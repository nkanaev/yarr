package server

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/nkanaev/yarr/src/server/auth"
	"github.com/nkanaev/yarr/src/server/router"
	"github.com/nkanaev/yarr/src/storage"
)

type FeverGroup struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

type FeverFeedsGroup struct {
	GroupID int64  `json:"group_id"`
	FeedIDs string `json:"feed_ids"`
}

type FeverFeed struct {
	ID          int64  `json:"id"`
	FaviconID   int64  `json:"favicon_id"`
	Title       string `json:"title"`
	Url         string `json:"url"`
	SiteUrl     string `json:"site_url"`
	IsSpark     int    `json:"is_spark"`
	LastUpdated int64  `json:"last_updated_on_time"`
}

type FeverItem struct {
	ID        int64  `json:"id"`
	FeedID    int64  `json:"feed_id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	HTML      string `json:"html"`
	Url       string `json:"url"`
	IsSaved   int    `json:"is_saved"`
	IsRead    int    `json:"is_read"`
	CreatedAt int64  `json:"created_on_time"`
}

type FeverFavicon struct {
	ID   int64  `json:"id"`
	Data string `json:"data"`
}

func writeFeverJSON(c *router.Context, data map[string]interface{}, lastRefreshed int64) {
	data["api_version"] = 1
	data["auth"] = 1
	data["last_refreshed_on_time"] = lastRefreshed
	c.JSON(http.StatusOK, data)
}

func getLastRefreshedOnTime(httpStates map[int64]storage.HTTPState) int64 {
	if len(httpStates) == 0 {
		return 0
	}

	var lastRefreshed int64
	for _, state := range httpStates {
		if state.LastRefreshed.Unix() > lastRefreshed {
			lastRefreshed = state.LastRefreshed.Unix()
		}
	}
	return lastRefreshed
}

func (s *Server) feverAuth(c *router.Context) bool {
	if s.Username != "" && s.Password != "" {
		apiKey := c.Req.FormValue("api_key")
		md5HashValue := md5.Sum([]byte(fmt.Sprintf("%s:%s", s.Username, s.Password)))
		hexMD5HashValue := fmt.Sprintf("%x", md5HashValue[:])
		if auth.StringsEqual(apiKey, hexMD5HashValue) {
			return false
		}
	}
	return true
}

func formHasValue(values url.Values, value string) bool {
	if _, ok := values[value]; ok {
		return true
	}
	return false
}

func (s *Server) handleFever(c *router.Context) {
	c.Req.ParseForm()
	if !s.feverAuth(c) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"api_version":            1,
			"auth":                   0,
			"last_refreshed_on_time": 0,
		})
		return
	}

	switch {
	case formHasValue(c.Req.Form, "groups"):
		s.feverGroupsHandler(c)
	case formHasValue(c.Req.Form, "feeds"):
		s.feverFeedsHandler(c)
	case formHasValue(c.Req.Form, "unread_item_ids"):
		s.feverUnreadItemIDsHandler(c)
	case formHasValue(c.Req.Form, "saved_item_ids"):
		s.feverSavedItemIDsHandler(c)
	case formHasValue(c.Req.Form, "favicons"):
		s.feverFaviconsHandler(c)
	case formHasValue(c.Req.Form, "items"):
		s.feverItemsHandler(c)
	case formHasValue(c.Req.Form, "links"):
		s.feverLinksHandler(c)
	case formHasValue(c.Req.Form, "mark"):
		s.feverMarkHandler(c)
	default:
		c.JSON(http.StatusOK, map[string]interface{}{
			"api_version":            1,
			"auth":                   1,
			"last_refreshed_on_time": 0,
		})
	}
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

func (s *Server) feverGroupsHandler(c *router.Context) {
	folders := s.db.ListFolders()
	groups := make([]*FeverGroup, len(folders))
	for i, folder := range folders {
		groups[i] = &FeverGroup{ID: folder.Id, Title: folder.Title}
	}
	writeFeverJSON(c, map[string]interface{}{
		"groups":       groups,
		"feeds_groups": feedGroups(s.db),
	}, getLastRefreshedOnTime(s.db.ListHTTPStates()))
}

func (s *Server) feverFeedsHandler(c *router.Context) {
	feeds := s.db.ListFeeds()
	httpStates := s.db.ListHTTPStates()

	feverFeeds := make([]*FeverFeed, len(feeds))
	for i, feed := range feeds {
		var lastUpdated int64
		if state, ok := httpStates[feed.Id]; ok {
			lastUpdated = state.LastRefreshed.Unix()
		}
		feverFeeds[i] = &FeverFeed{
			ID:          feed.Id,
			FaviconID:   feed.Id,
			Title:       feed.Title,
			Url:         feed.FeedLink,
			SiteUrl:     feed.Link,
			IsSpark:     0,
			LastUpdated: lastUpdated,
		}
	}
	writeFeverJSON(c, map[string]interface{}{
		"feeds":        feverFeeds,
		"feeds_groups": feedGroups(s.db),
	}, getLastRefreshedOnTime(httpStates))
}

func (s *Server) feverFaviconsHandler(c *router.Context) {
	feeds := s.db.ListFeeds()
	favicons := make([]*FeverFavicon, len(feeds))
	for i, feed := range feeds {
		data := "data:image/gif;base64,R0lGODlhAQABAAAAACw="
		if feed.HasIcon {
			icon := s.db.GetFeed(feed.Id).Icon
			data = fmt.Sprintf(
				"data:%s;base64,%s",
				http.DetectContentType(*icon),
				base64.StdEncoding.EncodeToString(*icon),
			)
		}
		favicons[i] = &FeverFavicon{ID: feed.Id, Data: data}
	}

	writeFeverJSON(c, map[string]interface{}{
		"favicons": favicons,
	}, getLastRefreshedOnTime(s.db.ListHTTPStates()))
}

// for memory pressure reasons, we only return a limited number of items
// documented at https://github.com/DigitalDJ/tinytinyrss-fever-plugin/blob/master/fever-api.md#items
const listLimit = 50

func (s *Server) feverItemsHandler(c *router.Context) {
	filter := storage.ItemFilter{}
	query := c.Req.URL.Query()

	switch {
	case query.Get("with_ids") != "":
		ids := make([]int64, 0)
		for _, idstr := range strings.Split(query.Get("with_ids"), ",") {
			if idnum, err := strconv.ParseInt(idstr, 10, 64); err == nil {
				ids = append(ids, idnum)
			}
		}
		filter.IDs = &ids
	case query.Get("since_id") != "":
		idstr := query.Get("since_id")
		if idnum, err := strconv.ParseInt(idstr, 10, 64); err == nil {
			filter.SinceID = &idnum
		}
	case query.Get("max_id") != "":
		idstr := query.Get("max_id")
		if idnum, err := strconv.ParseInt(idstr, 10, 64); err == nil {
			filter.MaxID = &idnum
		}
	}

	items := s.db.ListItems(filter, listLimit, true, true)

	feverItems := make([]FeverItem, len(items))
	for i, item := range items {
		date := item.Date
		time := date.Unix()

		isSaved := 0
		if item.Status == storage.STARRED {
			isSaved = 1
		}
		isRead := 0
		if item.Status == storage.READ {
			isRead = 1
		}
		feverItems[i] = FeverItem{
			ID:        item.Id,
			FeedID:    item.FeedId,
			Title:     item.Title,
			Author:    "",
			HTML:      item.Content,
			Url:       item.Link,
			IsSaved:   isSaved,
			IsRead:    isRead,
			CreatedAt: time,
		}
	}

	writeFeverJSON(c, map[string]interface{}{
		"items": feverItems,
	}, getLastRefreshedOnTime(s.db.ListHTTPStates()))
}

func (s *Server) feverLinksHandler(c *router.Context) {
	writeFeverJSON(c, map[string]interface{}{
		"links": make([]interface{}, 0),
	}, getLastRefreshedOnTime(s.db.ListHTTPStates()))
}

func (s *Server) feverUnreadItemIDsHandler(c *router.Context) {
	status := storage.UNREAD
	itemIds := make([]int64, 0)

	itemFilter := storage.ItemFilter{
		Status: &status,
	}
	for {
		items := s.db.ListItems(itemFilter, listLimit, true, false)
		if len(items) == 0 {
			break
		}
		for _, item := range items {
			itemIds = append(itemIds, item.Id)
		}
		itemFilter.After = &items[len(items)-1].Id
	}
	writeFeverJSON(c, map[string]interface{}{
		"unread_item_ids": joinInts(itemIds),
	}, getLastRefreshedOnTime(s.db.ListHTTPStates()))
}

func (s *Server) feverSavedItemIDsHandler(c *router.Context) {
	status := storage.STARRED
	itemIds := make([]int64, 0)

	itemFilter := storage.ItemFilter{
		Status: &status,
	}
	for {
		items := s.db.ListItems(itemFilter, listLimit, true, false)
		if len(items) == 0 {
			break
		}
		for _, item := range items {
			itemIds = append(itemIds, item.Id)
		}
		itemFilter.After = &items[len(items)-1].Id
	}
	writeFeverJSON(c, map[string]interface{}{
		"saved_item_ids": joinInts(itemIds),
	}, getLastRefreshedOnTime(s.db.ListHTTPStates()))
}

func (s *Server) feverMarkHandler(c *router.Context) {
	id, err := strconv.ParseInt(c.Req.Form.Get("id"), 10, 64)
	if err != nil {
		log.Print("invalid id:", err)
		return
	}

	switch c.Req.Form.Get("mark") {
	case "item":
		var status storage.ItemStatus
		switch c.Req.Form.Get("as") {
		case "read":
			status = storage.READ
		case "unread":
			status = storage.UNREAD
		case "saved":
			status = storage.STARRED
		case "unsaved":
			status = storage.READ
		default:
			c.Out.WriteHeader(http.StatusBadRequest)
			return
		}
		s.db.UpdateItemStatus(id, status)
	case "feed":
		if c.Req.Form.Get("as") != "read" {
			c.Out.WriteHeader(http.StatusBadRequest)
		}
		markFilter := storage.MarkFilter{FeedID: &id}
		x, _ := strconv.ParseInt(c.Req.Form.Get("before"), 10, 64)
		if x > 0 {
			before := time.Unix(x, 0)
			markFilter.Before = &before
		}
		s.db.MarkItemsRead(markFilter)
	case "group":
		if c.Req.Form.Get("as") != "read" {
			c.Out.WriteHeader(http.StatusBadRequest)
		}
		markFilter := storage.MarkFilter{FolderID: &id}
		x, _ := strconv.ParseInt(c.Req.Form.Get("before"), 10, 64)
		if x > 0 {
			before := time.Unix(x, 0)
			markFilter.Before = &before
		}
		s.db.MarkItemsRead(markFilter)
	default:
		c.Out.WriteHeader(http.StatusBadRequest)
		return
	}
}
