package scraper

import (
	"strings"

	"github.com/nkanaev/yarr/src/content/htmlutil"
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
			t := htmlutil.Attr(n, "type")
			for _, tt := range linkTypes {
				if tt == t {
					return true
				}
			}
		}
		return false
	}
	for _, node := range htmlutil.FindNodes(doc, isFeedLink) {
		href := htmlutil.Attr(node, "href")
		name := htmlutil.Attr(node, "title")
		link := htmlutil.AbsoluteUrl(href, base)
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
				href := strings.Trim(htmlutil.Attr(n, "href"), "/")
				for _, feedHref := range feedHrefs {
					if strings.HasSuffix(href, feedHref) {
						return true
					}
				}
				text := htmlutil.Text(n)
				for _, feedText := range feedTexts {
					if strings.EqualFold(text, feedText) {
						return true
					}
				}
			}
			return false
		}
		for _, node := range htmlutil.FindNodes(doc, isFeedHyperLink) {
			href := htmlutil.Attr(node, "href")
			link := htmlutil.AbsoluteUrl(href, base)
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
	for _, node := range htmlutil.FindNodes(doc, isLink) {
		rels := strings.Split(htmlutil.Attr(node, "rel"), " ")
		for _, rel := range rels {
			if strings.EqualFold(rel, "icon") {
				icons = append(icons, htmlutil.AbsoluteUrl(htmlutil.Attr(node, "href"), base))
			}
		}
	}
	return icons
}
