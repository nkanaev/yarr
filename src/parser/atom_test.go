package parser

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestAtom(t *testing.T) {
	have, _ := Parse(strings.NewReader(`
		<?xml version="1.0" encoding="utf-8"?>
		<feed xmlns="http://www.w3.org/2005/Atom">
			<title>Example Feed</title>
			<subtitle>A subtitle.</subtitle>
			<link href="http://example.org/feed/" rel="self" />
			<link href="http://example.org/" />
			<id>urn:uuid:60a76c80-d399-11d9-b91C-0003939e0af6</id>
			<updated>2003-12-13T18:30:02Z</updated>
			<entry>
				<title>Atom-Powered Robots Run Amok</title>
				<link href="http://example.org/2003/12/13/atom03" />
				<link rel="alternate" type="text/html" href="http://example.org/2003/12/13/atom03.html"/>
				<link rel="edit" href="http://example.org/2003/12/13/atom03/edit"/>
				<id>urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a</id>
				<updated>2003-12-13T18:30:02Z</updated>
				<summary>Some text.</summary>
				<content type="xhtml">
					<div xmlns="http://www.w3.org/1999/xhtml"><p>This is the entry content.</p></div>
				</content>
				<author>
					<name>John Doe</name>
					<email>johndoe@example.com</email>
				</author>
			</entry>
		</feed>
	`))
	want := &Feed{
		Title:   "Example Feed",
		SiteURL: "http://example.org/",
		Items: []Item{
			{
				GUID:     "urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a",
				Date:     time.Unix(1071340202, 0).UTC(),
				URL:      "http://example.org/2003/12/13/atom03.html",
				Title:    "Atom-Powered Robots Run Amok",
				Content:  `<div xmlns="http://www.w3.org/1999/xhtml"><p>This is the entry content.</p></div>`,
				ImageURL: "",
				AudioURL: "",
			},
		},
	}
	if !reflect.DeepEqual(want, have) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fatal("invalid atom")
	}
}

func TestAtomClashingNamespaces(t *testing.T) {
	have, err := Parse(strings.NewReader(`
		<?xml version="1.0" encoding="utf-8"?>
		<feed xmlns="http://www.w3.org/2005/Atom">
			<entry>
				<content>atom content</content>
				<media:content xmlns:media="http://search.yahoo.com/mrss/" />
			</entry>
		</feed>
	`))
	want := &Feed{Items: []Item{{Content: "atom content"}}}
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(want, have) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.FailNow()
	}
}

func TestAtomHTMLTitle(t *testing.T) {
	feed, _ := Parse(strings.NewReader(`
		<?xml version="1.0" encoding="utf-8"?>
		<feed xmlns="http://www.w3.org/2005/Atom">
			<entry><title type="html">say &lt;code&gt;what&lt;/code&gt;?</entry>
		</feed>
	`))
	have := feed.Items[0].Title
	want := "say what?"
	if !reflect.DeepEqual(want, have) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.FailNow()
	}
}

func TestAtomXHTMLTitle(t *testing.T) {
	feed, _ := Parse(strings.NewReader(`
		<?xml version="1.0" encoding="utf-8"?>
		<feed xmlns="http://www.w3.org/2005/Atom">
			<entry><title type="xhtml">say &lt;code&gt;what&lt;/code&gt;?</entry>
		</feed>
	`))
	have := feed.Items[0].Title
	want := "say what?"
	if !reflect.DeepEqual(want, have) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.FailNow()
	}
}

func TestAtomXHTMLNestedTitle(t *testing.T) {
	feed, _ := Parse(strings.NewReader(`
		<?xml version="1.0" encoding="utf-8"?>
		<feed xmlns="http://www.w3.org/2005/Atom">
			<entry>
				<title type="xhtml">
					<div xmlns="http://www.w3.org/1999/xhtml">
						<a href="https://example.com">Link to Example</a>
					</div>
				</title>
			</entry>
		</feed>
	`))
	have := feed.Items[0].Title
	want := "Link to Example"
	if !reflect.DeepEqual(want, have) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.FailNow()
	}
}

