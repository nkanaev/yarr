package server

import (
	"encoding/json"
	"github.com/nkanaev/yarr/src/assets"
	"github.com/nkanaev/yarr/src/auth"
	"github.com/nkanaev/yarr/src/router"
	"github.com/nkanaev/yarr/src/storage"
	"github.com/nkanaev/yarr/src/opml"
	"github.com/nkanaev/yarr/src/worker"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"reflect"
)

func (s *Server) handler() http.Handler {
	r := router.NewRouter()

	// TODO: auth, base, security
	if s.Username != "" && s.Password != "" {
		a := &authMiddleware{
			username: s.Username,
			password: s.Password,
			basepath: BasePath,
			public: BasePath + "/static",
		}
		r.Use(a.handler)
	}

	r.For("/", s.handleIndex)
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

	return r
}

func (s *Server) handleIndex(c *router.Context) {
	c.Out.Header().Set("Content-Type", "text/html")
	assets.Render("index.html", c.Out, nil)
}

func (s *Server) handleStatic(c *router.Context) {
	// TODO: gzip?
	http.StripPrefix(BasePath+"/static/", http.FileServer(http.FS(assets.FS))).ServeHTTP(c.Out, c.Req)
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
		s.worker.FetchAllFeeds()
		c.Out.WriteHeader(http.StatusOK)
	} else {
		c.Out.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleFeedErrors(c *router.Context) {
	errors := s.db.GetFeedErrors()
	c.JSON(http.StatusOK, errors)
}

func (s *Server) handleFeedIcon(c *router.Context) {
	id, err := c.VarInt64("id")
	if err != nil {
		c.Out.WriteHeader(http.StatusBadRequest)
		return
	}
	feed := s.db.GetFeed(id)
	if feed != nil && feed.Icon != nil {
		c.Out.Header().Set("Content-Type", http.DetectContentType(*feed.Icon))
		c.Out.Write(*feed.Icon)
	} else {
		c.Out.WriteHeader(http.StatusNotFound)
	}
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

		feed, sources, err := worker.DiscoverFeed(form.Url)
		if err != nil {
			log.Print(err)
			c.JSON(http.StatusOK, map[string]string{"status": "notfound"})
			return
		}

		if feed != nil {
			storedFeed := s.db.CreateFeed(
				feed.Title,
				feed.Description,
				feed.Link,
				feed.FeedLink,
				form.FolderID,
			)
			s.db.CreateItems(worker.ConvertItems(feed.Items, *storedFeed))

			icon, err := worker.FindFavicon(storedFeed.Link, storedFeed.FeedLink)
			if icon != nil {
				s.db.UpdateFeedIcon(storedFeed.Id, icon)
			}
			if err != nil {
				log.Printf("Failed to find favicon for %s (%d): %s", storedFeed.FeedLink, storedFeed.Id, err)
			}

			c.JSON(http.StatusOK, map[string]string{"status": "success"})
		} else if sources != nil {
			c.JSON(http.StatusOK, map[string]interface{}{"status": "multiple", "choice": sources})
		} else {
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
		c.Out.WriteHeader(http.StatusOK)
	} else if c.Req.Method == "DELETE" {
		s.db.DeleteFeed(id)
		c.Out.WriteHeader(http.StatusNoContent)
	} else {
		c.Out.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleItem(c *router.Context) {
	if c.Req.Method == "PUT" {
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
		c.Out.WriteHeader(http.StatusOK)
	} else {
		c.Out.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleItemList(c *router.Context) {
	if c.Req.Method == "GET" {
		perPage := 20
		curPage := 1
		query := c.Req.URL.Query()
		if page, err := c.QueryInt64("page"); err == nil {
			curPage = int(page)
		}
		filter := storage.ItemFilter{}
		if folderID, err := c.QueryInt64("folder_id"); err == nil {
			filter.FolderID = &folderID
		}
		if feedID, err := c.QueryInt64("feed_id"); err == nil {
			filter.FeedID = &feedID
		}
		if status := query.Get("status"); len(status) != 0 {
			statusValue := storage.StatusValues[status]
			filter.Status = &statusValue
		}
		if search := query.Get("search"); len(search) != 0 {
			filter.Search = &search
		}
		newestFirst := query.Get("oldest_first") != "true"
		items := s.db.ListItems(filter, (curPage-1)*perPage, perPage, newestFirst)
		count := s.db.CountItems(filter)
		c.JSON(http.StatusOK, map[string]interface{}{
			"page": map[string]int{
				"cur": curPage,
				"num": int(math.Ceil(float64(count) / float64(perPage))),
			},
			"list": items,
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
			return
		}
		for _, outline := range doc.Outlines {
			if outline.Type == "rss" {
				s.db.CreateFeed(outline.Title, outline.Description, outline.SiteURL, outline.FeedURL, nil)
			} else {
				folder := s.db.CreateFolder(outline.Title)
				for _, o := range outline.AllFeeds() {
					s.db.CreateFeed(o.Title, o.Description, o.SiteURL, o.FeedURL, &folder.Id)
				}
			}
		}
		s.worker.FetchAllFeeds()
		c.Out.WriteHeader(http.StatusOK)
	} else {
		c.Out.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleOPMLExport(c *router.Context) {
	if c.Req.Method == "GET" {
		c.Out.Header().Set("Content-Type", "application/xml; charset=utf-8")
		c.Out.Header().Set("Content-Disposition", `attachment; filename="subscriptions.opml"`)

		rootFeeds := make([]*storage.Feed, 0)
		feedsByFolderID := make(map[int64][]*storage.Feed)
		for _, feed := range s.db.ListFeeds() {
			feed := feed
			if feed.FolderId == nil {
				rootFeeds = append(rootFeeds, &feed)
			} else {
				id := *feed.FolderId
				if feedsByFolderID[id] == nil {
					feedsByFolderID[id] = make([]*storage.Feed, 0)
				}
				feedsByFolderID[id] = append(feedsByFolderID[id], &feed)
			}
		}
		builder := opml.NewBuilder()
		
		for _, feed := range rootFeeds {
			builder.AddFeed(feed.Title, feed.Description, feed.FeedLink, feed.Link)
		}
		for _, folder := range s.db.ListFolders() {
			folderFeeds := feedsByFolderID[folder.Id]
			if len(folderFeeds) == 0 {
				continue
			}
			feedFolder := builder.AddFolder(folder.Title)
			for _, feed := range folderFeeds {
				feedFolder.AddFeed(feed.Title, feed.Description, feed.FeedLink, feed.Link)
			}
		}

		c.Out.Write([]byte(builder.String()))
	}
}

func (s *Server) handlePageCrawl(c *router.Context) {
	query := c.Req.URL.Query()
	if url := query.Get("url"); len(url) > 0 {
		res, err := http.Get(url)
		if err == nil {
			body, err := ioutil.ReadAll(res.Body)
			if err == nil {
				c.Out.Write(body)
			}
		}
	}
}

func (s *Server) handleLogout(c *router.Context) {
	auth.Logout(c.Out)
	c.Out.WriteHeader(http.StatusNoContent)
}
