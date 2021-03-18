package worker

import (
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func getAttr(node *html.Node, key string) string {
	for _, a := range node.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

func getText(node *html.Node) string {
	text := make([]string, 0)
	isTextNode := func(n *html.Node) bool {
		return n.Type == html.TextNode
	}
	for _, n := range getNodes(node, isTextNode) {
		text = append(text, strings.TrimSpace(n.Data))
	}
	return strings.Join(text, " ")
}

func getNodes(node *html.Node, match func(*html.Node) bool) []*html.Node {
	nodes := make([]*html.Node, 0)

	queue := make([]*html.Node, 0)
	queue = append(queue, node)
	for len(queue) > 0 {
		queue, n := queue[1:], queue[0]
		if match(n) {
			nodes = append(nodes, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			queue = append(queue, c)
		}
	}
	return nodes
}

func FindFeeds(doc *html.Node, baseUrl *url.URL) []*FeedSource {
	candidates := make(map[string]string)

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
		link := baseUrl.ResolveReference(href).String()

		if href != "" {
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
				text := strings.Lower(getText(n))

				for _, feedHref := range feedHrefs {
					if strings.HasSuffix(href, feedHref) {
						return true
					}
				}
				for _, feedText := range feedTexts {
					if strings.Contains(text, feedText) {
						return true
					}
				}
			}
			return false
		}
		for _, node := range getNodes(doc, isFeedHyperLink) {
			href := getAttr(node, "href")
			link := baseUrl.ResolveReference(href).String()
			candidates[link] = ""
		}
	}

	sources := make([]*FeedSource, 0, len(candidates))
	for url, title := range candidates {
		sources = append(sources, &FeedSource{Title: title, Url: url})
	}
	return sources
}
