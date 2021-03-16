package server

import (
	"encoding/json"
	"fmt"
	"github.com/nkanaev/yarr/src/assets"
	"github.com/nkanaev/yarr/src/auth"
	"github.com/nkanaev/yarr/src/router"
	"github.com/nkanaev/yarr/src/storage"
	"html"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

func (s *Server) handler() http.Handler {
	r := router.NewRouter()

	// TODO: auth, base, security

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
	if s.requiresAuth() && !auth.IsAuthenticated(c.Req, s.Username, s.Password) {
		if c.Req.Method == "POST" {
			username := c.Req.FormValue("username")
			password := c.Req.FormValue("password")
			if auth.StringsEqual(username, s.Username) && auth.StringsEqual(password, s.Password) {
				auth.Authenticate(c.Out, username, password, BasePath)
				http.Redirect(c.Out, c.Req, c.Req.URL.Path, http.StatusFound)
				return
			}
		}

		c.Out.Header().Set("Content-Type", "text/html")
		assets.Render("login.html", c.Out, nil)
		return
	}
	c.Out.Header().Set("Content-Type", "text/html")
	assets.Render("index.html", c.Out, nil)
}

func (s *Server) handleStatic(c *router.Context) {
	// TODO: gzip?
	http.StripPrefix(BasePath+"/static/", http.FileServer(http.FS(assets.FS))).ServeHTTP(c.Out, c.Req)
}

func (s *Server) handleStatus(c *router.Context) {
	writeJSON(c.Out, map[string]interface{}{
		"running": *s.queueSize,
		"stats":   s.db.FeedStats(),
	})
}

func (s *Server) handleFolderList(c *router.Context) {
	if c.Req.Method == "GET" {
		list := s.db.ListFolders()
		writeJSON(c.Out, list)
	} else if c.Req.Method == "POST" {
		var body FolderCreateForm
		if err := json.NewDecoder(c.Req.Body).Decode(&body); err != nil {
			log.Print(err)
			c.Out.WriteHeader(http.StatusBadRequest)
			return
		}
		if len(body.Title) == 0 {
			c.Out.WriteHeader(http.StatusBadRequest)
			writeJSON(c.Out, map[string]string{"error": "Folder title missing."})
			return
		}
		folder := s.db.CreateFolder(body.Title)
		c.Out.WriteHeader(http.StatusCreated)
		writeJSON(c.Out, folder)
	} else {
		c.Out.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleFolder(c *router.Context) {
	id, err := strconv.ParseInt(c.Vars["id"], 10, 64)
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
		s.fetchAllFeeds()
		c.Out.WriteHeader(http.StatusOK)
	} else {
		c.Out.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleFeedErrors(c *router.Context) {
	errors := s.db.GetFeedErrors()
	writeJSON(c.Out, errors)
}

func (s *Server) handleFeedIcon(c *router.Context) {
	id, err := strconv.ParseInt(c.Vars["id"], 10, 64)
	if err != nil {
		c.Out.WriteHeader(http.StatusBadRequest)
		return
	}
	feed := s.db.GetFeed(id)
	if feed != nil && feed.Icon != nil {
		c.Out.Header().Set("Content-Type", http.DetectContentType(*feed.Icon))
		c.Out.Header().Set("Content-Length", strconv.Itoa(len(*feed.Icon)))
		c.Out.Write(*feed.Icon)
	} else {
		c.Out.WriteHeader(http.StatusNotFound)
	}
}

func (s *Server) handleFeedList(c *router.Context) {
	if c.Req.Method == "GET" {
		list := s.db.ListFeeds()
		writeJSON(c.Out, list)
	} else if c.Req.Method == "POST" {
		var form FeedCreateForm
		if err := json.NewDecoder(c.Req.Body).Decode(&form); err != nil {
			log.Print(err)
			c.Out.WriteHeader(http.StatusBadRequest)
			return
		}

		feed, sources, err := discoverFeed(form.Url)
		if err != nil {
			log.Print(err)
			writeJSON(c.Out, map[string]string{"status": "notfound"})
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
			s.db.CreateItems(convertItems(feed.Items, *storedFeed))

			icon, err := findFavicon(storedFeed.Link, storedFeed.FeedLink)
			if icon != nil {
				s.db.UpdateFeedIcon(storedFeed.Id, icon)
			}
			if err != nil {
				log.Printf("Failed to find favicon for %s (%d): %s", storedFeed.FeedLink, storedFeed.Id, err)
			}

			writeJSON(c.Out, map[string]string{"status": "success"})
		} else if sources != nil {
			writeJSON(c.Out, map[string]interface{}{"status": "multiple", "choice": sources})
		} else {
			writeJSON(c.Out, map[string]string{"status": "notfound"})
		}
	}
}

func (s *Server) handleFeed(c *router.Context) {
	id, err := strconv.ParseInt(c.Vars["id"], 10, 64)
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
		id, err := strconv.ParseInt(c.Vars["id"], 10, 64)
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
		if page, err := strconv.ParseInt(query.Get("page"), 10, 64); err == nil {
			curPage = int(page)
		}
		filter := storage.ItemFilter{}
		if folderID, err := strconv.ParseInt(query.Get("folder_id"), 10, 64); err == nil {
			filter.FolderID = &folderID
		}
		if feedID, err := strconv.ParseInt(query.Get("feed_id"), 10, 64); err == nil {
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
		writeJSON(c.Out, map[string]interface{}{
			"page": map[string]int{
				"cur": curPage,
				"num": int(math.Ceil(float64(count) / float64(perPage))),
			},
			"list": items,
		})
	} else if c.Req.Method == "PUT" {
		query := c.Req.URL.Query()
		filter := storage.MarkFilter{}
		if folderID, err := strconv.ParseInt(query.Get("folder_id"), 10, 64); err == nil {
			filter.FolderID = &folderID
		}
		if feedID, err := strconv.ParseInt(query.Get("feed_id"), 10, 64); err == nil {
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
		writeJSON(c.Out, s.db.GetSettings())
	} else if c.Req.Method == "PUT" {
		settings := make(map[string]interface{})
		if err := json.NewDecoder(c.Req.Body).Decode(&settings); err != nil {
			c.Out.WriteHeader(http.StatusBadRequest)
			return
		}
		if s.db.UpdateSettings(settings) {
			if _, ok := settings["refresh_rate"]; ok {
				s.refreshRate <- s.db.GetSettingsValueInt64("refresh_rate")
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
		doc, err := parseOPML(file)
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
		s.fetchAllFeeds()
		c.Out.WriteHeader(http.StatusOK)
	} else {
		c.Out.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleOPMLExport(c *router.Context) {
	if c.Req.Method == "GET" {
		c.Out.Header().Set("Content-Type", "application/xml; charset=utf-8")
		c.Out.Header().Set("Content-Disposition", `attachment; filename="subscriptions.opml"`)

		builder := strings.Builder{}

		line := func(s string, args ...string) {
			if len(args) > 0 {
				escapedargs := make([]interface{}, len(args))
				for idx, arg := range args {
					escapedargs[idx] = html.EscapeString(arg)
				}
				s = fmt.Sprintf(s, escapedargs...)
			}
			builder.WriteString(s)
			builder.WriteString("\n")
		}

		feedline := func(feed storage.Feed, indent int) {
			line(
				strings.Repeat(" ", indent)+
					`<outline type="rss" text="%s" description="%s" xmlUrl="%s" htmlUrl="%s"/>`,
				feed.Title, feed.Description,
				feed.FeedLink, feed.Link,
			)
		}
		line(`<?xml version="1.0" encoding="UTF-8"?>`)
		line(`<opml version="1.1">`)
		line(`<head>`)
		line(`  <title>subscriptions.opml</title>`)
		line(`</head>`)
		line(`<body>`)
		feedsByFolderID := make(map[int64][]storage.Feed)
		for _, feed := range s.db.ListFeeds() {
			var folderId = int64(0)
			if feed.FolderId != nil {
				folderId = *feed.FolderId
			}
			if feedsByFolderID[folderId] == nil {
				feedsByFolderID[folderId] = make([]storage.Feed, 0)
			}
			feedsByFolderID[folderId] = append(feedsByFolderID[folderId], feed)
		}
		for _, folder := range s.db.ListFolders() {
			line(`  <outline text="%s">`, folder.Title)
			for _, feed := range feedsByFolderID[folder.Id] {
				feedline(feed, 4)
			}
			line(`  </outline>`)
		}
		for _, feed := range feedsByFolderID[0] {
			feedline(feed, 2)
		}
		line(`</body>`)
		line(`</opml>`)
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
