package client

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	lmodFixture        = "Wed, 04 Mar 2026 05:06:07 GMT"
	lmodFixtureUpdated = "Thu, 07 May 2026 08:09:10 GMT"
)

type memoryStateStore struct {
	state  State
	getErr error
	putErr error
	gets   int
	puts   int
}

func (s *memoryStateStore) Get(ctx context.Context) (State, error) {
	s.gets++
	return s.state, s.getErr
}

func (s *memoryStateStore) Put(ctx context.Context, state State) error {
	s.puts++
	s.state = state
	return s.putErr
}

func statusResponder(status int, header http.Header) http.RoundTripper {
	return stubRoundTripper(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: status,
			Header:     header,
			Body:       http.NoBody,
			Request:    r,
		}, nil
	})
}

func validatorHeader(etag, lastModified string) http.Header {
	header := http.Header{}
	if etag != "" {
		header.Set("Etag", etag)
	}
	if lastModified != "" {
		header.Set("Last-Modified", lastModified)
	}
	return header
}

// roundTrip runs one request through the middleware chain and reports the
// request the innermost round tripper saw.
func roundTrip(t *testing.T, middleware Middleware, req *http.Request) *http.Request {
	t.Helper()
	var seen *http.Request
	capture := stubRoundTripper(func(r *http.Request) (*http.Response, error) {
		seen = r
		return &http.Response{StatusCode: http.StatusOK, Header: http.Header{}, Body: http.NoBody, Request: r}, nil
	})
	res, err := middleware(capture).RoundTrip(req)
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if seen == nil {
		t.Fatal("inner round tripper was not reached")
	}
	return seen
}

func TestUserAgentSetsAndOverwritesHeader(t *testing.T) {
	var seen string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = r.Header.Get("User-Agent")
	}))
	defer server.Close()

	client := NewBuilder().Middleware(UserAgent("Yarr/2.0")).Build()

	get(t, client, server.URL)
	if seen != "Yarr/2.0" {
		t.Errorf("User-Agent = %q; want %q", seen, "Yarr/2.0")
	}
}

func TestConditionalRequestsReadsFromStore(t *testing.T) {
	store := &memoryStateStore{state: State{ETag: "abc"}}

	req, _ := http.NewRequest(http.MethodGet, "http://example.com/", nil)
	roundTrip(t, ConditionalRequests(store), req)

	if store.gets != 1 {
		t.Errorf("gets = %d; want 1", store.gets)
	}
}

func TestConditionalRequestsPassesHeaders(t *testing.T) {
	store := &memoryStateStore{state: State{ETag: "abc", LastModified: lmodFixture}}
	req, _ := http.NewRequest(http.MethodGet, "http://example.com/", nil)
	seen := roundTrip(t, ConditionalRequests(store), req)
	if got := seen.Header.Get("If-None-Match"); got != "abc" {
		t.Errorf("If-None-Match = %q; want %q", got, "abc")
	}
	if got, want := seen.Header.Get("If-Modified-Since"), lmodFixture; got != want {
		t.Errorf("If-Modified-Since = %q; want %q", got, want)
	}
	if req.Header.Get("If-None-Match") != "" || req.Header.Get("If-Modified-Since") != "" {
		t.Error("original request was mutated; want it left unmodified")
	}
}

func TestConditionalRequestsStoresHeaders(t *testing.T) {
	t.Run("200 stores response validators", func(t *testing.T) {
		store := &memoryStateStore{state: State{ETag: "old", LastModified: lmodFixtureUpdated}}
		req, _ := http.NewRequest(http.MethodGet, "http://example.com/", nil)
		res, err := ConditionalRequests(store)(statusResponder(http.StatusOK, validatorHeader("new-etag", lmodFixture))).
			RoundTrip(req)
		if err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
		if store.puts != 1 {
			t.Errorf("puts = %d; want 1", store.puts)
		}
		if want := (State{ETag: "new-etag", LastModified: lmodFixture}); store.state != want {
			t.Errorf("stored state = %+v; want %+v, a 200 replaces the stored state", store.state, want)
		}
	})
}

func TestConditionalRequestsContinuesOnStoreGetError(t *testing.T) {
	store := &memoryStateStore{
		state:  State{ETag: "abc"},
		getErr: errors.New("store unavailable"),
	}

	req, _ := http.NewRequest(http.MethodGet, "http://example.com/", nil)
	seen := roundTrip(t, ConditionalRequests(store), req)

	if got := seen.Header.Get("If-None-Match"); got != "" {
		t.Errorf("If-None-Match = %q; want empty on a store read failure", got)
	}
	if store.puts != 0 {
		t.Errorf("puts = %d; want 0 when the read failed", store.puts)
	}
}

func TestConditionalRequestsContinuesOnStorePutError(t *testing.T) {
	store := &memoryStateStore{putErr: errors.New("store write failed")}
	req, _ := http.NewRequest(http.MethodGet, "http://example.com/", nil)
	res, err := ConditionalRequests(store)(statusResponder(http.StatusOK, validatorHeader("etag", ""))).
		RoundTrip(req)
	if err != nil {
		t.Fatalf("err = %v; want the already-received response returned uncached", err)
	}
	res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d; want %d", res.StatusCode, http.StatusOK)
	}
	if store.puts != 1 {
		t.Errorf("puts = %d; want 1", store.puts)
	}
}
