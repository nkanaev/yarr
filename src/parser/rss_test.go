package parser

import (
	"reflect"
	"strings"
	"testing"
)

func TestRSSFeed(t *testing.T) {
	have, _ := Parse(strings.NewReader(`
		<?xml version="1.0"?>
		<!DOCTYPE rss SYSTEM "http://my.netscape.com/publish/formats/rss-0.91.dtd">
		<rss version="0.91">
		<channel>
			<language>en</language>
			<description>???</description>
			<link>http://www.scripting.com/</link>
			<title>Scripting News</title>
			<item>
				<title>Title 1</title>
				<link>http://www.scripting.com/one/</link>
				<description>Description 1</description>
			</item>
			<item>
				<title>Title 2</title>
				<link>http://www.scripting.com/two/</link>
				<description>Description 2</description>
			</item>
		</channel>
		</rss>
	`))
	want := &Feed{
		Title:   "Scripting News",
		SiteURL: "http://www.scripting.com/",
		Items: []Item{
			{
				GUID:    "http://www.scripting.com/one/",
				URL:     "http://www.scripting.com/one/",
				Title:   "Title 1",
				Content: "Description 1",
			},
			{
				GUID:    "http://www.scripting.com/two/",
				URL:     "http://www.scripting.com/two/",
				Title:   "Title 2",
				Content: "Description 2",
			},
		},
	}

	if !reflect.DeepEqual(want, have) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fatal("invalid rss")
	}
}

func TestRSSMediaContentThumbnail(t *testing.T) {
	// see: https://vimeo.com/channels/staffpicks/videos/rss
	feed, _ := Parse(strings.NewReader(`
		<?xml version="1.0" encoding="UTF-8"?>
		<rss version="2.0"
			xmlns:atom="http://www.w3.org/2005/Atom"
			xmlns:media="http://search.yahoo.com/mrss/" xml:lang="en-US">
			<channel>
				<item>
					<title></title>
					<media:content>
						<media:player url="https://player.vimeo.com/video/527877676"/>
						<media:credit role="author" scheme="urn:ebu"></media:credit>
						<media:thumbnail height="540" width="960" url="https://i.vimeocdn.com/video/1092705247_960.jpg"/>
						<media:title></media:title>
					</media:content>
				</item>
			</channel>
		</rss>
	`))
	have := feed.Items[0].ImageURL
	want := "https://i.vimeocdn.com/video/1092705247_960.jpg"
	if have != want {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.FailNow()
	}
}
