package server

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
	"github.com/nkanaev/yarr/storage"
	"io/ioutil"
	"net/http"
	"net/url"
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

func FindFeeds(r *http.Response) ([]FeedSource, error) {
	sources := make([]FeedSource, 0, 0)
	doc, err := goquery.NewDocumentFromResponse(r)
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
			url := doc.Url.ResolveReference(feedUrl).String()
			sources = append(sources, FeedSource{Title: title, Url: url})
		}
	})
	return sources, nil
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

	client := http.Client{}

	imageTypes := [4]string{
		"image/x-icon",
		"image/png",
		"image/jpeg",
		"image/gif",
	}
	for _, url := range candidateUrls {
		if res, err := client.Get(url); err == nil && res.StatusCode == 200 {
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
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(f.FeedLink)
	if err != nil {
		return nil, err
	}
	return convertItems(feed.Items, f), nil
}

func createFeed(s *storage.Storage, url string, folderId *int64) error {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return err
	}
	feedLink := feed.FeedLink
	if len(feedLink) == 0 {
		feedLink = url
	}
	storedFeed := s.CreateFeed(
		feed.Title,
		feed.Description,
		feed.Link,
		feedLink,
		folderId,
	)
	s.CreateItems(convertItems(feed.Items, *storedFeed))
	return nil
}
