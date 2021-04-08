package storage

import (
	"reflect"
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
		{GUID: "item112", FeedId: feed11.Id, Title: "title112", Date: now.Add(time.Hour * 24 * 2)},
		{GUID: "item113", FeedId: feed11.Id, Title: "title113", Date: now.Add(time.Hour * 24 * 3)},
		// feed12
		{GUID: "item121", FeedId: feed12.Id, Title: "title121", Date: now.Add(time.Hour * 24 * 4)},
		{GUID: "item122", FeedId: feed12.Id, Title: "title122", Date: now.Add(time.Hour * 24 * 5)},
		// feed21
		{GUID: "item211", FeedId: feed21.Id, Title: "title211", Date: now.Add(time.Hour * 24 * 6)},
		{GUID: "item212", FeedId: feed21.Id, Title: "title212", Date: now.Add(time.Hour * 24 * 7)},
		// feed01
		{GUID: "item011", FeedId: feed01.Id, Title: "title011", Date: now.Add(time.Hour * 24 * 8)},
		{GUID: "item012", FeedId: feed01.Id, Title: "title012", Date: now.Add(time.Hour * 24 * 9)},
		{GUID: "item013", FeedId: feed01.Id, Title: "title013", Date: now.Add(time.Hour * 24 * 10)},
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

func testItemGuids(items []Item) []string {
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

	have := testItemGuids(db.ListItems(ItemFilter{FolderID: &scope.folder1.Id}, 0, 10, false))
	want := []string{"item111", "item112", "item113", "item121", "item122"}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}

	have = testItemGuids(db.ListItems(ItemFilter{FolderID: &scope.folder2.Id}, 0, 10, false))
	want = []string{"item211", "item212"}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}

	// filter by feed_id

	have = testItemGuids(db.ListItems(ItemFilter{FeedID: &scope.feed11.Id}, 0, 10, false))
	want = []string{"item111", "item112", "item113"}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}

	have = testItemGuids(db.ListItems(ItemFilter{FeedID: &scope.feed01.Id}, 0, 10, false))
	want = []string{"item011", "item012", "item013"}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}

	// filter by status

	var starred ItemStatus = STARRED
	have = testItemGuids(db.ListItems(ItemFilter{Status: &starred}, 0, 10, false))
	want = []string{"item113", "item212", "item013"}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}

	var unread ItemStatus = UNREAD
	have = testItemGuids(db.ListItems(ItemFilter{Status: &unread}, 0, 10, false))
	want = []string{"item111", "item121", "item011"}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}

	// filter by offset,limit

	have = testItemGuids(db.ListItems(ItemFilter{}, 0, 2, false))
	want = []string{"item111", "item112"}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}

	have = testItemGuids(db.ListItems(ItemFilter{}, 2, 3, false))
	want = []string{"item113", "item121", "item122"}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}

	// filter by search
	db.SyncSearch()
	search1 := "title111"
	have = testItemGuids(db.ListItems(ItemFilter{Search: &search1}, 0, 4, true))
	want = []string{"item111"}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}

	// sort by date
	have = testItemGuids(db.ListItems(ItemFilter{}, 0, 4, true))
	want = []string{"item013", "item012", "item011", "item212"}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}
}

func TestCountItems(t *testing.T) {
	db := testDB()
	scope := testItemsSetup(db)

	have := db.CountItems(ItemFilter{})
	want := int64(10)
	if have != want {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}
	
	// folders

	have = db.CountItems(ItemFilter{FolderID: &scope.folder1.Id})
	want = int64(5)
	if have != want {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}

	have = db.CountItems(ItemFilter{FolderID: &scope.folder2.Id})
	want = int64(2)
	if have != want {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}

	// feeds

	have = db.CountItems(ItemFilter{FeedID: &scope.feed21.Id})
	want = int64(2)
	if have != want {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}

	have = db.CountItems(ItemFilter{FeedID: &scope.feed01.Id})
	want = int64(3)
	if have != want {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}

	// statuses

	var unread ItemStatus = UNREAD
	have = db.CountItems(ItemFilter{Status: &unread})
	want = int64(3)
	if have != want {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fail()
	}

	// search

	db.SyncSearch()
	search := "title0"
	have = db.CountItems(ItemFilter{Search: &search})
	want = int64(3)
	if have != want {
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
	have := testItemGuids(db1.ListItems(ItemFilter{Status: &read}, 0, 10, false))
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
	have = testItemGuids(db2.ListItems(ItemFilter{Status: &read}, 0, 10, false))
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
	have = testItemGuids(db3.ListItems(ItemFilter{Status: &read}, 0, 10, false))
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
