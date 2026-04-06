package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/nkanaev/yarr/src/assets"
	"github.com/nkanaev/yarr/src/content/htmlutil"
	"github.com/nkanaev/yarr/src/content/sanitizer"
	"github.com/nkanaev/yarr/src/server/router"
	"github.com/nkanaev/yarr/src/storage"
)

// Template data structures for server-rendered partials

type FeedInFolder struct {
	Feed      storage.Feed
	FeedCount int64
	Error     string
	Hidden    bool
}

type FolderWithFeeds struct {
	Folder      storage.Folder
	Feeds       []FeedInFolder
	FolderCount int64
	Hidden      bool
}

type FeedListData struct {
	FoldersWithFeeds []FolderWithFeeds
	FeedSelected     string
	Filter           string
	TotalCount       int64
	Loading          int32
}

type ItemViewData struct {
	Id          int64
	Title       string
	Link        string
	FeedTitle   string
	FolderTitle string
	StatusStr   string
	DateRepr    string
	FeedId      int64
}

type ItemListData struct {
	Items        []ItemViewData
	HasMore      bool
	NextPageVals template.HTMLAttr
	FeedError    string
	SearchActive bool
	Search       string
}

type MediaLinkView struct {
	URL         string
	Description template.HTML
}

type ItemContentData struct {
	Item          ItemViewData
	FeedTitle     string
	FormattedDate string
	Content       template.HTML
	Images        []MediaLinkView
	Audios        []MediaLinkView
	Videos        []MediaLinkView
	ThemeFont     string
	ThemeSize     float64
}

// Helper functions

func dateRepr(d time.Time) string {
	sec := time.Since(d).Seconds()
	neg := sec < 0
	if neg {
		sec = -sec
	}
	var out string
	switch {
	case sec < 2700: // 45 min
		out = fmt.Sprintf("%dm", int(math.Round(sec/60)))
	case sec < 86400: // 24h
		out = fmt.Sprintf("%dh", int(math.Round(sec/3600)))
	case sec < 604800: // 7d
		out = fmt.Sprintf("%dd", int(math.Round(sec/86400)))
	default:
		out = d.Format("January 2, 2006")
	}
	if neg {
		return "-" + out
	}
	return out
}

func statusStr(s storage.ItemStatus) string {
	return storage.StatusRepresentations[s]
}

