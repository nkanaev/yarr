package parser

import "time"

type Feed struct {
	Title   string
	SiteURL string
	Items   []Item
}

type Item struct {
	GUID  string
	Date  time.Time
	URL   string
	Title string

	Content    string
	MediaLinks []MediaLink
}

type MediaLink struct {
	URL         string
	Type        string
	Description string
}
