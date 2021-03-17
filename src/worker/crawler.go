package worker

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
	"github.com/nkanaev/yarr/src/storage"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
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

func (c *Client) getConditional(url, lastModified, etag string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("If-Modified-Since", lastModified)
	req.Header.Set("If-None-Match", etag)
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

	// feed {url: title} map
	feeds := make(map[string]string)

	doc.Find(feedLinks).Each(func(i int, s *goquery.Selection) {
		// Unlikely to happen, but don't get more than N links
		if len(feeds) > 10 {
			return
		}
		if href, ok := s.Attr("href"); ok {
			feedUrl, err := url.Parse(href)
			if err != nil {
				return
			}

			title := s.AttrOr("title", "")
			url := base.ResolveReference(feedUrl).String()

			if _, alreadyExists := feeds[url]; alreadyExists {
				if feeds[url] == "" {
					feeds[url] = title
				}
			} else {
				feeds[url] = title
			}
		}
	})

	for url, title := range feeds {
		sources = append(sources, FeedSource{Title: title, Url: url})
	}
	return sources, nil
}

func DiscoverFeed(candidateUrl string) (*gofeed.Feed, *[]FeedSource, error) {
	// Query URL
	res, err := defaultClient.get(candidateUrl)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		errmsg := fmt.Sprintf("Failed to fetch feed %s (status: %d)", candidateUrl, res.StatusCode)
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
			feed.FeedLink = candidateUrl
		}

		// WILD: resolve relative links (path, without host)
		base, _ := url.Parse(candidateUrl)
		if link, err := url.Parse(feed.Link); err == nil && link.Host == "" {
			feed.Link = base.ResolveReference(link).String()
		}
		if link, err := url.Parse(feed.FeedLink); err == nil && link.Host == "" {
			feed.FeedLink = base.ResolveReference(link).String()
		}

		return feed, nil, nil
	}

	// Possibly an html link. Search for feed links
	sources, err := searchFeedLinks(content, candidateUrl)
	if err != nil {
		return nil, nil, err
	} else if len(sources) == 0 {
		return nil, nil, errors.New("No feeds found at the given url")
	} else if len(sources) == 1 {
		if sources[0].Url == candidateUrl {
			return nil, nil, errors.New("Recursion!")
		}
		return DiscoverFeed(sources[0].Url)
	}
	return nil, &sources, nil
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
		base, err := url.Parse(websiteUrl)
		if err != nil {
			return nil, err
		}
		res, err := defaultClient.get(websiteUrl)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			return nil, err
		}
		doc.Find(`link[rel=icon]`).EachWithBreak(func(i int, s *goquery.Selection) bool {
			if href, ok := s.Attr("href"); ok {
				if hrefUrl, err := url.Parse(href); err == nil {
					faviconUrl := base.ResolveReference(hrefUrl).String()
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

func ConvertItems(items []*gofeed.Item, feed storage.Feed) []storage.Item {
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
		var podcastUrl *string
		if item.Enclosures != nil {
			for _, enclosure := range item.Enclosures {
				if strings.ToLower(enclosure.Type) == "audio/mpeg" {
					podcastUrl = &enclosure.URL
				}
			}
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
			PodcastURL:  podcastUrl,
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
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode/100 == 4 || res.StatusCode/100 == 5 {
		errmsg := fmt.Sprintf("Failed to list feed items for %s (status: %d)", f.FeedLink, res.StatusCode)
		return nil, errors.New(errmsg)
	}

	if res.StatusCode == 304 {
		return nil, nil
	}

	lastModified := res.Header.Get("Last-Modified")
	etag := res.Header.Get("Etag")
	if lastModified != "" || etag != "" {
		db.SetHTTPState(f.Id, lastModified, etag)
	}

	feedparser := gofeed.NewParser()
	feed, err := feedparser.Parse(res.Body)
	if err != nil {
		return nil, err
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
