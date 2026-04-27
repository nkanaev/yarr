package storage

import (
	"database/sql"
	"log"
	"reflect"
	"strconv"
	"testing"
	"testing/synctest"
	"time"
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
	feed11, feed12   *Feed
	feed21, feed01   *Feed
	folder1, folder2 *Folder
}

func testItemsSetup(db *Storage) testItemScope {
	folder1 := db.CreateFolder("folder1")
	folder2 := db.CreateFolder("folder2")

	feed11 := db.CreateFeed("feed11", "", "", "http://test.com/feed11.xml", &folder1.Id)
	feed12 := db.CreateFeed("feed12", "", "", "http://test.com/feed12.xml", &folder1.Id)
	feed21 := db.CreateFeed("feed21", "", "", "http://test.com/feed21.xml", &folder2.Id)
	feed01 := db.CreateFeed("feed01", "", "", "http://test.com/feed01.xml", nil)

	now := time.Now()
	db.CreateItems([]Item{
		// feed11
		{GUID: "item111", FeedId: feed11.Id, Title: "title111", Date: now.Add(time.Hour * 24 * 1)},
		{
			GUID:   "item112",
			FeedId: feed11.Id,
			Title:  "title112",
			Date:   now.Add(time.Hour * 24 * 2),
		}, // read
		{
			GUID:   "item113",
			FeedId: feed11.Id,
			Title:  "title113",
			Date:   now.Add(time.Hour * 24 * 3),
		}, // starred
		// feed12
		{GUID: "item121", FeedId: feed12.Id, Title: "title121", Date: now.Add(time.Hour * 24 * 4)},
		{
			GUID:   "item122",
			FeedId: feed12.Id,
			Title:  "title122",
			Date:   now.Add(time.Hour * 24 * 5),
		}, // read
		// feed21
		{
			GUID:   "item211",
			FeedId: feed21.Id,
			Title:  "title211",
			Date:   now.Add(time.Hour * 24 * 6),
		}, // read
		{
			GUID:   "item212",
			FeedId: feed21.Id,
			Title:  "title212",
			Date:   now.Add(time.Hour * 24 * 7),
		}, // starred
		// feed01
		{GUID: "item011", FeedId: feed01.Id, Title: "title011", Date: now.Add(time.Hour * 24 * 8)},
		{
			GUID:   "item012",
			FeedId: feed01.Id,
			Title:  "title012",
			Date:   now.Add(time.Hour * 24 * 9),
		}, // read
		{
			GUID:   "item013",
			FeedId: feed01.Id,
			Title:  "title013",
			Date:   now.Add(time.Hour * 24 * 10),
		}, // starred
	})
	db.db.Exec(
		`update items set status = :status where guid in ("item112", "item122", "item211", "item012")`,
		sql.Named("status", READ),
	)
	db.db.Exec(
		`update items set status = :status where guid in ("item113", "item212", "item013")`,
		sql.Named("status", STARRED),
	)

	return testItemScope{
		feed11:  feed11,
		feed12:  feed12,
		feed21:  feed21,
		feed01:  feed01,
		folder1: folder1,
		folder2: folder2,
	}
}

