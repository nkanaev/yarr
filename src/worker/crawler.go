package worker

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/nkanaev/yarr/src/crawler"
	feedparser "github.com/nkanaev/yarr/src/feed"
	"github.com/nkanaev/yarr/src/storage"
	"golang.org/x/net/html/charset"
)

type FeedSource struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

type Client struct {
	httpClient *http.Client
	userAgent  string
}

func (c *Client) get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.userAgent)
	return c.httpClient.Do(req)
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

var defaultClient *Client

func searchFeedLinks(html []byte, siteurl string) ([]FeedSource, error) {
	sources := make([]FeedSource, 0, 0)
	for url, title := range crawler.FindFeeds(string(html), siteurl) {
		sources = append(sources, FeedSource{Title: title, Url: url})
	}
	return sources, nil
}

func DiscoverFeed(candidateUrl string) (*feedparser.Feed, string, *[]FeedSource, error) {
	// Query URL
	res, err := defaultClient.get(candidateUrl)
	if err != nil {
		return nil, "", nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		errmsg := fmt.Sprintf("Failed to fetch feed %s (status: %d)", candidateUrl, res.StatusCode)
		return nil, "", nil, errors.New(errmsg)
	}
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, "", nil, err
	}

	// Try to feed into parser
	feed, err := feedparser.Parse(bytes.NewReader(content))
	if err == nil {
		/*
		// WILD: feeds may not always have link to themselves
		if len(feed.FeedLink) == 0 {
			feed.FeedLink = candidateUrl
		}
		*/

		// WILD: resolve relative links (path, without host)
		/*
		base, _ := url.Parse(candidateUrl)
		if link, err := url.Parse(feed.Link); err == nil && link.Host == "" {
			feed.Link = base.ResolveReference(link).String()
		}
		if link, err := url.Parse(feed.FeedLink); err == nil && link.Host == "" {
			feed.FeedLink = base.ResolveReference(link).String()
		}
		*/
		err := feed.TranslateURLs(candidateUrl)
		if err != nil {
			log.Printf("Failed to translate feed urls: %s", err)
		}

		return feed, candidateUrl, nil, nil
	}

	// Possibly an html link. Search for feed links
	sources, err := searchFeedLinks(content, candidateUrl)
	if err != nil {
		return nil, "", nil, err
	} else if len(sources) == 0 {
		return nil, "", nil, errors.New("No feeds found at the given url")
	} else if len(sources) == 1 {
		if sources[0].Url == candidateUrl {
			return nil, "", nil, errors.New("Recursion!")
		}
		return DiscoverFeed(sources[0].Url)
	}
	return nil, "", &sources, nil
}

func FindFavicon(websiteUrl, feedUrl string) (*[]byte, error) {
	candidateUrls := make([]string, 0)

	favicon := func(link string) string {
		u, err := url.Parse(link)
		if err != nil {
			return ""
		}
		return fmt.Sprintf("%s://%s/favicon.ico", u.Scheme, u.Host)
	}

	if len(websiteUrl) != 0 {
		res, err := defaultClient.get(websiteUrl)
		if err != nil {
			return nil, err
		}
		body, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			return nil, err
		}
		candidateUrls = append(candidateUrls, crawler.FindIcons(string(body), websiteUrl)...)
		if c := favicon(websiteUrl); len(c) != 0 {
			candidateUrls = append(candidateUrls, c)
		}
	}
	if c := favicon(feedUrl); len(c) != 0 {
		candidateUrls = append(candidateUrls, c)
	}

	imageTypes := [4]string{
		"image/x-icon",
		"image/png",
		"image/jpeg",
		"image/gif",
	}
	for _, url := range candidateUrls {
		res, err := defaultClient.get(url)
		if err != nil {
			continue
		}
		defer res.Body.Close()
		if res.StatusCode == 200 {
			if content, err := ioutil.ReadAll(res.Body); err == nil {
				ctype := http.DetectContentType(content)
				for _, itype := range imageTypes {
					if ctype == itype {
						return &content, nil
					}
				}
			}
		}
	}
	return nil, nil
}

func ConvertItems(items []feedparser.Item, feed storage.Feed) []storage.Item {
	result := make([]storage.Item, len(items))
	for i, item := range items {
		podcastUrl := item.PodcastURL

		/*
		var podcastUrl *string
		if item.Enclosures != nil {
			for _, enclosure := range item.Enclosures {
				if strings.ToLower(enclosure.Type) == "audio/mpeg" {
					podcastUrl = &enclosure.URL
				}
			}
		}
		*/
		result[i] = storage.Item{
			GUID:        item.GUID,
			FeedId:      feed.Id,
			Title:       item.Title,
			Link:        item.URL,
			Description: "",
			Content:     item.Content,
			Author:      "",
			Date:        &item.Date,
			Status:      storage.UNREAD,
			Image:       item.ImageURL,
			PodcastURL:  &podcastUrl,
		}
	}
	return result
}

func listItems(f storage.Feed, db *storage.Storage) ([]storage.Item, error) {
	var res *http.Response
	var err error

	httpState := db.GetHTTPState(f.Id)
	if httpState != nil {
		res, err = defaultClient.getConditional(f.FeedLink, httpState.LastModified, httpState.Etag)
	} else {
		res, err = defaultClient.get(f.FeedLink)
	}

	if err != nil {
		return nil, fmt.Errorf("unable to get: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode/100 == 4 || res.StatusCode/100 == 5 {
		return nil, fmt.Errorf("status code %d", res.StatusCode)
	}

	if res.StatusCode == 304 {
		return nil, nil
	}

	body, err := charset.NewReader(res.Body, res.Header.Get("Content-Type"))
	if err != nil {
		return nil, fmt.Errorf("failed to init response body: %s", err)
	}
	feed, err := feedparser.Parse(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse: %s", err)
	}

	lastModified := res.Header.Get("Last-Modified")
	etag := res.Header.Get("Etag")
	if lastModified != "" || etag != "" {
		db.SetHTTPState(f.Id, lastModified, etag)
	}
	return ConvertItems(feed.Items, f), nil
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
	defaultClient = &Client{
		httpClient: httpClient,
		userAgent:  "Yarr/1.0",
	}
}
