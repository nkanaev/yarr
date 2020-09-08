package server

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
	"github.com/nkanaev/yarr/storage"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type FeedSource struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

const feedLinks = `
	link[type='application/rss+xml'],
	link[type='application/atom+xml'],
	a[href$="/feed"],
	a[href$="/feed/"],
	a[href$="feed.xml"],
	a[href$="atom.xml"],
	a[href$="rss.xml"],
	a:contains("rss"),
	a:contains("RSS"),
	a:contains("feed"),
	a:contains("FEED")
`

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

var defaultClient *Client

func searchFeedLinks(html []byte, siteurl string) ([]FeedSource, error) {
	sources := make([]FeedSource, 0, 0)

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return sources, err
	}
	base, err := url.Parse(siteurl)
	if err != nil {
		return sources, err
	}

	doc.Find(feedLinks).Each(func(i int, s *goquery.Selection) {
		if href, ok := s.Attr("href"); ok {
			feedUrl, err := url.Parse(href)
			if err != nil {
				return
			}
			title := s.AttrOr("title", "")
			url := base.ResolveReference(feedUrl).String()
			sources = append(sources, FeedSource{Title: title, Url: url})
		}
	})
	return sources, nil
}

func discoverFeed(url string) (*gofeed.Feed, *[]FeedSource, error) {
	// Query URL
	res, err := defaultClient.get(url)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		errmsg := fmt.Sprintf("Failed to fetch feed %s (status: %d)", url, res.StatusCode)
		return nil, nil, errors.New(errmsg)
	}
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}

	// Try to feed into parser
	feedparser := gofeed.NewParser()
	feed, err := feedparser.Parse(bytes.NewReader(content))
	if err == nil {
		// WILD: feeds may not always have link to themselves
		if len(feed.FeedLink) == 0 {
			feed.FeedLink = url
		}
		return feed, nil, nil
	}

	// Possibly an html link. Search for feed links
	sources, err := searchFeedLinks(content, url)
	if err != nil {
		return nil, nil, err
	} else if len(sources) == 0 {
		return nil, nil, errors.New("No feeds found at the given url")
	} else if len(sources) == 1 {
		if sources[0].Url == url {
			return nil, nil, errors.New("Recursion!")
		}
		return discoverFeed(sources[0].Url)
	}
	return nil, &sources, nil
}

func findFavicon(websiteUrl, feedUrl string) (*[]byte, error) {
	candidateUrls := make([]string, 0)

	favicon := func(link string) string {
		u, err := url.Parse(link)
		if err != nil {
			return ""
		}
		return fmt.Sprintf("%s://%s/favicon.ico", u.Scheme, u.Host)
	}

	if len(websiteUrl) != 0 {
		doc, err := goquery.NewDocument(websiteUrl)
		if err != nil {
			return nil, err
		}
		doc.Find(`link[rel=icon]`).EachWithBreak(func(i int, s *goquery.Selection) bool {
			if href, ok := s.Attr("href"); ok {
				if hrefUrl, err := url.Parse(href); err == nil {
					faviconUrl := doc.Url.ResolveReference(hrefUrl).String()
					candidateUrls = append(candidateUrls, faviconUrl)
				}
			}
			return true
		})

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

func convertItems(items []*gofeed.Item, feed storage.Feed) []storage.Item {
	result := make([]storage.Item, len(items))
	for i, item := range items {
		imageURL := ""
		if item.Image != nil {
			imageURL = item.Image.URL
		}
		author := ""
		if item.Author != nil {
			author = item.Author.Name
		}
		result[i] = storage.Item{
			GUID:        item.GUID,
			FeedId:      feed.Id,
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
			Content:     item.Content,
			Author:      author,
			Date:        item.PublishedParsed,
			DateUpdated: item.UpdatedParsed,
			Status:      storage.UNREAD,
			Image:       imageURL,
		}
	}
	return result
}

func listItems(f storage.Feed) ([]storage.Item, error) {
	res, err := defaultClient.get(f.FeedLink)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	feedparser := gofeed.NewParser()
	feed, err := feedparser.Parse(res.Body)
	if err != nil {
		return nil, err
	}
	return convertItems(feed.Items, f), nil
}

func init() {
	defaultClient = &Client{
		httpClient: &http.Client{Timeout: time.Second * 5},
		userAgent:  "Yarr/1.0",
	}
}
