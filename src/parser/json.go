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
	Authors       []jsonAuthor     `json:"authors"`
}

type jsonAttachment struct {
	URL      string `json:"url"`
	MimeType string `json:"mime_type"`
	Title    string `json:"title"`
	Size     int64  `json:"size_in_bytes"`
	Duration int    `json:"duration_in_seconds"`
}

type jsonAuthor struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Avatar string `json:"avatar"`
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
		authors := []string{}
		for _, v := range srcitem.Authors {
			authors = append(authors, v.Name)
		}
		dstfeed.Items = append(dstfeed.Items, Item{
			GUID:    firstNonEmpty(srcitem.ID, srcitem.URL),
			Date:    dateParse(firstNonEmpty(srcitem.DatePublished, srcitem.DateModified)),
			URL:     srcitem.URL,
			Title:   srcitem.Title,
			Content: firstNonEmpty(srcitem.HTML, srcitem.Text, srcitem.Summary),
			Author:  strings.Join(authors, ","),
		})
	}
	return dstfeed, nil
}
