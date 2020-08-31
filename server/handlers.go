package server

import (
	"bytes"
	"encoding/json"
	"encoding/base64"
	"compress/gzip"
	"fmt"
	"github.com/nkanaev/yarr/storage"
	"html"
	"html/template"
	"io"
	"io/ioutil"
	"math"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

var routes []Route = []Route{
	p("/", IndexHandler),
	p("/static/*path", StaticHandler),
	p("/api/status", StatusHandler),
	p("/api/folders", FolderListHandler),
	p("/api/folders/:id", FolderHandler),
	p("/api/feeds", FeedListHandler),
	p("/api/feeds/find", FeedHandler),
	p("/api/feeds/refresh", FeedRefreshHandler),
	p("/api/feeds/:id/icon", FeedIconHandler),
	p("/api/feeds/:id", FeedHandler),
	p("/api/items", ItemListHandler),
	p("/api/items/:id", ItemHandler),
	p("/api/settings", SettingsHandler),
	p("/opml/import", OPMLImportHandler),
	p("/opml/export", OPMLExportHandler),
	p("/page", PageCrawlHandler),
}

type asset struct {
	etag    string
	body    string  // base64(gzip(content))
	gzipped *[]byte
	decoded *string
}

func (a *asset) gzip() *[]byte {
	if a.gzipped == nil {
		gzipped, _ := base64.StdEncoding.DecodeString(a.body)
		a.gzipped = &gzipped
	}
	return a.gzipped
}

func (a *asset) text() *string {
	if a.decoded == nil {
		gzipped, _ := base64.StdEncoding.DecodeString(a.body)
		reader, _ := gzip.NewReader(bytes.NewBuffer(gzipped))
		decoded, _ := ioutil.ReadAll(reader)
		reader.Close()

		decoded_string := string(decoded)
		a.decoded = &decoded_string
	}
	return a.decoded
}

var assets map[string]asset

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
	if assets != nil {
		asset := assets["index.html"]

		rw.Header().Set("Content-Type", "text/html")
		rw.Header().Set("Content-Encoding", "gzip")
		rw.Write(*asset.gzip())
	} else {
		t := template.Must(template.New("index.html").Delims("{%", "%}").Funcs(template.FuncMap{
			"inline": func(svg string) template.HTML {
				content, _ := ioutil.ReadFile("assets/graphicarts/" + svg)
				return template.HTML(content)
			},
		}).ParseFiles("assets/index.html"))
		rw.Header().Set("Content-Type", "text/html")
		t.Execute(rw, nil)
	}
}

func StaticHandler(rw http.ResponseWriter, req *http.Request) {
	path := Vars(req)["path"]
	ctype := mime.TypeByExtension(filepath.Ext(path))

	if assets != nil {
		if asset, ok := assets[path]; ok {
			if req.Header.Get("if-none-match") == asset.etag {
				rw.WriteHeader(http.StatusNotModified)
				return
			}
			rw.Header().Set("Content-Type", ctype)
			rw.Header().Set("Content-Encoding", "gzip")
			rw.Header().Set("Etag", asset.etag)
			rw.Write(*asset.gzip())
		}
	}

	f, err := os.Open("assets/" + path)
	if err != nil {
		return
	}
	defer f.Close()
	rw.Header().Set("Content-Type", ctype)
	io.Copy(rw, f)
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
		var feed FeedCreateForm
		if err := json.NewDecoder(req.Body).Decode(&feed); err != nil {
			handler(req).log.Print(err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		feedUrl := feed.Url
		feedreq, _ := http.NewRequest("GET", feedUrl, nil)
		feedreq.Header.Set("user-agent", req.Header.Get("user-agent"))
		feedclient := &http.Client{}
		res, err := feedclient.Do(feedreq)
		if err != nil {
			handler(req).log.Print(err)
			writeJSON(rw, map[string]string{"status": "notfound"})
			return
		} else if res.StatusCode != 200 {
			handler(req).log.Printf("Failed to fetch %s (status: %d)", feedUrl, res.StatusCode)
			body, err := ioutil.ReadAll(res.Body)
			handler(req).log.Print(string(body), err)
			writeJSON(rw, map[string]string{"status": "notfound"})
			return
		}

		contentType := res.Header.Get("Content-Type")
		if strings.HasPrefix(contentType, "text/html") || contentType == "" {
			sources, err := FindFeeds(res)
			if err != nil {
				handler(req).log.Print(err)
				writeJSON(rw, map[string]string{"status": "notfound"})
				return
			}
			if len(sources) == 0 {
				writeJSON(rw, map[string]string{"status": "notfound"})
			} else if len(sources) > 1 {
				writeJSON(rw, map[string]interface{}{
					"status": "multiple",
					"choice": sources,
				})
			} else if len(sources) == 1 {
				feedUrl = sources[0].Url
				err = createFeed(db(req), feedUrl, feed.FolderID)
				if err != nil {
					handler(req).log.Print(err)
					rw.WriteHeader(http.StatusBadRequest)
					return
				}
				writeJSON(rw, map[string]string{"status": "success"})
			}
		} else if strings.Contains(contentType, "xml") || strings.Contains(contentType, "json") {
			// text/xml, application/xml, application/rss+xml, application/atom+xml
			err = createFeed(db(req), feedUrl, feed.FolderID)
			if err == nil {
				writeJSON(rw, map[string]string{"status": "success"})
			}
		} else {
			writeJSON(rw, map[string]string{"status": "notfound"})
			return
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
		filter := storage.ItemFilter{}
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
