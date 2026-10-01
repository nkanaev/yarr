package server

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/nkanaev/yarr/src/assets"
	"github.com/nkanaev/yarr/src/storage"
	"github.com/nkanaev/yarr/src/storage/model"
)

func testServer() *Server {
	db, _ := storage.New(":memory:")

	server := NewServer("127.0.0.1:8000")
	server.StaticFS = assets.StaticFS()
	server.Template = assets.Templates()
	server.Storage = NewLocalStorage(db)

	return server
}

func TestStatic(t *testing.T) {
	handler := testServer().Handler()
	url := "/static/bundle.js"

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", url, nil)
	handler.ServeHTTP(recorder, request)
	if recorder.Result().StatusCode != 200 {
		t.FailNow()
	}
}

func TestStaticWithBase(t *testing.T) {
	server := testServer()
	server.BasePath = "/sub"

	handler := server.Handler()
	url := "/sub/static/bundle.js"

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", url, nil)
	handler.ServeHTTP(recorder, request)
	if recorder.Result().StatusCode != 200 {
		t.FailNow()
	}
}

func TestStaticBanTemplates(t *testing.T) {
	server := testServer()
	handler := server.Handler()
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
	log.SetOutput(os.Stderr)
	handler := testServer().Handler()
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

type fakeIngestor struct {
	result AddFeedResult
	err    error
}

func (f fakeIngestor) AddFeed(ctx context.Context, store storage.Storage, feed FeedLink) (AddFeedResult, error) {
	return f.result, f.err
}

func (f fakeIngestor) AddOPML(ctx context.Context, store storage.Storage, file io.Reader) error {
	return f.err
}

func newFeedCreateServer(t *testing.T, ingestor FeedIngestor) http.Handler {
	t.Helper()
	db, err := storage.New(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	server := NewServer("127.0.0.1:8000")
	server.Storage = NewLocalStorage(db)
	server.Ingestor = ingestor
	return server.Handler()
}

func feedCreateRequest(t *testing.T, handler http.Handler, body string) *httptest.ResponseRecorder {
	t.Helper()
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/api/feeds", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	handler.ServeHTTP(recorder, request)
	return recorder
}

func TestFeedCreateSuccess(t *testing.T) {
	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stderr)

	feed := &model.Feed{Id: 1, Title: "RSS Title", FeedLink: "http://example.com/feed"}

	handler := newFeedCreateServer(t, fakeIngestor{result: AddFeedResult{Feed: feed}})
	recorder := feedCreateRequest(t, handler, `{"url":"http://example.com"}`)

	if recorder.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Result().StatusCode)
	}
	var resp map[string]any
	json.NewDecoder(recorder.Result().Body).Decode(&resp)
	if resp["status"] != "success" {
		t.Fatalf("expected success, got %v", resp["status"])
	}
}

func TestFeedCreateMultiple(t *testing.T) {
	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stderr)

	choices := []FeedLink{{URL: "http://example.com/rss", Title: "RSS"}}
	handler := newFeedCreateServer(t, fakeIngestor{result: AddFeedResult{Choices: choices}})
	recorder := feedCreateRequest(t, handler, `{"url":"http://example.com"}`)

	if recorder.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Result().StatusCode)
	}
	var resp map[string]any
	json.NewDecoder(recorder.Result().Body).Decode(&resp)
	if resp["status"] != "multiple" {
		t.Fatalf("expected multiple, got %v", resp["status"])
	}
	if len(resp["choice"].([]any)) != 1 {
		t.Fatalf("expected 1 choice, got %v", resp["choice"])
	}
}

func TestFeedCreateError(t *testing.T) {
	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stderr)

	handler := newFeedCreateServer(t, fakeIngestor{err: errors.New("discovery failed")})
	recorder := feedCreateRequest(t, handler, `{"url":"http://example.com"}`)

	if recorder.Result().StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", recorder.Result().StatusCode)
	}
}
