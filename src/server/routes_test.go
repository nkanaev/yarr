package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/nkanaev/yarr/src/storage"
	"github.com/nkanaev/yarr/src/storage/model"
)

func TestStatic(t *testing.T) {
	handler := NewServer(nil, "127.0.0.1:8000").handler()
	url := "/static/bundle.js"

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", url, nil)
	handler.ServeHTTP(recorder, request)
	if recorder.Result().StatusCode != 200 {
		t.FailNow()
	}
}

func TestStaticWithBase(t *testing.T) {
	server := NewServer(nil, "127.0.0.1:8000")
	server.BasePath = "/sub"

	handler := server.handler()
	url := "/sub/static/bundle.js"

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", url, nil)
	handler.ServeHTTP(recorder, request)
	if recorder.Result().StatusCode != 200 {
		t.FailNow()
	}
}

func TestStaticBanTemplates(t *testing.T) {
	handler := NewServer(nil, "127.0.0.1:8000").handler()
	url := "/static/login.html"

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", url, nil)
	handler.ServeHTTP(recorder, request)
	if recorder.Result().StatusCode != 404 {
		t.FailNow()
	}
}

func TestIndexGzipped(t *testing.T) {
	log.SetOutput(io.Discard)
	db, _ := storage.New(":memory:")
	log.SetOutput(os.Stderr)
	handler := NewServer(db, "127.0.0.1:8000").handler()
	url := "/"

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", url, nil)
	request.Header.Set("accept-encoding", "gzip")
	handler.ServeHTTP(recorder, request)
	response := recorder.Result()
	if response.StatusCode != 200 {
		t.FailNow()
	}
	if response.Header.Get("content-encoding") != "gzip" {
		t.Errorf("invalid content-encoding header: %#v", response.Header.Get("content-encoding"))
	}
	if response.Header.Get("content-type") != "text/html" {
		t.Errorf("invalid content-type header: %#v", response.Header.Get("content-type"))
	}
}

func TestFeedCreateWithTitleOverride(t *testing.T) {
	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stderr)

	feedSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml")
		w.Write([]byte(`<?xml version="1.0"?>
			<rss version="2.0">
				<channel>
					<title>RSS Title</title>
					<link>http://example.com</link>
					<item>
						<title>Item 1</title>
						<link>http://example.com/1</link>
					</item>
				</channel>
			</rss>
		`))
	}))
	defer feedSrv.Close()

	db, err := storage.New(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	server := NewServer(db, "127.0.0.1:8000")
	handler := server.handler()

	t.Run("override title", func(t *testing.T) {
		body := fmt.Sprintf(`{"url":%q,"title_override":"Override Title"}`, feedSrv.URL)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/api/feeds", strings.NewReader(body))
		request.Header.Set("Content-Type", "application/json")
		handler.ServeHTTP(recorder, request)

		if recorder.Result().StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", recorder.Result().StatusCode)
		}

		var resp map[string]any
		json.NewDecoder(recorder.Result().Body).Decode(&resp)
		if resp["status"] != "success" {
			t.Fatalf("expected success, got %v", resp["status"])
		}
		feed := resp["feed"].(map[string]any)
		if feed["title"] != "Override Title" {
			t.Fatalf("expected 'Override Title', got %v", feed["title"])
		}
	})

	t.Run("no override uses rss title", func(t *testing.T) {
		body := fmt.Sprintf(`{"url":%q}`, feedSrv.URL)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/api/feeds", strings.NewReader(body))
		request.Header.Set("Content-Type", "application/json")
		handler.ServeHTTP(recorder, request)

		if recorder.Result().StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", recorder.Result().StatusCode)
		}

		var resp map[string]any
		json.NewDecoder(recorder.Result().Body).Decode(&resp)
		if resp["status"] != "success" {
			t.Fatalf("expected success, got %v", resp["status"])
		}
		feed := resp["feed"].(map[string]any)
		if feed["title"] != "RSS Title" {
			t.Fatalf("expected 'RSS Title', got %v", feed["title"])
		}
	})
}

func TestFeedIcons(t *testing.T) {
	log.SetOutput(io.Discard)
	db, _ := storage.New(":memory:")
	icon := []byte("test")
	feed := db.CreateFeed(model.CreateFeedParams{})
	db.UpdateFeed(feed.Id, model.UpdateFeedParams{Icon: model.SetNullable(&icon)})
	log.SetOutput(os.Stderr)

	recorder := httptest.NewRecorder()
	url := fmt.Sprintf("/api/feeds/%d/icon", feed.Id)
	request := httptest.NewRequest("GET", url, nil)

	handler := NewServer(db, "127.0.0.1:8000").handler()
	handler.ServeHTTP(recorder, request)
	response := recorder.Result()

	if response.StatusCode != http.StatusOK {
		t.Fatal()
	}
	body, _ := io.ReadAll(response.Body)
	if !reflect.DeepEqual(body, icon) {
		t.Fatal()
	}
	if response.Header.Get("Etag") == "" {
		t.Fatal()
	}

	recorder2 := httptest.NewRecorder()
	request2 := httptest.NewRequest("GET", url, nil)
	request2.Header.Set("If-None-Match", response.Header.Get("Etag"))
	handler.ServeHTTP(recorder2, request2)
	response2 := recorder2.Result()

	if response2.StatusCode != http.StatusNotModified {
		t.Fatal("got", response2.StatusCode)
	}
}
