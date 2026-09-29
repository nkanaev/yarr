package worker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nkanaev/yarr/src/client"
	"github.com/nkanaev/yarr/src/server"
)

const feedFixture = `<?xml version="1.0"?>
	<rss version="2.0">
		<channel>
			<title>RSS Title</title>
			<link>http://example.com</link>
			<item>
				<title>Item 1</title>
				<link>http://example.com/1</link>
			</item>
		</channel>
	</rss>`

func TestAddFeedTitle(t *testing.T) {
	for _, tc := range []struct {
		name     string
		override string
		want     string
	}{
		{name: "rss title", want: "RSS Title"},
		{name: "override", override: "Override Title", want: "Override Title"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			db := newTestDB(t)
			feedSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/rss+xml")
				w.Write([]byte(feedFixture))
			}))
			defer feedSrv.Close()

			w := NewWorker(db, client.NewBuilder().Build())
			result, err := w.AddFeed(context.Background(), db, server.FeedLink{
				URL:   feedSrv.URL,
				Title: tc.override,
			})
			if err != nil {
				t.Fatal(err)
			}
			if result.Feed == nil {
				t.Fatal("expected a feed, got nil")
			}
			if result.Feed.Title != tc.want {
				t.Errorf("title = %q; want %q", result.Feed.Title, tc.want)
			}
			if items := db.CountItems(); items != 1 {
				t.Errorf("stored items = %d; want 1", items)
			}
		})
	}
}

func TestAddFeedChoices(t *testing.T) {
	db := newTestDB(t)
	pageSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<html><head>
			<link rel="alternate" type="application/rss+xml" title="First" href="http://example.com/one.xml">
			<link rel="alternate" type="application/atom+xml" title="Second" href="http://example.com/two.xml">
		</head></html>`))
	}))
	defer pageSrv.Close()

	w := NewWorker(db, client.NewBuilder().Build())
	result, err := w.AddFeed(context.Background(), db, server.FeedLink{URL: pageSrv.URL})
	if err != nil {
		t.Fatal(err)
	}
	if result.Feed != nil {
		t.Errorf("expected no feed, got %+v", result.Feed)
	}
	if len(result.Choices) != 2 {
		t.Fatalf("choices = %d; want 2", len(result.Choices))
	}
}