func getItem(db *Storage, guid string) *Item {
	i := &Item{}
	err := db.db.QueryRow(`
		select
			i.id, i.guid, i.feed_id, i.title, i.link, i.content,
			i.date, i.status, i.media_links
		from items i
		where i.guid = :guid
	`, sql.Named("guid", guid)).Scan(
		&i.Id, &i.GUID, &i.FeedId, &i.Title, &i.Link, &i.Content,
		&i.Date, &i.Status, &i.MediaLinks,
	)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

func getItemGuids(items []Item) []string {
	guids := make([]string, 0)
	for _, item := range items {
		guids = append(guids, item.GUID)
	}
	return guids
}

func TestListItems(t *testing.T) {
	db := testDB()
	scope := testItemsSetup(db)

	// filter by folder_id

	have := getItemGuids(db.ListItems(ItemFilter{FolderID: &scope.folder1.Id}, 10, false, false))
	want := []string{"item111", "item112", "item113", "item121", "item122"}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}

	have = getItemGuids(db.ListItems(ItemFilter{FolderID: &scope.folder2.Id}, 10, false, false))
	want = []string{"item211", "item212"}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}

	// filter by feed_id

	have = getItemGuids(db.ListItems(ItemFilter{FeedID: &scope.feed11.Id}, 10, false, false))
	want = []string{"item111", "item112", "item113"}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}

	have = getItemGuids(db.ListItems(ItemFilter{FeedID: &scope.feed01.Id}, 10, false, false))
	want = []string{"item011", "item012", "item013"}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}

	// filter by status

	var starred ItemStatus = STARRED
	have = getItemGuids(db.ListItems(ItemFilter{Status: &starred}, 10, false, false))
	want = []string{"item113", "item212", "item013"}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}

	var unread ItemStatus = UNREAD
	have = getItemGuids(db.ListItems(ItemFilter{Status: &unread}, 10, false, false))
	want = []string{"item111", "item121", "item011"}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}

	// limit

	have = getItemGuids(db.ListItems(ItemFilter{}, 2, false, false))
	want = []string{"item111", "item112"}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}

	// filter by search
	db.SyncSearch()
	search1 := "title111"
	have = getItemGuids(db.ListItems(ItemFilter{Search: &search1}, 4, true, false))
	want = []string{"item111"}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}

	// sort by date
	have = getItemGuids(db.ListItems(ItemFilter{}, 4, true, false))
	want = []string{"item013", "item012", "item011", "item212"}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}
}

func TestListItemsPaginated(t *testing.T) {
	db := testDB()
	testItemsSetup(db)

	item012 := getItem(db, "item012")
	item121 := getItem(db, "item121")

	// all, newest first
	have := getItemGuids(db.ListItems(ItemFilter{After: &item012.Id}, 3, true, false))
	want := []string{"item011", "item212", "item211"}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}

	// unread, newest first
	unread := UNREAD
	have = getItemGuids(
		db.ListItems(ItemFilter{After: &item012.Id, Status: &unread}, 3, true, false),
	)
	want = []string{"item011", "item121", "item111"}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}

	// starred, oldest first
	starred := STARRED
	have = getItemGuids(
		db.ListItems(ItemFilter{After: &item121.Id, Status: &starred}, 3, false, false),
	)
	want = []string{"item212", "item013"}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}
}

func TestMarkItemsRead(t *testing.T) {
	// NOTE: starred items must not be marked as read
	var read ItemStatus = READ

	db1 := testDB()
	testItemsSetup(db1)
	db1.MarkItemsRead(MarkFilter{})
	have := getItemGuids(db1.ListItems(ItemFilter{Status: &read}, 10, false, false))
	want := []string{
		"item111", "item112", "item121", "item122",
		"item211", "item011", "item012",
	}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}

	db2 := testDB()
	scope2 := testItemsSetup(db2)
	db2.MarkItemsRead(MarkFilter{FolderID: &scope2.folder1.Id})
	have = getItemGuids(db2.ListItems(ItemFilter{Status: &read}, 10, false, false))
	want = []string{
		"item111", "item112", "item121", "item122",
		"item211", "item012",
	}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}

	db3 := testDB()
	scope3 := testItemsSetup(db3)
	db3.MarkItemsRead(MarkFilter{FeedID: &scope3.feed11.Id})
	have = getItemGuids(db3.ListItems(ItemFilter{Status: &read}, 10, false, false))
	want = []string{
		"item111", "item112", "item122",
		"item211", "item012",
	}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}
}

