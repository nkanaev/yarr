package parser

import (
	"reflect"
	"strings"
	"testing"
)

func TestSniff(t *testing.T) {
	testcases := [][2]string{
		{
			`<?xml version="1.0"?><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"></rdf:RDF>`,
			"rdf",
		},
		{
			`<?xml version="1.0" encoding="ISO-8859-1"?><rss version="2.0"><channel></channel></rss>`,
			"rss",
		},
		{
			`<?xml version="1.0"?><rss version="2.0"><channel></channel></rss>`,
			"rss",
		},
		{
			`<?xml version="1.0" encoding="utf-8"?><feed xmlns="http://www.w3.org/2005/Atom"></feed>`,
			"atom",
		},
		{
			`{}`,
			"json",
		},
		{
			`<!DOCTYPE html><html><head><title></title></head><body></body></html>`,
			"",
		},
	}
	for _, testcase := range testcases {
		have, _ := sniff(testcase[0])
		want := testcase[1]
		if want != have {
			t.Log(testcase[0])
			t.Errorf("Invalid format: want=%#v have=%#v", want, have)
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