func (s *Server) buildFeedListData(settings map[string]interface{}) FeedListData {
	feeds := s.db.ListFeeds()
	folders := s.db.ListFolders()
	stats := s.db.FeedStats()
	errors := s.db.GetFeedErrors()

	filter := ""
	if f, ok := settings["filter"].(string); ok {
		filter = f
	}
	feedSelected := ""
	if f, ok := settings["feed"].(string); ok {
		feedSelected = f
	}

	// Build stats maps
	statsByFeed := make(map[int64]storage.FeedStat)
	for _, stat := range stats {
		statsByFeed[stat.FeedId] = stat
	}

	// Group feeds by folder
	feedsByFolder := make(map[int64][]storage.Feed)
	var orphanFeeds []storage.Feed
	for _, feed := range feeds {
		if feed.FolderId != nil {
			feedsByFolder[*feed.FolderId] = append(feedsByFolder[*feed.FolderId], feed)
		} else {
			orphanFeeds = append(orphanFeeds, feed)
		}
	}

	var totalCount int64

	var foldersWithFeeds []FolderWithFeeds

	for _, folder := range folders {
		fwf := FolderWithFeeds{Folder: folder}
		var folderCount int64

		for _, feed := range feedsByFolder[folder.Id] {
			fif := FeedInFolder{Feed: feed}
			if stat, ok := statsByFeed[feed.Id]; ok {
				switch filter {
				case "unread":
					fif.FeedCount = stat.UnreadCount
				case "starred":
					fif.FeedCount = stat.StarredCount
				}
			}
			if filter == "archived" {
				if feed.Archived {
					fif.FeedCount = 1
				}
			}
			if err, ok := errors[feed.Id]; ok {
				fif.Error = err
			}

			// Hide logic
			if filter != "" && filter != "archived" && feed.Archived {
				fif.Hidden = true
			}
			if filter == "archived" && !feed.Archived {
				fif.Hidden = true
			}
			if filter != "" && !fif.Hidden && fif.FeedCount == 0 {
				isCurrent := feedSelected == fmt.Sprintf("feed:%d", feed.Id)
				if !isCurrent {
					fif.Hidden = true
				}
			}

			folderCount += fif.FeedCount
			fwf.Feeds = append(fwf.Feeds, fif)
		}

		fwf.FolderCount = folderCount
		totalCount += folderCount

		// Hide folder if empty in filtered mode
		if filter != "" && folderCount == 0 {
			isCurrent := feedSelected == fmt.Sprintf("folder:%d", folder.Id)
			if !isCurrent {
				fwf.Hidden = true
			}
		}

		foldersWithFeeds = append(foldersWithFeeds, fwf)
	}

	// Orphan feeds (no folder)
	if len(orphanFeeds) > 0 {
		orphan := FolderWithFeeds{Folder: storage.Folder{}} // Id=0 means no folder
		var orphanCount int64
		for _, feed := range orphanFeeds {
			fif := FeedInFolder{Feed: feed}
			if stat, ok := statsByFeed[feed.Id]; ok {
				switch filter {
				case "unread":
					fif.FeedCount = stat.UnreadCount
				case "starred":
					fif.FeedCount = stat.StarredCount
				}
			}
			if filter == "archived" {
				if feed.Archived {
					fif.FeedCount = 1
				}
			}
			if err, ok := errors[feed.Id]; ok {
				fif.Error = err
			}
			if filter != "" && filter != "archived" && feed.Archived {
				fif.Hidden = true
			}
			if filter == "archived" && !feed.Archived {
				fif.Hidden = true
			}
			if filter != "" && !fif.Hidden && fif.FeedCount == 0 {
				isCurrent := feedSelected == fmt.Sprintf("feed:%d", feed.Id)
				if !isCurrent {
					fif.Hidden = true
				}
			}

			orphanCount += fif.FeedCount
			orphan.Feeds = append(orphan.Feeds, fif)
		}
		orphan.FolderCount = orphanCount
		totalCount += orphanCount
		foldersWithFeeds = append(foldersWithFeeds, orphan)
	}

	return FeedListData{
		FoldersWithFeeds: foldersWithFeeds,
		FeedSelected:     feedSelected,
		Filter:           filter,
		TotalCount:       totalCount,
		Loading:          s.worker.FeedsPending(),
	}
}

