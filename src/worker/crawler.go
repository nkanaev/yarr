package worker

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/nkanaev/yarr/src/scraper"
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

	body, err := charset.NewReader(res.Body, res.Header.Get("Content-Type"))
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
		res, err := client.get(websiteUrl)
		if err != nil {
			return nil, err
		}
		body, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			return nil, err
		}
		candidateUrls = append(candidateUrls, scraper.FindIcons(string(body), websiteUrl)...)
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
		res, err := client.get(url)
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

func ConvertItems(items []parser.Item, feed storage.Feed) []storage.Item {
	result := make([]storage.Item, len(items))
	for i, item := range items {
		item := item
		var podcastURL *string = nil
		if item.AudioURL != "" {
			podcastURL = &item.AudioURL
		}
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
			PodcastURL:  podcastURL,
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
		return nil, fmt.Errorf("unable to get: %s", err)
	}
	defer res.Body.Close()

	switch {
	case res.StatusCode < 200 || res.StatusCode > 399:
		return nil, fmt.Errorf("status code %d", res.StatusCode)
	case res.StatusCode == http.StatusNotModified:
		return nil, nil
	}

	body, err := charset.NewReader(res.Body, res.Header.Get("Content-Type"))
	if err != nil {
		return nil, fmt.Errorf("failed to init response body: %s", err)
	}

	feed, err := parser.Parse(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse: %s", err)
	}

	lmod = res.Header.Get("Last-Modified")
	etag = res.Header.Get("Etag")
	if lmod != "" || etag != "" {
		db.SetHTTPState(f.Id, lmod, etag)
	}
	feed.TranslateURLs(f.FeedLink)
	return ConvertItems(feed.Items, f), nil
}
