package server

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/nkanaev/yarr/src/server/router"
)

// handleAiProxy forwards /api/ai/* requests to the Python AI service.
// SSE streams are passed through transparently.
func (s *Server) handleAiProxy(c *router.Context) {
	if s.AiServiceURL == "" {
		c.Out.WriteHeader(http.StatusServiceUnavailable)
		c.Out.Write([]byte(`{"error":"AI service not configured"}`))
		return
	}

	// Strip /api/ai prefix to get the path for the Python service
	path := c.Req.URL.Path
	prefix := s.BasePath + "/api/ai"
	aiPath := strings.TrimPrefix(path, prefix)
	if aiPath == "" {
		aiPath = "/"
	}

	targetURL := strings.TrimRight(s.AiServiceURL, "/") + aiPath
	if c.Req.URL.RawQuery != "" {
		targetURL += "?" + c.Req.URL.RawQuery
	}

	// Build proxy request
	proxyReq, err := http.NewRequestWithContext(c.Req.Context(), c.Req.Method, targetURL, c.Req.Body)
	if err != nil {
		log.Printf("AI proxy: failed to create request: %v", err)
		c.Out.WriteHeader(http.StatusBadGateway)
		return
	}

	// Copy relevant headers
	for _, h := range []string{"Content-Type", "Accept", "HX-Request"} {
		if v := c.Req.Header.Get(h); v != "" {
			proxyReq.Header.Set(h, v)
		}
	}

	// Use a client with appropriate timeout for streaming
	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Do(proxyReq)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		log.Printf("AI proxy: request failed: %v", err)
		c.Out.WriteHeader(http.StatusBadGateway)
		c.Out.Write([]byte(`{"error":"AI service unavailable"}`))
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for k, vs := range resp.Header {
		for _, v := range vs {
			c.Out.Header().Add(k, v)
		}
	}
	c.Out.WriteHeader(resp.StatusCode)

	// Stream the response body (important for SSE)
	if f, ok := c.Out.(http.Flusher); ok {
		buf := make([]byte, 4096)
		for {
			n, err := resp.Body.Read(buf)
			if n > 0 {
				c.Out.Write(buf[:n])
				f.Flush()
			}
			if err != nil {
				break
			}
		}
	} else {
		io.Copy(c.Out, resp.Body)
	}
}
