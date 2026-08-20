package gzip

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"
)

// gzipResponseWriter wraps http.ResponseWriter to route writes through gzip.Writer
type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// Pool gzip writers to reduce memory allocations under high traffic
var gzipPool = sync.Pool{
	New: func() any {
		return gzip.NewWriter(io.Discard)
	},
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Check if client accepts gzip encoding
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// 2. Set headers
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Vary", "Accept-Encoding")

		// 3. Acquire gzip writer from pool
		gz := gzipPool.Get().(*gzip.Writer)
		defer gzipPool.Put(gz)

		gz.Reset(w)
		defer gz.Close()

		// 4. Wrap response writer and execute next handler
		gzw := gzipResponseWriter{
			ResponseWriter: w,
			Writer:         gz,
		}

		next.ServeHTTP(gzw, r)
	})
}
