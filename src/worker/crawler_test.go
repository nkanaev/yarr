package worker

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/nkanaev/yarr/src/storage/model"
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

func TestListItemsRejectsOversizedBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml")
		fmt.Fprint(w, `<?xml version="1.0"?><rss version="2.0"><channel>`)

		chunk := bytes.Repeat([]byte(" "), 32*1024)
		for {
			if _, err := w.Write(chunk); err != nil {
				return
			}
		}
	}))
	defer srv.Close()

	errCh := make(chan error, 1)
	go func() {
		_, err := listItems(srv.Client(), model.Feed{FeedLink: srv.URL})
		errCh <- err
	}()

	select {
	case err := <-errCh:
		var maxErr *http.MaxBytesError
		if !errors.As(err, &maxErr) {
			t.Fatalf("listItems error = %v; want *http.MaxBytesError", err)
		}
	case <-time.After(5 * time.Second):
		srv.CloseClientConnections()
		t.Fatal("listItems did not stop reading oversized response")
	}
}
