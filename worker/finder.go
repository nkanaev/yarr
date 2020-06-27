package worker

import (
	"log"
	//"net/http"
	"net/url"
	"github.com/PuerkitoBio/goquery"
)

type FeedSource struct {
	Title string
	Url *url.URL
}


func FindFeeds(u string) []FeedSource {
	doc, err := goquery.NewDocument(u)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(doc.Url)
	// Find the review items
	sources := make([]FeedSource, 0, 0)
	doc.Find("link[type='application/rss+xml'],link[type='application/atom+xml']").Each(func(i int, s *goquery.Selection) {
		if href, ok := s.Attr("href"); ok {
			feedUrl, feedErr := url.Parse(href)
			if feedErr != nil {
				log.Fatal(err)
			}
			title := s.AttrOr("title", "")
			feedUrl = doc.Url.ResolveReference(feedUrl)
			sources = append(sources, FeedSource{Title: title, Url: feedUrl})
		}
	})
	return sources
}
