package client

import (
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

type stubRoundTripper func(*http.Request) (*http.Response, error)

func (f stubRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func get(t *testing.T, client *http.Client, url string) {
	t.Helper()
	res, err := client.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	io.Copy(io.Discard, res.Body)
	res.Body.Close()
}

func buildTransport(t *testing.T, builder *ClientBuilder) *http.Transport {
	t.Helper()
	client := builder.Build()
	transport, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("Transport = %T; want *http.Transport with no middleware", client.Transport)
	}
	return transport
}

func TestNewBuilderUsesDefaults(t *testing.T) {
	wantTransport := DefaultTransport()
	wantTimeout := DefaultTimeout()
	client := NewBuilder().Build()
	transport := client.Transport.(*http.Transport)

	if transport.MaxIdleConns != wantTransport.MaxIdleConns {
		t.Errorf("MaxIdleConns = %d; want %d from DefaultTransport()", transport.MaxIdleConns, wantTransport.MaxIdleConns)
	}
	if transport.MaxIdleConnsPerHost != wantTransport.MaxIdleConnsPerHost {
		t.Errorf("MaxIdleConnsPerHost = %d; want %d from DefaultTransport()", transport.MaxIdleConnsPerHost, wantTransport.MaxIdleConnsPerHost)
	}
	if transport.MaxConnsPerHost != wantTransport.MaxConnsPerHost {
		t.Errorf("MaxConnsPerHost = %d; want %d from DefaultTransport()", transport.MaxConnsPerHost, wantTransport.MaxConnsPerHost)
	}
	if transport.IdleConnTimeout != wantTransport.IdleConnTimeout {
		t.Errorf("IdleConnTimeout = %v; want %v from DefaultTransport()", transport.IdleConnTimeout, wantTransport.IdleConnTimeout)
	}
	if transport.DisableKeepAlives {
		t.Error("DisableKeepAlives = true; want false")
	}
	if !transport.ForceAttemptHTTP2 {
		t.Error("ForceAttemptHTTP2 = false; want true")
	}
	if transport.Proxy == nil {
		t.Error("Proxy = nil; want non-nil")
	}

	if client.Timeout != wantTimeout.Request {
		t.Errorf("Timeout = %v; want %v from DefaultTimeout()", client.Timeout, wantTimeout.Request)
	}
	if transport.TLSHandshakeTimeout != wantTimeout.TLSHandshake {
		t.Errorf("TLSHandshakeTimeout = %v; want %v from DefaultTimeout()", transport.TLSHandshakeTimeout, wantTimeout.TLSHandshake)
	}
	if transport.ResponseHeaderTimeout != wantTimeout.ResponseHeader {
		t.Errorf("ResponseHeaderTimeout = %v; want %v from DefaultTimeout()", transport.ResponseHeaderTimeout, wantTimeout.ResponseHeader)
	}
	if transport.ExpectContinueTimeout != wantTimeout.ExpectContinue {
		t.Errorf("ExpectContinueTimeout = %v; want %v from DefaultTimeout()", transport.ExpectContinueTimeout, wantTimeout.ExpectContinue)
	}
	if transport.DialContext == nil {
		t.Error("DialContext = nil; want a dialer carrying the default dial timeout")
	}
}

func TestBuildPassesThroughTransport(t *testing.T) {
	custom := DefaultTransport()
	custom.MaxIdleConns = 10
	custom.MaxIdleConnsPerHost = 3
	custom.MaxConnsPerHost = 5
	custom.IdleConnTimeout = time.Minute

	transport := buildTransport(t, NewBuilder().Transport(custom))

	if transport.MaxIdleConns != 10 {
		t.Errorf("MaxIdleConns = %d; want 10", transport.MaxIdleConns)
	}
	if transport.MaxIdleConnsPerHost != 3 {
		t.Errorf("MaxIdleConnsPerHost = %d; want 3", transport.MaxIdleConnsPerHost)
	}
	if transport.MaxConnsPerHost != 5 {
		t.Errorf("MaxConnsPerHost = %d; want 5", transport.MaxConnsPerHost)
	}
	if transport.IdleConnTimeout != time.Minute {
		t.Errorf("IdleConnTimeout = %v; want %v", transport.IdleConnTimeout, time.Minute)
	}
}

func TestBuildWithoutMiddleware(t *testing.T) {
	if _, ok := NewBuilder().Build().Transport.(*http.Transport); !ok {
		t.Fatalf("Transport = %T; want *http.Transport with no middleware", NewBuilder().Build().Transport)
	}
}

func TestMiddlewareFirstRegisteredIsOutermost(t *testing.T) {
	var order []string

	recording := func(name string, order *[]string) Middleware {
		return func(next http.RoundTripper) http.RoundTripper {
			return stubRoundTripper(func(req *http.Request) (*http.Response, error) {
				*order = append(*order, name)
				return next.RoundTrip(req)
			})
		}
	}

	client := NewBuilder().
		Middleware(recording("a", &order), recording("b", &order)).
		Middleware(nil). // must be ignored
		Middleware(recording("c", &order)).
		Build()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer server.Close()

	get(t, client, server.URL)

	if want := []string{"a", "b", "c"}; !reflect.DeepEqual(order, want) {
		t.Errorf("order = %v; want %v", order, want)
	}
}

func TestBuildReturnsIndependentClients(t *testing.T) {
	builder := NewBuilder()
	first := builder.Build()
	second := builder.Build()

	if first == second {
		t.Error("Build returned the same *http.Client twice; want independent instances")
	}
	if first.Transport == second.Transport {
		t.Error("Build shared a Transport; want independent transports")
	}
}
