package opml

import (
	"reflect"
	"testing"
)

func TestOPML(t *testing.T) {
	have := (Folder{
		Title: "",
		Feeds: []Feed{
			{
				Title:       "title1",
				FeedUrl:     "https://baz.com/feed.xml",
				SiteUrl:     "https://baz.com/",
				CustomOrder: "",
			},
		},
		Folders: []Folder{
			{
				Title: "sub",
				Feeds: []Feed{
					{
						Title:       "subtitle1",
						FeedUrl:     "https://foo.com/feed.xml",
						SiteUrl:     "https://foo.com/",
						CustomOrder: "123",
					},
					{
						Title:       "&>",
						FeedUrl:     "https://bar.com/feed.xml",
						SiteUrl:     "https://bar.com/",
						CustomOrder: "456",
					},
				},
				Folders: []Folder{},
			},
		},
	}).OPML()
	want := `<?xml version="1.0" encoding="UTF-8"?>
<opml version="1.1">
<head><title>subscriptions</title></head>
<body>
  <outline text="sub">
    <outline type="rss" text="subtitle1" xmlUrl="https://foo.com/feed.xml" htmlUrl="https://foo.com/" customOrder="123"/>
    <outline type="rss" text="&amp;&gt;" xmlUrl="https://bar.com/feed.xml" htmlUrl="https://bar.com/" customOrder="456"/>
  </outline>
  <outline type="rss" text="title1" xmlUrl="https://baz.com/feed.xml" htmlUrl="https://baz.com/" customOrder=""/>
</body>
</opml>
`
	if !reflect.DeepEqual(want, have) {
		t.Logf("want: %s", want)
		t.Logf("have: %s", have)
		t.Fatal("invalid opml")
	}
}
