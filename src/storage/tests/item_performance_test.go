package tests

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/nkanaev/yarr/src/storage"
	"github.com/nkanaev/yarr/src/storage/model"
)

var loremIpsum = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. " +
	"Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. " +
	"Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris " +
	"nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in " +
	"reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla " +
	"pariatur. Excepteur sint occaecat cupidatat non proident, sunt in " +
	"culpa qui officia deserunt mollit anim id est laborum. " +
	"Lorem ipsum dolor sit amet, consectetur adipiscing elit. " +
	"Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. " +
	"Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris " +
	"nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in " +
	"reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla " +
	"pariatur. Excepteur sint occaecat cupidatat non proident, sunt in " +
	"culpa qui officia deserunt mollit anim id est laborum."

func perfDB(t testing.TB) storage.Storage {
	t.Helper()
	db, err := storage.New(filepath.Join(t.TempDir(), "perf.db"))
	if err != nil {
		t.Fatalf("failed to create perf db: %v", err)
	}
	return db
}

func createBenchFeed(db storage.Storage, tag string, n int) *model.Feed {
	feed := db.CreateFeed(model.CreateFeedParams{FeedLink: fmt.Sprintf("http://%s.xml", tag)})
	now := time.Now()
	items := make([]model.Item, n)
	for i := range items {
		items[i] = model.Item{
			GUID:   fmt.Sprintf("b-%s-%d", tag, i),
			FeedId: feed.Id,
			Title:  "t",
			Date:   now.Add(time.Duration(i) * time.Second),
		}
	}
	db.CreateItems(items)
	return feed
}

func BenchmarkCreateItems_Batch10(b *testing.B) {
	db := perfDB(b)
	feed := db.CreateFeed(model.CreateFeedParams{FeedLink: "http://b10.xml"})

	for b.Loop() {
		items := make([]model.Item, 10)
		for j := range items {
			items[j] = model.Item{
				GUID:    fmt.Sprintf("b10-%d", j),
				FeedId:  feed.Id,
				Title:   "t",
				Content: loremIpsum,
				Date:    time.Now(),
			}
		}
		db.CreateItems(items)
	}
}

func BenchmarkCreateItems_Batch100(b *testing.B) {
	db := perfDB(b)
	feed := db.CreateFeed(model.CreateFeedParams{FeedLink: "http://b100.xml"})

	for b.Loop() {
		items := make([]model.Item, 100)
		for j := range items {
			items[j] = model.Item{
				GUID:    fmt.Sprintf("b100-%d", j),
				FeedId:  feed.Id,
				Title:   "t",
				Content: loremIpsum,
				Date:    time.Now(),
			}
		}
		db.CreateItems(items)
	}
}

func BenchmarkCreateItems_Batch1000(b *testing.B) {
	db := perfDB(b)
	feed := db.CreateFeed(model.CreateFeedParams{FeedLink: "http://b1000.xml"})

	for b.Loop() {
		items := make([]model.Item, 1000)
		for j := range items {
			items[j] = model.Item{
				GUID:    fmt.Sprintf("b1000-%d", j),
				FeedId:  feed.Id,
				Title:   "t",
				Content: loremIpsum,
				Date:    time.Now(),
			}
		}
		db.CreateItems(items)
	}
}

func BenchmarkCreateItems_Upsert(b *testing.B) {
	db := perfDB(b)
	feed := db.CreateFeed(model.CreateFeedParams{FeedLink: "http://upsert.xml"})

	items := make([]model.Item, 100)
	for j := range items {
		items[j] = model.Item{
			GUID:   fmt.Sprintf("u-%d", j),
			FeedId: feed.Id,
			Title:  "t",
			Date:   time.Now(),
		}
	}
	db.CreateItems(items)

	for b.Loop() {
		db.CreateItems(items)
	}
}

func BenchmarkFeedStats_Empty(b *testing.B) {
	db := perfDB(b)

	for b.Loop() {
		db.FeedStats()
	}
}

