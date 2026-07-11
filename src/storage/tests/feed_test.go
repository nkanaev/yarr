package tests

import (
	"reflect"
	"testing"

	"github.com/nkanaev/yarr/src/storage"
	"github.com/nkanaev/yarr/src/storage/model"
)

func TestCreateFeed(t *testing.T) {
	dbtest(t, func(t *testing.T, db storage.Storage) {
		feed1 := db.CreateFeed(model.CreateFeedParams{Title: "title", Link: "http://example.com", FeedLink: "http://example.com/feed.xml"})
		if feed1 == nil || feed1.Id == 0 {
			t.Fatal("expected feed")
		}
		feed2 := db.GetFeed(feed1.Id)
		if feed2 == nil || !reflect.DeepEqual(feed1, feed2) {
			t.Fatal("invalid feed")
		}
	})
}

func TestCreateFeedSameLink(t *testing.T) {
	dbtest(t, func(t *testing.T, db storage.Storage) {
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
	})
}

func TestReadFeed(t *testing.T) {
	dbtest(t, func(t *testing.T, db storage.Storage) {
		if db.GetFeed(100500) != nil {
			t.Fatal("cannot get nonexistent feed")
		}

		feed1 := db.CreateFeed(model.CreateFeedParams{Title: "feed 1", Link: "http://example1.com", FeedLink: "http://example1.com/feed.xml"})
		feed2 := db.CreateFeed(model.CreateFeedParams{Title: "feed 2", Link: "http://example2.com", FeedLink: "http://example2.com/feed.xml"})
		feeds := db.ListFeeds()
		if !reflect.DeepEqual(feeds, []model.Feed{*feed1, *feed2}) {
			t.Fatalf("invalid feed list: %#v", feeds)
		}
	})
}

func TestUpdateFeed(t *testing.T) {
	dbtest(t, func(t *testing.T, db storage.Storage) {
		feed1 := db.CreateFeed(model.CreateFeedParams{Title: "feed 1", Link: "http://example1.com", FeedLink: "http://example1.com/feed.xml"})
		folder := db.CreateFolder("test")
		icon := model.Icon("icon")

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
		if feed2.Icon == nil || string(*feed2.Icon) != "icon" {
			t.Error("invalid icon")
		}
	})
}

func TestFeedStats(t *testing.T) {
	dbtest(t, func(t *testing.T, db storage.Storage) {
		// empty
		stats := db.FeedStats()
		if len(stats) != 0 {
			t.Errorf("expected 0 stats, got %d", len(stats))
		}

		scope := testItemsSetup(db)

		stats = db.FeedStats()
		statByFeed := make(map[int64]model.FeedStat)
		for _, s := range stats {
			statByFeed[s.FeedId] = s
		}

		for _, tc := range []struct {
			feedID  int64
			unread  int64
			starred int64
		}{
			{scope.feed11.Id, 1, 1},
			{scope.feed12.Id, 1, 0},
			{scope.feed21.Id, 0, 1},
			{scope.feed01.Id, 1, 1},
		} {
			s, ok := statByFeed[tc.feedID]
			if !ok {
				t.Errorf("feed %d missing from stats", tc.feedID)
				continue
			}
			if s.UnreadCount != tc.unread {
				t.Errorf("feed %d unread: expected %d, got %d", tc.feedID, tc.unread, s.UnreadCount)
			}
			if s.StarredCount != tc.starred {
				t.Errorf("feed %d starred: expected %d, got %d", tc.feedID, tc.starred, s.StarredCount)
			}
		}

		// mark feed11 read, verify stats update
		db.MarkItemsRead(model.MarkFilter{FeedID: &scope.feed11.Id})
		stats = db.FeedStats()
		statByFeed = make(map[int64]model.FeedStat)
		for _, s := range stats {
			statByFeed[s.FeedId] = s
		}
		if s := statByFeed[scope.feed11.Id]; s.UnreadCount != 0 {
			t.Errorf("feed11 unread after mark-read: expected 0, got %d", s.UnreadCount)
		}
		if s := statByFeed[scope.feed11.Id]; s.StarredCount != 1 {
			t.Errorf("feed11 starred after mark-read: expected 1, got %d", s.StarredCount)
		}
	})
}

func TestDeleteFeed(t *testing.T) {
	dbtest(t, func(t *testing.T, db storage.Storage) {
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
	})
}