func TestDeleteOldItems(t *testing.T) {
	now := time.Now().UTC()
	starred := STARRED

	t.Run("keeps at least 50 items", func(t *testing.T) {
		db := testDB()
		feed := db.CreateFeed("f", "", "", "http://f.xml", nil)
		items := make([]Item, 100)
		for i := range 100 {
			items[i] = Item{GUID: strconv.Itoa(i), FeedId: feed.Id, Date: now.Add(time.Duration(i) * time.Hour * 24)}
		}
		db.CreateItems(items)

		// // Set 1 recent (latest), 100 old (100 days ago)
		db.db.Exec(`update items set last_arrived = :la where guid = "99"`, sql.Named("la", now))
		db.db.Exec(`update items set last_arrived = :la where guid != "99"`, sql.Named("la", now.Add(-time.Hour*24*100)))

		db.DeleteOldItems()
		have := db.CountItems(ItemFilter{FeedID: &feed.Id})
		if have != 50 {
			t.Errorf("expected 50 items, have %d", have)
		}
	})

	t.Run("keeps all less than 90 days old", func(t *testing.T) {
		db := testDB()
		feed := db.CreateFeed("f", "", "", "http://f.xml", nil)
		items := make([]Item, 100)
		for i := 0; i < 100; i++ {
			items[i] = Item{GUID: strconv.Itoa(i), FeedId: feed.Id, Date: now.Add(time.Duration(i) * time.Second)}
		}
		db.CreateItems(items)

		// Latest item at "now"
		// All others at 80 days ago (keep)
		db.db.Exec(`update items set last_arrived = :la where guid = "99"`, sql.Named("la", now))
		db.db.Exec(`update items set last_arrived = :la where guid != "99"`, sql.Named("la", now.Add(-time.Hour*24*80)))

		db.DeleteOldItems()
		have := db.CountItems(ItemFilter{FeedID: &feed.Id})
		if have != 100 {
			t.Errorf("expected 100 items, have %d", have)
		}
	})

	t.Run("keeps starred", func(t *testing.T) {
		db := testDB()
		feed := db.CreateFeed("f", "", "", "http://f.xml", nil)
		items := make([]Item, 100)
		for i := 0; i < 100; i++ {
			items[i] = Item{GUID: strconv.Itoa(i), FeedId: feed.Id, Date: now.Add(time.Duration(i) * time.Second)}
		}
		db.CreateItems(items)

		// Set all to 100 days ago, except one recent
		db.db.Exec(`update items set last_arrived = :la`, sql.Named("la", now.Add(-time.Hour*24*100)))
		db.db.Exec(`update items set last_arrived = :la where guid = "99"`, sql.Named("la", now))
		// Star 10 old items that would otherwise be deleted (rn > 50 and old)
		db.db.Exec(`update items set status = :s where cast(guid as integer) < 10`, sql.Named("s", starred))

		db.DeleteOldItems()
		have := db.CountItems(ItemFilter{FeedID: &feed.Id})
		// 50 (limit) + 10 (starred) = 60 items should remain.
		if have != 60 {
			t.Errorf("expected 60 items, have %d", have)
		}
	})
}



func TestCreateItemsLastArrived(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		db := testDB()
		defer db.db.Close()
		feed := db.CreateFeed("test feed", "", "", "http://example.com/feed", nil)

		item := Item{
			GUID:   "item1",
			FeedId: feed.Id,
			Title:  "Title 1",
			Date:   time.Now(),
		}

		// 1. Initial creation
		db.CreateItems([]Item{item})

		var lastArrived1 time.Time
		err := db.db.QueryRow("select last_arrived from items where guid = ?", item.GUID).Scan(&lastArrived1)
		if err != nil {
			t.Fatal(err)
		}

		time.Sleep(time.Second * 10)

		// 2. Update on conflict
		db.CreateItems([]Item{item})

		var lastArrived2 time.Time
		err = db.db.QueryRow("select last_arrived from items where guid = ?", item.GUID).Scan(&lastArrived2)
		if err != nil {
			t.Fatal(err)
		}

		if !lastArrived2.After(lastArrived1) {
			t.Errorf("expected last_arrived to be updated. old: %v, new: %v", lastArrived1, lastArrived2)
		}
	})
}
