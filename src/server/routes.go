package server

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/nkanaev/yarr/src/assets"
	"github.com/nkanaev/yarr/src/content/htmlutil"
	"github.com/nkanaev/yarr/src/content/readability"
	"github.com/nkanaev/yarr/src/content/sanitizer"
	"github.com/nkanaev/yarr/src/content/silo"
	"github.com/nkanaev/yarr/src/server/auth"
	"github.com/nkanaev/yarr/src/server/gzip"
	"github.com/nkanaev/yarr/src/server/opml"
	"github.com/nkanaev/yarr/src/server/router"
	"github.com/nkanaev/yarr/src/storage"
	"github.com/nkanaev/yarr/src/worker"
)

func (s *Server) handler() http.Handler {
	r := router.NewRouter(s.BasePath)

	r.Use(gzip.Middleware)

	if s.Username != "" && s.Password != "" {
		a := &auth.Middleware{
			BasePath: s.BasePath,
			Username: s.Username,
			Password: s.Password,
			Public:   []string{"/static", "/fever"},
			DB:       s.db,
		}
		r.Use(a.Handler)
	}

	r.For("/", s.handleIndex)
	r.For("/manifest.json", s.handleManifest)
	r.For("/static/*path", s.handleStatic)
	r.For("/api/status", s.handleStatus)
	r.For("/api/folders", s.handleFolderList)
	r.For("/api/folders/:id", s.handleFolder)
	r.For("/api/feeds", s.handleFeedList)
	r.For("/api/feeds/refresh", s.handleFeedRefresh)
	r.For("/api/feeds/errors", s.handleFeedErrors)
	r.For("/api/feeds/:id/icon", s.handleFeedIcon)
	r.For("/api/feeds/:id", s.handleFeed)
	r.For("/api/items", s.handleItemList)
	r.For("/api/items/:id", s.handleItem)
	r.For("/api/settings", s.handleSettings)
	r.For("/opml/import", s.handleOPMLImport)
	r.For("/opml/export", s.handleOPMLExport)
	r.For("/page", s.handlePageCrawl)
	r.For("/logout", s.handleLogout)
	r.For("/fever/", s.handleFever)

	return r
}

func (s *Server) handleIndex(c *router.Context) {
	c.HTML(http.StatusOK, assets.Template("index.html"), map[string]interface{}{
		"settings":      s.db.GetSettings(),
		"authenticated": s.Username != "" && s.Password != "",
	})
}

func (s *Server) handleStatic(c *router.Context) {
	// don't serve templates
	dir, name := filepath.Split(c.Vars["path"])
	if dir == "" && strings.HasSuffix(name, ".html") {
		c.Out.WriteHeader(http.StatusNotFound)
		return
	}
	http.StripPrefix(s.BasePath+"/static/", http.FileServer(http.FS(assets.FS))).ServeHTTP(c.Out, c.Req)
}

func (s *Server) handleManifest(c *router.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"$schema":     "https://json.schemastore.org/web-manifest-combined.json",
		"name":        "yarr!",
		"short_name":  "yarr",
		"description": "yet another rss reader",
		"display":     "standalone",
		"start_url":   "/" + strings.TrimPrefix(s.BasePath, "/"),
		"icons": []map[string]interface{}{
			{
				"src":   s.BasePath + "/static/graphicarts/favicon.png",
				"sizes": "64x64",
				"type":  "image/png",
			},
		},
	})
}

func (s *Server) handleStatus(c *router.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"running": s.worker.FeedsPending(),
		"stats":   s.db.FeedStats(),
	})
}

