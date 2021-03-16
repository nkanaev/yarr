package opml

import "testing"

var sample = `<?xml version="1.0" encoding="UTF-8"?>
<opml version="1.1">
<head><title>Subscriptions</title></head>
<body>
  <outline text="sub">
    <outline type="rss" text="subtitle1" description="sub1" xmlUrl="https://foo.com/feed.xml" htmlUrl="https://foo.com/"/>
    <outline type="rss" text="&amp;&gt;" description="&lt;&gt;" xmlUrl="https://bar.com/feed.xml" htmlUrl="https://bar.com/"/>
  </outline>
  <outline type="rss" text="title1" description="desc1" xmlUrl="https://example.com/feed.xml" htmlUrl="https://example.com/"/>
</body>
</opml>
`

func TestOPMLBuilder(t *testing.T) {
	builder := NewBuilder()
	builder.AddFeed("title1", "desc1", "https://example.com/feed.xml", "https://example.com/")

	folder := builder.AddFolder("sub")
	folder.AddFeed("subtitle1", "sub1", "https://foo.com/feed.xml", "https://foo.com/")
	folder.AddFeed("&>", "<>", "https://bar.com/feed.xml", "https://bar.com/")

	output := builder.String()
	if output != sample {
		t.Errorf("\n=== expected:\n%s\n=== got:\n%s\n===", sample, output)
	}
}
