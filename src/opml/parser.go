package opml

import (
	"encoding/xml"
	"io"
)

type OPML struct {
	XMLName  xml.Name  `xml:"opml"`
	Version  string    `xml:"version,attr"`
	Outlines []Outline `xml:"body>outline"`
}

type Outline struct {
	Type        string    `xml:"type,attr,omitempty"`
	Title       string    `xml:"text,attr"`
	FeedURL     string    `xml:"xmlUrl,attr,omitempty"`
	SiteURL     string    `xml:"htmlUrl,attr,omitempty"`
	Description string    `xml:"description,attr,omitempty"`
	Outlines    []Outline `xml:"outline,omitempty"`
}

func (o Outline) AllFeeds() []Outline {
	result := make([]Outline, 0)
	for _, sub := range o.Outlines {
		if sub.Type == "rss" {
			result = append(result, sub)
		} else {
			result = append(result, sub.AllFeeds()...)
		}
	}
	return result
}

func Parse(r io.Reader) (*OPML, error) {
	feeds := new(OPML)
	decoder := xml.NewDecoder(r)
	decoder.Entity = xml.HTMLEntity
	decoder.Strict = false
	err := decoder.Decode(&feeds)
	return feeds, err
}
