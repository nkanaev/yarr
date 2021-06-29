package worker

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/nkanaev/yarr/src/content/scraper"
	"github.com/nkanaev/yarr/src/parser"
	"github.com/nkanaev/yarr/src/storage"
	"golang.org/x/net/html/charset"
)

type FeedSource struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

type DiscoverResult struct {
	Feed     *parser.Feed
	FeedLink string
	Sources  []FeedSource
}

func DiscoverFeed(candidateUrl string) (*DiscoverResult, error) {
	result := &DiscoverResult{}
	// Query URL
	res, err := client.get(candidateUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code %d", res.StatusCode)
	}

	body, err := httpBody(res)
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	// Try to feed into parser
	feed, err := parser.Parse(bytes.NewReader(content))
	if err == nil {
		feed.TranslateURLs(candidateUrl)
		feed.SetMissingDatesTo(time.Now())
		result.Feed = feed
		result.FeedLink = candidateUrl
		return result, nil
	}

	// Possibly an html link. Search for feed links
	sources := make([]FeedSource, 0)
	for url, title := range scraper.FindFeeds(string(content), candidateUrl) {
		sources = append(sources, FeedSource{Title: title, Url: url})
	}
	switch {
	case len(sources) == 0:
		return nil, errors.New("No feeds found at the given url")
	case len(sources) == 1:
		if sources[0].Url == candidateUrl {
			return nil, errors.New("Recursion!")
		}
		return DiscoverFeed(sources[0].Url)
	}

	result.Sources = sources
	return result, nil
}

var emptyIcon = make([]byte, 0)
var imageTypes = map[string]bool{
	"image/x-icon": true,
	"image/png":    true,
	"image/jpeg":   true,
	"image/gif":    true,
}

func findFavicon(siteUrl, feedUrl string) (*[]byte, error) {
	urls := make([]string, 0)

	favicon := func(link string) string {
		u, err := url.Parse(link)
		if err != nil {
			return ""
		}
		return fmt.Sprintf("%s://%s/favicon.ico", u.Scheme, u.Host)
	}

	if siteUrl != "" {
		if res, err := client.get(siteUrl); err == nil {
			defer res.Body.Close()
			if body, err := ioutil.ReadAll(res.Body); err == nil {
				urls = append(urls, scraper.FindIcons(string(body), siteUrl)...)
				if c := favicon(siteUrl); c != "" {
					urls = append(urls, c)
				}
			}
		}
	}

	if c := favicon(feedUrl); c != "" {
		urls = append(urls, c)
	}

	for _, u := range urls {
		res, err := client.get(u)
		if err != nil {
			continue
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			continue
		}

		content, err := ioutil.ReadAll(res.Body)
		if err != nil {
			continue
		}

		ctype := http.DetectContentType(content)
		if imageTypes[ctype] {
			return &content, nil
		}
	}
	return &emptyIcon, nil
}

func ConvertItems(items []parser.Item, feed storage.Feed) []storage.Item {
	result := make([]storage.Item, len(items))
	for i, item := range items {
		item := item
		var audioURL *string = nil
		if item.AudioURL != "" {
			audioURL = &item.AudioURL
		}
		var imageURL *string = nil
		if item.ImageURL != "" {
			imageURL = &item.ImageURL
		}
		result[i] = storage.Item{
			GUID:     item.GUID,
			FeedId:   feed.Id,
			Title:    item.Title,
			Link:     item.URL,
			Content:  item.Content,
			Date:     item.Date,
			Status:   storage.UNREAD,
			ImageURL: imageURL,
			AudioURL: audioURL,
		}
	}
	return result
}

func listItems(f storage.Feed, db *storage.Storage) ([]storage.Item, error) {
	lmod := ""
	etag := ""
	if state := db.GetHTTPState(f.Id); state != nil {
		lmod = state.LastModified
		etag = state.Etag
	}

	res, err := client.getConditional(f.FeedLink, lmod, etag)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	switch {
	case res.StatusCode < 200 || res.StatusCode > 399:
		if res.StatusCode == 404 {
			return nil, fmt.Errorf("feed not found")
		}
		return nil, fmt.Errorf("status code %d", res.StatusCode)
	case res.StatusCode == http.StatusNotModified:
		return nil, nil
	}

	body, err := httpBody(res)
	if err != nil {
		return nil, err
	}

	feed, err := parser.Parse(body)
	if err != nil {
		return nil, err
	}

	lmod = res.Header.Get("Last-Modified")
	etag = res.Header.Get("Etag")
	if lmod != "" || etag != "" {
		db.SetHTTPState(f.Id, lmod, etag)
	}
	feed.TranslateURLs(f.FeedLink)
	feed.SetMissingDatesTo(time.Now())
	return ConvertItems(feed.Items, f), nil
}

func httpBody(res *http.Response) (io.Reader, error) {
	ctype := res.Header.Get("Content-Type")
	if strings.Contains(ctype, "charset") {
		return charset.NewReader(res.Body, ctype)
	}
	return res.Body, nil
}
