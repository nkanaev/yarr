// Parser for RSS versions:
// - 0.91 netscape
// - 0.91 userland
// - 2.0
package parser

import (
	"encoding/xml"
	"html"
	"io"
	"path"
	"strings"
)

type rssFeed struct {
	XMLName xml.Name  `xml:"rss"`
	Version string    `xml:"version,attr"`
	Title   string    `xml:"channel>title"`
	Link    string    `xml:"channel>link"`
	Items   []rssItem `xml:"channel>item"`
}

type rssItem struct {
	GUID        rssGuid        `xml:"guid"`
	Title       string         `xml:"title"`
	Link        string         `xml:"rss link"`
	Description string         `xml:"rss description"`
	PubDate     string         `xml:"pubDate"`
	Enclosures  []rssEnclosure `xml:"enclosure"`

	DublinCoreDate string `xml:"http://purl.org/dc/elements/1.1/ date"`
	ContentEncoded string `xml:"http://purl.org/rss/1.0/modules/content/ encoded"`

	OrigLink          string `xml:"http://rssnamespace.org/feedburner/ext/1.0 origLink"`
	OrigEnclosureLink string `xml:"http://rssnamespace.org/feedburner/ext/1.0 origEnclosureLink"`

	media
}

type rssGuid struct {
	GUID        string `xml:",chardata"`
	IsPermaLink string `xml:"isPermaLink,attr"`
}

type rssLink struct {
	XMLName xml.Name
	Data    string `xml:",chardata"`
	Href    string `xml:"href,attr"`
	Rel     string `xml:"rel,attr"`
}

type rssTitle struct {
	XMLName xml.Name
	Data    string `xml:",chardata"`
	Inner   string `xml:",innerxml"`
}

type rssEnclosure struct {
	URL    string `xml:"url,attr"`
	Type   string `xml:"type,attr"`
	Length string `xml:"length,attr"`
}

func ParseRSS(r io.Reader) (*Feed, error) {
	srcfeed := rssFeed{}

	decoder := xmlDecoder(r)
	decoder.DefaultSpace = "rss"
	if err := decoder.Decode(&srcfeed); err != nil {
		return nil, err
	}

	dstfeed := &Feed{
		Title:   srcfeed.Title,
		SiteURL: srcfeed.Link,
	}
	for _, srcitem := range srcfeed.Items {
		mediaLinks := srcitem.mediaLinks()
		for _, e := range srcitem.Enclosures {
			if strings.HasPrefix(e.Type, "audio/") {
				podcastURL := e.URL
				if srcitem.OrigEnclosureLink != "" && strings.Contains(podcastURL, path.Base(srcitem.OrigEnclosureLink)) {
					podcastURL = srcitem.OrigEnclosureLink
				}
				mediaLinks = append(mediaLinks, MediaLink{URL: podcastURL, Type: "audio"})
				break
			}
		}

		for _, e := range srcitem.Enclosures {
			if strings.HasPrefix(e.Type, "image/") {
				mediaLinks = append(mediaLinks, MediaLink{URL: e.URL, Type: "image"})
			}
		}

		if isLinkPossiblyAImage(srcitem.Link) {
			mediaLinks = append(mediaLinks, MediaLink{URL: srcitem.Link, Type: "image"})
		}

		content := firstNonEmpty(srcitem.ContentEncoded, srcitem.Description, srcitem.firstMediaDescription())
		if contentImage := findImageInContent(html.UnescapeString(content)); contentImage != nil {
			mediaLinks = append(mediaLinks, MediaLink{URL: *contentImage, Type: "image"})
		}

		permalink := ""
		if srcitem.GUID.IsPermaLink == "true" {
			permalink = srcitem.GUID.GUID
		}

		if len(mediaLinks) <= 0 {
			mediaLinks = nil
		}

		dstfeed.Items = append(dstfeed.Items, Item{
			GUID:       firstNonEmpty(srcitem.GUID.GUID, srcitem.Link),
			Date:       dateParse(firstNonEmpty(srcitem.DublinCoreDate, srcitem.PubDate)),
			URL:        firstNonEmpty(srcitem.OrigLink, srcitem.Link, permalink),
			Title:      srcitem.Title,
			Content:    content,
			MediaLinks: mediaLinks,
		})
	}
	return dstfeed, nil
}
