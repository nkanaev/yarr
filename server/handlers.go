package server

import (
	"github.com/nkanaev/yarr/storage"
	"github.com/mmcdole/gofeed"
	"net/http"
	"encoding/json"
	"os"
	"log"
	"io"
	"mime"
	"strings"
	"path/filepath"
	"strconv"
)

func IndexHandler(rw http.ResponseWriter, req *http.Request) {
	f, err := os.Open("template/index.html")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	rw.Header().Set("Content-Type", "text/html")
	io.Copy(rw, f)

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
		"running": handler(req).fetchRunning,
		"stats": map[string]int64{},
	})
}

type NewFolder struct {
	Title string `json:"title"`
}

func FolderListHandler(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		list := db(req).ListFolders()
		json.NewEncoder(rw).Encode(list)
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
		json.NewEncoder(rw).Encode(folder)
	} else {
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}


type UpdateFolder struct {
	Title *string `json:"title,omitempty"`
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
		rw.WriteHeader(http.StatusOK)
	} else if req.Method == "DELETE" {
		db(req).DeleteFolder(id)
		rw.WriteHeader(http.StatusNoContent)
	}
}

type NewFeed struct {
	Url string	    `json:"url"`
	FolderID *int64 `json:"folder_id,omitempty"`
}

type UpdateFeed struct {
	Title *string `json:"title,omitempty"`
	FolderID *int64 `json:"folder_id,omitempty"`
}

func FeedListHandler(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		list := db(req).ListFeeds()
		json.NewEncoder(rw).Encode(list)
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
		} else if strings.HasPrefix(contentType, "text/xml") || strings.HasPrefix(contentType, "application/xml") {
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
			GUID: item.GUID,
			FeedId: feed.Id,
			Title: item.Title,
			Link: item.Link,
			Description: item.Description,
			Content: item.Content,
			Author: author,
			Date: item.PublishedParsed,
			DateUpdated: item.UpdatedParsed,
			Status: storage.UNREAD,
			Image: imageURL,
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
		"",
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
		var body UpdateFeed
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			log.Print(err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		if body.Title != nil {
			db(req).RenameFeed(id, *body.Title)
		}
		if body.FolderID != nil {
			db(req).UpdateFeedFolder(id, *body.FolderID)
		}
		rw.WriteHeader(http.StatusOK)
	} else if req.Method == "DELETE" {
		db(req).DeleteFeed(id)
		rw.WriteHeader(http.StatusNoContent)
	} else {
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func findItems(db *storage.Storage, filter storage.ItemFilter) []storage.Item {
	statusFilter := db.GetSettingsValue("filter")
	if statusFilter != nil && len(statusFilter.(string)) != 0 {
		status := storage.StatusValues[statusFilter.(string)]
		filter.Status = &status
	}
	return db.ListItems(filter)	
}

func FeedItemsHandler(rw http.ResponseWriter, req *http.Request) {
	id, err := strconv.ParseInt(Vars(req)["id"], 10, 64)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusOK)
	items := findItems(db(req), storage.ItemFilter{FeedID: &id})
	writeJSON(rw, items)
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

func FolderItemsHandler(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		id, err := strconv.ParseInt(Vars(req)["id"], 10, 64)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		items := findItems(db(req), storage.ItemFilter{FolderID: &id})
		writeJSON(rw, items)
	} else {
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func ItemListHandler(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		items := findItems(db(req), storage.ItemFilter{})
		writeJSON(rw, items)
	} else {
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func SettingsHandler(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		rw.WriteHeader(http.StatusOK)
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