func (s *Server) handleFolderList(c *router.Context) {
	if c.Req.Method == "GET" {
		list := s.db.ListFolders()
		c.JSON(http.StatusOK, list)
	} else if c.Req.Method == "POST" {
		var body FolderCreateForm
		if err := json.NewDecoder(c.Req.Body).Decode(&body); err != nil {
			log.Print(err)
			c.Out.WriteHeader(http.StatusBadRequest)
			return
		}
		if len(body.Title) == 0 {
			c.JSON(http.StatusBadRequest, map[string]string{"error": "Folder title missing."})
			return
		}
		folder := s.db.CreateFolder(body.Title)
		c.JSON(http.StatusCreated, folder)
	} else {
		c.Out.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleFolder(c *router.Context) {
	id, err := c.VarInt64("id")
	if err != nil {
		c.Out.WriteHeader(http.StatusBadRequest)
		return
	}
	if c.Req.Method == "PUT" {
		var body FolderUpdateForm
		if err := json.NewDecoder(c.Req.Body).Decode(&body); err != nil {
			log.Print(err)
			c.Out.WriteHeader(http.StatusBadRequest)
			return
		}
		if body.Title != nil {
			s.db.RenameFolder(id, *body.Title)
		}
		if body.IsExpanded != nil {
			s.db.ToggleFolderExpanded(id, *body.IsExpanded)
		}
		c.Out.WriteHeader(http.StatusOK)
	} else if c.Req.Method == "DELETE" {
		s.db.DeleteFolder(id)
		c.Out.WriteHeader(http.StatusNoContent)
	}
}

func (s *Server) handleFeedRefresh(c *router.Context) {
	if c.Req.Method == "POST" {
		s.worker.RefreshFeeds()
		c.Out.WriteHeader(http.StatusOK)
	} else {
		c.Out.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleFeedErrors(c *router.Context) {
	errors := s.db.GetFeedErrors()
	c.JSON(http.StatusOK, errors)
}

type feedicon struct {
	ctype string
	bytes []byte
	etag  string
}

func (s *Server) handleFeedIcon(c *router.Context) {
	id, err := c.VarInt64("id")
	if err != nil {
		c.Out.WriteHeader(http.StatusBadRequest)
		return
	}

	cachekey := "icon:" + strconv.FormatInt(id, 10)
	s.cache_mutex.Lock()
	cachedat := s.cache[cachekey]
	s.cache_mutex.Unlock()
	if cachedat == nil {
		feed := s.db.GetFeed(id)
		if feed == nil || feed.Icon == nil {
			c.Out.WriteHeader(http.StatusNotFound)
			return
		}

		hash := md5.New()
		hash.Write(*feed.Icon)

		etag := fmt.Sprintf("%x", hash.Sum(nil))[:16]

		cachedat = feedicon{
			ctype: http.DetectContentType(*feed.Icon),
			bytes: *(*feed).Icon,
			etag:  etag,
		}
		s.cache_mutex.Lock()
		s.cache[cachekey] = cachedat
		s.cache_mutex.Unlock()
	}

	icon := cachedat.(feedicon)

	if c.Req.Header.Get("If-None-Match") == icon.etag {
		c.Out.WriteHeader(http.StatusNotModified)
		return
	}

	c.Out.Header().Set("Content-Type", icon.ctype)
	c.Out.Header().Set("Etag", icon.etag)
	c.Out.Write(icon.bytes)
}

func (s *Server) handleFeedList(c *router.Context) {
	if c.Req.Method == "GET" {
		list := s.db.ListFeeds()
		c.JSON(http.StatusOK, list)
	} else if c.Req.Method == "POST" {
		var form FeedCreateForm
		if err := json.NewDecoder(c.Req.Body).Decode(&form); err != nil {
			log.Print(err)
			c.Out.WriteHeader(http.StatusBadRequest)
			return
		}

		result, err := worker.DiscoverFeed(form.Url)
		switch {
		case err != nil:
			log.Printf("Faild to discover feed for %s: %s", form.Url, err)
			c.JSON(http.StatusOK, map[string]string{"status": "notfound"})
		case len(result.Sources) > 0:
			c.JSON(http.StatusOK, map[string]interface{}{"status": "multiple", "choice": result.Sources})
		case result.Feed != nil:
			feed := s.db.CreateFeed(
				result.Feed.Title,
				"",
				result.Feed.SiteURL,
				result.FeedLink,
				form.FolderID,
			)
			items := worker.ConvertItems(result.Feed.Items, *feed)
			if len(items) > 0 {
				s.db.CreateItems(items)
				s.db.SetFeedSize(feed.Id, len(items))
				s.db.SyncSearch()
			}
			s.worker.FindFeedFavicon(*feed)

			c.JSON(http.StatusOK, map[string]interface{}{
				"status": "success",
				"feed":   feed,
			})
		default:
			c.JSON(http.StatusOK, map[string]string{"status": "notfound"})
		}
	}
}

func (s *Server) handleFeed(c *router.Context) {
	id, err := c.VarInt64("id")
	if err != nil {
		c.Out.WriteHeader(http.StatusBadRequest)
		return
	}
	if c.Req.Method == "PUT" {
		feed := s.db.GetFeed(id)
		if feed == nil {
			c.Out.WriteHeader(http.StatusBadRequest)
			return
		}
		body := make(map[string]interface{})
		if err := json.NewDecoder(c.Req.Body).Decode(&body); err != nil {
			log.Print(err)
			c.Out.WriteHeader(http.StatusBadRequest)
			return
		}
		if title, ok := body["title"]; ok {
			if reflect.TypeOf(title).Kind() == reflect.String {
				s.db.RenameFeed(id, title.(string))
			}
		}
		if f_id, ok := body["folder_id"]; ok {
			if f_id == nil {
				s.db.UpdateFeedFolder(id, nil)
			} else if reflect.TypeOf(f_id).Kind() == reflect.Float64 {
				folderId := int64(f_id.(float64))
				s.db.UpdateFeedFolder(id, &folderId)
			}
		}
		if link, ok := body["feed_link"]; ok {
			if reflect.TypeOf(link).Kind() == reflect.String {
				s.db.UpdateFeedLink(id, link.(string))
			}
		}
		c.Out.WriteHeader(http.StatusOK)
	} else if c.Req.Method == "DELETE" {
		s.db.DeleteFeed(id)
		c.Out.WriteHeader(http.StatusNoContent)
	} else {
		c.Out.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleItem(c *router.Context) {
	id, err := c.VarInt64("id")
	if err != nil {
		c.Out.WriteHeader(http.StatusBadRequest)
		return
	}
	if c.Req.Method == "GET" {
		item := s.db.GetItem(id)
		if item == nil {
			c.Out.WriteHeader(http.StatusBadRequest)
			return
		}

		// runtime fix for relative links
		if !htmlutil.IsAPossibleLink(item.Link) {
			if feed := s.db.GetFeed(item.FeedId); feed != nil {
				item.Link = htmlutil.AbsoluteUrl(item.Link, feed.Link)
			}
		}

		item.Content = sanitizer.Sanitize(item.Link, item.Content)

		c.JSON(http.StatusOK, item)
	} else if c.Req.Method == "PUT" {
		var body ItemUpdateForm
		if err := json.NewDecoder(c.Req.Body).Decode(&body); err != nil {
			log.Print(err)
			c.Out.WriteHeader(http.StatusBadRequest)
			return
		}
		if body.Status != nil {
			s.db.UpdateItemStatus(id, *body.Status)
		}
		c.Out.WriteHeader(http.StatusOK)
	} else {
		c.Out.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleItemList(c *router.Context) {
	if c.Req.Method == "GET" {
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

		items := s.db.ListItems(filter, perPage+1, newestFirst, true)
		hasMore := false
		if len(items) == perPage+1 {
			hasMore = true
			items = items[:perPage]
		}

		for i, item := range items {
			if item.Title == "" {
				text := htmlutil.ExtractText(item.Content)
				items[i].Title = htmlutil.TruncateText(text, 140)
			}
		}
		c.JSON(http.StatusOK, map[string]interface{}{
			"list":     items,
			"has_more": hasMore,
		})
	} else if c.Req.Method == "PUT" {
		filter := storage.MarkFilter{}

		if folderID, err := c.QueryInt64("folder_id"); err == nil {
			filter.FolderID = &folderID
		}
		if feedID, err := c.QueryInt64("feed_id"); err == nil {
			filter.FeedID = &feedID
		}
		s.db.MarkItemsRead(filter)
		c.Out.WriteHeader(http.StatusOK)
	} else {
		c.Out.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleSettings(c *router.Context) {
	if c.Req.Method == "GET" {
		c.JSON(http.StatusOK, s.db.GetSettings())
	} else if c.Req.Method == "PUT" {
		settings := make(map[string]interface{})
		if err := json.NewDecoder(c.Req.Body).Decode(&settings); err != nil {
			c.Out.WriteHeader(http.StatusBadRequest)
			return
		}
		if s.db.UpdateSettings(settings) {
			if _, ok := settings["refresh_rate"]; ok {
				s.worker.SetRefreshRate(s.db.GetSettingsValueInt64("refresh_rate"))
			}
			c.Out.WriteHeader(http.StatusOK)
		} else {
			c.Out.WriteHeader(http.StatusBadRequest)
		}
	}
}

func (s *Server) handleOPMLImport(c *router.Context) {
	if c.Req.Method == "POST" {
		file, _, err := c.Req.FormFile("opml")
		if err != nil {
			log.Print(err)
			return
		}
		doc, err := opml.Parse(file)
		if err != nil {
			log.Print(err)
			c.Out.WriteHeader(http.StatusBadRequest)
			return
		}
		for _, f := range doc.Feeds {
			s.db.CreateFeed(f.Title, "", f.SiteUrl, f.FeedUrl, nil)
		}
		for _, f := range doc.Folders {
			folder := s.db.CreateFolder(f.Title)
			for _, ff := range f.AllFeeds() {
				s.db.CreateFeed(ff.Title, "", ff.SiteUrl, ff.FeedUrl, &folder.Id)
			}
		}

		s.worker.FindFavicons()
		s.worker.RefreshFeeds()

		c.Out.WriteHeader(http.StatusOK)
	} else {
		c.Out.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleOPMLExport(c *router.Context) {
	if c.Req.Method == "GET" {
		c.Out.Header().Set("Content-Type", "application/xml; charset=utf-8")
		c.Out.Header().Set("Content-Disposition", `attachment; filename="subscriptions.opml"`)

		doc := opml.Folder{}

		feedsByFolderID := make(map[int64][]*storage.Feed)
		for _, feed := range s.db.ListFeeds() {
			feed := feed
			if feed.FolderId == nil {
				doc.Feeds = append(doc.Feeds, opml.Feed{
					Title:   feed.Title,
					FeedUrl: feed.FeedLink,
					SiteUrl: feed.Link,
				})
			} else {
				id := *feed.FolderId
				feedsByFolderID[id] = append(feedsByFolderID[id], &feed)
			}
		}

		for _, folder := range s.db.ListFolders() {
			folderFeeds := feedsByFolderID[folder.Id]
			if len(folderFeeds) == 0 {
				continue
			}
			opmlfolder := opml.Folder{Title: folder.Title}
			for _, feed := range folderFeeds {
				opmlfolder.Feeds = append(opmlfolder.Feeds, opml.Feed{
					Title:   feed.Title,
					FeedUrl: feed.FeedLink,
					SiteUrl: feed.Link,
				})
			}
			doc.Folders = append(doc.Folders, opmlfolder)
		}

		c.Out.Write([]byte(doc.OPML()))
	}
}

func (s *Server) handlePageCrawl(c *router.Context) {
	url := c.Req.URL.Query().Get("url")

	if newUrl := silo.RedirectURL(url); newUrl != "" {
		url = newUrl
	}
	if content := silo.VideoIFrame(url); content != "" {
		c.JSON(http.StatusOK, map[string]string{
			"content": sanitizer.Sanitize(url, content),
		})
		return
	}

	body, err := worker.GetBody(url)
	if err != nil {
		log.Print(err)
		c.Out.WriteHeader(http.StatusBadRequest)
		return
	}
	content, err := readability.ExtractContent(strings.NewReader(body))
	if err != nil {
		c.JSON(http.StatusOK, map[string]string{
			"content": "error: " + err.Error(),
		})
		return
	}
	content = sanitizer.Sanitize(url, content)
	c.JSON(http.StatusOK, map[string]string{
		"content": content,
	})
}

func (s *Server) handleLogout(c *router.Context) {
	auth.Logout(c.Out, s.BasePath)
	c.Out.WriteHeader(http.StatusNoContent)
}
