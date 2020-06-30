package server

import (
	"encoding/json"
	"context"
	"regexp"
	"net/http"
	"github.com/nkanaev/yarr/storage"
	"log"
)

type Route struct {
	url string
	urlRegex *regexp.Regexp
	handler func(http.ResponseWriter, *http.Request)
}

type Handler struct {
	db *storage.Storage
}

func p(path string, handler func(http.ResponseWriter, *http.Request)) Route {
	var urlRegexp string
	urlRegexp = regexp.MustCompile(`[\*\:]\w+`).ReplaceAllStringFunc(path, func(m string) string {
		if m[0:1] == `*` {
			return "(?P<" + m[1:] + ">.+)"
		}
		return "(?P<" + m[1:] + ">[^/]+)"
	})
	urlRegexp = "^" + urlRegexp + "$"
	return Route{
		url: path,
		urlRegex: regexp.MustCompile(urlRegexp),
		handler: handler,
	}
}

var routes []Route = []Route{
	p("/", IndexHandler),
	p("/static/*path", StaticHandler),
	p("/api/folders", FolderListHandler),
	p("/api/folders/:id", FolderHandler),
	p("/api/feeds", FeedListHandler),
	p("/api/feeds/:id", FeedHandler),
	p("/api/feeds/find", FeedHandler),
}

func Vars(req *http.Request) map[string]string {
	if rv := req.Context().Value(ctxVars); rv != nil {
		return rv.(map[string]string)
	}
	return nil
}

func db(req *http.Request) *storage.Storage {
	if rv := req.Context().Value(ctxDB); rv != nil {
		return rv.(*storage.Storage)
	}
	return nil
}

const (
	ctxDB = 1
	ctxVars = 2
)

func (h Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for _, route := range routes {
		if route.urlRegex.MatchString(req.URL.Path) {
			ctx := context.WithValue(req.Context(), ctxDB, h.db)
			if route.urlRegex.NumSubexp() > 0 {
				vars := make(map[string]string)
				matches := route.urlRegex.FindStringSubmatchIndex(req.URL.Path)
				for i, key := range route.urlRegex.SubexpNames()[1:] {
					vars[key] = req.URL.Path[matches[i*2+2]:matches[i*2+3]]
				}
				ctx = context.WithValue(ctx, ctxVars, vars)
			}
			route.handler(rw, req.WithContext(ctx))
			return
		}
	}
	rw.WriteHeader(http.StatusNotFound)
}

func writeJSON(rw http.ResponseWriter, data interface{}) {
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	reply, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	rw.Write(reply)
	rw.Write([]byte("\n"))
}

func New() *http.Server {
	db, _ := storage.New()
	h := Handler{db: db}
	s := &http.Server{Addr: "127.0.0.1:8000", Handler: h}
	return s
}
