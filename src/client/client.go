package client

import (
	"net"
	"net/http"
	"net/url"
	"slices"
	"time"
)

type Transport struct {
	Proxy               func(*http.Request) (*url.URL, error)
	DisableKeepAlives   bool
	ForceAttemptHTTP2   bool
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	MaxConnsPerHost     int
	IdleConnTimeout     time.Duration
}

type Middleware func(http.RoundTripper) http.RoundTripper

type Timeout struct {
	Request        time.Duration
	Dial           time.Duration
	TLSHandshake   time.Duration
	ResponseHeader time.Duration
	ExpectContinue time.Duration
}

type ClientBuilder struct {
	transport  Transport
	middleware []Middleware
	timeout    Timeout
}

func DefaultTransport() Transport {
	return Transport{
		Proxy:               http.ProxyFromEnvironment,
		DisableKeepAlives:   false,
		ForceAttemptHTTP2:   true,
		MaxIdleConns:        64,
		MaxIdleConnsPerHost: 1,
		MaxConnsPerHost:     0,
		IdleConnTimeout:     90 * time.Second,
	}
}

func DefaultTimeout() Timeout {
	return Timeout{
		Request:        30 * time.Second,
		Dial:           10 * time.Second,
		TLSHandshake:   10 * time.Second,
		ResponseHeader: 15 * time.Second,
		ExpectContinue: time.Second,
	}
}

func NewBuilder() *ClientBuilder {
	return &ClientBuilder{
		transport: DefaultTransport(),
		timeout:   DefaultTimeout(),
	}
}

func (b *ClientBuilder) Transport(transport Transport) *ClientBuilder {
	b.transport = transport
	return b
}

func (b *ClientBuilder) Timeout(timeout Timeout) *ClientBuilder {
	b.timeout = timeout
	return b
}

func (b *ClientBuilder) Middleware(middleware ...Middleware) *ClientBuilder {
	for _, m := range middleware {
		if m != nil {
			b.middleware = append(b.middleware, m)
		}
	}
	return b
}

func (b *ClientBuilder) Build() *http.Client {
	transport := &http.Transport{
		Proxy:               b.transport.Proxy,
		DialContext:         (&net.Dialer{Timeout: b.timeout.Dial}).DialContext,
		DisableKeepAlives:   b.transport.DisableKeepAlives,
		ForceAttemptHTTP2:   b.transport.ForceAttemptHTTP2,
		MaxIdleConns:        b.transport.MaxIdleConns,
		MaxIdleConnsPerHost: b.transport.MaxIdleConnsPerHost,
		MaxConnsPerHost:     b.transport.MaxConnsPerHost,
		IdleConnTimeout:     b.transport.IdleConnTimeout,
		TLSHandshakeTimeout: b.timeout.TLSHandshake,

		ResponseHeaderTimeout: b.timeout.ResponseHeader,
		ExpectContinueTimeout: b.timeout.ExpectContinue,
	}

	var roundTripper http.RoundTripper = transport
	for _, middleware := range slices.Backward(b.middleware) {
		roundTripper = middleware(roundTripper)
	}

	return &http.Client{
		Timeout:   b.timeout.Request,
		Transport: roundTripper,
	}
}
