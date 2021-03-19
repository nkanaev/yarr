package worker

import (
	"github.com/nkanaev/yarr/src/storage"
	"log"
	"runtime"
	"sync/atomic"
	"time"
)

type Worker struct {
	db *storage.Storage

	feedQueue   chan storage.Feed
	queueSize   *int32
	refreshRate chan int64
}

func NewWorker(db *storage.Storage) *Worker {
	queueSize := int32(0)
	return &Worker{
		db:          db,
		feedQueue:   make(chan storage.Feed, 3000),
		queueSize:   &queueSize,
		refreshRate: make(chan int64),
	}
}

func (w *Worker) Start() {
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
			case feed := <-w.feedQueue:
				items, err := listItems(feed, w.db)
				atomic.AddInt32(w.queueSize, -1)
				if err != nil {
					log.Printf("Failed to fetch %s (%d): %s", feed.FeedLink, feed.Id, err)
					w.db.SetFeedError(feed.Id, err)
					continue
				}
				w.db.CreateItems(items)
				syncSearch()
				if !feed.HasIcon {
					icon, err := FindFavicon(feed.Link, feed.FeedLink)
					if icon != nil {
						w.db.UpdateFeedIcon(feed.Id, icon)
					}
					if err != nil {
						log.Printf("Failed to search favicon for %s (%s): %s", feed.Link, feed.FeedLink, err)
					}
				}
			case <-delTicker.C:
				w.db.DeleteOldItems()
			case <-syncSearchChannel:
				w.db.SyncSearch()
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
	go w.db.DeleteOldItems()
	go w.db.SyncSearch()

	go func() {
		var refreshTicker *time.Ticker
		refreshTick := make(<-chan time.Time)
		for {
			select {
			case <-refreshTick:
				w.FetchAllFeeds()
			case val := <-w.refreshRate:
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
	refreshRate := w.db.GetSettingsValueInt64("refresh_rate")
	w.refreshRate <- refreshRate
	if refreshRate > 0 {
		w.FetchAllFeeds()
	}
}

func (w *Worker) FetchAllFeeds() {
	log.Print("Refreshing all feeds")
	w.db.ResetFeedErrors()
	for _, feed := range w.db.ListFeeds() {
		w.fetchFeed(feed)
	}
}

func (w *Worker) fetchFeed(feed storage.Feed) {
	atomic.AddInt32(w.queueSize, 1)
	w.feedQueue <- feed
}

func (w *Worker) FeedsPending() int32 {
	return *w.queueSize
}

func (w *Worker) SetRefreshRate(val int64) {
	w.refreshRate <- val
}
