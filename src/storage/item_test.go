package storage

import (
	"log"
	"reflect"
	"strconv"
	"testing"
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
		{GUID: "item112", FeedId: feed11.Id, Title: "title112", Date: now.Add(time.Hour * 24 * 2)}, // read
		{GUID: "item113", FeedId: feed11.Id, Title: "title113", Date: now.Add(time.Hour * 24 * 3)}, // starred
		// feed12
		{GUID: "item121", FeedId: feed12.Id, Title: "title121", Date: now.Add(time.Hour * 24 * 4)},
		{GUID: "item122", FeedId: feed12.Id, Title: "title122", Date: now.Add(time.Hour * 24 * 5)}, // read
		// feed21
		{GUID: "item211", FeedId: feed21.Id, Title: "title211", Date: now.Add(time.Hour * 24 * 6)}, // read
		{GUID: "item212", FeedId: feed21.Id, Title: "title212", Date: now.Add(time.Hour * 24 * 7)}, // starred
		// feed01
		{GUID: "item011", FeedId: feed01.Id, Title: "title011", Date: now.Add(time.Hour * 24 * 8)},
		{GUID: "item012", FeedId: feed01.Id, Title: "title012", Date: now.Add(time.Hour * 24 * 9)},  // read
		{GUID: "item013", FeedId: feed01.Id, Title: "title013", Date: now.Add(time.Hour * 24 * 10)}, // starred
	})
	db.db.Exec(`update items set status = ? where guid in ("item112", "item122", "item211", "item012")`, READ)
	db.db.Exec(`update items set status = ? where guid in ("item113", "item212", "item013")`, STARRED)

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
			i.date, i.status, i.image, i.podcast_url
		from items i
		where i.guid = ?
	`, guid).Scan(
		&i.Id, &i.GUID, &i.FeedId, &i.Title, &i.Link, &i.Content,
		&i.Date, &i.Status, &i.ImageURL, &i.AudioURL,
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

	have := getItemGuids(db.ListItems(ItemFilter{FolderID: &scope.folder1.Id}, 10, false))
	want := []string{"item111", "item112", "item113", "item121", "item122"}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}

	have = getItemGuids(db.ListItems(ItemFilter{FolderID: &scope.folder2.Id}, 10, false))
	want = []string{"item211", "item212"}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}

	// filter by feed_id

	have = getItemGuids(db.ListItems(ItemFilter{FeedID: &scope.feed11.Id}, 10, false))
	want = []string{"item111", "item112", "item113"}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}

	have = getItemGuids(db.ListItems(ItemFilter{FeedID: &scope.feed01.Id}, 10, false))
	want = []string{"item011", "item012", "item013"}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}

	// filter by status

	var starred ItemStatus = STARRED
	have = getItemGuids(db.ListItems(ItemFilter{Status: &starred}, 10, false))
	want = []string{"item113", "item212", "item013"}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}

	var unread ItemStatus = UNREAD
	have = getItemGuids(db.ListItems(ItemFilter{Status: &unread}, 10, false))
	want = []string{"item111", "item121", "item011"}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}

	// limit

	have = getItemGuids(db.ListItems(ItemFilter{}, 2, false))
	want = []string{"item111", "item112"}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}

	// filter by search
	db.SyncSearch()
	search1 := "title111"
	have = getItemGuids(db.ListItems(ItemFilter{Search: &search1}, 4, true))
	want = []string{"item111"}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}

	// sort by date
	have = getItemGuids(db.ListItems(ItemFilter{}, 4, true))
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
	have := getItemGuids(db.ListItems(ItemFilter{After: &item012.Id}, 3, true))
	want := []string{"item011", "item212", "item211"}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}

	// unread, newest first
	unread := UNREAD
	have = getItemGuids(db.ListItems(ItemFilter{After: &item012.Id, Status: &unread}, 3, true))
	want = []string{"item011", "item121", "item111"}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}

	// starred, oldest first
	starred := STARRED
	have = getItemGuids(db.ListItems(ItemFilter{After: &item121.Id, Status: &starred}, 3, false))
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
	have := getItemGuids(db1.ListItems(ItemFilter{Status: &read}, 10, false))
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
	have = getItemGuids(db2.ListItems(ItemFilter{Status: &read}, 10, false))
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
	have = getItemGuids(db3.ListItems(ItemFilter{Status: &read}, 10, false))
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
	extraItems := 10

	now := time.Now()
	db := testDB()
	feed := db.CreateFeed("feed", "", "", "http://test.com/feed11.xml", nil)

	items := make([]Item, 0)
	for i := 0; i < itemsKeepSize+extraItems; i++ {
		istr := strconv.Itoa(i)
		items = append(items, Item{
			GUID:   istr,
			FeedId: feed.Id,
			Title:  istr,
			Date:   now.Add(time.Hour * time.Duration(i)),
		})
	}
	db.CreateItems(items)
	
	db.SetFeedSize(feed.Id, itemsKeepSize)
	var feedSize int
	err := db.db.QueryRow(
		`select size from feed_sizes where feed_id = ?`, feed.Id,
	).Scan(&feedSize)
	if err != nil {
		t.Fatal(err)
	}
	if feedSize != itemsKeepSize {
		t.Fatalf(
			"expected feed size to get updated\nwant: %d\nhave: %d", 
			itemsKeepSize+extraItems,
			feedSize,
		)
	}

	// expire only the first 3 articles
	_, err = db.db.Exec(
		`update items set date_arrived = ?
		where id in (select id from items limit 3)`,
		now.Add(-time.Hour*time.Duration(itemsKeepDays*24)),
	)
	if err != nil {
		t.Fatal(err)
	}

	db.DeleteOldItems()
	feedItems := db.ListItems(ItemFilter{FeedID: &feed.Id}, 1000, false)
	if len(feedItems) != len(items)-3 {
		t.Fatalf(
			"invalid number of old items kept\nwant: %d\nhave: %d",
			len(items)-3,
			len(feedItems),
		)
	}
}
