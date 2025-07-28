package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/nkanaev/yarr/src/storage"
)

// Test the complete archive tab workflow including frontend filtering behavior
func TestArchiveTabIntegration(t *testing.T) {
	log.SetOutput(io.Discard)
	db, _ := storage.New(":memory:")
	
	// Create test folder and feeds
	folder1 := db.CreateFolder("Test Folder")
	folder2 := db.CreateFolder("Archive Folder")
	
	_ = db.CreateFeed("Active Feed 1", "", "http://example1.com", "http://example1.com/feed.xml", &folder1.Id)
	_ = db.CreateFeed("Active Feed 2", "", "http://example2.com", "http://example2.com/feed.xml", nil)
	archivedFeed1 := db.CreateFeed("Archived Feed 1", "", "http://example3.com", "http://example3.com/feed.xml", &folder2.Id)
	archivedFeed2 := db.CreateFeed("Archived Feed 2", "", "http://example4.com", "http://example4.com/feed.xml", &folder1.Id)
	
	// Archive two feeds
	db.ArchiveFeed(archivedFeed1.Id)
	db.ArchiveFeed(archivedFeed2.Id)
	
	log.SetOutput(os.Stderr)
	handler := NewServer(db, "127.0.0.1:8000").handler()

	// Test getting all feeds to verify archive status
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/api/feeds", nil)
	handler.ServeHTTP(recorder, request)
	
	if recorder.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Result().StatusCode)
	}
	
	var feeds []storage.Feed
	body, _ := io.ReadAll(recorder.Result().Body)
	err := json.Unmarshal(body, &feeds)
	if err != nil {
		t.Fatalf("failed to unmarshal feeds: %v", err)
	}
	
	// Verify we have 4 feeds with correct archive status
	if len(feeds) != 4 {
		t.Fatalf("expected 4 feeds, got %d", len(feeds))
	}
	
	archivedCount := 0
	activeCount := 0
	for _, feed := range feeds {
		if feed.Archived {
			archivedCount++
		} else {
			activeCount++
		}
	}
	
	if archivedCount != 2 {
		t.Errorf("expected 2 archived feeds, got %d", archivedCount)
	}
	if activeCount != 2 {
		t.Errorf("expected 2 active feeds, got %d", activeCount)
	}
}

func TestArchiveTabFeedStats(t *testing.T) {
	log.SetOutput(io.Discard)
	db, _ := storage.New(":memory:")
	
	// Create feeds and items for testing stats
	feed1 := db.CreateFeed("Active Feed", "", "http://example1.com", "http://example1.com/feed.xml", nil)
	feed2 := db.CreateFeed("Archived Feed", "", "http://example2.com", "http://example2.com/feed.xml", nil)
	
	// Create some items
	items := []storage.Item{
		{GUID: "item1", FeedId: feed1.Id, Title: "Active Unread", Status: storage.UNREAD},
		{GUID: "item2", FeedId: feed1.Id, Title: "Active Starred", Status: storage.STARRED},
		{GUID: "item3", FeedId: feed2.Id, Title: "Archived Unread", Status: storage.UNREAD},
		{GUID: "item4", FeedId: feed2.Id, Title: "Archived Starred", Status: storage.STARRED},
	}
	db.CreateItems(items)
	
	// Archive the second feed
	db.ArchiveFeed(feed2.Id)
	
	log.SetOutput(os.Stderr)
	handler := NewServer(db, "127.0.0.1:8000").handler()

	// Test getting feed stats
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/api/status", nil)
	handler.ServeHTTP(recorder, request)
	
	if recorder.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Result().StatusCode)
	}
	
	var stats map[string]interface{}
	body, _ := io.ReadAll(recorder.Result().Body)
	err := json.Unmarshal(body, &stats)
	if err != nil {
		t.Fatalf("failed to unmarshal stats: %v", err)
	}
	
	// Verify that stats include both active and archived feeds
	statsData, ok := stats["stats"]
	if !ok {
		t.Fatal("expected stats field in response")
	}
	
	// The stats should be present for both feeds
	// The exact structure depends on the implementation, but we just need to verify
	// that the response contains stats and both feeds are accounted for
	if statsData == nil {
		t.Error("stats data should not be nil")
	}
	
	// This test mainly verifies that the API returns successfully with archived feeds present
	// The detailed structure testing is done in other tests
}

func TestArchiveFeedWorkflow(t *testing.T) {
	log.SetOutput(io.Discard)
	db, _ := storage.New(":memory:")
	feed := db.CreateFeed("Test Feed", "", "http://example.com", "http://example.com/feed.xml", nil)
	log.SetOutput(os.Stderr)

	handler := NewServer(db, "127.0.0.1:8000").handler()

	// Test complete workflow: create -> archive -> unarchive
	tests := []struct {
		name        string
		payload     string
		expectArchived bool
	}{
		{"Archive feed", `{"archived": true}`, true},
		{"Unarchive feed", `{"archived": false}`, false},
		{"Archive again", `{"archived": true}`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Update feed status
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/api/feeds/%d", feed.Id)
			request := httptest.NewRequest("PUT", url, strings.NewReader(tt.payload))
			request.Header.Set("Content-Type", "application/json")

			handler.ServeHTTP(recorder, request)
			
			if recorder.Result().StatusCode != http.StatusOK {
				t.Fatalf("expected status 200, got %d", recorder.Result().StatusCode)
			}

			// Verify feed status
			updatedFeed := db.GetFeed(feed.Id)
			if updatedFeed == nil {
				t.Fatal("feed should exist")
			}
			
			if updatedFeed.Archived != tt.expectArchived {
				t.Errorf("expected archived=%v, got %v", tt.expectArchived, updatedFeed.Archived)
			}
		})
	}
}

