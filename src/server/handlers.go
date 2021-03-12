package server

import (
	"encoding/json"
	"fmt"
	"github.com/nkanaev/yarr/src/storage"
	"github.com/nkanaev/yarr/src/assets"
	"html"
	"io/ioutil"
	"math"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

var routes []Route = []Route{
	p("/", IndexHandler).ManualAuth(),
	p("/static/*path", StaticHandler).ManualAuth(),

	p("/api/status", StatusHandler),
	p("/api/folders", FolderListHandler),
	p("/api/folders/:id", FolderHandler),
	p("/api/feeds", FeedListHandler),
	p("/api/feeds/find", FeedHandler),
	p("/api/feeds/refresh", FeedRefreshHandler),
	p("/api/feeds/errors", FeedErrorsHandler),
	p("/api/feeds/:id/icon", FeedIconHandler),
	p("/api/feeds/:id", FeedHandler),
	p("/api/items", ItemListHandler),
	p("/api/items/:id", ItemHandler),
	p("/api/settings", SettingsHandler),
	p("/opml/import", OPMLImportHandler),
	p("/opml/export", OPMLExportHandler),
	p("/page", PageCrawlHandler),
	p("/logout", LogoutHandler),
}

type FolderCreateForm struct {
	Title string `json:"title"`
}

type FolderUpdateForm struct {
	Title      *string `json:"title,omitempty"`
	IsExpanded *bool   `json:"is_expanded,omitempty"`
}

type FeedCreateForm struct {
	Url      string `json:"url"`
	FolderID *int64 `json:"folder_id,omitempty"`
}

type ItemUpdateForm struct {
	Status *storage.ItemStatus `json:"status,omitempty"`
}

func IndexHandler(rw http.ResponseWriter, req *http.Request) {
	h := handler(req)
	if h.requiresAuth() && !userIsAuthenticated(req, h.Username, h.Password) {
		if req.Method == "POST" {
			username := req.FormValue("username")
			password := req.FormValue("password")
			if stringsEqual(username, h.Username) && stringsEqual(password, h.Password) {
				userAuthenticate(rw, username, password)
				http.Redirect(rw, req, req.URL.Path, http.StatusFound)
				return
			}
		}

		rw.Header().Set("Content-Type", "text/html")
		assets.Render("login.html", rw, nil)
		return
	}
	rw.Header().Set("Content-Type", "text/html")
	assets.Render("index.html", rw, nil)
}

func StaticHandler(rw http.ResponseWriter, req *http.Request) {
	// TODO: gzip?
	http.StripPrefix(BasePath+"/static/", http.FileServer(http.FS(assets.FS))).ServeHTTP(rw, req)
}

func StatusHandler(rw http.ResponseWriter, req *http.Request) {
	writeJSON(rw, map[string]interface{}{
		"running": *handler(req).queueSize,
		"stats":   db(req).FeedStats(),
	})
}

func FolderListHandler(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		list := db(req).ListFolders()
		writeJSON(rw, list)
	} else if req.Method == "POST" {
		var body FolderCreateForm
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			handler(req).log.Print(err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		if len(body.Title) == 0 {
			rw.WriteHeader(http.StatusBadRequest)
			writeJSON(rw, map[string]string{"error": "Folder title missing."})
			return
		}
		folder := db(req).CreateFolder(body.Title)
		rw.WriteHeader(http.StatusCreated)
		writeJSON(rw, folder)
	} else {
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func FolderHandler(rw http.ResponseWriter, req *http.Request) {
	id, err := strconv.ParseInt(Vars(req)["id"], 10, 64)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.Method == "PUT" {
		var body FolderUpdateForm
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			handler(req).log.Print(err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		if body.Title != nil {
			db(req).RenameFolder(id, *body.Title)
		}
		if body.IsExpanded != nil {
			db(req).ToggleFolderExpanded(id, *body.IsExpanded)
		}
		rw.WriteHeader(http.StatusOK)
	} else if req.Method == "DELETE" {
		db(req).DeleteFolder(id)
		rw.WriteHeader(http.StatusNoContent)
	}
}

func FeedRefreshHandler(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		handler(req).fetchAllFeeds()
		rw.WriteHeader(http.StatusOK)
	} else {
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func FeedErrorsHandler(rw http.ResponseWriter, req *http.Request) {
	errors := db(req).GetFeedErrors()
	writeJSON(rw, errors)
}

func FeedIconHandler(rw http.ResponseWriter, req *http.Request) {
	id, err := strconv.ParseInt(Vars(req)["id"], 10, 64)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	feed := db(req).GetFeed(id)
	if feed != nil && feed.Icon != nil {
		rw.Header().Set("Content-Type", http.DetectContentType(*feed.Icon))
		rw.Header().Set("Content-Length", strconv.Itoa(len(*feed.Icon)))
		rw.Write(*feed.Icon)
	} else {
		rw.WriteHeader(http.StatusNotFound)
	}
}

func FeedListHandler(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		list := db(req).ListFeeds()
		writeJSON(rw, list)
	} else if req.Method == "POST" {
		var form FeedCreateForm
		if err := json.NewDecoder(req.Body).Decode(&form); err != nil {
			handler(req).log.Print(err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		feed, sources, err := discoverFeed(form.Url)
		if err != nil {
			handler(req).log.Print(err)
			writeJSON(rw, map[string]string{"status": "notfound"})
			return
		}

		if feed != nil {
			storedFeed := db(req).CreateFeed(
				feed.Title,
				feed.Description,
				feed.Link,
				feed.FeedLink,
				form.FolderID,
			)
			db(req).CreateItems(convertItems(feed.Items, *storedFeed))

			icon, err := findFavicon(storedFeed.Link, storedFeed.FeedLink)
			if icon != nil {
				db(req).UpdateFeedIcon(storedFeed.Id, icon)
			}
			if err != nil {
				handler(req).log.Printf("Failed to find favicon for %s (%d): %s", storedFeed.FeedLink, storedFeed.Id, err)
			}

			writeJSON(rw, map[string]string{"status": "success"})
		} else if sources != nil {
			writeJSON(rw, map[string]interface{}{"status": "multiple", "choice": sources})
		} else {
			writeJSON(rw, map[string]string{"status": "notfound"})
		}
	}
}

func FeedHandler(rw http.ResponseWriter, req *http.Request) {
	id, err := strconv.ParseInt(Vars(req)["id"], 10, 64)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.Method == "PUT" {
		feed := db(req).GetFeed(id)
		if feed == nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		body := make(map[string]interface{})
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			handler(req).log.Print(err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		if title, ok := body["title"]; ok {
			if reflect.TypeOf(title).Kind() == reflect.String {
				db(req).RenameFeed(id, title.(string))
			}
		}
		if f_id, ok := body["folder_id"]; ok {
			if f_id == nil {
				db(req).UpdateFeedFolder(id, nil)
			} else if reflect.TypeOf(f_id).Kind() == reflect.Float64 {
				folderId := int64(f_id.(float64))
				db(req).UpdateFeedFolder(id, &folderId)
			}
		}
		rw.WriteHeader(http.StatusOK)
	} else if req.Method == "DELETE" {
		db(req).DeleteFeed(id)
		rw.WriteHeader(http.StatusNoContent)
	} else {
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func ItemHandler(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "PUT" {
		id, err := strconv.ParseInt(Vars(req)["id"], 10, 64)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		var body ItemUpdateForm
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			handler(req).log.Print(err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		if body.Status != nil {
			db(req).UpdateItemStatus(id, *body.Status)
		}
		rw.WriteHeader(http.StatusOK)
	} else {
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func ItemListHandler(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		perPage := 20
		curPage := 1
		query := req.URL.Query()
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
		items := db(req).ListItems(filter, (curPage-1)*perPage, perPage, newestFirst)
		count := db(req).CountItems(filter)
		writeJSON(rw, map[string]interface{}{
			"page": map[string]int{
				"cur": curPage,
				"num": int(math.Ceil(float64(count) / float64(perPage))),
			},
			"list": items,
		})
	} else if req.Method == "PUT" {
		query := req.URL.Query()
		filter := storage.MarkFilter{}
		if folderID, err := strconv.ParseInt(query.Get("folder_id"), 10, 64); err == nil {
			filter.FolderID = &folderID
		}
		if feedID, err := strconv.ParseInt(query.Get("feed_id"), 10, 64); err == nil {
			filter.FeedID = &feedID
		}
		db(req).MarkItemsRead(filter)
		rw.WriteHeader(http.StatusOK)
	} else {
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func SettingsHandler(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		writeJSON(rw, db(req).GetSettings())
	} else if req.Method == "PUT" {
		settings := make(map[string]interface{})
		if err := json.NewDecoder(req.Body).Decode(&settings); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		if db(req).UpdateSettings(settings) {
			if _, ok := settings["refresh_rate"]; ok {
				handler(req).refreshRate <- db(req).GetSettingsValueInt64("refresh_rate")
			}
			rw.WriteHeader(http.StatusOK)
		} else {
			rw.WriteHeader(http.StatusBadRequest)
		}
	}
}

func OPMLImportHandler(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		file, _, err := req.FormFile("opml")
		if err != nil {
			handler(req).log.Print(err)
			return
		}
		doc, err := parseOPML(file)
		if err != nil {
			handler(req).log.Print(err)
			return
		}
		for _, outline := range doc.Outlines {
			if outline.Type == "rss" {
				db(req).CreateFeed(outline.Title, outline.Description, outline.SiteURL, outline.FeedURL, nil)
			} else {
				folder := db(req).CreateFolder(outline.Title)
				for _, o := range outline.AllFeeds() {
					db(req).CreateFeed(o.Title, o.Description, o.SiteURL, o.FeedURL, &folder.Id)
				}
			}
		}
		handler(req).fetchAllFeeds()
		rw.WriteHeader(http.StatusOK)
	} else {
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func OPMLExportHandler(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		rw.Header().Set("Content-Type", "application/xml; charset=utf-8")
		rw.Header().Set("Content-Disposition", `attachment; filename="subscriptions.opml"`)

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
		for _, feed := range db(req).ListFeeds() {
			var folderId = int64(0)
			if feed.FolderId != nil {
				folderId = *feed.FolderId
			}
			if feedsByFolderID[folderId] == nil {
				feedsByFolderID[folderId] = make([]storage.Feed, 0)
			}
			feedsByFolderID[folderId] = append(feedsByFolderID[folderId], feed)
		}
		for _, folder := range db(req).ListFolders() {
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
		rw.Write([]byte(builder.String()))
	}
}

func PageCrawlHandler(rw http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	if url := query.Get("url"); len(url) > 0 {
		res, err := http.Get(url)
		if err == nil {
			body, err := ioutil.ReadAll(res.Body)
			if err == nil {
				rw.Write(body)
			}
		}
	}
}

func LogoutHandler(rw http.ResponseWriter, req *http.Request) {
	userLogout(rw)
	rw.WriteHeader(http.StatusNoContent)
}
