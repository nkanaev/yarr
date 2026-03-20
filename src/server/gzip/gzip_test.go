package gzip

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/nkanaev/yarr/src/server/router"
)

func TestMiddleware_CompressesResponse(t *testing.T) {
	body := "Hello, this is a test response body that should be compressed."

	r := router.NewRouter("")
	r.Use(Middleware)
	r.For("/test", func(c *router.Context) {
		c.Out.Write([]byte(body))
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/test", nil)
	request.Header.Set("Accept-Encoding", "gzip")
	r.ServeHTTP(recorder, request)

	result := recorder.Result()
	if result.Header.Get("Content-Encoding") != "gzip" {
		t.Fatal("expected Content-Encoding: gzip")
	}

	gr, err := gzip.NewReader(result.Body)
	if err != nil {
		t.Fatalf("failed to create gzip reader: %v", err)
	}
	defer gr.Close()

	decoded, err := io.ReadAll(gr)
	if err != nil {
		t.Fatalf("failed to decompress: %v", err)
	}
	if string(decoded) != body {
		t.Errorf("decompressed body mismatch:\ngot:  %q\nwant: %q", string(decoded), body)
	}
}

func TestMiddleware_NoCompressWithoutHeader(t *testing.T) {
	body := "Uncompressed response body"

	r := router.NewRouter("")
	r.Use(Middleware)
	r.For("/test", func(c *router.Context) {
		c.Out.Write([]byte(body))
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/test", nil)
	// No Accept-Encoding header
	r.ServeHTTP(recorder, request)

	result := recorder.Result()
	if result.Header.Get("Content-Encoding") == "gzip" {
		t.Fatal("should not compress without Accept-Encoding: gzip")
	}

	got, _ := io.ReadAll(result.Body)
	if string(got) != body {
		t.Errorf("body mismatch:\ngot:  %q\nwant: %q", string(got), body)
	}
}

func TestMiddleware_CompressesLargeBody(t *testing.T) {
	// Large body should compress well
	body := bytes.Repeat([]byte("abcdefghijklmnop"), 1000)

	r := router.NewRouter("")
	r.Use(Middleware)
	r.For("/test", func(c *router.Context) {
		c.Out.Write(body)
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/test", nil)
	request.Header.Set("Accept-Encoding", "gzip")
	r.ServeHTTP(recorder, request)

	result := recorder.Result()
	compressed, _ := io.ReadAll(result.Body)
	if len(compressed) >= len(body) {
		t.Errorf("compressed size (%d) should be smaller than original (%d)", len(compressed), len(body))
	}

	gr, _ := gzip.NewReader(bytes.NewReader(compressed))
	decoded, _ := io.ReadAll(gr)
	if !bytes.Equal(decoded, body) {
		t.Error("decompressed content doesn't match original")
	}
}

func TestMiddleware_AcceptEncodingPartialMatch(t *testing.T) {
	r := router.NewRouter("")
	r.Use(Middleware)
	r.For("/test", func(c *router.Context) {
		c.Out.Write([]byte("test"))
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/test", nil)
	request.Header.Set("Accept-Encoding", "deflate, gzip, br")
	r.ServeHTTP(recorder, request)

	if recorder.Result().Header.Get("Content-Encoding") != "gzip" {
		t.Fatal("should compress when gzip is among accepted encodings")
	}
}
