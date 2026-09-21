// JSON 1.0 parser
package parser

import (
	"encoding/json"
	"io"
	"strings"
)

type jsonFeed struct {
	Version string     `json:"version"`
	Title   string     `json:"title"`
	SiteURL string     `json:"home_page_url"`
	Items   []jsonItem `json:"items"`
}

type jsonItem struct {
	ID            string           `json:"id"`
	URL           string           `json:"url"`
	Title         string           `json:"title"`
	Summary       string           `json:"summary"`
	Text          string           `json:"content_text"`
	HTML          string           `json:"content_html"`
	DatePublished string           `json:"date_published"`
	DateModified  string           `json:"date_modified"`
	Attachments   []jsonAttachment `json:"attachments"`
}

type jsonAttachment struct {
	URL      string `json:"url"`
	MimeType string `json:"mime_type"`
	Title    string `json:"title"`
	Size     int64  `json:"size_in_bytes"`
	Duration int    `json:"duration_in_seconds"`
}

func (item *jsonItem) mediaLinks() []MediaLink {
	links := make([]MediaLink, 0)
	for _, a := range item.Attachments {
		if a.URL == "" {
			continue
		}
		var typ string
		switch {
		case strings.HasPrefix(a.MimeType, "image/"):
			typ = "image"
		case strings.HasPrefix(a.MimeType, "audio/"):
			typ = "audio"
		case strings.HasPrefix(a.MimeType, "video/"):
			typ = "video"
		default:
			continue
		}
		links = append(links, MediaLink{URL: a.URL, Type: typ, Description: a.Title})
	}
	if len(links) == 0 {
		return nil
	}
	return links
}

func ParseJSON(data io.Reader) (*Feed, error) {
	srcfeed := new(jsonFeed)
	decoder := json.NewDecoder(data)
	if err := decoder.Decode(&srcfeed); err != nil {
		return nil, err
	}

	dstfeed := &Feed{
		Title:   srcfeed.Title,
		SiteURL: srcfeed.SiteURL,
	}
	for _, srcitem := range srcfeed.Items {
		dstfeed.Items = append(dstfeed.Items, Item{
			GUID:       firstNonEmpty(srcitem.ID, srcitem.URL),
			Date:       dateParse(firstNonEmpty(srcitem.DatePublished, srcitem.DateModified)),
			URL:        srcitem.URL,
			Title:      srcitem.Title,
			Content:    firstNonEmpty(srcitem.HTML, srcitem.Text, srcitem.Summary),
			MediaLinks: srcitem.mediaLinks(),
		})
	}
	return dstfeed, nil
}
