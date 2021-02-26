package server

import (
	"context"
	"log"
	"net/http"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	"github.com/nkanaev/yarr/src/storage"
)

type Handler struct {
	Addr        string
	db          *storage.Storage
	log         *log.Logger
	feedQueue   chan storage.Feed
	queueSize   *int32
	refreshRate chan int64
	// auth
	Username string
	Password string
	// https
	CertFile string
	KeyFile  string
}

func New(db *storage.Storage, logger *log.Logger, addr string) *Handler {
	queueSize := int32(0)
	return &Handler{
		db:          db,
		log:         logger,
		feedQueue:   make(chan storage.Feed, 3000),
		queueSize:   &queueSize,
		Addr:        addr,
		refreshRate: make(chan int64),
	}
}

func (h *Handler) GetAddr() string {
	proto := "http"
	if h.CertFile != "" && h.KeyFile != "" {
		proto = "https"
	}
	return proto + "://" + h.Addr + BasePath
}

func (h *Handler) Start() {
	h.startJobs()
	s := &http.Server{Addr: h.Addr, Handler: h}

	var err error
	if h.CertFile != "" && h.KeyFile != "" {
		err = s.ListenAndServeTLS(h.CertFile, h.KeyFile)
	} else {
		err = s.ListenAndServe()
	}
	if err != http.ErrServerClosed {
		h.log.Fatal(err)
	}
}

func unsafeMethod(method string) bool {
	return method == "POST" || method == "PUT" || method == "DELETE"
}

func (h Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	reqPath := req.URL.Path
	if BasePath != "" {
		if !strings.HasPrefix(reqPath, BasePath) {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		reqPath = strings.TrimPrefix(req.URL.Path, BasePath)
		if reqPath == "" {
			http.Redirect(rw, req, BasePath+"/", http.StatusFound)
			return
		}
	}
	route, vars := getRoute(reqPath)
	if route == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	if h.requiresAuth() && !route.manualAuth {
		if unsafeMethod(req.Method) && req.Header.Get("X-Requested-By") != "yarr" {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}
		if !userIsAuthenticated(req, h.Username, h.Password) {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	ctx := context.WithValue(req.Context(), ctxHandler, &h)
	ctx = context.WithValue(ctx, ctxVars, vars)
	route.handler(rw, req.WithContext(ctx))
}

func (h *Handler) startJobs() {
	delTicker := time.NewTicker(time.Hour * 24)

	syncSearchChannel := make(chan bool, 10)
	var syncSearchTimer *time.Timer // TODO: should this be atomic?

	syncSearch := func() {
		if syncSearchTimer == nil {
			syncSearchTimer = time.AfterFunc(time.Second*2, func() {
				syncSearchChannel <- true
			})
		} else {
			syncSearchTimer.Reset(time.Second * 2)
		}
	}

	worker := func() {
		for {
			select {
			case feed := <-h.feedQueue:
				items, err := listItems(feed, h.db)
				atomic.AddInt32(h.queueSize, -1)
				if err != nil {
					h.log.Printf("Failed to fetch %s (%d): %s", feed.FeedLink, feed.Id, err)
					h.db.SetFeedError(feed.Id, err)
					continue
				}
				h.db.CreateItems(items)
				syncSearch()
				if !feed.HasIcon {
					icon, err := findFavicon(feed.Link, feed.FeedLink)
					if icon != nil {
						h.db.UpdateFeedIcon(feed.Id, icon)
					}
					if err != nil {
						h.log.Printf("Failed to search favicon for %s (%s): %s", feed.Link, feed.FeedLink, err)
					}
				}
			case <-delTicker.C:
				h.db.DeleteOldItems()
			case <-syncSearchChannel:
				h.db.SyncSearch()
			}
		}
	}

	num := runtime.NumCPU() - 1
	if num < 1 {
		num = 1
	}
	for i := 0; i < num; i++ {
		go worker()
	}
	go h.db.DeleteOldItems()
	go h.db.SyncSearch()

	go func() {
		var refreshTicker *time.Ticker
		refreshTick := make(<-chan time.Time)
		for {
			select {
			case <-refreshTick:
				h.fetchAllFeeds()
			case val := <-h.refreshRate:
				if refreshTicker != nil {
					refreshTicker.Stop()
					if val == 0 {
						refreshTick = make(<-chan time.Time)
					}
				}
				if val > 0 {
					refreshTicker = time.NewTicker(time.Duration(val) * time.Minute)
					refreshTick = refreshTicker.C
				}
			}
		}
	}()
	refreshRate := h.db.GetSettingsValueInt64("refresh_rate")
	h.refreshRate <- refreshRate
	if refreshRate > 0 {
		h.fetchAllFeeds()
	}
}

func (h Handler) requiresAuth() bool {
	return h.Username != "" && h.Password != ""
}

func (h *Handler) fetchAllFeeds() {
	h.log.Print("Refreshing all feeds")
	h.db.ResetFeedErrors()
	for _, feed := range h.db.ListFeeds() {
		h.fetchFeed(feed)
	}
}

func (h *Handler) fetchFeed(feed storage.Feed) {
	atomic.AddInt32(h.queueSize, 1)
	h.feedQueue <- feed
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
