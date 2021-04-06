package worker

import (
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nkanaev/yarr/src/storage"
)

type Worker struct {
	db      *storage.Storage
	pending *int32
	refresh *time.Ticker
	reflock sync.Mutex
	stopper chan bool
}

func NewWorker(db *storage.Storage) *Worker {
	pending := int32(0)
	return &Worker{db: db, pending: &pending}
}

func (w *Worker) FeedsPending() int32 {
	return *w.pending
}

func (w *Worker) StartFeedCleaner() {
	go w.db.DeleteOldItems()
	ticker := time.NewTicker(time.Hour * 24)
	go func() {
		for {
			<-ticker.C
			w.db.DeleteOldItems()
		}
	}()
}

func (w *Worker) FindFavicons() {
	go func() {
		for _, feed := range w.db.ListFeeds() {
			if !feed.HasIcon {
				w.FindFeedFavicon(feed)
			}
		}
	}()
}

func (w *Worker) FindFeedFavicon(feed storage.Feed) {
	icon, err := findFavicon(feed.Link, feed.FeedLink)
	if err != nil {
		log.Printf("Failed to find favicon for %s (%s): %s", feed.FeedLink, feed.Link, err)
	}
	if icon != nil {
		w.db.UpdateFeedIcon(feed.Id, icon)
	}
}

func (w *Worker) SetRefreshRate(minute int64) {
	if w.stopper != nil {
		w.refresh.Stop()
		w.refresh = nil
		w.stopper <- true
		w.stopper = nil
	}

	if minute == 0 {
		return
	}

	w.stopper = make(chan bool)
	w.refresh = time.NewTicker(time.Minute * time.Duration(minute))

	go func(fire <-chan time.Time, stop <-chan bool, m int64) {
		log.Printf("auto-refresh %dm: starting", m)
		for {
			select {
			case <-fire:
				log.Printf("auto-refresh %dm: firing", m)
				w.RefreshFeeds()
			case <-stop:
				log.Printf("auto-refresh %dm: stopping", m)
				return
			}
		}
	}(w.refresh.C, w.stopper, minute)
}

func (w *Worker) RefreshFeeds() {
	log.Print("Refreshing feeds")
	go w.refresher()
}

func (w *Worker) refresher() {
	w.reflock.Lock()

	w.db.ResetFeedErrors()

	feeds := w.db.ListFeeds()
	if len(feeds) == 0 {
		return
	}

	atomic.StoreInt32(w.pending, int32(len(feeds)))

	srcqueue := make(chan storage.Feed, len(feeds))
	dstqueue := make(chan []storage.Item)

	// hardcoded to 4 workers ;)
	go w.worker(srcqueue, dstqueue)
	go w.worker(srcqueue, dstqueue)
	go w.worker(srcqueue, dstqueue)
	go w.worker(srcqueue, dstqueue)

	for _, feed := range feeds {
		srcqueue <- feed
	}
	for i := 0; i < len(feeds); i++ {
		w.db.CreateItems(<-dstqueue)
		atomic.AddInt32(w.pending, -1)
	}
	close(srcqueue)
	close(dstqueue)

	w.db.SyncSearch()
	log.Printf("Finished refreshing %d feeds", len(feeds))

	w.reflock.Unlock()
}

func (w *Worker) worker(srcqueue <-chan storage.Feed, dstqueue chan<- []storage.Item) {
	for feed := range srcqueue {
		items, err := listItems(feed, w.db)
		if err != nil {
			w.db.SetFeedError(feed.Id, err)
		}
		dstqueue <- items
	}
}
