package tests

import (
	"reflect"
	"testing"

	"github.com/nkanaev/yarr/src/storage/model"
)

func TestCreateFeed(t *testing.T) {
	db := testDB()
	feed1 := db.CreateFeed(model.CreateFeedParams{Title: "title", Link: "http://example.com", FeedLink: "http://example.com/feed.xml"})
	if feed1 == nil || feed1.Id == 0 {
		t.Fatal("expected feed")
	}
	feed2 := db.GetFeed(feed1.Id)
	if feed2 == nil || !reflect.DeepEqual(feed1, feed2) {
		t.Fatal("invalid feed")
	}
}

func TestCreateFeedSameLink(t *testing.T) {
	db := testDB()
	feed1 := db.CreateFeed(model.CreateFeedParams{Title: "title", FeedLink: "http://example1.com/feed.xml"})
	if feed1 == nil || feed1.Id == 0 {
		t.Fatal("expected feed")
	}

	for range 10 {
		db.CreateFeed(model.CreateFeedParams{Title: "title", FeedLink: "http://example2.com/feed.xml"})
	}

	feed2 := db.CreateFeed(model.CreateFeedParams{Title: "title", Link: "http://example.com", FeedLink: "http://example1.com/feed.xml"})
	if feed1.Id != feed2.Id {
		t.Fatalf("expected the same feed.\nwant: %#v\nhave: %#v", feed1, feed2)
	}
}

func TestReadFeed(t *testing.T) {
	db := testDB()
	if db.GetFeed(100500) != nil {
		t.Fatal("cannot get nonexistent feed")
	}

	feed1 := db.CreateFeed(model.CreateFeedParams{Title: "feed 1", Link: "http://example1.com", FeedLink: "http://example1.com/feed.xml"})
	feed2 := db.CreateFeed(model.CreateFeedParams{Title: "feed 2", Link: "http://example2.com", FeedLink: "http://example2.com/feed.xml"})
	feeds := db.ListFeeds()
	if !reflect.DeepEqual(feeds, []model.Feed{*feed1, *feed2}) {
		t.Fatalf("invalid feed list: %#v", feeds)
	}
}

func TestUpdateFeed(t *testing.T) {
	db := testDB()
	feed1 := db.CreateFeed(model.CreateFeedParams{Title: "feed 1", Link: "http://example1.com", FeedLink: "http://example1.com/feed.xml"})
	folder := db.CreateFolder("test")
	icon := []byte("icon")

	title := "newtitle"
	db.UpdateFeed(feed1.Id, model.UpdateFeedParams{
		Title:    &title,
		FolderID: model.SetNullable(&folder.Id),
		Icon:     model.SetNullable(&icon),
	})

	feed2 := db.GetFeed(feed1.Id)
	if feed2.Title != "newtitle" {
		t.Error("invalid title")
	}
	if feed2.FolderId == nil || *feed2.FolderId != folder.Id {
		t.Error("invalid folder")
	}
	if !feed2.HasIcon || string(*feed2.Icon) != "icon" {
		t.Error("invalid icon")
	}
}

func TestDeleteFeed(t *testing.T) {
	db := testDB()
	feed1 := db.CreateFeed(model.CreateFeedParams{Title: "title", Link: "http://example.com", FeedLink: "http://example.com/feed.xml"})

	if db.DeleteFeed(100500) {
		t.Error("cannot delete what does not exist")
	}

	if !db.DeleteFeed(feed1.Id) {
		t.Fatal("did not delete existing feed")
	}
	if db.GetFeed(feed1.Id) != nil {
		t.Fatal("feed still exists")
	}
}