func TestAtomImageLink(t *testing.T) {
	feed, _ := Parse(strings.NewReader(`
		<?xml version="1.0" encoding="UTF-8"?>
		<feed xmlns="http://www.w3.org/2005/Atom" xmlns:media="http://search.yahoo.com/mrss/">
			<entry>
				<media:thumbnail url="https://example.com/image.png?width=100&height=100" />
			</entry>
		</feed>
	`))
	have := feed.Items[0].ImageURL
	want := `https://example.com/image.png?width=100&height=100`
	if want != have {
		t.Fatalf("item.image_url doesn't match\nwant: %#v\nhave: %#v\n", want, have)
	}
}

// found in: https://www.reddit.com/r/funny.rss
// items come with thumbnail urls which are also present in the content
func TestAtomImageLinkDuplicated(t *testing.T) {
	feed, _ := Parse(strings.NewReader(`
		<?xml version="1.0" encoding="utf-8"?>
		<feed xmlns="http://www.w3.org/2005/Atom" xmlns:media="http://search.yahoo.com/mrss/">
			<entry>
				<content type="html">&lt;img src="https://example.com/image.png?width=100&amp;height=100"&gt;</content>
				<media:thumbnail url="https://example.com/image.png?width=100&height=100" />
			</entry>
		</feed>
	`))
	have := feed.Items[0].Content
	want := `<img src="https://example.com/image.png?width=100&height=100">`
	if want != have {
		t.Fatalf("want: %#v\nhave: %#v\n", want, have)
	}
	if feed.Items[0].ImageURL != "" {
		t.Fatal("item.image_url must be unset if present in the content")
	}
}

func TestAtomLinkInID(t *testing.T) {
	feed, _ := Parse(strings.NewReader(`
		<?xml version="1.0" encoding="utf-8"?>
		<feed xmlns="http://www.w3.org/2005/Atom" xmlns:media="http://search.yahoo.com/mrss/">
			<entry>
                <title>one updated</title>
                <id>https://example.com/posts/1</id>
                <updated>2003-12-13T09:17:51</updated>
			</entry>
			<entry>
                <title>two</title>
                <id>urn:uuid:60a76c80-d399-11d9-b93C-0003939e0af6</id>
			</entry>
			<entry>
                <title>one</title>
                <id>https://example.com/posts/1</id>
			</entry>
		</feed>
	`))
	have := feed.Items
	want := []Item{
		Item{
			GUID:  "https://example.com/posts/1::2003-12-13T09:17:51",
			Date:  time.Date(2003, time.December, 13, 9, 17, 51, 0, time.UTC),
			URL:   "https://example.com/posts/1",
			Title: "one updated",
		},
		Item{
			GUID: "urn:uuid:60a76c80-d399-11d9-b93C-0003939e0af6",
			Date: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), URL: "",
			Title: "two",
		},
		Item{
			GUID:    "https://example.com/posts/1::",
			Date:    time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			URL:     "https://example.com/posts/1",
			Title:   "one",
			Content: "",
		},
	}
	if !reflect.DeepEqual(want, have) {
		t.Fatalf("\nwant: %#v\nhave: %#v\n", want, have)
	}
}

func TestAtomDoesntEscapeHTMLTags(t *testing.T) {
	feed, _ := Parse(strings.NewReader(`
		<?xml version="1.0" encoding="utf-8"?>
		<feed xmlns="http://www.w3.org/2005/Atom">
			<entry><summary type="html">&amp;lt;script&amp;gt;alert(1);&amp;lt;/script&amp;gt;</summary></entry>
		</feed>
	`))
	have := feed.Items[0].Content
	want := "&lt;script&gt;alert(1);&lt;/script&gt;"
	if !reflect.DeepEqual(want, have) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.FailNow()
	}
}
