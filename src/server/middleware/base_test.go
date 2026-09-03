package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBaseRedirect(t *testing.T) {
	handler := Base("/base")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/base", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusFound {
		t.Errorf("status = %d; want %d", rec.Code, http.StatusFound)
	}
	if loc := rec.Header().Get("Location"); loc != "/base/" {
		t.Errorf("location = %q; want %q", loc, "/base/")
	}
}

func TestBaseOutsidePath(t *testing.T) {
	called := false
	handler := Base("/base")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))

	req := httptest.NewRequest("GET", "/other", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d; want %d", rec.Code, http.StatusNotFound)
	}
	if called {
		t.Error("next handler was called")
	}
}

func TestBaseStripPrefix(t *testing.T) {
	var gotPath string
	handler := Base("/base")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
	}))

	req := httptest.NewRequest("GET", "/base/api/status", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if gotPath != "/api/status" {
		t.Errorf("path = %q; want %q", gotPath, "/api/status")
	}
}

func TestBaseEmpty(t *testing.T) {
	handler := Base("")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/status", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d; want %d", rec.Code, http.StatusOK)
	}
}
