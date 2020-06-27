package server

import (
	"net/http"
	"os"
	"log"
	"io"
	"fmt"
	"mime"
)

func Index(rw http.ResponseWriter, req *http.Request) {
	fmt.Println(os.Getwd())
	f, err := os.Open("template/index.html")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	rw.Header().Set("Content-Type", "text/html")
	io.Copy(rw, f)

}

func Static(rw http.ResponseWriter, req *http.Request) {
	path := "template/static/" + Vars(req)["path"]
	f, err := os.Open(path)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	rw.Header().Set("Content-Type", mime.TypeByExtension(path))
	io.Copy(rw, f)
}

func FolderList(rw http.ResponseWriter, req *http.Request) {
}

func Folder(rw http.ResponseWriter, req *http.Request) {
}

func FeedList(rw http.ResponseWriter, req *http.Request) {
}

func Feed(rw http.ResponseWriter, req *http.Request) {
}
