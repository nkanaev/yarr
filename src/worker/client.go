package worker

import (
	"net"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	userAgent  string
}

func (c *Client) get(url string) (*http.Response, error) {
	return c.getConditional(url, "", "")
}

func (c *Client) getConditional(url, lastModified, etag string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.userAgent)
	if lastModified != "" {
		req.Header.Set("If-Modified-Since", lastModified)
	}
	if etag != "" {
		req.Header.Set("If-None-Match", etag)
	}
	return c.httpClient.Do(req)
}

var client *Client

func SetVersion(num string) {
	client.userAgent = "Yarr/" + num
}

func init() {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).DialContext,
		DisableKeepAlives:   true,
		TLSHandshakeTimeout: time.Second * 10,
	}
	httpClient := &http.Client{
		Timeout:   time.Second * 30,
		Transport: transport,
	}
	client = &Client{
		httpClient: httpClient,
		userAgent:  "Yarr/1.0",
	}
}