func BenchmarkFeedStats_MultipleFeeds(b *testing.B) {
	db := perfDB(b)

	numFeeds := 1000
	var feeds []*model.Feed
	for k := 0; k < numFeeds; k++ {
		feeds = append(feeds, db.CreateFeed(model.CreateFeedParams{
			FeedLink: fmt.Sprintf("http://f%d.xml", k),
		}))
	}

	now := time.Now()
	var all []model.Item
	for i := 0; i < 100_000; i++ {
		all = append(all, model.Item{
			GUID:   fmt.Sprintf("i-%d", i),
			FeedId: feeds[i%numFeeds].Id,
			Title:  "t",
			Date:   now.Add(time.Duration(i) * time.Second),
			Status: model.ItemStatus(i % 3),
		})
	}
	db.CreateItems(all)

	for b.Loop() {
		db.FeedStats()
	}
}

func BenchmarkListItems_All(b *testing.B) {
	db := perfDB(b)
	createBenchFeed(db, "all", 10000)

	for b.Loop() {
		db.ListItems(model.ItemFilter{}, 50, false, false)
	}
}

func BenchmarkListItems_ByFeed(b *testing.B) {
	db := perfDB(b)
	feed := createBenchFeed(db, "feed", 10000)

	for b.Loop() {
		db.ListItems(model.ItemFilter{FeedID: &feed.Id}, 50, false, false)
	}
}

func BenchmarkListItems_ByStatus(b *testing.B) {
	db := perfDB(b)
	createBenchFeed(db, "status", 10000)
	starred := model.STARRED

	for b.Loop() {
		db.ListItems(model.ItemFilter{Status: &starred}, 50, false, false)
	}
}

func BenchmarkListItems_Search(b *testing.B) {
	db := perfDB(b)
	feed := db.CreateFeed(model.CreateFeedParams{FeedLink: "http://search.xml"})
	now := time.Now()
	var all []model.Item
	for i := 0; i < 10000; i++ {
		all = append(all, model.Item{
			GUID:    fmt.Sprintf("s-%d", i),
			FeedId:  feed.Id,
			Title:   fmt.Sprintf("searchable title %d", i),
			Content: "common text for full-text search indexing",
			Date:    now.Add(time.Duration(i) * time.Second),
		})
	}
	db.CreateItems(all)
	query := "searchable"

	for b.Loop() {
		db.ListItems(model.ItemFilter{Search: &query}, 50, false, false)
	}
}

func BenchmarkListItems_Paginated(b *testing.B) {
	db := perfDB(b)
	feed := createBenchFeed(db, "page", 10000)

	all := db.ListItems(model.ItemFilter{FeedID: &feed.Id}, 10000, false, false)
	cursor := all[len(all)/2].Id

	for b.Loop() {
		db.ListItems(model.ItemFilter{After: &cursor}, 50, false, false)
	}
}

func BenchmarkMarkItemsRead_All(b *testing.B) {
	db := perfDB(b)
	createBenchFeed(db, "all", 1_000_000)

	for b.Loop() {
		db.MarkItemsRead(model.MarkFilter{})
	}
}

func BenchmarkMarkItemsRead_ByFeed(b *testing.B) {
	db := perfDB(b)
	createBenchFeed(db, "feed0", 500_000)
	feed := createBenchFeed(db, "feed1", 500_000)
	for b.Loop() {
		db.MarkItemsRead(model.MarkFilter{FeedID: &feed.Id})
	}
}

func BenchmarkMarkItemsRead_ByFolder(b *testing.B) {
	db := perfDB(b)

	folder := db.CreateFolder("perf")
	now := time.Now()
	var all []model.Item
	for k := 0; k < 5; k++ {
		feed := db.CreateFeed(model.CreateFeedParams{
			FeedLink: fmt.Sprintf("http://f%d.xml", k),
			FolderID: &folder.Id,
		})
		for i := 0; i < 1_000_000 / 5; i++ {
			all = append(all, model.Item{
				GUID:   fmt.Sprintf("f%d-i%d", k, i),
				FeedId: feed.Id,
				Title:  "t",
				Date:   now.Add(time.Duration(len(all)) * time.Second),
			})
		}
	}
	db.CreateItems(all)

	for b.Loop() {
		db.MarkItemsRead(model.MarkFilter{FolderID: &folder.Id})
	}
}
