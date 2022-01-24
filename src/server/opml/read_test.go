package opml

import (
	"os"
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
			{
				Title:   "title1",
				FeedUrl: "https://baz.com/feed.xml",
				SiteUrl: "https://baz.com/",
			},
		},
		Folders: []Folder{
			{
				Title: "sub",
				Feeds: []Feed{
					{
						Title:   "subtitle1",
						FeedUrl: "https://foo.com/feed.xml",
						SiteUrl: "https://foo.com/",
					},
					{
						Title:   "&>",
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

func TestParseFallback(t *testing.T) {
	// as reported in https://github.com/nkanaev/yarr/pull/56
	// the feed below comes without `outline[text]` & `outline[type=rss]` attributes
	have, _ := Parse(strings.NewReader(`
		<?xml version="1.0" encoding="utf-8"?>
		<opml xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" version="1.0">
		  <head>
			<title>Newsflow</title>
		  </head>
		  <body>
			<outline title="foldertitle">
				<outline htmlUrl="https://example.com" text="feedtext" title="feedtitle" xmlUrl="https://example.com/feed.xml" />
			</outline>
		  </body>
		</opml>
	`))
	want := Folder{
		Folders: []Folder{{
			Title: "foldertitle",
			Feeds: []Feed{
				{Title: "feedtext", FeedUrl: "https://example.com/feed.xml", SiteUrl: "https://example.com"},
			},
		}},
	}
	if !reflect.DeepEqual(want, have) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fatal("invalid opml")
	}
}

func TestParseWithEncoding(t *testing.T) {
	file, err := os.Open("sample_win1251.xml")
	if err != nil {
		t.Fatal(err)
	}
	have, err := Parse(file)
	if err != nil {
		t.Fatal(err)
	}
	want := Folder{
		Title: "",
		Feeds: []Feed{
			{
				Title:   "пример1",
				FeedUrl: "https://baz.com/feed.xml",
				SiteUrl: "https://baz.com/",
			},
		},
		Folders: []Folder{
			{
				Title: "папка",
				Feeds: []Feed{
					{
						Title:   "пример2",
						FeedUrl: "https://foo.com/feed.xml",
						SiteUrl: "https://foo.com/",
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
