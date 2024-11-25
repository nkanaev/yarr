package parser

import (
	"bytes"
	"crypto/sha256"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/nkanaev/yarr/src/content/htmlutil"
	"golang.org/x/net/html/charset"
)

var UnknownFormat = errors.New("unknown feed format")

type feedProbe struct {
	feedType string
	callback func(r io.Reader) (*Feed, error)
	encoding string
}

func sniff(lookup string) (out feedProbe) {
	lookup = strings.TrimSpace(lookup)
	lookup = strings.TrimLeft(lookup, "\x00\xEF\xBB\xBF\xFE\xFF")

	if len(lookup) == 0 {
		return
	}

	switch lookup[0] {
	case '<':
		decoder := xmlDecoder(strings.NewReader(lookup))
		for {
			token, _ := decoder.Token()
			if token == nil {
				break
			}

			// check <?xml encoding="ENCODING" ?>
			if el, ok := token.(xml.ProcInst); ok && el.Target == "xml" {
				out.encoding = strings.ToLower(procInst("encoding", string(el.Inst)))
			}

			if el, ok := token.(xml.StartElement); ok {
				switch el.Name.Local {
				case "rss":
					out.feedType = "rss"
					out.callback = ParseRSS
					return
				case "RDF":
					out.feedType = "rdf"
					out.callback = ParseRDF
					return
				case "feed":
					out.feedType = "atom"
					out.callback = ParseAtom
					return
				}
			}
		}
	case '{':
		out.feedType = "json"
		out.callback = ParseJSON
		return
	}
	return
}

func Parse(r io.Reader) (*Feed, error) {
	return ParseWithEncoding(r, "")
}

func ParseWithEncoding(r io.Reader, fallbackEncoding string) (*Feed, error) {
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

	out := sniff(string(lookup))
	if out.feedType == "" {
		return nil, UnknownFormat
	}

	if out.encoding == "" && fallbackEncoding != "" {
		r, err = charset.NewReaderLabel(fallbackEncoding, r)
		if err != nil {
			return nil, err
		}
	}

	if (out.feedType != "json") && (out.encoding == "" || out.encoding == "utf-8") {
		// XML decoder will not rely on custom CharsetReader (see `xmlDecoder`)
		// to handle invalid xml characters.
		// Assume input is already UTF-8 and do the cleanup here.
		r = NewSafeXMLReader(r)
	}

	feed, err := out.callback(r)
	if feed != nil {
		feed.cleanup()
	}
	return feed, err
}

func ParseAndFix(r io.Reader, baseURL, fallbackEncoding string) (*Feed, error) {
	feed, err := ParseWithEncoding(r, fallbackEncoding)
	if err != nil {
		return nil, err
	}
	feed.TranslateURLs(baseURL)
	feed.SetMissingDatesTo(time.Now())
	feed.SetMissingGUIDs()
	return feed, nil
}

func (feed *Feed) cleanup() {
	feed.Title = strings.TrimSpace(feed.Title)
	feed.SiteURL = strings.TrimSpace(feed.SiteURL)

	for i, item := range feed.Items {
		feed.Items[i].GUID = strings.TrimSpace(item.GUID)
		feed.Items[i].URL = strings.TrimSpace(item.URL)
		feed.Items[i].Title = strings.TrimSpace(htmlutil.ExtractText(item.Title))
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

func (feed *Feed) SetMissingGUIDs() {
	for i, item := range feed.Items {
		if item.GUID == "" {
			id := strings.Join([]string{item.Title, item.Date.Format(time.RFC3339), item.URL}, ";;")
			feed.Items[i].GUID = fmt.Sprintf("%x", sha256.Sum256([]byte(id)))
		}
	}
}
