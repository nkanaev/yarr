package parser

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestRDFFeed(t *testing.T) {
	have, _ := Parse(strings.NewReader(`<?xml version="1.0"?>
		<rdf:RDF
		xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
		xmlns="http://channel.netscape.com/rdf/simple/0.9/">

		  <channel>
			<title>Mozilla Dot Org</title>
			<link>http://www.mozilla.org</link>
			<description>the Mozilla Organization
			  web site</description>
		  </channel>

		  <image>
			<title>Mozilla</title>
			<url>http://www.mozilla.org/images/moz.gif</url>
			<link>http://www.mozilla.org</link>
		  </image>

		  <item>
			<title>New Status Updates</title>
			<link>http://www.mozilla.org/status/</link>
		  </item>

		  <item>
			<title>Bugzilla Reorganized</title>
			<link>http://www.mozilla.org/bugs/</link>
		  </item>

		</rdf:RDF>
	`))
	want := &Feed{
		Title:   "Mozilla Dot Org",
		SiteURL: "http://www.mozilla.org",
		Items: []Item{
			{GUID: "http://www.mozilla.org/status/", URL: "http://www.mozilla.org/status/", Title: "New Status Updates"},
			{GUID: "http://www.mozilla.org/bugs/", URL: "http://www.mozilla.org/bugs/", Title: "Bugzilla Reorganized"},
		},
	}

	if !reflect.DeepEqual(want, have) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fatal("invalid rdf")
	}
}

func TestRDFExtensions(t *testing.T) {
	have, _ := Parse(strings.NewReader(`
		<?xml version="1.0" encoding="utf-8"?>
		<rdf:RDF xmlns="http://purl.org/rss/1.0/"
				xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
				xmlns:dc="http://purl.org/dc/elements/1.1/"
				xmlns:content="http://purl.org/rss/1.0/modules/content/">
			<item>
				<dc:date>2006-01-02T15:04:05-07:00</dc:date>
				<content:encoded><![CDATA[test]]></content:encoded>
			</item>
		</rdf:RDF>
	`))
	date, _ := time.Parse(time.RFC1123Z, time.RFC1123Z)
	want := &Feed{
		Items: []Item{
			{Content: "test", Date: date},
		},
	}
	if !reflect.DeepEqual(want, have) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.FailNow()
	}
}
