package worker

import (
	"context"
	"testing"
	"time"

	"github.com/nkanaev/yarr/src/client"
	"github.com/nkanaev/yarr/src/storage"
	"github.com/nkanaev/yarr/src/storage/model"
)

func newTestDB(t *testing.T) storage.Storage {
	t.Helper()
	db, err := storage.New(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestFeedStateStoreGet(t *testing.T) {
	db := newTestDB(t)
	feed := db.CreateFeed(model.CreateFeedParams{Title: "Test", FeedLink: "http://example.com"})
	lmod := "Wed, 04 Mar 2026 05:06:07 GMT"
	etag := "v1"
	db.UpdateFeedState(feed.Id, model.UpdateFeedStateParams{HTTPLastModified: &lmod, HTTPEtag: &etag})

	state, err := NewFeedStateStore(db).Get(contextSetFeedID(context.Background(), feed.Id))
	if err != nil {
		t.Fatal(err)
	}
	if want := (client.State{ETag: etag, LastModified: lmod}); state != want {
		t.Errorf("state = %+v; want %+v", state, want)
	}
}

func TestFeedStateStoreGetMissingState(t *testing.T) {
	db := newTestDB(t)
	feed := db.CreateFeed(model.CreateFeedParams{Title: "Test", FeedLink: "http://example.com"})

	state, err := NewFeedStateStore(db).Get(contextSetFeedID(context.Background(), feed.Id))
	if err != nil {
		t.Fatal(err)
	}
	if state != (client.State{}) {
		t.Errorf("state = %+v; want empty", state)
	}
}

func TestFeedStateStoreMissingFeedID(t *testing.T) {
	store := NewFeedStateStore(newTestDB(t))
	if _, err := store.Get(context.Background()); err == nil {
		t.Error("Get: expected error without a feed id in context")
	}
	if err := store.Put(context.Background(), client.State{ETag: "x"}); err == nil {
		t.Error("Put: expected error without a feed id in context")
	}
}

func TestFeedStateStorePut(t *testing.T) {
	db := newTestDB(t)
	feed := db.CreateFeed(model.CreateFeedParams{Title: "Test", FeedLink: "http://example.com"})

	before := time.Now().UTC()
	err := NewFeedStateStore(db).Put(contextSetFeedID(context.Background(), feed.Id), client.State{
		ETag:         "v2",
		LastModified: "Thu, 07 May 2026 08:09:10 GMT",
	})
	if err != nil {
		t.Fatal(err)
	}

	state, err := db.GetFeedState(feed.Id)
	if err != nil {
		t.Fatal(err)
	}
	if state == nil {
		t.Fatal("expected state, got nil")
	}
	if state.HTTPEtag != "v2" {
		t.Errorf("HTTPEtag = %q; want %q", state.HTTPEtag, "v2")
	}
	if state.HTTPLastModified != "Thu, 07 May 2026 08:09:10 GMT" {
		t.Errorf("HTTPLastModified = %q", state.HTTPLastModified)
	}
	if state.LastRefreshed.Before(before) {
		t.Errorf("LastRefreshed = %v; want >= %v", state.LastRefreshed, before)
	}
}
