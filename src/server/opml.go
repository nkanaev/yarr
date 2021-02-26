package server

import (
	"encoding/xml"
	"io"
)

type opml struct {
	XMLName  xml.Name  `xml:"opml"`
	Version  string    `xml:"version,attr"`
	Outlines []outline `xml:"body>outline"`
}

type outline struct {
	Type        string    `xml:"type,attr,omitempty"`
	Title       string    `xml:"text,attr"`
	FeedURL     string    `xml:"xmlUrl,attr,omitempty"`
	SiteURL     string    `xml:"htmlUrl,attr,omitempty"`
	Description string    `xml:"description,attr,omitempty"`
	Outlines    []outline `xml:"outline,omitempty"`
}

func (o outline) AllFeeds() []outline {
	result := make([]outline, 0)
	for _, sub := range o.Outlines {
		if sub.Type == "rss" {
			result = append(result, sub)
		} else {
			result = append(result, sub.AllFeeds()...)
		}
	}
	return result
}

func parseOPML(r io.Reader) (*opml, error) {
	feeds := new(opml)
	decoder := xml.NewDecoder(r)
	decoder.Entity = xml.HTMLEntity
	decoder.Strict = false
	err := decoder.Decode(&feeds)
	return feeds, err
}
