package opml

import (
	"reflect"
	"strings"
	"testing"
)


func TestParse(t *testing.T) {
	have, _ := Parse(strings.NewReader(`
		<?xml version="1.0" encoding="UTF-8"?>
		<opml version="1.1">
		<head><title>Subscriptions</title></head>
		<body>
			<outline text="sub">
				<outline type="rss" text="subtitle1" description="sub1"
						 xmlUrl="https://foo.com/feed.xml" htmlUrl="https://foo.com/"/>
				<outline type="rss" text="&amp;&gt;" description="&lt;&gt;"
						 xmlUrl="https://bar.com/feed.xml" htmlUrl="https://bar.com/"/>
			</outline>
			<outline type="rss" text="title1" description="desc1"
					 xmlUrl="https://baz.com/feed.xml" htmlUrl="https://baz.com/"/>
		</body>
		</opml>
	`))
	want := Folder{
		Title: "",
		Feeds: []Feed{
			Feed{
				Title: "title1",
				FeedUrl: "https://baz.com/feed.xml",
				SiteUrl: "https://baz.com/",
			},
		},
		Folders: []Folder{
			Folder{
				Title: "sub",
				Feeds: []Feed{
					Feed{
						Title: "subtitle1",
						FeedUrl: "https://foo.com/feed.xml",
						SiteUrl: "https://foo.com/",
					},
					Feed{
						Title: "&>",
						FeedUrl: "https://bar.com/feed.xml",
						SiteUrl: "https://bar.com/",
					},
				},
			},
		},
	}	
	if !reflect.DeepEqual(want, have) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fatal("invalid opml")
	}
}
