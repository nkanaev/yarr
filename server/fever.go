package server

import (
	"encoding/base64"
	"fmt"
	"github.com/nkanaev/yarr/storage"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var feverHandlers = map[string]func(rw http.ResponseWriter, req *http.Request){
	"groups":          FeverGroupsHandler,
	"feeds":           FeverFeedsHandler,
	"unread_item_ids": FeverFilteredItemIDsHandler,
	"saved_item_ids":  FeverFilteredItemIDsHandler,

	"favicons": FeverFaviconsHandler,
	"items":    FeverItemsHandler,
	"links":    FeverLinksHandler,
	"mark":     FeverMarkHandler,
}

type FeverGroup struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

type FeverFeedsGroup struct {
	GroupID int64  `json:"group_id"`
	FeedIDs string `json:"feed_ids"`
}

type FeverFeed struct {
	ID                int64  `json:"id"`
	FaviconID         int64  `json:"favicon_id"`
	Title             string `json:"title"`
	Url               string `json:"url"`
	SiteUrl           string `json:"site_url"`
	IsSpark           int    `json:"is_spark"`
	LastUpdatedOnTime int64  `json:"last_updated_on_time"`
}

type FeverItem struct {
	ID            int64  `json:"id"`
	FeedID        int64  `json:"feed_id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	HTML          string `json:"html"`
	Url           string `json:"url"`
	IsSaved       int    `json:"is_saved"`
	IsRead        int    `json:"is_read"`
	CreatedOnTime int64  `json:"created_on_time"`
}

type FeverFavicon struct {
	ID   int64  `json:"id"`
	Data string `json:"data"`
}

func writeFeverJSON(rw http.ResponseWriter, data map[string]interface{}) {
	data["api_version"] = 1
	data["auth"] = 1
	writeJSON(rw, data)
}

func FeverHandler(rw http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	fmt.Println(req.URL.String())
	for key, handler := range feverHandlers {
		if _, ok := query[key]; ok {
			handler(rw, req)
			return
		}
	}
	writeJSON(rw, map[string]interface{}{
		"api_version": 1,
		"auth":        1,
	})
}

func joinInts(values []int64) string {
	var result strings.Builder
	for i, val := range values {
		fmt.Fprintf(&result, "%d", val)
		if i != len(values)-1 {
			result.WriteString(",")
		}
	}
	return result.String()
}

func feedGroups(db *storage.Storage) []*FeverFeedsGroup {
	feeds := db.ListFeeds()

	groupFeeds := make(map[int64][]int64)
	for _, feed := range feeds {
		// TODO: what about top-level feeds?
		if feed.FolderId == nil {
			continue
		}
		groupFeeds[*feed.FolderId] = append(groupFeeds[*feed.FolderId], feed.Id)
	}
	result := make([]*FeverFeedsGroup, 0)
	for groupId, feedIds := range groupFeeds {
		result = append(result, &FeverFeedsGroup{
			GroupID: groupId,
			FeedIDs: joinInts(feedIds),
		})
	}
	return result
}

func FeverGroupsHandler(rw http.ResponseWriter, req *http.Request) {
	folders := db(req).ListFolders()
	groups := make([]*FeverGroup, len(folders))
	for i, folder := range folders {
		groups[i] = &FeverGroup{ID: folder.Id, Title: folder.Title}
	}
	writeFeverJSON(rw, map[string]interface{}{
		"groups":       groups,
		"feeds_groups": feedGroups(db(req)),
	})
}

func FeverFeedsHandler(rw http.ResponseWriter, req *http.Request) {
	feeds := db(req).ListFeeds()

	feverFeeds := make([]*FeverFeed, len(feeds))
	for i, feed := range feeds {
		// TODO: check url/siteurl
		// TODO: store last updated on time?
		feverFeeds[i] = &FeverFeed{
			ID:                feed.Id,
			FaviconID:         feed.Id,
			Title:             feed.Title,
			Url:               feed.FeedLink,
			SiteUrl:           feed.Link,
			IsSpark:           0,
			LastUpdatedOnTime: 1,
		}
	}
	writeFeverJSON(rw, map[string]interface{}{
		"feeds":        feverFeeds,
		"feeds_groups": feedGroups(db(req)),
	})
}

func FeverFilteredItemIDsHandler(rw http.ResponseWriter, req *http.Request) {
	var status storage.ItemStatus
	var filter string
	if _, ok := req.URL.Query()["unread_item_ids"]; ok {
		status = storage.UNREAD
		filter = "unread_item_ids"
	} else {
		status = storage.STARRED
		filter = "saved_item_ids"
	}

	itemIds := make([]int64, 0, 4000)
	batch := 1000
	index := 0
	for {
		items := db(req).ListItems(storage.ItemFilter{Status: &status}, index*batch, batch, true)
		if len(items) == 0 {
			break
		}
		for _, item := range items {
			itemIds = append(itemIds, item.Id)
		}
		index += 1
	}
	writeFeverJSON(rw, map[string]interface{}{
		filter: joinInts(itemIds),
	})
}

func FeverFaviconsHandler(rw http.ResponseWriter, req *http.Request) {
	feeds := db(req).ListFeeds()
	favicons := make([]*FeverFavicon, len(feeds))
	for i, feed := range feeds {
		data := "data:image/gif;base64,R0lGODlhAQABAAAAACw="
		if feed.HasIcon {
			icon := db(req).GetFeed(feed.Id).Icon
			data = fmt.Sprintf(
				"data:%s;base64,%s",
				http.DetectContentType(*icon),
				base64.StdEncoding.EncodeToString(*icon),
			)
		}
		favicons[i] = &FeverFavicon{ID: feed.Id, Data: data}
	}

	writeFeverJSON(rw, map[string]interface{}{
		"favicons": favicons,
	})
}

func FeverItemsHandler(rw http.ResponseWriter, req *http.Request) {
	filter := storage.ItemFilter{}
	query := req.URL.Query()
	// TODO: must be switch case?
	if _, ok := query["with_ids"]; ok {
		ids := make([]int64, 0)
		for _, idstr := range strings.Split(query.Get("with_ids"), ",") {
			if idnum, err := strconv.ParseInt(idstr, 10, 64); err == nil {
				ids = append(ids, idnum)
			}
		}
		filter.IDs = &ids
	}

	if _, ok := query["since_id"]; ok {
		idstr := query.Get("since_id")
		if idnum, err := strconv.ParseInt(idstr, 10, 64); err == nil {
			filter.SinceID = &idnum
		}
	}

	items := db(req).ListItems(filter, 0, 50, true)

	feverItems := make([]FeverItem, len(items))
	for i, item := range items {
		date := item.Date
		if date == nil {
			date = item.DateUpdated
		}
		time := int64(0)
		if date != nil {
			time = date.UnixNano() / 1000_000_000
		}

		isSaved := 0
		if item.Status == storage.STARRED {
			isSaved = 1
		}
		isRead := 0
		if item.Status == storage.READ {
			isRead = 1
		}
		feverItems[i] = FeverItem{
			ID:            item.Id,
			FeedID:        item.FeedId,
			Title:         item.Title,
			Author:        item.Author,
			HTML:          item.Content,
			Url:           item.Link,
			IsSaved:       isSaved,
			IsRead:        isRead,
			CreatedOnTime: time,
		}
	}

	writeFeverJSON(rw, map[string]interface{}{
		"items": feverItems,
	})
}

func FeverLinksHandler(rw http.ResponseWriter, req *http.Request) {
	writeFeverJSON(rw, map[string]interface{}{
		"links": make([]interface{}, 0),
	})
}

func FeverMarkHandler(rw http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()

	id, err := strconv.ParseInt(query.Get("id"), 10, 0)
	if err != nil {
		handler(req).log.Print("invalid id:", err)
		return
	}

	switch query.Get("mark") {
	case "item":
		var status storage.ItemStatus
		switch query.Get("as") {
		case "read":
			status = storage.READ
		case "unread":
			status = storage.UNREAD
		case "saved":
			status = storage.STARRED
		case "unsaved":
			status = storage.READ
		default:
			fmt.Println("TODO: handle")
		}
		db(req).UpdateItemStatus(id, status)
	case "feed":
		x, _ := strconv.ParseInt(query.Get("before"), 10, 0)
		before := time.Unix(x, 0)
		db(req).MarkItemsRead(storage.MarkFilter{FeedID: &id, Before: &before})
	case "group":
		x, _ := strconv.ParseInt(query.Get("before"), 10, 0)
		before := time.Unix(x, 0)
		db(req).MarkItemsRead(storage.MarkFilter{FolderID: &id, Before: &before})
	default:
		fmt.Println("TODO: handle")
	}
}
