package tests

import (
	"database/sql"
	"fmt"
	"log"
	"maps"
	"reflect"
	"slices"
	"strconv"
	"testing"
	"testing/synctest"
	"time"

	"github.com/nkanaev/yarr/src/storage"
	"github.com/nkanaev/yarr/src/storage/model"
)

/*
- folder1
  - feed11
	- item111 (unread)
	- item112 (read)
	- item113 (starred)
  - feed12
	- item121 (unread)
	- item122 (read)
- folder2
  - feed21
    - item211 (read)
	- item212 (starred)
- feed01
  - item011 (unread)
  - item012 (read)
  - item013 (starred)
*/

type testItemScope struct {
	feed11, feed12   *model.Feed
	feed21, feed01   *model.Feed
	folder1, folder2 *model.Folder
	items map[string]model.Item
}

func MustGet[K comparable, V any](m map[K]V, key K) V {
    value, ok := m[key]
    if !ok {
        panic(fmt.Sprintf("key %v not found in map", key))
    }
    return value
}

func testItemsSetup(db storage.Storage) testItemScope {
	folder1 := db.CreateFolder("folder1")
	folder2 := db.CreateFolder("folder2")

	feed11 := db.CreateFeed(model.CreateFeedParams{Title: "feed11", FeedLink: "http://test.com/feed11.xml", FolderID: &folder1.Id})
	feed12 := db.CreateFeed(model.CreateFeedParams{Title: "feed12", FeedLink: "http://test.com/feed12.xml", FolderID: &folder1.Id})
	feed21 := db.CreateFeed(model.CreateFeedParams{Title: "feed21", FeedLink: "http://test.com/feed21.xml", FolderID: &folder2.Id})
	feed01 := db.CreateFeed(model.CreateFeedParams{Title: "feed01", FeedLink: "http://test.com/feed01.xml"})

	now := time.Now()
	items := map[string]model.Item{
		// feed11
		"item111": {
			GUID: "item111",
			FeedId: feed11.Id,
			Title: "title111",
			Date: now.Add(time.Hour * 24 * 1),
		},
		"item112": {
			GUID:   "item112",
			FeedId: feed11.Id,
			Title:  "title112",
			Date:   now.Add(time.Hour * 24 * 2),
			Status: model.READ,
		}, // read
		"item113": {
			GUID:   "item113",
			FeedId: feed11.Id,
			Title:  "title113",
			Date:   now.Add(time.Hour * 24 * 3),
			Status: model.STARRED,
		}, // starred
		// feed12
		"item121": {
			GUID: "item121",
			FeedId: feed12.Id,
			Title: "title121",
			Date: now.Add(time.Hour * 24 * 4),
		},
		"item122": {
			GUID:   "item122",
			FeedId: feed12.Id,
			Title:  "title122",
			Date:   now.Add(time.Hour * 24 * 5),
			Status: model.READ,
		}, // read
		// feed21
		"item211": {
			GUID:   "item211",
			FeedId: feed21.Id,
			Title:  "title211",
			Date:   now.Add(time.Hour * 24 * 6),
			Status: model.READ,
		}, // read
		"item212": {
			GUID:   "item212",
			FeedId: feed21.Id,
			Title:  "title212",
			Date:   now.Add(time.Hour * 24 * 7),
			Status: model.STARRED,
		}, // starred
		// feed01
		"item011": {
			GUID: "item011",
			FeedId: feed01.Id,
			Title: "title011",
			Date: now.Add(time.Hour * 24 * 8),
		},
		"item012": {
			GUID:   "item012",
			FeedId: feed01.Id,
			Title:  "title012",
			Date:   now.Add(time.Hour * 24 * 9),
			Status: model.READ,
		}, // read
		"item013": {
			GUID:   "item013",
			FeedId: feed01.Id,
			Title:  "title013",
			Date:   now.Add(time.Hour * 24 * 10),
			Status: model.STARRED,
		}, // starred
	}

	db.CreateItems(slices.Collect(maps.Values(items)))

	return testItemScope{
		feed11:  feed11,
		feed12:  feed12,
		feed21:  feed21,
		feed01:  feed01,
		folder1: folder1,
		folder2: folder2,
		items: items,
	}
}

