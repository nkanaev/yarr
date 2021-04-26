package gzip

import (
	"compress/gzip"
	"net/http"
	"strings"

	"github.com/nkanaev/yarr/src/server/router"
)

type gzipResponseWriter struct {
	http.ResponseWriter

	out *gzip.Writer
	src http.ResponseWriter
}

func (rw *gzipResponseWriter) Header() http.Header {
	return rw.src.Header()
}

func (rw *gzipResponseWriter) Write(x []byte) (int, error) {
	return rw.out.Write(x)
}

func (rw *gzipResponseWriter) WriteHeader(statusCode int) {
	rw.src.WriteHeader(statusCode)
}

func Middleware(c *router.Context) {
	if !strings.Contains(c.Req.Header.Get("Accept-Encoding"), "gzip") {
		c.Next()
		return
	}

	gz := &gzipResponseWriter{out: gzip.NewWriter(c.Out), src: c.Out}
	defer gz.out.Close()

	c.Out.Header().Set("Content-Encoding", "gzip")
	c.Out = gz

	c.Next()
}
