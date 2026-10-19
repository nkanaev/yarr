package server

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
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

func TestIndexSettingsWhenUnauthenticated(t *testing.T) {
	server := testServer()
	db := server.Storage.(*localStorage).storage
	db.UpdateSettings(model.UpdateSettingsParams{
		Language:      new("fr"),
		ThemeName:     new("night"),
		RefreshRate:   new(int64(15)),
		FeedListWidth: new(400),
	})
	server.Auth = NewLocalAuthProvider("u", "p", "")

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	server.Handler().ServeHTTP(recorder, request)
	body := recorder.Body.String()

	if !strings.Contains(body, `"language":"fr"`) {
		t.Error("expected language to be preserved")
	}
	if !strings.Contains(body, `"theme_name":"night"`) {
		t.Error("expected theme_name to be preserved")
	}
	if strings.Contains(body, `"refresh_rate":15`) {
		t.Error("expected refresh_rate to be hidden")
	}
	if strings.Contains(body, `"feed_list_width":400`) {
		t.Error("expected feed_list_width to be hidden")
	}
}

func TestIndexSettingsWhenAuthenticated(t *testing.T) {
	server := testServer()
	db := server.Storage.(*localStorage).storage
	db.UpdateSettings(model.UpdateSettingsParams{
		Language:      new("fr"),
		ThemeName:     new("night"),
		RefreshRate:   new(int64(15)),
		FeedListWidth: new(400),
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	server.Handler().ServeHTTP(recorder, request)
	body := recorder.Body.String()

	if !strings.Contains(body, `"refresh_rate":15`) {
		t.Error("expected refresh_rate to be exposed")
	}
	if !strings.Contains(body, `"feed_list_width":400`) {
		t.Error("expected feed_list_width to be exposed")
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

type crawlTestTransport struct {
	target *url.URL
}

func (t *crawlTestTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	req.URL.Scheme = t.target.Scheme
	req.URL.Host = t.target.Host
	return http.DefaultTransport.RoundTrip(req)
}

func newCrawlTestServer(t *testing.T, handler http.HandlerFunc) *Server {
	t.Helper()

	target := httptest.NewServer(handler)
	t.Cleanup(target.Close)

	targetURL, err := url.Parse(target.URL)
	if err != nil {
		t.Fatal(err)
	}

	server := testServer()
	server.crawlClient = &http.Client{
		Transport: &crawlTestTransport{target: targetURL},
	}
	return server
}

func crawlPage(t *testing.T, server *Server, pageURL, userAgent string) *httptest.ResponseRecorder {
	t.Helper()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/page?url="+url.QueryEscape(pageURL), nil)
	if userAgent != "" {
		request.Header.Set("User-Agent", userAgent)
	}
	server.handlePageCrawl(recorder, request)
	return recorder
}

func TestPageCrawlSuccess(t *testing.T) {
	server := newCrawlTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("User-Agent"); got != "test-agent" {
			t.Errorf("expected user agent %q, got %q", "test-agent", got)
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		io.WriteString(w, `<!doctype html>
<html>
<head><title>Test article</title></head>
<body>
<article>
<h1>Test article</h1>
<p>Hello readability world. This is the first paragraph with enough text for extraction.</p>
<p>This is the second paragraph with more text to make the article look real.</p>
</article>
</body>
</html>`)
	})

	recorder := crawlPage(t, server, "http://example.com/article", "test-agent")
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}

	var response map[string]string
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(response["content"], "Hello readability world") {
		t.Fatalf("expected extracted content, got %q", response["content"])
	}
}

func TestPageCrawlEmptyBody(t *testing.T) {
	log.SetOutput(io.Discard)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })

	server := newCrawlTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	})

	recorder := crawlPage(t, server, "http://example.com/empty", "")
	if recorder.Code != http.StatusBadGateway {
		t.Fatalf("expected 502, got %d", recorder.Code)
	}
}

func TestPageCrawlUpstreamError(t *testing.T) {
	log.SetOutput(io.Discard)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })

	server := newCrawlTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	recorder := crawlPage(t, server, "http://example.com/missing", "")
	if recorder.Code != http.StatusBadGateway {
		t.Fatalf("expected 502, got %d", recorder.Code)
	}
}

func TestPageCrawlRejectsInternalRedirect(t *testing.T) {
	log.SetOutput(io.Discard)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })

	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://127.0.0.1/secret", http.StatusFound)
	}))
	t.Cleanup(target.Close)

	targetURL, err := url.Parse(target.URL)
	if err != nil {
		t.Fatal(err)
	}

	server := testServer()
	server.crawlClient.Transport = &crawlTestTransport{target: targetURL}

	recorder := crawlPage(t, server, "http://example.com/article", "")
	if recorder.Code != http.StatusBadGateway {
		t.Fatalf("expected 502, got %d", recorder.Code)
	}
}
