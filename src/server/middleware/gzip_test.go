package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGzipPassThrough(t *testing.T) {
	payload := "hello world"
	handler := Gzip(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, payload)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if got := rec.Header().Get("Content-Encoding"); got != "" {
		t.Errorf("content-encoding = %q; want empty", got)
	}
	if got := rec.Body.String(); got != payload {
		t.Errorf("body = %q; want %q", got, payload)
	}
}

func TestGzipCompressed(t *testing.T) {
	payload := "hello world"
	handler := Gzip(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, payload)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if got := rec.Header().Get("Content-Encoding"); got != "gzip" {
		t.Errorf("content-encoding = %q; want %q", got, "gzip")
	}
	reader, err := gzip.NewReader(rec.Body)
	if err != nil {
		t.Fatal(err)
	}
	body, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	if got := string(body); got != payload {
		t.Errorf("body = %q; want %q", got, payload)
	}
}
