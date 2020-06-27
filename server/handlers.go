package server

import "net/http"

func Index(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("index"))
}

func Static(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("static:" + Vars(req)["path"]))
}

func FolderList(rw http.ResponseWriter, req *http.Request) {
}

func Folder(rw http.ResponseWriter, req *http.Request) {
}

func FeedList(rw http.ResponseWriter, req *http.Request) {
}

func Feed(rw http.ResponseWriter, req *http.Request) {
}
