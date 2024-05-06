package opml

import (
	"fmt"
	"html"
	"strings"
)

type Folder struct {
	Title   string
	Folders []Folder
	Feeds   []Feed
}

type Feed struct {
	Title       string
	FeedUrl     string
	SiteUrl     string
	CustomOrder string
}

func (f Folder) AllFeeds() []Feed {
	feeds := make([]Feed, 0)
	feeds = append(feeds, f.Feeds...)
	for _, subfolder := range f.Folders {
		feeds = append(feeds, subfolder.AllFeeds()...)
	}
	return feeds
}

var e = html.EscapeString
var indent = "  "
var nl = "\n"

func (f Folder) outline(level int) string {
	builder := strings.Builder{}
	prefix := strings.Repeat(indent, level)

	if level > 0 {
		builder.WriteString(prefix + fmt.Sprintf(`<outline text="%s">`+nl, e(f.Title)))
	}
	for _, folder := range f.Folders {
		builder.WriteString(folder.outline(level + 1))
	}
	for _, feed := range f.Feeds {
		builder.WriteString(feed.outline(level + 1))
	}
	if level > 0 {
		builder.WriteString(prefix + `</outline>` + nl)
	}
	return builder.String()
}

func (f Feed) outline(level int) string {
	return strings.Repeat(indent, level) + fmt.Sprintf(
		`<outline type="rss" text="%s" xmlUrl="%s" htmlUrl="%s" customOrder="%s"/>`+nl,
		e(f.Title), e(f.FeedUrl), e(f.SiteUrl), e(f.CustomOrder),
	)
}

func (f Folder) OPML() string {
	builder := strings.Builder{}
	builder.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + nl)
	builder.WriteString(`<opml version="1.1">` + nl)
	builder.WriteString(`<head><title>subscriptions</title></head>` + nl)
	builder.WriteString(`<body>` + nl)
	builder.WriteString(f.outline(0))
	builder.WriteString(`</body>` + nl)
	builder.WriteString(`</opml>` + nl)
	return builder.String()
}
