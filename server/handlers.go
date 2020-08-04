package server

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/nkanaev/yarr/storage"
	"html"
	"html/template"
	"io"
	"io/ioutil"
	"log"
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

func IndexHandler(rw http.ResponseWriter, req *http.Request) {
	t := template.Must(template.New("index.html").Delims("{%", "%}").Funcs(template.FuncMap{
		"inline": func(svg string) template.HTML {
			content, _ := ioutil.ReadFile("template/static/images/" + svg)
			return template.HTML(content)
		},
	}).ParseFiles("template/index.html"))
	rw.Header().Set("Content-Type", "text/html")
	t.Execute(rw, nil)
}

func StaticHandler(rw http.ResponseWriter, req *http.Request) {
	path := "template/static/" + Vars(req)["path"]
	f, err := os.Open(path)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	rw.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(path)))
	io.Copy(rw, f)
}

func StatusHandler(rw http.ResponseWriter, req *http.Request) {
	writeJSON(rw, map[string]interface{}{
		"running": *handler(req).queueSize,
		"stats":   db(req).FeedStats(),
	})
}

type NewFolder struct {
	Title string `json:"title"`
}

func FolderListHandler(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		list := db(req).ListFolders()
		writeJSON(rw, list)
	} else if req.Method == "POST" {
		var body NewFolder
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			log.Print(err)
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

type UpdateFolder struct {
	Title      *string `json:"title,omitempty"`
	IsExpanded *bool   `json:"is_expanded,omitempty"`
}

func FolderHandler(rw http.ResponseWriter, req *http.Request) {
	id, err := strconv.ParseInt(Vars(req)["id"], 10, 64)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.Method == "PUT" {
		var body UpdateFolder
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			log.Print(err)
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

type NewFeed struct {
	Url      string `json:"url"`
	FolderID *int64 `json:"folder_id,omitempty"`
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
		var feed NewFeed
		if err := json.NewDecoder(req.Body).Decode(&feed); err != nil {
			log.Print(err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		feedUrl := feed.Url
		res, err := http.Get(feedUrl)
		if err != nil {
			log.Print(err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		} else if res.StatusCode != 200 {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		contentType := res.Header.Get("Content-Type")
		if strings.HasPrefix(contentType, "text/html") || contentType == "" {
			sources, err := FindFeeds(res)
			if err != nil {
				log.Print(err)
				rw.WriteHeader(http.StatusBadRequest)
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
					log.Print(err)
					rw.WriteHeader(http.StatusBadRequest)
					return
				}
				writeJSON(rw, map[string]string{"status": "success"})
			}
		} else if strings.HasPrefix(contentType, "text/xml") || strings.HasPrefix(contentType, "application/xml") || strings.HasPrefix(contentType, "application/rss+xml") {
			err = createFeed(db(req), feedUrl, feed.FolderID)
			if err == nil {
				writeJSON(rw, map[string]string{"status": "success"})
			}
		} else {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}

func convertItems(items []*gofeed.Item, feed storage.Feed) []storage.Item {
	result := make([]storage.Item, len(items))
	for _, item := range items {
		imageURL := ""
		if item.Image != nil {
			imageURL = item.Image.URL
		}
		author := ""
		if item.Author != nil {
			author = item.Author.Name
		}
		result = append(result, storage.Item{
			GUID:        item.GUID,
			FeedId:      feed.Id,
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
			Content:     item.Content,
			Author:      author,
			Date:        item.PublishedParsed,
			DateUpdated: item.UpdatedParsed,
			Status:      storage.UNREAD,
			Image:       imageURL,
		})
	}
	return result
}

func listItems(f storage.Feed) []storage.Item {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(f.FeedLink)
	if err != nil {
		log.Print(err)
		return nil
	}
	return convertItems(feed.Items, f)
}

func createFeed(s *storage.Storage, url string, folderId *int64) error {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return err
	}
	feedLink := feed.FeedLink
	if len(feedLink) == 0 {
		feedLink = url
	}
	storedFeed := s.CreateFeed(
		feed.Title,
		feed.Description,
		feed.Link,
		feedLink,
		folderId,
	)
	s.CreateItems(convertItems(feed.Items, *storedFeed))
	return nil
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
			log.Print(err)
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

type UpdateItem struct {
	Status *storage.ItemStatus `json:"status,omitempty"`
}

func ItemHandler(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "PUT" {
		id, err := strconv.ParseInt(Vars(req)["id"], 10, 64)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		var body UpdateItem
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			log.Print(err)
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

type opml struct {
	XMLName  xml.Name  `xml:"opml"`
	Version  string    `xml:"version,attr"`
	Outlines []outline `xml:"body>outline"`
}

type outline struct {
	Type        string    `xml:"type,attr,omitempty"`
	Title       string    `xml:"text,attr"`
	FeedURL     string    `xml:"xmlUrl,attr,omitempty"`
	SiteURL     string    `xml:"htmlUrl,attr,omitempty"`
	Description string    `xml:"description,attr,omitempty"`
	Outlines    []outline `xml:"outline,omitempty"`
}

func (o outline) AllFeeds() []outline {
	result := make([]outline, 0)
	for _, sub := range o.Outlines {
		if sub.Type == "rss" {
			result = append(result, sub)
		} else {
			result = append(result, sub.AllFeeds()...)
		}
	}
	return result
}

func OPMLImportHandler(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		file, _, err := req.FormFile("opml")
		if err != nil {
			log.Print(err)
			return
		}
		feeds := new(opml)
		decoder := xml.NewDecoder(file)
		decoder.Entity = xml.HTMLEntity
		decoder.Strict = false
		err = decoder.Decode(&feeds)
		if err != nil {
			log.Print(err)
			return
		}
		for _, outline := range feeds.Outlines {
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
