package server

import (
	"context"
	"github.com/nkanaev/yarr/storage"
	"log"
	"net/http"
)

type Handler struct {
	db           *storage.Storage
	log          *log.Logger
	fetchRunning bool
	feedQueue    chan storage.Feed
	counter      chan int
	queueSize    int
}

func New(db *storage.Storage, logger *log.Logger) *Handler {
	return &Handler{
		db:        db,
		log:       logger,
		feedQueue: make(chan storage.Feed),
		counter:   make(chan int),
	}
}

func (h *Handler) Start(addr string) {
	h.startJobs()
	s := &http.Server{Addr: addr, Handler: h}
	s.ListenAndServe()
}

func (h Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	route, vars := getRoute(req)
	if route == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	ctx := context.WithValue(req.Context(), ctxHandler, &h)
	ctx = context.WithValue(ctx, ctxVars, vars)
	route.handler(rw, req.WithContext(ctx))
}

func (h *Handler) startJobs() {
	h.db.DeleteOldItems()
	go func() {
		for {
			feed := <-h.feedQueue
			items := listItems(feed)
			h.db.CreateItems(items)
		}
	}()
	go func() {
		for {
			val := <-h.counter
			h.queueSize += val
		}
	}()
	go h.db.SyncSearch()
	h.fetchAllFeeds()
}

func (h *Handler) fetchFeed(feed storage.Feed) {
	h.queueSize += 1
	h.feedQueue <- feed
}

func (h *Handler) fetchAllFeeds() {
	for _, feed := range h.db.ListFeeds() {
		h.fetchFeed(feed)
	}
}

func Vars(req *http.Request) map[string]string {
	if rv := req.Context().Value(ctxVars); rv != nil {
		return rv.(map[string]string)
	}
	return nil
}

func db(req *http.Request) *storage.Storage {
	if h := handler(req); h != nil {
		return h.db
	}
	return nil
}

func handler(req *http.Request) *Handler {
	return req.Context().Value(ctxHandler).(*Handler)
}

const (
	ctxVars    = 2
	ctxHandler = 3
)
