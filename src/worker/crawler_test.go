package worker

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestDiscoverFeedLinkCycle(t *testing.T) {
	var requests atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Stop the infinite recursion so the test itself terminates.
		if requests.Add(1) > 5 {
			http.Error(w, "request limit reached", http.StatusInternalServerError)
			return
		}

		target := "/a"
		if r.URL.Path == "/a" {
			target = "/b"
		}
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<html><head><link rel="alternate" type="application/rss+xml" href="%s"></head></html>`, target)
	}))
	defer srv.Close()

	_, err := DiscoverFeed(srv.Client(), srv.URL+"/a")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if got := requests.Load(); got > 2 {
		t.Fatalf("DiscoverFeed made %d requests; want cycle detected after 2", got)
	}
}
