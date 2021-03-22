package feed

import (
	"encoding/xml"
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
