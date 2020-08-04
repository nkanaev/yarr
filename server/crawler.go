package server

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"net/url"
)

type FeedSource struct {
	Title string `json:"title"`
	Url   string `json:"url"`
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
