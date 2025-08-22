// Parser for RSS versions:
// - 0.90
// - 1.0
package parser

import (
	"encoding/xml"
	"html"
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

	DublinCoreDate string `xml:"http://purl.org/dc/elements/1.1/ date"`
	ContentEncoded string `xml:"http://purl.org/rss/1.0/modules/content/ encoded"`

	media
}

func ParseRDF(r io.Reader) (*Feed, error) {
	srcfeed := rdfFeed{}

	decoder := xmlDecoder(r)
	if err := decoder.Decode(&srcfeed); err != nil {
		return nil, err
	}

	dstfeed := &Feed{
		Title:   srcfeed.Title,
		SiteURL: srcfeed.Link,
	}
	for _, srcitem := range srcfeed.Items {
		mediaLinks := []MediaLink{}

		if isLinkPossiblyAImage(srcitem.Link) {
			mediaLinks = append(mediaLinks, MediaLink{URL: srcitem.Link, Type: "image"})
		}

		content := firstNonEmpty(srcitem.ContentEncoded, srcitem.Description)
		if contentImage := findImageInContent(html.UnescapeString(content)); contentImage != nil {
			mediaLinks = append(mediaLinks, MediaLink{URL: *contentImage, Type: "image"})
		}

		if len(mediaLinks) <= 0 {
			mediaLinks = nil
		}

		dstfeed.Items = append(dstfeed.Items, Item{
			GUID:       srcitem.Link,
			URL:        srcitem.Link,
			Date:       dateParse(srcitem.DublinCoreDate),
			Title:      srcitem.Title,
			Content:    content,
			MediaLinks: mediaLinks,
		})
	}
	return dstfeed, nil
}