func getItem(db storage.Storage, guid string) *model.Item {
	i := &model.Item{}
	err := db.db.QueryRow(`
		select
			i.id, i.guid, i.feed_id, i.title, i.link, i.content,
			i.date, i.status, i.media_links
		from items i
		where i.guid = :guid
	`, sql.Named("guid", guid)).Scan(
		&i.Id, &i.GUID, &i.FeedId, &i.Title, &i.Link, &i.Content,
		&i.Date, &i.Status, (*model.MediaLinks)(&i.MediaLinks),
	)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

func getItemGuids(items []model.Item) []string {
	guids := make([]string, 0)
	for _, item := range items {
		guids = append(guids, item.GUID)
	}
	return guids
}

func TestListItems(t *testing.T) {
	dbtest(t, func(t *testing.T, db storage.Storage) {
		scope := testItemsSetup(db)

		// filter by folder_id

		have := getItemGuids(db.ListItems(model.ItemFilter{FolderID: &scope.folder1.Id}, 10, false, false))
		want := []string{"item111", "item112", "item113", "item121", "item122"}
		if !reflect.DeepEqual(have, want) {
			t.Logf("want: %#v", want)
			t.Logf("have: %#v", have)
			t.Fail()
		}

		have = getItemGuids(db.ListItems(model.ItemFilter{FolderID: &scope.folder2.Id}, 10, false, false))
		want = []string{"item211", "item212"}
		if !reflect.DeepEqual(have, want) {
			t.Logf("want: %#v", want)
			t.Logf("have: %#v", have)
			t.Fail()
		}

		// filter by feed_id

		have = getItemGuids(db.ListItems(model.ItemFilter{FeedID: &scope.feed11.Id}, 10, false, false))
		want = []string{"item111", "item112", "item113"}
		if !reflect.DeepEqual(have, want) {
			t.Logf("want: %#v", want)
			t.Logf("have: %#v", have)
			t.Fail()
		}

		have = getItemGuids(db.ListItems(model.ItemFilter{FeedID: &scope.feed01.Id}, 10, false, false))
		want = []string{"item011", "item012", "item013"}
		if !reflect.DeepEqual(have, want) {
			t.Logf("want: %#v", want)
			t.Logf("have: %#v", have)
			t.Fail()
		}

		// filter by status

		var starred model.ItemStatus = model.STARRED
		have = getItemGuids(db.ListItems(model.ItemFilter{Status: &starred}, 10, false, false))
		want = []string{"item113", "item212", "item013"}
		if !reflect.DeepEqual(have, want) {
			t.Logf("want: %#v", want)
			t.Logf("have: %#v", have)
			t.Fail()
		}

		var unread model.ItemStatus = model.UNREAD
		have = getItemGuids(db.ListItems(model.ItemFilter{Status: &unread}, 10, false, false))
		want = []string{"item111", "item121", "item011"}
		if !reflect.DeepEqual(have, want) {
			t.Logf("want: %#v", want)
			t.Logf("have: %#v", have)
			t.Fail()
		}

		// limit

		have = getItemGuids(db.ListItems(model.ItemFilter{}, 2, false, false))
		want = []string{"item111", "item112"}
		if !reflect.DeepEqual(have, want) {
			t.Logf("want: %#v", want)
			t.Logf("have: %#v", have)
			t.Fail()
		}

		// filter by search
		search1 := "title111"
		have = getItemGuids(db.ListItems(model.ItemFilter{Search: &search1}, 4, true, false))
		want = []string{"item111"}
		if !reflect.DeepEqual(have, want) {
			t.Logf("want: %#v", want)
			t.Logf("have: %#v", have)
			t.Fail()
		}

		// sort by date
		have = getItemGuids(db.ListItems(model.ItemFilter{}, 4, true, false))
		want = []string{"item013", "item012", "item011", "item212"}
		if !reflect.DeepEqual(have, want) {
			t.Logf("want: %#v", want)
			t.Logf("have: %#v", have)
			t.Fail()
		}
	})
}

func TestListItemsPaginated(t *testing.T) {
	dbtest(t, func(t *testing.T, db storage.Storage) {
		scope := testItemsSetup(db)

		item012 := MustGet(scope.items, "item012")
		item121 := MustGet(scope.items, "item121")

		// all, newest first
		have := getItemGuids(db.ListItems(model.ItemFilter{After: &item012.Id}, 3, true, false))
		want := []string{"item011", "item212", "item211"}
		if !reflect.DeepEqual(have, want) {
			t.Logf("want: %#v", want)
			t.Logf("have: %#v", have)
			t.Fail()
		}

		// unread, newest first
		unread := model.UNREAD
		have = getItemGuids(
			db.ListItems(model.ItemFilter{After: &item012.Id, Status: &unread}, 3, true, false),
		)
		want = []string{"item011", "item121", "item111"}
		if !reflect.DeepEqual(have, want) {
			t.Logf("want: %#v", want)
			t.Logf("have: %#v", have)
			t.Fail()
		}

		// starred, oldest first
		starred := model.STARRED
		have = getItemGuids(
			db.ListItems(model.ItemFilter{After: &item121.Id, Status: &starred}, 3, false, false),
		)
		want = []string{"item212", "item013"}
		if !reflect.DeepEqual(have, want) {
			t.Logf("want: %#v", want)
			t.Logf("have: %#v", have)
			t.Fail()
		}
	})
}

func TestMarkItemsRead(t *testing.T) {
	// NOTE: starred items must not be marked as read
	var read model.ItemStatus = model.READ

	dbtest(t, func(t *testing.T, db1 storage.Storage) {
		testItemsSetup(db1)
		db1.MarkItemsRead(model.MarkFilter{})
		have := getItemGuids(db1.ListItems(model.ItemFilter{Status: &read}, 10, false, false))
		want := []string{
			"item111", "item112", "item121", "item122",
			"item211", "item011", "item012",
		}
		if !reflect.DeepEqual(have, want) {
			t.Logf("want: %#v", want)
			t.Logf("have: %#v", have)
			t.Fail()
		}
	})

	dbtest(t, func(t *testing.T, db2 storage.Storage) {
		scope2 := testItemsSetup(db2)
		db2.MarkItemsRead(model.MarkFilter{FolderID: &scope2.folder1.Id})
		have := getItemGuids(db2.ListItems(model.ItemFilter{Status: &read}, 10, false, false))
		want := []string{
			"item111", "item112", "item121", "item122",
			"item211", "item012",
		}
		if !reflect.DeepEqual(have, want) {
			t.Logf("want: %#v", want)
			t.Logf("have: %#v", have)
			t.Fail()
		}
	})

	dbtest(t, func(t *testing.T, db3 storage.Storage) {
		scope3 := testItemsSetup(db3)
		db3.MarkItemsRead(model.MarkFilter{FeedID: &scope3.feed11.Id})
		have := getItemGuids(db3.ListItems(model.ItemFilter{Status: &read}, 10, false, false))
		want := []string{
			"item111", "item112", "item122",
			"item211", "item012",
		}
		if !reflect.DeepEqual(have, want) {
			t.Logf("want: %#v", want)
			t.Logf("have: %#v", have)
			t.Fail()
		}
	})
}

func TestDeleteOldItems(t *testing.T) {
	now := time.Now().UTC()
	starred := model.STARRED
	dbtest(t, func(t *testing.T, db storage.Storage) {
		t.Run("keeps at least 50 items", func(t *testing.T) {
			feed := db.CreateFeed(model.CreateFeedParams{Title: "f", FeedLink: "http://f.xml"})
			items := make([]model.Item, 100)
			for i := range 100 {
				items[i] = model.Item{GUID: strconv.Itoa(i), FeedId: feed.Id, Date: now.Add(time.Duration(i) * time.Hour * 24)}
			}
			db.CreateItems(items)

			// // Set 1 recent (latest), 100 old (100 days ago)
			db.db.Exec(`update items set last_arrived = :la where guid = "99"`, sql.Named("la", now))
			db.db.Exec(`update items set last_arrived = :la where guid != "99"`, sql.Named("la", now.Add(-time.Hour*24*100)))

			db.DeleteOldItems()
			var have int
			db.db.QueryRow("select count(*) from items where feed_id = ?", feed.Id).Scan(&have)
			if have != 50 {
				t.Errorf("expected 50 items, have %d", have)
			}
		})

		t.Run("keeps all less than 90 days old", func(t *testing.T) {
			feed := db.CreateFeed(model.CreateFeedParams{Title: "f", FeedLink: "http://f.xml"})
			items := make([]model.Item, 100)
			for i := 0; i < 100; i++ {
				items[i] = model.Item{GUID: strconv.Itoa(i), FeedId: feed.Id, Date: now.Add(time.Duration(i) * time.Second)}
			}
			db.CreateItems(items)

			// Latest item at "now"
			// All others at 80 days ago (keep)
			db.db.Exec(`update items set last_arrived = :la where guid = "99"`, sql.Named("la", now))
			db.db.Exec(`update items set last_arrived = :la where guid != "99"`, sql.Named("la", now.Add(-time.Hour*24*80)))

			db.DeleteOldItems()
			var have int
			db.db.QueryRow("select count(*) from items where feed_id = ?", feed.Id).Scan(&have)
			if have != 100 {
				t.Errorf("expected 100 items, have %d", have)
			}
		})

		t.Run("keeps starred", func(t *testing.T) {
			feed := db.CreateFeed(model.CreateFeedParams{Title: "f", FeedLink: "http://f.xml"})
			items := make([]model.Item, 100)
			for i := 0; i < 100; i++ {
				items[i] = model.Item{GUID: strconv.Itoa(i), FeedId: feed.Id, Date: now.Add(time.Duration(i) * time.Second)}
			}
			db.CreateItems(items)

			// Set all to 100 days ago, except one recent
			db.db.Exec(`update items set last_arrived = :la`, sql.Named("la", now.Add(-time.Hour*24*100)))
			db.db.Exec(`update items set last_arrived = :la where guid = "99"`, sql.Named("la", now))
			// Star 10 old items that would otherwise be deleted (rn > 50 and old)
			db.db.Exec(`update items set status = :s where cast(guid as integer) < 10`, sql.Named("s", starred))

			db.DeleteOldItems()
			var have int
			db.db.QueryRow("select count(*) from items where feed_id = ?", feed.Id).Scan(&have)
			// 50 (limit) + 10 (starred) = 60 items should remain.
			if have != 60 {
				t.Errorf("expected 60 items, have %d", have)
			}
		})
	})
}

func TestCreateItemsLastArrived(t *testing.T) {
	dbtest(t, func(t *testing.T, db storage.Storage) {
		synctest.Test(t, func(t *testing.T) {
			feed := db.CreateFeed(model.CreateFeedParams{Title: "test feed", FeedLink: "http://example.com/feed"})

			item := model.Item{
				GUID:   "item1",
				FeedId: feed.Id,
				Title:  "Title 1",
				Date:   time.Now(),
			}

			// 1. Initial creation
			db.CreateItems([]model.Item{item})

			var lastArrived1 time.Time
			err := db.db.QueryRow("select last_arrived from items where guid = ?", item.GUID).Scan(&lastArrived1)
			if err != nil {
				t.Fatal(err)
			}

			time.Sleep(time.Second * 10)

			// 2. Update on conflict
			db.CreateItems([]model.Item{item})

			var lastArrived2 time.Time
			err = db.db.QueryRow("select last_arrived from items where guid = ?", item.GUID).Scan(&lastArrived2)
			if err != nil {
				t.Fatal(err)
			}

			if !lastArrived2.After(lastArrived1) {
				t.Errorf("expected last_arrived to be updated. old: %v, new: %v", lastArrived1, lastArrived2)
			}
		})
	})
}

func TestSearch(t *testing.T) {
	dbtest(t, func(t *testing.T, db storage.Storage) {
		feed := db.CreateFeed(model.CreateFeedParams{Title: "f", FeedLink: "http://f.xml"})

		db.CreateItems([]model.Item{
			{
				GUID:    "i1",
				FeedId:  feed.Id,
				Title:   "Hello World",
				Content: "This is a <b>test</b> of the <i>emergency</i> broadcast system.",
			},
			{
				GUID:    "i2",
				FeedId:  feed.Id,
				Title:   "FTS5 Unicode",
				Content: "Unicode support with characters like: Привет, 世界, 🚀",
			},
			{
				GUID:    "i3",
				FeedId:  feed.Id,
				Title:   "Hidden Tag",
				Content: `<div class="secret-class">Don't find me by my class name</div>`,
			},
		})

		// 1. Basic search
		s1 := "emergency"
		have := getItemGuids(db.ListItems(model.ItemFilter{Search: &s1}, 10, true, false))
		if !reflect.DeepEqual(have, []string{"i1"}) {
			t.Errorf("basic search failed: expected [i1], got %v", have)
		}

		// 2. HTML stripping: Should find text, but NOT the tags
		s2 := "test"
		have = getItemGuids(db.ListItems(model.ItemFilter{Search: &s2}, 10, true, false))
		if !reflect.DeepEqual(have, []string{"i1"}) {
			t.Errorf("html text search failed: expected [i1], got %v", have)
		}

		s3 := "secret-class"
		have = getItemGuids(db.ListItems(model.ItemFilter{Search: &s3}, 10, true, false))
		if len(have) > 0 {
			t.Errorf("html tag search should have failed but found: %v", have)
		}

		// 3. Multi-word (AND)
		s4 := "broadcast system"
		have = getItemGuids(db.ListItems(model.ItemFilter{Search: &s4}, 10, true, false))
		if !reflect.DeepEqual(have, []string{"i1"}) {
			t.Errorf("multi-word search failed: expected [i1], got %v", have)
		}

		// 4. Unicode
		s5 := "Привет"
		have = getItemGuids(db.ListItems(model.ItemFilter{Search: &s5}, 10, true, false))
		if !reflect.DeepEqual(have, []string{"i2"}) {
			t.Errorf("unicode search failed: expected [i2], got %v", have)
		}

		s6 := "世界"
		have = getItemGuids(db.ListItems(model.ItemFilter{Search: &s6}, 10, true, false))
		if !reflect.DeepEqual(have, []string{"i2"}) {
			t.Errorf("unicode search (CJK) failed: expected [i2], got %v", have)
		}

		// 5. Trigger: Update
		db.db.Exec("update items set title = 'Updated Title' where guid = 'i1'")
		s7 := "Updated"
		have = getItemGuids(db.ListItems(model.ItemFilter{Search: &s7}, 10, true, false))
		if !reflect.DeepEqual(have, []string{"i1"}) {
			t.Errorf("update trigger failed: expected [i1], got %v", have)
		}

		// 6. Trigger: Delete
		db.db.Exec("delete from items where guid = 'i1'")
		have = getItemGuids(db.ListItems(model.ItemFilter{Search: &s7}, 10, true, false))
		if len(have) > 0 {
			t.Errorf("delete trigger failed: found deleted item: %v", have)
		}
	})
}