func (s *Server) buildItemListData(c *router.Context) ItemListData {
	perPage := 20
	query := c.Req.URL.Query()

	filter := storage.ItemFilter{}
	if folderID, err := c.QueryInt64("folder_id"); err == nil {
		filter.FolderID = &folderID
	}
	if feedID, err := c.QueryInt64("feed_id"); err == nil {
		filter.FeedID = &feedID
	}
	if after, err := c.QueryInt64("after"); err == nil {
		filter.After = &after
	}
	if status := query.Get("status"); len(status) != 0 {
		statusValue := storage.StatusValues[status]
		filter.Status = &statusValue
	}
	if search := query.Get("search"); len(search) != 0 {
		filter.Search = &search
	}
	newestFirst := query.Get("oldest_first") != "true"

	items := s.db.ListItems(filter, perPage+1, newestFirst, false)
	hasMore := false
	if len(items) == perPage+1 {
		hasMore = true
		items = items[:perPage]
	}

	// Build feed title and folder title maps
	feeds := s.db.ListFeeds()
	feedMap := make(map[int64]string)
	feedFolderMap := make(map[int64]int64)
	for _, f := range feeds {
		feedMap[f.Id] = f.Title
		if f.FolderId != nil {
			feedFolderMap[f.Id] = *f.FolderId
		}
	}
	folders := s.db.ListFolders()
	folderMap := make(map[int64]string)
	for _, f := range folders {
		folderMap[f.Id] = f.Title
	}

	var viewItems []ItemViewData
	for _, item := range items {
		title := item.Title
		if title == "" {
			text := htmlutil.ExtractText(item.Content)
			title = htmlutil.TruncateText(text, 140)
		}
		folderTitle := ""
		if fid, ok := feedFolderMap[item.FeedId]; ok {
			folderTitle = folderMap[fid]
		}
		viewItems = append(viewItems, ItemViewData{
			Id:          item.Id,
			Title:       title,
			FeedTitle:   feedMap[item.FeedId],
			FolderTitle: folderTitle,
			StatusStr:   statusStr(item.Status),
			DateRepr:    dateRepr(item.Date),
			FeedId:      item.FeedId,
		})
	}

	// Build next page vals
	var nextPageVals string
	if hasMore && len(items) > 0 {
		lastItem := items[len(items)-1]
		vals := make(map[string]string)
		vals["after"] = fmt.Sprintf("%d", lastItem.Id)
		if filter.FeedID != nil {
			vals["feed_id"] = fmt.Sprintf("%d", *filter.FeedID)
		}
		if filter.FolderID != nil {
			vals["folder_id"] = fmt.Sprintf("%d", *filter.FolderID)
		}
		if filter.Status != nil {
			vals["status"] = statusStr(*filter.Status)
		}
		if filter.Search != nil {
			vals["search"] = *filter.Search
		}
		if !newestFirst {
			vals["oldest_first"] = "true"
		}
		valsJSON, _ := json.Marshal(vals)
		nextPageVals = string(valsJSON)
	}

	// Feed error
	var feedError string
	if filter.FeedID != nil {
		errors := s.db.GetFeedErrors()
		if e, ok := errors[*filter.FeedID]; ok {
			feedError = e
		}
	}

	searchActive := filter.Search != nil && *filter.Search != ""
	searchStr := ""
	if searchActive {
		searchStr = *filter.Search
	}

	return ItemListData{
		Items:        viewItems,
		HasMore:      hasMore,
		NextPageVals: template.HTMLAttr(nextPageVals),
		FeedError:    feedError,
		SearchActive: searchActive,
		Search:       searchStr,
	}
}

// Partial handlers

func (s *Server) handlePartialFeedList(c *router.Context) {
	settings := s.db.GetSettings()
	data := s.buildFeedListData(settings)
	s.renderPartial(c.Out, "feed_list", data)
}

func (s *Server) handlePartialItems(c *router.Context) {
	data := s.buildItemListData(c)
	s.renderPartial(c.Out, "item_list", data)
}

func (s *Server) handlePartialItemContent(c *router.Context) {
	id, err := c.VarInt64("id")
	if err != nil {
		c.Out.WriteHeader(http.StatusBadRequest)
		return
	}

	if c.Req.Method == "PUT" {
		var body ItemUpdateForm
		if err := json.NewDecoder(c.Req.Body).Decode(&body); err != nil {
			log.Print(err)
			c.Out.WriteHeader(http.StatusBadRequest)
			return
		}
		if body.Status != nil {
			s.db.UpdateItemStatus(id, *body.Status)
		}
	}

	item := s.db.GetItem(id)
	if item == nil {
		c.Out.WriteHeader(http.StatusNotFound)
		return
	}

	// Mark as read on GET
	if c.Req.Method == "GET" && item.Status == storage.UNREAD {
		s.db.UpdateItemStatus(id, storage.READ)
		item.Status = storage.READ
	}

	// Fix relative links
	feed := s.db.GetFeed(item.FeedId)
	feedTitle := ""
	if feed != nil {
		feedTitle = feed.Title
		if !htmlutil.IsAPossibleLink(item.Link) {
			item.Link = htmlutil.AbsoluteUrl(item.Link, feed.Link)
		}
	}

	sanitizedContent := sanitizer.Sanitize(item.Link, item.Content)

	settings := s.db.GetSettings()
	themeFont := ""
	if f, ok := settings["theme_font"].(string); ok {
		themeFont = f
	}
	themeSize := 1.0
	if sz, ok := settings["theme_size"].(float64); ok {
		themeSize = sz
	}

	var images, audios, videos []MediaLinkView
	for _, link := range item.MediaLinks {
		mlv := MediaLinkView{
			URL:         link.URL,
			Description: template.HTML(sanitizer.Sanitize(item.Link, link.Description)),
		}
		switch link.Type {
		case "image":
			images = append(images, mlv)
		case "audio":
			audios = append(audios, mlv)
		case "video":
			videos = append(videos, mlv)
		}
	}

	data := ItemContentData{
		Item: ItemViewData{
			Id:        item.Id,
			Title:     item.Title,
			Link:      item.Link,
			StatusStr: statusStr(item.Status),
			FeedId:    item.FeedId,
		},
		FeedTitle:     feedTitle,
		FormattedDate: item.Date.Format("January 2, 2006 3:04 PM"),
		Content:       template.HTML(sanitizedContent),
		Images:        images,
		Audios:        audios,
		Videos:        videos,
		ThemeFont:     themeFont,
		ThemeSize:     themeSize,
	}
	// Also set the item Link on the view data for template use
	data.Item.Title = item.Title
	// We need Link accessible in template — add it to ItemViewData or use a wrapper
	s.renderPartialWithLink(c.Out, "item_content", data, item.Link)
}

