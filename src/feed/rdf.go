package feed

import (
	"encoding/xml"
	"io"
)

type rdfFeed struct {
	XMLName xml.Name  `xml:"RDF"`
	Title   string    `xml:"channel>title"`
	Link    string    `xml:"channel>link"`
	Items   []rdfItem `xml:"item"`
}

type rdfItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`

	DublinCoreDate    string `xml:"http://purl.org/dc/elements/1.1/ date"`
	DublinCoreContent string `xml:"http://purl.org/rss/1.0/modules/content/ encoded"`
}

func ParseRDF(r io.Reader) (*Feed, error) {
	f := rdfFeed{}

	decoder := xml.NewDecoder(r)
	if err := decoder.Decode(&f); err != nil {
		return nil, err
	}

	feed := &Feed{
		Title:   f.Title,
		SiteURL: f.Link,
	}
	for _, e := range f.Items {
		feed.Items = append(feed.Items, Item{
			GUID:  e.Link,
			URL:   e.Link,
			Title: e.Title,
		})
	}
	return feed, nil
}
