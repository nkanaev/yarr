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

func TestRSSWithLotsOfSpaces(t *testing.T) {
	// https://pxlnv.com/: https://feedpress.me/pxlnv
	feed, err := Parse(strings.NewReader(strings.ReplaceAll(`
		<?xml version="1.0" encoding="UTF-8"?>
		<?xml-stylesheet type="text/xsl" media="screen" href="/~files/feed-premium.xsl"?>
		<lotsofspaces>
		<rss xmlns:content="http://purl.org/rss/1.0/modules/content/"
		     xmlns:wfw="http://wellformedweb.org/CommentAPI/"
			 xmlns:dc="http://purl.org/dc/elements/1.1/"
			 xmlns:atom="http://www.w3.org/2005/Atom"
			 xmlns:sy="http://purl.org/rss/1.0/modules/syndication/"
			 xmlns:slash="http://purl.org/rss/1.0/modules/slash/"
			 xmlns:feedpress="https://feed.press/xmlns"
			 xmlns:media="http://search.yahoo.com/mrss/"
			 version="2.0">
			<channel>
				<title>finally</title>
			</channel>
		</rss>
	`, "<lotsofspaces>", strings.Repeat(" ", 500))))
	if err != nil {
		t.Fatal(err)
	}
	have := feed.Title
	want := "finally"
	if have != want {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.FailNow()
	}
}

func TestRSSPodcast(t *testing.T) {
	feed, _ := Parse(strings.NewReader(`
		<?xml version="1.0" encoding="UTF-8"?>
		<rss version="2.0">
			<channel>
				<item>
					<enclosure length="100500" type="audio/x-m4a" url="http://example.com/audio.ext"/>
				</item>
			</channel>
		</rss>
	`))
	have := feed.Items[0].AudioURL
	want := "http://example.com/audio.ext"
	if want != have {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.FailNow()
	}
}

func TestRSSOpusPodcast(t *testing.T) {
	feed, _ := Parse(strings.NewReader(`
		<?xml version="1.0" encoding="UTF-8"?>
		<rss version="2.0">
			<channel>
				<item>
					<enclosure length="100500" type="audio/opus" url="http://example.com/audio.ext"/>
				</item>
			</channel>
		</rss>
	`))
	have := feed.Items[0].AudioURL
	want := "http://example.com/audio.ext"
	if want != have {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.FailNow()
	}
}

// found in: https://podcast.cscript.site/podcast.xml
func TestRSSPodcastDuplicated(t *testing.T) {
	feed, _ := Parse(strings.NewReader(`
		<?xml version="1.0" encoding="UTF-8"?>
		<rss version="2.0" xmlns:content="http://purl.org/rss/1.0/modules/content/">
			<channel>
				<item>
				<content:encoded>
					<![CDATA[ <audio src="http://example.com/audio.ext"></audio> ]]>
				</content:encoded>
					<enclosure length="100500" type="audio/x-m4a" url="http://example.com/audio.ext"/>
				</item>
			</channel>
		</rss>
	`))
	have := feed.Items[0].Content
	want := `<audio src="http://example.com/audio.ext"></audio>`
	if want != have {
		t.Fatalf("content doesn't match\nwant: %#v\nhave: %#v\n", want, have)
	}
	if feed.Items[0].AudioURL != "" {
		t.Fatal("item.audio_url must be unset if present in the content")
	}
}

func TestRSSTitleHTMLTags(t *testing.T) {
	feed, _ := Parse(strings.NewReader(`
		<?xml version="1.0" encoding="UTF-8"?>
		<rss version="2.0" xmlns:content="http://purl.org/rss/1.0/modules/content/">
			<channel>
				<item>
					<title>&lt;p&gt;title in p&lt;/p&gt;</title>
				</item>
				<item>
					<title>very &lt;strong&gt;strong&lt;/strong&gt; title</title>
				</item>
			</channel>
		</rss>
	`))
	have := []string{feed.Items[0].Title, feed.Items[1].Title}
	want := []string{"title in p", "very strong title"}
	for i := 0; i < len(want); i++ {
		if want[i] != have[i] {
			t.Errorf("title doesn't match\nwant: %#v\nhave: %#v\n", want[i], have[i])
		}
	}
}

func TestRSSIsPermalink(t *testing.T) {
	feed, _ := Parse(strings.NewReader(`
		<?xml version="1.0" encoding="UTF-8"?>
		<rss version="2.0" xmlns:content="http://purl.org/rss/1.0/modules/content/">
			<channel>
				<item>
                    <guid isPermaLink="true">http://example.com/posts/1</guid>
				</item>
			</channel>
		</rss>
	`))
	have := feed.Items
	want := []Item{
        {
            GUID:    "http://example.com/posts/1",
            URL:     "http://example.com/posts/1",
        },
    }
	for i := 0; i < len(want); i++ {
		if want[i] != have[i] {
			t.Errorf("Failed to handle isPermalink\nwant: %#v\nhave: %#v\n", want[i], have[i])
		}
	}
}
