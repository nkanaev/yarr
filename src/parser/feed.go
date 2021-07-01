package parser

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"
)

var UnknownFormat = errors.New("unknown feed format")

type processor func(r io.Reader) (*Feed, error)

func sniff(lookup string) (string, processor) {
	lookup = strings.TrimSpace(lookup)
	lookup = strings.TrimLeft(lookup, "\x00\xEF\xBB\xBF\xFE\xFF")

	if len(lookup) < 0 {
		return "", nil
	}

	switch lookup[0] {
	case '<':
		decoder := xmlDecoder(strings.NewReader(lookup))
		for {
			token, _ := decoder.Token()
			if token == nil {
				break
			}
			if el, ok := token.(xml.StartElement); ok {
				switch el.Name.Local {
				case "rss":
					return "rss", ParseRSS
				case "RDF":
					return "rdf", ParseRDF
				case "feed":
					return "atom", ParseAtom
				}
			}
		}
	case '{':
		return "json", ParseJSON
	}
	return "", nil
}

func Parse(r io.Reader) (*Feed, error) {
	lookup := make([]byte, 2048)
	n, err := io.ReadFull(r, lookup)
	switch {
	case err == io.ErrUnexpectedEOF:
		lookup = lookup[:n]
		r = bytes.NewReader(lookup)
	case err != nil:
		return nil, err
	default:
		r = io.MultiReader(bytes.NewReader(lookup), r)
	}

	_, callback := sniff(string(lookup))
	if callback == nil {
		return nil, UnknownFormat
	}

	feed, err := callback(r)
	if feed != nil {
		feed.cleanup()
	}
	return feed, err
}

func (feed *Feed) cleanup() {
	feed.Title = strings.TrimSpace(feed.Title)
	feed.SiteURL = strings.TrimSpace(feed.SiteURL)

	for i, item := range feed.Items {
		feed.Items[i].GUID = strings.TrimSpace(item.GUID)
		feed.Items[i].URL = strings.TrimSpace(item.URL)
		feed.Items[i].Title = strings.TrimSpace(item.Title)
		feed.Items[i].Content = strings.TrimSpace(item.Content)

		if item.ImageURL != "" && strings.Contains(item.Content, item.ImageURL) {
			feed.Items[i].ImageURL = ""
		}
		if item.AudioURL != "" && strings.Contains(item.Content, item.AudioURL) {
			feed.Items[i].AudioURL = ""
		}
	}
}

func (feed *Feed) SetMissingDatesTo(newdate time.Time) {
	for i, item := range feed.Items {
		if item.Date.IsZero() {
			feed.Items[i].Date = newdate
		}
	}
}

func (feed *Feed) TranslateURLs(base string) error {
	baseUrl, err := url.Parse(base)
	if err != nil {
		return fmt.Errorf("failed to parse base url: %#v", base)
	}
	siteUrl, err := url.Parse(feed.SiteURL)
	if err != nil {
		return fmt.Errorf("failed to parse feed url: %#v", feed.SiteURL)
	}
	feed.SiteURL = baseUrl.ResolveReference(siteUrl).String()
	for _, item := range feed.Items {
		itemUrl, err := url.Parse(item.URL)
		if err != nil {
			return fmt.Errorf("failed to parse item url: %#v", item.URL)
		}
		item.URL = siteUrl.ResolveReference(itemUrl).String()
	}
	return nil
}
