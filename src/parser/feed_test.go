package parser

import (
	"reflect"
	"strings"
	"testing"
)

func TestSniff(t *testing.T) {
	testcases := []struct {
		input string
		want  feedProbe
	}{
		{
			`<?xml version="1.0"?><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"></rdf:RDF>`,
			feedProbe{feedType: "rdf", callback: ParseRDF},
		},
		{
			`<?xml version="1.0" encoding="ISO-8859-1"?><rss version="2.0"><channel></channel></rss>`,
			feedProbe{feedType: "rss", callback: ParseRSS, encoding: "iso-8859-1"},
		},
		{
			`<?xml version="1.0"?><rss version="2.0"><channel></channel></rss>`,
			feedProbe{feedType: "rss", callback: ParseRSS},
		},
		{
			`<?xml version="1.0" encoding="utf-8"?><feed xmlns="http://www.w3.org/2005/Atom"></feed>`,
			feedProbe{feedType: "atom", callback: ParseAtom, encoding: "utf-8"},
		},
		{
			`{}`,
			feedProbe{feedType: "json", callback: ParseJSON},
		},
		{
			`<!DOCTYPE html><html><head><title></title></head><body></body></html>`,
			feedProbe{},
		},
	}
	for _, testcase := range testcases {
		want := testcase.want
		have := sniff(testcase.input)
		if want.encoding != have.encoding || want.feedType != have.feedType {
			t.Errorf("Invalid output\n---\n%s\n---\n\nwant=%#v\nhave=%#v", testcase.input, want, have)
		}
	}
}

func TestParse(t *testing.T) {
	have, _ := Parse(strings.NewReader(`
		<?xml version="1.0"?>
		<rss version="2.0">
		   <channel>
			  <title>
				 Title
			  </title>
			  <item>
				 <title>
				  Item 1
				 </title>
				 <description>
					<![CDATA[<div>content</div>]]>
				 </description>
			  </item>
		   </channel>
		</rss>
	`))
	want := &Feed{
		Title: "Title",
		Items: []Item{
			{
				Title:   "Item 1",
				Content: "<div>content</div>",
			},
		},
	}
	if !reflect.DeepEqual(want, have) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fatal("invalid content")
	}
}

func TestParseShortFeed(t *testing.T) {
	have, err := Parse(strings.NewReader(
		`<?xml version="1.0"?><feed xmlns="http://www.w3.org/2005/Atom"></feed>`,
	))
	want := &Feed{}
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(want, have) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.FailNow()
	}
}

func TestParseFeedWithBOM(t *testing.T) {
	have, err := Parse(strings.NewReader(
		"\xEF\xBB\xBF" + `<?xml version="1.0"?><feed xmlns="http://www.w3.org/2005/Atom"></feed>`,
	))
	want := &Feed{}
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(want, have) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.FailNow()
	}
}

func TestParseCleanIllegalCharsInUTF8(t *testing.T) {
	data := `
		<?xml version="1.0" encoding="UTF-8"?>
		<rss version="2.0" xmlns:content="http://purl.org/rss/1.0/modules/content/">
			<channel>
				<item>
					<title>` + "\a" + `title</title>
				</item>
			</channel>
		</rss>
	`
	feed, err := Parse(strings.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if len(feed.Items) != 1 || feed.Items[0].Title != "title" {
		t.Fatalf("invalid feed, got: %v", feed)
	}
}

func TestParseCleanIllegalCharsInNonUTF8(t *testing.T) {
	// echo привет | iconv -f utf8 -t cp1251 | hexdump -C
	data := `
		<?xml version="1.0" encoding="windows-1251"?>
		<rss version="2.0" xmlns:content="http://purl.org/rss/1.0/modules/content/">
			<channel>
				<item>
					<title>` + "\a \xef\xf0\xe8\xe2\xe5\xf2\x0a \a" + `</title>
				</item>
			</channel>
		</rss>
	`
	feed, err := Parse(strings.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if len(feed.Items) != 1 || feed.Items[0].Title != "привет" {
		t.Fatalf("invalid feed, got: %v", feed)
	}
}

func TestParseMissingGUID(t *testing.T) {
	data := `
		<?xml version="1.0" encoding="windows-1251"?>
		<rss version="2.0" xmlns:content="http://purl.org/rss/1.0/modules/content/">
			<channel>
				<item>
					<title>foo</title>
				</item>
				<item>
					<title>bar</title>
				</item>
			</channel>
		</rss>
	`
	feed, err := ParseAndFix(strings.NewReader(data), "", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(feed.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(feed.Items))
	}
	if feed.Items[0].GUID == "" || feed.Items[1].GUID == "" {
		t.Fatalf("item GUIDs are missing, got %#v", feed.Items)
	}
	if feed.Items[0].GUID == feed.Items[1].GUID {
		t.Fatalf("item GUIDs are not unique, got %#v", feed.Items)
	}
}
