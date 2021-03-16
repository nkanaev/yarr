package opml

import (
	"fmt"
	"html"
	"strings"
)

type OPMLBuilder struct {
	rootfolder *OPMLFolder
	folders    []*OPMLFolder
}

type OPMLFolder struct {
	title string
	feeds []*OPMLFeed
}

type OPMLFeed struct {
	title, description, feedUrl, siteUrl string
}

func NewBuilder() *OPMLBuilder {
	return &OPMLBuilder{
		rootfolder: &OPMLFolder{feeds: make([]*OPMLFeed, 0)},
		folders:    make([]*OPMLFolder, 0),
	}	
}

func (b *OPMLBuilder) AddFeed(title, description, feedUrl, siteUrl string) {
	b.rootfolder.AddFeed(title, description, feedUrl, siteUrl)
}

func (b *OPMLBuilder) AddFolder(title string) *OPMLFolder {
	folder := &OPMLFolder{title: title}
	b.folders = append(b.folders, folder)
	return folder
}

func (f *OPMLFolder) AddFeed(title, description, feedUrl, siteUrl string) {
	f.feeds = append(f.feeds, &OPMLFeed{title, description, feedUrl, siteUrl})
}

func (b *OPMLBuilder) String() string {
	builder := strings.Builder{}

	line := func(s string, args ...string) {
		if len(args) > 0 {
			escapedargs := make([]interface{}, len(args))
			for idx, arg := range args {
				escapedargs[idx] = html.EscapeString(arg)
			}
			s = fmt.Sprintf(s, escapedargs...)
		}
		builder.WriteString(s)
		builder.WriteRune('\n')
	}
	feedline := func(feed *OPMLFeed, indent int) {
		line(
			strings.Repeat(" ", indent) + `<outline type="rss" text="%s" description="%s" xmlUrl="%s" htmlUrl="%s"/>`,
			feed.title, feed.description,
			feed.feedUrl, feed.siteUrl,
		)
	}
	line(`<?xml version="1.0" encoding="UTF-8"?>`)
	line(`<opml version="1.1">`)
	line(`<head><title>Subscriptions</title></head>`)
	line(`<body>`)
	for _, folder := range b.folders {
		line(`  <outline text="%s">`, folder.title)
		for _, feed := range folder.feeds {
			feedline(feed, 4)
		}
		line(`  </outline>`)
	}
	for _, feed := range b.rootfolder.feeds {
		feedline(feed, 2)
	}
	line(`</body>`)
	line(`</opml>`)

	return builder.String()
}
