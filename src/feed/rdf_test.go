package feed

import (
	"reflect"
	"strings"
	"testing"
)

func TestRDFFeed(t *testing.T) {
	have, _ := ParseRDF(strings.NewReader(`<?xml version="1.0"?>
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
		Title: "Mozilla Dot Org",
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
