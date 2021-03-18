package crawler

import (
	"strings"

	"golang.org/x/net/html"
)

func FindFeeds(body string, base string) map[string]string {
	candidates := make(map[string]string)

	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		return candidates
	}

	linkTypes := []string{"application/atom+xml", "application/rss+xml", "application/json"}
	isFeedLink := func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == "link" {
			t := getAttr(n, "type")
			for _, tt := range linkTypes {
				if tt == t {
					return true
				}
			}
		}
		return false
	}
	for _, node := range getNodes(doc, isFeedLink) {
		href := getAttr(node, "href")
		name := getAttr(node, "title")
		link := absoluteUrl(href, base)
		if link != "" {
			candidates[link] = name
		}
	}

	if len(candidates) == 0 {
		// guess by hyperlink properties:
		// - a[href="feed"]
		// - a:contains("rss")
		// ...etc
		feedHrefs := []string{"feed", "feed.xml", "rss.xml", "atom.xml"}
		feedTexts := []string{"rss", "feed"}
		isFeedHyperLink := func(n *html.Node) bool {
			if n.Type == html.ElementNode && n.Data == "a" {
				href := strings.Trim(getAttr(n, "href"), "/")
				text := getText(n)

				for _, feedHref := range feedHrefs {
					if strings.HasSuffix(href, feedHref) {
						return true
					}
				}
				for _, feedText := range feedTexts {
					if strings.EqualFold(text, feedText) {
						return true
					}
				}
			}
			return false
		}
		for _, node := range getNodes(doc, isFeedHyperLink) {
			href := getAttr(node, "href")
			link := absoluteUrl(href, base)
			if link != "" {
				candidates[link] = ""
			}
		}
	}

	return candidates
}
