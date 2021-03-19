package feed

import "time"

type Feed struct {
	Title   string
	SiteURL string
	FeedURL string
	Items   []Item
}

type Item struct {
	GUID string
	Date time.Time
	URL  string
	Title string

	Content string
	ImageURL string
	PodcastURL string
}
