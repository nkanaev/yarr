package feed

import (
	"encoding/json"
	"io"
)

type jsonFeed struct {
	Version string     `json:"version"`
	Title   string     `json:"title"`
	SiteURL string     `json:"home_page_url"`
	FeedURL string     `json:"feed_url"`
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

func first(vals ...string) string {
	for _, val := range vals {
		if len(val) > 0 {
			return val
		}
	}
	return ""
}

func (f *jsonFeed) convert() *Feed {
	feed := &Feed{
		Title: f.Title,
		SiteURL: f.SiteURL,
		FeedURL: f.FeedURL,
	}
	for _, item := range f.Items {
		date, _ := dateParse(first(item.DatePublished, item.DateModified))
		content := first(item.HTML, item.Text, item.Summary)
		imageUrl := ""
		podcastUrl := ""
	
		feed.Items = append(feed.Items, Item{
			GUID: item.ID,
			Date: date,
			URL:  item.URL,
			Title: item.Title,
			Content: content,
			ImageURL: imageUrl,
			PodcastURL: podcastUrl,
		})
	}
	return feed
}

func ParseJSON(data io.Reader) (*Feed, error) {
	feed := new(jsonFeed)
	decoder := json.NewDecoder(data)
	if err := decoder.Decode(&feed); err != nil {
		return nil, err
	}
	return feed.convert(), nil
}
