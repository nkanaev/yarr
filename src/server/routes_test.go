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
)

func TestStatic(t *testing.T) {
	handler := NewServer(nil, "127.0.0.1:8000").handler()
	url := "/static/javascripts/app.js"

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
	url := "/sub/static/javascripts/app.js"

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

func TestFeedIcons(t *testing.T) {
	log.SetOutput(io.Discard)
	db, _ := storage.New(":memory:")
	icon := []byte("test")
	feed := db.CreateFeed("", "", "", "", nil)
	db.UpdateFeedIcon(feed.Id, &icon)
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

func TestHealthEndpoint(t *testing.T) {
	log.SetOutput(io.Discard)
	db, _ := storage.New(":memory:")
	log.SetOutput(os.Stderr)

	handler := NewServer(db, "127.0.0.1:8000").handler()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/up", nil)
	handler.ServeHTTP(recorder, request)

	if recorder.Result().StatusCode != 200 {
		t.Fatalf("expected 200, got %d", recorder.Result().StatusCode)
	}
	body, _ := io.ReadAll(recorder.Result().Body)
	if string(body) != "OK" {
		t.Fatalf("expected body 'OK', got %q", string(body))
	}
}

func TestHealthEndpoint_NotAuthGated(t *testing.T) {
	log.SetOutput(io.Discard)
	db, _ := storage.New(":memory:")
	log.SetOutput(os.Stderr)

	srv := NewServer(db, "127.0.0.1:8000")
	srv.Username = "admin"
	srv.Password = "secret"
	srv.SecureCookie = true
	handler := srv.handler()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/up", nil)
	handler.ServeHTTP(recorder, request)

	if recorder.Result().StatusCode != 200 {
		t.Fatalf("health endpoint should not require auth, got %d", recorder.Result().StatusCode)
	}
}

func TestServiceWorker(t *testing.T) {
	handler := NewServer(nil, "127.0.0.1:8000").handler()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/sw.js", nil)
	handler.ServeHTTP(recorder, request)
	response := recorder.Result()

	if response.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", response.StatusCode)
	}
	if ct := response.Header.Get("Content-Type"); ct != "application/javascript" {
		t.Errorf("expected Content-Type application/javascript, got %q", ct)
	}
	if cc := response.Header.Get("Cache-Control"); cc != "no-cache" {
		t.Errorf("expected Cache-Control no-cache, got %q", cc)
	}
	if swa := response.Header.Get("Service-Worker-Allowed"); swa != "/" {
		t.Errorf("expected Service-Worker-Allowed /, got %q", swa)
	}
	body, _ := io.ReadAll(response.Body)
	if !strings.Contains(string(body), "CACHE_NAME") {
		t.Error("service worker body missing expected content")
	}
}

func TestServiceWorkerWithBase(t *testing.T) {
	srv := NewServer(nil, "127.0.0.1:8000")
	srv.BasePath = "/sub"
	handler := srv.handler()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/sub/sw.js", nil)
	handler.ServeHTTP(recorder, request)
	response := recorder.Result()

	if response.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", response.StatusCode)
	}
	if swa := response.Header.Get("Service-Worker-Allowed"); swa != "/sub/" {
		t.Errorf("expected Service-Worker-Allowed /sub/, got %q", swa)
	}
}

func TestServiceWorker_NotAuthGated(t *testing.T) {
	log.SetOutput(io.Discard)
	db, _ := storage.New(":memory:")
	log.SetOutput(os.Stderr)

	srv := NewServer(db, "127.0.0.1:8000")
	srv.Username = "admin"
	srv.Password = "secret"
	srv.SecureCookie = true
	handler := srv.handler()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/sw.js", nil)
	handler.ServeHTTP(recorder, request)

	if recorder.Result().StatusCode != 200 {
		t.Fatalf("sw.js should not require auth, got %d", recorder.Result().StatusCode)
	}
}

func TestManifestPWAFields(t *testing.T) {
	handler := NewServer(nil, "127.0.0.1:8000").handler()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/manifest.json", nil)
	handler.ServeHTTP(recorder, request)
	response := recorder.Result()

	if response.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", response.StatusCode)
	}

	var manifest map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&manifest); err != nil {
		t.Fatalf("failed to decode manifest: %v", err)
	}

	// Required PWA fields
	if manifest["scope"] == nil {
		t.Error("manifest missing scope")
	}
	if manifest["background_color"] != "#1a1a2e" {
		t.Errorf("expected background_color #1a1a2e, got %v", manifest["background_color"])
	}
	if manifest["theme_color"] != "#ffffff" {
		t.Errorf("expected theme_color #ffffff, got %v", manifest["theme_color"])
	}
	if manifest["display"] != "standalone" {
		t.Errorf("expected display standalone, got %v", manifest["display"])
	}

	// Check icons contain required PWA sizes
	icons, ok := manifest["icons"].([]interface{})
	if !ok {
		t.Fatal("manifest icons is not an array")
	}
	requiredSizes := map[string]bool{"192x192": false, "512x512": false, "180x180": false}
	for _, icon := range icons {
		entry, _ := icon.(map[string]interface{})
		if size, ok := entry["sizes"].(string); ok {
			if _, needed := requiredSizes[size]; needed {
				requiredSizes[size] = true
			}
		}
	}
	for size, found := range requiredSizes {
		if !found {
			t.Errorf("manifest missing required icon size %s", size)
		}
	}
}

func TestManifestWithBase(t *testing.T) {
	srv := NewServer(nil, "127.0.0.1:8000")
	srv.BasePath = "/yarr"
	handler := srv.handler()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/yarr/manifest.json", nil)
	handler.ServeHTTP(recorder, request)
	response := recorder.Result()

	if response.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", response.StatusCode)
	}

	var manifest map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&manifest); err != nil {
		t.Fatalf("failed to decode manifest: %v", err)
	}

	// Scope and start_url should include base path
	if manifest["scope"] != "/yarr" {
		t.Errorf("expected scope /yarr, got %v", manifest["scope"])
	}
	if manifest["start_url"] != "/yarr" {
		t.Errorf("expected start_url /yarr, got %v", manifest["start_url"])
	}

	// Icon paths should include base path
	icons := manifest["icons"].([]interface{})
	for _, icon := range icons {
		entry := icon.(map[string]interface{})
		src := entry["src"].(string)
		if !strings.HasPrefix(src, "/yarr/") {
			t.Errorf("icon src %q should be prefixed with /yarr/", src)
		}
	}
}
