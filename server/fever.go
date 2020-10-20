package server

import (
	"fmt"
	"github.com/nkanaev/yarr/storage"
	"net/http"
	"strings"
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

}

func FeverItemsHandler(rw http.ResponseWriter, req *http.Request) {

}

func FeverLinksHandler(rw http.ResponseWriter, req *http.Request) {

}

func FeverMarkHandler(rw http.ResponseWriter, req *http.Request) {

}
