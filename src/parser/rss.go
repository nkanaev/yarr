// Parser for RSS versions:
// - 0.91 netscape
// - 0.91 userland
// - 2.0
package parser

import (
	"encoding/xml"
	"io"
)

type rssFeed struct {
	XMLName xml.Name  `xml:"rss"`
	Version string    `xml:"version,attr"`
	Title   string    `xml:"channel>title"`
	Link    string    `xml:"channel>link"`
	Items   []rssItem `xml:"channel>item"`
}

type rssItem struct {
	GUID        string         `xml:"guid"`
	Title       string         `xml:"title"`
	Link        string         `xml:"link"`
	Description string         `xml:"rss description"`
	PubDate     string         `xml:"pubDate"`
	Enclosures  []rssEnclosure `xml:"enclosure"`

	DublinCoreDate string `xml:"http://purl.org/dc/elements/1.1/ date"`
	ContentEncoded string `xml:"http://purl.org/rss/1.0/modules/content/ encoded"`

	FeedBurnerLink          string `xml:"http://rssnamespace.org/feedburner/ext/1.0 origLink"`
	FeedBurnerEnclosureLink string `xml:"http://rssnamespace.org/feedburner/ext/1.0 origEnclosureLink"`

	ItunesSubtitle    string `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd subtitle"`
	ItunesSummary     string `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd summary"`
	GoogleDescription string `xml:"http://www.google.com/schemas/play-podcasts/1.0 description"`
	media
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
		podcastURL := ""
		for _, e := range srcitem.Enclosures {
			if e.Type == "audio/mpeg" || e.Type == "audio/x-m4a" {
				podcastURL = e.URL
				break
			}
		}

		dstfeed.Items = append(dstfeed.Items, Item{
			GUID:     firstNonEmpty(srcitem.GUID, srcitem.Link),
			Date:     dateParse(firstNonEmpty(srcitem.DublinCoreDate, srcitem.PubDate)),
			URL:      srcitem.Link,
			Title:    srcitem.Title,
			Content:  firstNonEmpty(srcitem.ContentEncoded, srcitem.Description),
			AudioURL: podcastURL,
			ImageURL: srcitem.firstMediaThumbnail(),
		})
	}
	return dstfeed, nil
}
