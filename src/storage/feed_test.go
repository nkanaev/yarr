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

func TestArchiveFeed(t *testing.T) {
	db := testDB()
	feed := db.CreateFeed("test feed", "", "http://example.com", "http://example.com/feed.xml", nil)

	// Initially, feed should not be archived
	if feed.Archived {
		t.Error("new feed should not be archived")
	}

	// Archive the feed
	if !db.ArchiveFeed(feed.Id) {
		t.Fatal("failed to archive feed")
	}

	// Verify feed is archived
	archivedFeed := db.GetFeed(feed.Id)
	if archivedFeed == nil {
		t.Fatal("archived feed should still exist")
	}
	if !archivedFeed.Archived {
		t.Error("feed should be archived")
	}

	// Test archiving non-existent feed (SQL UPDATE succeeds even if no rows affected)
	// This behavior is consistent with other operations in the codebase
	if !db.ArchiveFeed(999999) {
		t.Error("archive operation should succeed even for non-existent feed")
	}
}

func TestUnarchiveFeed(t *testing.T) {
	db := testDB()
	feed := db.CreateFeed("test feed", "", "http://example.com", "http://example.com/feed.xml", nil)

	// Archive the feed first
	if !db.ArchiveFeed(feed.Id) {
		t.Fatal("failed to archive feed")
	}

	// Verify it's archived
	archivedFeed := db.GetFeed(feed.Id)
	if !archivedFeed.Archived {
		t.Fatal("feed should be archived")
	}

	// Unarchive the feed
	if !db.UnarchiveFeed(feed.Id) {
		t.Fatal("failed to unarchive feed")
	}

	// Verify feed is unarchived
	unarchivedFeed := db.GetFeed(feed.Id)
	if unarchivedFeed == nil {
		t.Fatal("unarchived feed should still exist")
	}
	if unarchivedFeed.Archived {
		t.Error("feed should not be archived")
	}

	// Test unarchiving non-existent feed (SQL UPDATE succeeds even if no rows affected)
	// This behavior is consistent with other operations in the codebase
	if !db.UnarchiveFeed(999999) {
		t.Error("unarchive operation should succeed even for non-existent feed")
	}
}

func TestListFeedsWithArchived(t *testing.T) {
	db := testDB()
	feed1 := db.CreateFeed("active feed", "", "http://example1.com", "http://example1.com/feed.xml", nil)
	feed2 := db.CreateFeed("archived feed", "", "http://example2.com", "http://example2.com/feed.xml", nil)

	// Archive the second feed
	db.ArchiveFeed(feed2.Id)

	// List all feeds
	feeds := db.ListFeeds()
	if len(feeds) != 2 {
		t.Fatalf("expected 2 feeds, got %d", len(feeds))
	}

	// Find feeds in the list
	var activeFeed, archivedFeed *Feed
	for i := range feeds {
		if feeds[i].Id == feed1.Id {
			activeFeed = &feeds[i]
		} else if feeds[i].Id == feed2.Id {
			archivedFeed = &feeds[i]
		}
	}

	if activeFeed == nil {
		t.Fatal("active feed not found in list")
	}
	if archivedFeed == nil {
		t.Fatal("archived feed not found in list")
	}

	// Check archived status
	if activeFeed.Archived {
		t.Error("active feed should not be archived")
	}
	if !archivedFeed.Archived {
		t.Error("archived feed should be archived")
	}
}