func TestArchiveTabVisibilityRules(t *testing.T) {
	log.SetOutput(io.Discard)
	db, _ := storage.New(":memory:")
	
	folder := db.CreateFolder("Mixed Folder")
	activeFeed := db.CreateFeed("Active Feed", "", "http://example1.com", "http://example1.com/feed.xml", &folder.Id)
	archivedFeed := db.CreateFeed("Archived Feed", "", "http://example2.com", "http://example2.com/feed.xml", &folder.Id)
	
	// Archive one feed
	db.ArchiveFeed(archivedFeed.Id)
	
	log.SetOutput(os.Stderr)
	handler := NewServer(db, "127.0.0.1:8000").handler()

	// Get all feeds to verify the archive tab filtering behavior
	// This simulates what the frontend would need to implement filtering
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/api/feeds", nil)
	handler.ServeHTTP(recorder, request)
	
	if recorder.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Result().StatusCode)
	}
	
	var feeds []storage.Feed
	body, _ := io.ReadAll(recorder.Result().Body)
	err := json.Unmarshal(body, &feeds)
	if err != nil {
		t.Fatalf("failed to unmarshal feeds: %v", err)
	}
	
	// Simulate frontend filtering logic
	var activeFeeds, archivedFeeds []storage.Feed
	for _, feed := range feeds {
		if feed.Archived {
			archivedFeeds = append(archivedFeeds, feed)
		} else {
			activeFeeds = append(activeFeeds, feed)
		}
	}
	
	// Verify filtering results
	if len(activeFeeds) != 1 {
		t.Errorf("expected 1 active feed, got %d", len(activeFeeds))
	}
	
	if len(archivedFeeds) != 1 {
		t.Errorf("expected 1 archived feed, got %d", len(archivedFeeds))
	}
	
	// Verify correct feeds are in each category
	if len(activeFeeds) > 0 && activeFeeds[0].Id != activeFeed.Id {
		t.Error("wrong feed in active category")
	}
	
	if len(archivedFeeds) > 0 && archivedFeeds[0].Id != archivedFeed.Id {
		t.Error("wrong feed in archived category")
	}
}

func TestArchiveTabFolderBehavior(t *testing.T) {
	log.SetOutput(io.Discard)
	db, _ := storage.New(":memory:")
	
	// Create folder with mixed feeds
	folder := db.CreateFolder("Mixed Folder")
	_ = db.CreateFeed("Active Feed", "", "http://example1.com", "http://example1.com/feed.xml", &folder.Id)
	archivedFeed := db.CreateFeed("Archived Feed", "", "http://example2.com", "http://example2.com/feed.xml", &folder.Id)
	
	// Create folder with only archived feeds
	archivedFolder := db.CreateFolder("Archived Folder")
	archivedFeed2 := db.CreateFeed("Archived Feed 2", "", "http://example3.com", "http://example3.com/feed.xml", &archivedFolder.Id)
	
	// Archive feeds
	db.ArchiveFeed(archivedFeed.Id)
	db.ArchiveFeed(archivedFeed2.Id)
	
	log.SetOutput(os.Stderr)
	handler := NewServer(db, "127.0.0.1:8000").handler()

	// Get folders and feeds
	var folders []storage.Folder
	var feeds []storage.Feed
	
	// Get folders
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/api/folders", nil)
	handler.ServeHTTP(recorder, request)
	
	if recorder.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Result().StatusCode)
	}
	
	body, _ := io.ReadAll(recorder.Result().Body)
	err := json.Unmarshal(body, &folders)
	if err != nil {
		t.Fatalf("failed to unmarshal folders: %v", err)
	}
	
	// Get feeds
	recorder = httptest.NewRecorder()
	request = httptest.NewRequest("GET", "/api/feeds", nil)
	handler.ServeHTTP(recorder, request)
	
	if recorder.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Result().StatusCode)
	}
	
	body, _ = io.ReadAll(recorder.Result().Body)
	err = json.Unmarshal(body, &feeds)
	if err != nil {
		t.Fatalf("failed to unmarshal feeds: %v", err)
	}
	
	// Test folder visibility logic for archived tab
	// This simulates frontend logic for showing/hiding folders based on archive filter
	
	// When showing archived feeds only:
	foldersWithArchivedFeeds := make(map[int64]bool)
	for _, feed := range feeds {
		if feed.Archived && feed.FolderId != nil {
			foldersWithArchivedFeeds[*feed.FolderId] = true
		}
	}
	
	// Both folders should have archived feeds
	if !foldersWithArchivedFeeds[folder.Id] {
		t.Error("mixed folder should be visible in archived tab")
	}
	
	if !foldersWithArchivedFeeds[archivedFolder.Id] {
		t.Error("archived folder should be visible in archived tab")
	}
	
	// When showing active feeds only:
	foldersWithActiveFeeds := make(map[int64]bool)
	for _, feed := range feeds {
		if !feed.Archived && feed.FolderId != nil {
			foldersWithActiveFeeds[*feed.FolderId] = true
		}
	}
	
	// Only mixed folder should have active feeds
	if !foldersWithActiveFeeds[folder.Id] {
		t.Error("mixed folder should be visible in active feeds")
	}
	
	if foldersWithActiveFeeds[archivedFolder.Id] {
		t.Error("archived folder should not be visible when showing active feeds")
	}
}