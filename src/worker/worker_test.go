package worker

import (
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/nkanaev/yarr/src/parser"
	"github.com/nkanaev/yarr/src/storage"
)

func TestConvertItems_Empty(t *testing.T) {
	feed := storage.Feed{Id: 1}
	result := ConvertItems(nil, feed)
	if len(result) != 0 {
		t.Fatalf("expected 0 items, got %d", len(result))
	}
}

func TestConvertItems_BasicFields(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	feed := storage.Feed{Id: 42}

	items := []parser.Item{
		{
			GUID:    "guid-1",
			Title:   "Test Article",
			URL:     "https://example.com/article",
			Content: "<p>Hello</p>",
			Date:    now,
		},
	}

	result := ConvertItems(items, feed)
	if len(result) != 1 {
		t.Fatalf("expected 1 item, got %d", len(result))
	}

	got := result[0]
	if got.GUID != "guid-1" {
		t.Errorf("GUID: got %q, want %q", got.GUID, "guid-1")
	}
	if got.FeedId != 42 {
		t.Errorf("FeedId: got %d, want 42", got.FeedId)
	}
	if got.Title != "Test Article" {
		t.Errorf("Title: got %q", got.Title)
	}
	if got.Link != "https://example.com/article" {
		t.Errorf("Link: got %q", got.Link)
	}
	if got.Content != "<p>Hello</p>" {
		t.Errorf("Content: got %q", got.Content)
	}
	if !got.Date.Equal(now) {
		t.Errorf("Date: got %v, want %v", got.Date, now)
	}
	if got.Status != storage.UNREAD {
		t.Errorf("Status: got %d, want UNREAD (%d)", got.Status, storage.UNREAD)
	}
}

func TestConvertItems_MediaLinks(t *testing.T) {
	feed := storage.Feed{Id: 1}
	items := []parser.Item{
		{
			GUID: "guid-media",
			MediaLinks: []parser.MediaLink{
				{URL: "https://example.com/audio.mp3", Type: "audio/mpeg", Description: "Episode 1"},
				{URL: "https://example.com/video.mp4", Type: "video/mp4", Description: ""},
			},
		},
	}

	result := ConvertItems(items, feed)
	if len(result[0].MediaLinks) != 2 {
		t.Fatalf("expected 2 media links, got %d", len(result[0].MediaLinks))
	}

	want := storage.MediaLinks{
		{URL: "https://example.com/audio.mp3", Type: "audio/mpeg", Description: "Episode 1"},
		{URL: "https://example.com/video.mp4", Type: "video/mp4", Description: ""},
	}
	if !reflect.DeepEqual(result[0].MediaLinks, want) {
		t.Errorf("MediaLinks mismatch:\ngot:  %+v\nwant: %+v", result[0].MediaLinks, want)
	}
}

func TestConvertItems_MultipleItems(t *testing.T) {
	feed := storage.Feed{Id: 5}
	items := []parser.Item{
		{GUID: "a", Title: "First"},
		{GUID: "b", Title: "Second"},
		{GUID: "c", Title: "Third"},
	}

	result := ConvertItems(items, feed)
	if len(result) != 3 {
		t.Fatalf("expected 3 items, got %d", len(result))
	}
	for i, item := range result {
		if item.FeedId != 5 {
			t.Errorf("item[%d].FeedId: got %d, want 5", i, item.FeedId)
		}
		if item.Status != storage.UNREAD {
			t.Errorf("item[%d].Status: got %d, want UNREAD", i, item.Status)
		}
	}
}

func TestConvertItems_NoMediaLinks(t *testing.T) {
	feed := storage.Feed{Id: 1}
	items := []parser.Item{
		{GUID: "no-media"},
	}

	result := ConvertItems(items, feed)
	if result[0].MediaLinks == nil {
		t.Fatal("MediaLinks should be empty slice, not nil")
	}
	if len(result[0].MediaLinks) != 0 {
		t.Fatalf("expected 0 media links, got %d", len(result[0].MediaLinks))
	}
}

func TestGetCharset(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		want        string
	}{
		{"utf-8", "text/xml; charset=utf-8", "utf-8"},
		{"iso-8859-1", "text/html; charset=iso-8859-1", "iso-8859-1"},
		{"no charset", "text/html", ""},
		{"empty", "", ""},
		{"invalid charset", "text/html; charset=bogus-encoding-xyz", ""},
		{"windows-1252", "text/xml; charset=windows-1252", "windows-1252"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := &http.Response{
				Header: http.Header{"Content-Type": []string{tt.contentType}},
			}
			got := getCharset(res)
			if got != tt.want {
				t.Errorf("getCharset(%q) = %q, want %q", tt.contentType, got, tt.want)
			}
		})
	}
}

func TestSetVersion(t *testing.T) {
	SetVersion("2.6")
	if client.userAgent != "Yarr/2.6" {
		t.Errorf("expected user agent 'Yarr/2.6', got %q", client.userAgent)
	}
	// Restore default
	SetVersion("1.0")
}
