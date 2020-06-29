package worker

import (
	"net/url"
	"net/http"
	"github.com/PuerkitoBio/goquery"
)

type FeedSource struct {
	Title string `json:"title"`
	Url string `json:"url"`
}

const feedLinks = `link[type='application/rss+xml'],link[type='application/atom+xml']`


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
