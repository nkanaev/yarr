package scraper

import (
	"maps"
	"net/url"
	"slices"
	"strings"

	"github.com/nkanaev/yarr/src/content/htmlutil"
	"golang.org/x/net/html"
)

type FeedLink struct {
	URL           string `json:"url"`
	Title         string `json:"title"`
}

func FindFeeds(body string, base string) []FeedLink {
	candidates := make(map[string]FeedLink)

	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		return nil
	}

	// find direct links
	// css: link[type=application/atom+xml]
	linkTypes := []string{"application/atom+xml", "application/rss+xml", "application/json", "application/feed+json"}
	isFeedLink := func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == "link" {
			t := htmlutil.Attr(n, "type")
			if slices.Contains(linkTypes, t) {
				return true
			}
		}
		return false
	}
	for _, node := range htmlutil.FindNodes(doc, isFeedLink) {
		href := htmlutil.Attr(node, "href")
		name := htmlutil.Attr(node, "title")
		link := htmlutil.AbsoluteUrl(href, base)
		if link != "" {
			candidates[link] = FeedLink{URL: link, Title: name}

			l, err := url.Parse(link)
			if err == nil && l.Host == "www.youtube.com" && l.Path == "/feeds/videos.xml" {
				// https://wiki.archiveteam.org/index.php/YouTube/Technical_details#Playlists
				channelID, found := strings.CutPrefix(l.Query().Get("channel_id"), "UC")
				if found {
					const baseURL string = "https://www.youtube.com/feeds/videos.xml?playlist_id="

					ytTitle := name
					isOG := func(n *html.Node) bool {
						return n.Type == html.ElementNode && n.Data == "meta" &&
							htmlutil.Attr(n, "property") == "og:title"
					}
					for _, n := range htmlutil.FindNodes(doc, isOG) {
						ytTitle = htmlutil.Attr(n, "content")
						break
					}

					candidates[link] = FeedLink{
						URL:   link,
						Title: ytTitle + " - All",
					}
					candidates[baseURL+"UULF"+channelID] = FeedLink{
						URL:           baseURL + "UULF" + channelID,
						Title:         ytTitle + " - Videos",
					}
					candidates[baseURL+"UULV"+channelID] = FeedLink{
						URL:           baseURL + "UULV" + channelID,
						Title:         ytTitle + " - Live Streams",
					}
					candidates[baseURL+"UUSH"+channelID] = FeedLink{
						URL:           baseURL + "UUSH" + channelID,
						Title:         ytTitle + " - Short videos",
					}
				}
			}
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
				candidates[link] = FeedLink{URL: link, Title: ""}
			}
		}
	}

	result := slices.Collect(maps.Values(candidates))
	slices.SortFunc(result, func(a, b FeedLink) int {
		if a.Title != b.Title {
			return strings.Compare(a.Title, b.Title)
		}
		return strings.Compare(a.URL, b.URL)
	})
	return result
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
		rels := strings.SplitSeq(htmlutil.Attr(node, "rel"), " ")
		for rel := range rels {
			if strings.EqualFold(rel, "icon") {
				icons = append(icons, htmlutil.AbsoluteUrl(htmlutil.Attr(node, "href"), base))
			}
		}
	}
	return icons
}
