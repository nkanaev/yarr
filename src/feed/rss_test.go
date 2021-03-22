package feed

import (
	"reflect"
	"strings"
	"testing"
)

func TestRSSFeed(t *testing.T) {
	have, _ := ParseRSS(strings.NewReader(`
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
		Title: "Scripting News",
		SiteURL: "http://www.scripting.com/",
		Items: []Item{
			{
				GUID: "http://www.scripting.com/one/",
				URL: "http://www.scripting.com/one/",
				Title: "Title 1",
				Content: "Description 1",
			},
			{
				GUID: "http://www.scripting.com/two/",
				URL: "http://www.scripting.com/two/",
				Title: "Title 2",
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