func (s *Server) renderPartial(w http.ResponseWriter, name string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl := assets.PartialTemplate()
	if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
		log.Printf("Error rendering partial %s: %v", name, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) renderPartialWithLink(w http.ResponseWriter, name string, data ItemContentData, link string) {
	// Wrap data to include link for template
	wrapped := struct {
		ItemContentData
		Link string
	}{data, link}
	_ = wrapped
	// For now just render the content data directly - Link is on Item
	s.renderPartial(w, name, data)
}

// handlePartialItemUpdate handles PUT /partials/items/:id for status changes
// Returns an OOB swap for the item row + updated content
func (s *Server) handlePartialItemUpdate(c *router.Context) {
	id, err := c.VarInt64("id")
	if err != nil {
		c.Out.WriteHeader(http.StatusBadRequest)
		return
	}

	var body ItemUpdateForm
	if err := json.NewDecoder(c.Req.Body).Decode(&body); err != nil {
		log.Print(err)
		c.Out.WriteHeader(http.StatusBadRequest)
		return
	}
	if body.Status != nil {
		s.db.UpdateItemStatus(id, *body.Status)
	}

	// Re-render the item content with updated status
	s.handlePartialItemContent(c)
}

// Feed menu data
type FeedMenuFeed struct {
	Id          int64
	Title       string
	Link        string
	FeedLink    string
	Archived    bool
	FolderId    *int64
	FolderIdVal int64 // 0 if nil, for template comparison
}

type FeedMenuData struct {
	Type    string // "feed", "folder", or ""
	Feed    FeedMenuFeed
	Folder  storage.Folder
	Folders []storage.Folder
}

func (s *Server) handlePartialFeedMenu(c *router.Context) {
	selection := c.Req.URL.Query().Get("sel")
	if selection == "" {
		c.Out.WriteHeader(http.StatusBadRequest)
		return
	}

	data := FeedMenuData{}
	parts := strings.SplitN(selection, ":", 2)
	if len(parts) != 2 {
		c.Out.WriteHeader(http.StatusBadRequest)
		return
	}

	switch parts[0] {
	case "feed":
		id, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			c.Out.WriteHeader(http.StatusBadRequest)
			return
		}
		feed := s.db.GetFeed(id)
		if feed == nil {
			c.Out.WriteHeader(http.StatusNotFound)
			return
		}
		data.Type = "feed"
		data.Feed = FeedMenuFeed{
			Id:       feed.Id,
			Title:    feed.Title,
			Link:     feed.Link,
			FeedLink: feed.FeedLink,
			Archived: feed.Archived,
			FolderId: feed.FolderId,
		}
		if feed.FolderId != nil {
			data.Feed.FolderIdVal = *feed.FolderId
		}
		data.Folders = s.db.ListFolders()
	case "folder":
		id, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			c.Out.WriteHeader(http.StatusBadRequest)
			return
		}
		folders := s.db.ListFolders()
		for _, f := range folders {
			if f.Id == id {
				data.Type = "folder"
				data.Folder = f
				break
			}
		}
		if data.Type == "" {
			c.Out.WriteHeader(http.StatusNotFound)
			return
		}
	}

	s.renderPartial(c.Out, "feed_menu", data)
}

