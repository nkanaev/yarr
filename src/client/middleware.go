package client

import (
	"context"
	"net/http"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func UserAgent(userAgent string) Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			if userAgent != "" {
				req = req.Clone(req.Context())
				req.Header.Set("User-Agent", userAgent)
			}
			return next.RoundTrip(req)
		})
	}
}

type State struct {
	ETag         string
	LastModified string
}

type StateStore interface {
	Get(ctx context.Context) (State, error)
	Put(ctx context.Context, state State) error
}

func ConditionalRequests(store StateStore) Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return roundTripperFunc(func(r *http.Request) (*http.Response, error) {
			if r.Method != http.MethodGet && r.Method != http.MethodHead {
				return next.RoundTrip(r)
			}
			ctx := r.Context()

			state, err := store.Get(ctx)
			if err != nil {
				return next.RoundTrip(r)
			}

			if state.ETag != "" || state.LastModified != "" {
				r = r.Clone(ctx)
				if state.ETag != "" {
					r.Header.Set("If-None-Match", state.ETag)
				}
				if state.LastModified != "" {
					r.Header.Set("If-Modified-Since", state.LastModified)
				}
			}

			resp, err := next.RoundTrip(r)
			if err != nil {
				return nil, err
			}

			if resp.StatusCode == http.StatusOK {
				etag := resp.Header.Get("ETag")
				lmod := resp.Header.Get("Last-Modified")
				if etag != "" || lmod != "" {
					_ = store.Put(ctx, State{ETag: etag, LastModified: lmod})
				}
			}
			return resp, nil
		})
	}
}
