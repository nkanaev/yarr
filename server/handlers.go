package server

import (
	"github.com/nkanaev/yarr/worker"
	"net/http"
	"encoding/json"
	"os"
	"log"
	"io"
	"fmt"
	"mime"
	"strings"
	"path/filepath"
)

func IndexHandler(rw http.ResponseWriter, req *http.Request) {
	fmt.Println(os.Getwd())
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

func FolderListHandler(rw http.ResponseWriter, req *http.Request) {
}

func FolderHandler(rw http.ResponseWriter, req *http.Request) {
}

type NewFeed struct {
	Url string	   `json:"url"`
	FolderID int64 `json:"folder_id,omitempty"`
}

func FeedListHandler(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
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
		if strings.HasPrefix(contentType, "text/html") {
			sources, err := worker.FindFeeds(res)
			if err != nil {
				log.Print(err)
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
				fmt.Println("feedUrl:", feedUrl)
				writeJSON(rw, map[string]string{"status": "success"})
			}
			fmt.Println("got html url", sources, feedUrl)
		} else if strings.HasPrefix(contentType, "text/xml") {
			log.Print("got rss feed")
		}
		log.Print(res.Header.Get("Content-Type"))
	}
}

func FeedHandler(rw http.ResponseWriter, req *http.Request) {
}