// buildInitialItemList builds the item list for the initial page load
func (s *Server) buildInitialItemList(settings map[string]interface{}) ItemListData {
	perPage := 20
	filter := storage.ItemFilter{}

	feedSelected := ""
	if f, ok := settings["feed"].(string); ok {
		feedSelected = f
	}

	// If no feed selected, return empty
	if feedSelected == "" {
		return ItemListData{}
	}

	// Parse feed selection
	parts := strings.SplitN(feedSelected, ":", 2)
	if len(parts) == 2 {
		switch parts[0] {
		case "feed":
			if id, err := strconv.ParseInt(parts[1], 10, 64); err == nil {
				filter.FeedID = &id
			}
		case "folder":
			if id, err := strconv.ParseInt(parts[1], 10, 64); err == nil {
				filter.FolderID = &id
			}
		}
	}

	if filterStr, ok := settings["filter"].(string); ok && filterStr != "" {
		if sv, exists := storage.StatusValues[filterStr]; exists {
			filter.Status = &sv
		}
	}

	newestFirst := true
	if sn, ok := settings["sort_newest_first"].(bool); ok {
		newestFirst = sn
	}

	items := s.db.ListItems(filter, perPage+1, newestFirst, false)
	hasMore := false
	if len(items) == perPage+1 {
		hasMore = true
		items = items[:perPage]
	}

	feeds := s.db.ListFeeds()
	feedMap := make(map[int64]string)
	feedFolderMap := make(map[int64]int64)
	for _, f := range feeds {
		feedMap[f.Id] = f.Title
		if f.FolderId != nil {
			feedFolderMap[f.Id] = *f.FolderId
		}
	}
	folders := s.db.ListFolders()
	folderMap := make(map[int64]string)
	for _, f := range folders {
		folderMap[f.Id] = f.Title
	}

	var viewItems []ItemViewData
	for _, item := range items {
		title := item.Title
		if title == "" {
			title = "untitled"
		}
		folderTitle := ""
		if fid, ok := feedFolderMap[item.FeedId]; ok {
			folderTitle = folderMap[fid]
		}
		viewItems = append(viewItems, ItemViewData{
			Id:          item.Id,
			Title:       title,
			FeedTitle:   feedMap[item.FeedId],
			FolderTitle: folderTitle,
			StatusStr:   statusStr(item.Status),
			DateRepr:    dateRepr(item.Date),
			FeedId:      item.FeedId,
		})
	}

	var nextPageVals string
	if hasMore && len(items) > 0 {
		lastItem := items[len(items)-1]
		vals := make(map[string]string)
		vals["after"] = fmt.Sprintf("%d", lastItem.Id)
		if filter.FeedID != nil {
			vals["feed_id"] = fmt.Sprintf("%d", *filter.FeedID)
		}
		if filter.FolderID != nil {
			vals["folder_id"] = fmt.Sprintf("%d", *filter.FolderID)
		}
		if filter.Status != nil {
			vals["status"] = statusStr(*filter.Status)
		}
		if !newestFirst {
			vals["oldest_first"] = "true"
		}
		valsJSON, _ := json.Marshal(vals)
		nextPageVals = string(valsJSON)
	}

	return ItemListData{
		Items:        viewItems,
		HasMore:      hasMore,
		NextPageVals: template.HTMLAttr(nextPageVals),
	}
}

// isHTMXRequest checks if request comes from HTMX
func isHTMXRequest(r *http.Request) bool {
	return strings.ToLower(r.Header.Get("HX-Request")) == "true"
}
