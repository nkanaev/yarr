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

	// find direct links
	// css: link[type=application/atom+xml]
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

	// guess by hyperlink properties
	if len(candidates) == 0 {
		// css: a[href="feed"]
		// css: a:contains("rss")
		feedHrefs := []string{"feed", "feed.xml", "rss.xml", "atom.xml"}
		feedTexts := []string{"rss", "feed"}
		isFeedHyperLink := func(n *html.Node) bool {
			if n.Type == html.ElementNode && n.Data == "a" {
				if any(feedHrefs, strings.Trim(getAttr(n, "href"), "/"), strings.HasSuffix) {
					return true
				}
				if any(feedTexts, getText(n), strings.EqualFold) {
					return true
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

func FindIcons(body string, base string) []string {
	icons := make([]string, 0)

	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		return icons
	}

	// css: link[rel=icon]
	isLink := func(n *html.Node) bool {
		return n.Type == html.ElementNode && n.Data == "link"
	}
	for _, node := range getNodes(doc, isLink) {
		if any(strings.Split(getAttr(node, "rel"), " "), "icon", strings.EqualFold) {
			icons = append(icons, absoluteUrl(getAttr(node, "href"), base))
		}
	}
	return icons
}
