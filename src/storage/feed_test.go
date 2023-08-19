package storage

import (
	"reflect"
	"testing"
)

func TestCreateFeed(t *testing.T) {
	db := testDB()
	feed1 := db.CreateFeed("title", "", "http://example.com", "http://example.com/feed.xml", nil)
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
	feed1 := db.CreateFeed("title", "", "", "http://example1.com/feed.xml", nil)
	if feed1 == nil || feed1.Id == 0 {
		t.Fatal("expected feed")
	}

    for i := 0; i < 10; i++ {
	    db.CreateFeed("title", "", "", "http://example2.com/feed.xml", nil)
    }

	feed2 := db.CreateFeed("title", "", "http://example.com", "http://example1.com/feed.xml", nil)
	if feed1.Id != feed2.Id {
        t.Fatalf("expected the same feed.\nwant: %#v\nhave: %#v", feed1, feed2)
	}
}

func TestReadFeed(t *testing.T) {
	db := testDB()
	if db.GetFeed(100500) != nil {
		t.Fatal("cannot get nonexistent feed")
	}

	feed1 := db.CreateFeed("feed 1", "", "http://example1.com", "http://example1.com/feed.xml", nil)
	feed2 := db.CreateFeed("feed 2", "", "http://example2.com", "http://example2.com/feed.xml", nil)
	feeds := db.ListFeeds()
	if !reflect.DeepEqual(feeds, []Feed{*feed1, *feed2}) {
		t.Fatalf("invalid feed list: %#v", feeds)
	}
}

func TestUpdateFeed(t *testing.T) {
	db := testDB()
	feed1 := db.CreateFeed("feed 1", "", "http://example1.com", "http://example1.com/feed.xml", nil)
	folder := db.CreateFolder("test")
	icon := []byte("icon")

	db.RenameFeed(feed1.Id, "newtitle")
	db.UpdateFeedFolder(feed1.Id, &folder.Id)
	db.UpdateFeedIcon(feed1.Id, &icon)

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
	feed1 := db.CreateFeed("title", "", "http://example.com", "http://example.com/feed.xml", nil)

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
