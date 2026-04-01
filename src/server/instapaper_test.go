package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInstapaperAdd_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != "test@example.com" || pass != "secret" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		if r.URL.Query().Get("url") == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer srv.Close()

	origURL := instapaperAddURL
	instapaperAddURL = srv.URL
	defer func() { instapaperAddURL = origURL }()

	err := InstapaperAdd("test@example.com", "secret", "https://example.com/article", "Test Article")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestInstapaperAdd_BadCredentials(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	origURL := instapaperAddURL
	instapaperAddURL = srv.URL
	defer func() { instapaperAddURL = origURL }()

	err := InstapaperAdd("wrong", "creds", "https://example.com", "")
	if err == nil {
		t.Fatal("expected error for bad credentials")
	}
}
