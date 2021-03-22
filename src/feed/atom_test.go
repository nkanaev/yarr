package feed

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestAtom(t *testing.T) {
	have, _ := ParseAtom(strings.NewReader(`
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
				GUID:       "urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a",
				Date:       time.Unix(1071340202, 0).UTC(),
				URL:        "http://example.org/2003/12/13/atom03.html",
				Title:      "Atom-Powered Robots Run Amok",
				Content:    `<div xmlns="http://www.w3.org/1999/xhtml"><p>This is the entry content.</p></div>`,
				ImageURL:   "",
				PodcastURL: "",
			},
		},
	}
	if !reflect.DeepEqual(want, have) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fatal("invalid atom")
	}
}
