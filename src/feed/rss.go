package feed

import (
	"encoding/xml"
)

type rssFeed struct {
	XMLName     xml.Name  `xml:"rss"`
	Version     string    `xml:"version,attr"`
	Title       string    `xml:"channel>title"`
	Links       []rssLink `xml:"channel>link"`
	Language    string    `xml:"channel>language"`
	Description string    `xml:"channel>description"`
	PubDate     string    `xml:"channel>pubDate"`
	Items       []rssItem `xml:"channel>item"`
}

type rssItem struct {
	GUID           string         `xml:"guid"`
	Title          []rssTitle     `xml:"title"`
	Links          []rssLink      `xml:"link"`
	Description    string         `xml:"description"`
	PubDate        string         `xml:"pubDate"`
	EnclosureLinks []rssEnclosure `xml:"enclosure"`

	DublinCoreDate    string `xml:"http://purl.org/dc/elements/1.1/ date"`
	DublinCoreContent string `xml:"http://purl.org/rss/1.0/modules/content/ encoded"`

	FeedBurnerLink          string `xml:"http://rssnamespace.org/feedburner/ext/1.0 origLink"`
	FeedBurnerEnclosureLink string `xml:"http://rssnamespace.org/feedburner/ext/1.0 origEnclosureLink"`

	ItunesSubtitle    string `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd subtitle"`
	ItunesSummary     string `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd summary"`
	GoogleDescription string `xml:"http://www.google.com/schemas/play-podcasts/1.0 description"`
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
