package parser

import (
	"reflect"
	"strings"
	"testing"
)

func TestJSONFeed(t *testing.T) {
	have, _ := Parse(strings.NewReader(`{
		"version": "https://jsonfeed.org/version/1",
		"title": "My Example Feed",
		"home_page_url": "https://example.org/",
		"feed_url": "https://example.org/feed.json",
		"items": [
			{
				"id": "2",
				"content_text": "This is a second item.",
				"url": "https://example.org/second-item"
			},
			{
				"id": "1",
				"content_html": "<p>Hello, world!</p>",
				"url": "https://example.org/initial-post"
			}
		]
	}`))
	want := &Feed{
		Title:   "My Example Feed",
		SiteURL: "https://example.org/",
		Items: []Item{
			{GUID: "2", Content: "This is a second item.", URL: "https://example.org/second-item"},
			{GUID: "1", Content: "<p>Hello, world!</p>", URL: "https://example.org/initial-post"},
		},
	}

	if !reflect.DeepEqual(want, have) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fatal("invalid json")
	}
}

func TestJSONFeedItemAttachementsHasImage(t *testing.T) {
	feed, _ := Parse(strings.NewReader(`{
		"items": [
			{
				"attachments": [
					{
						"url": "http://example.org/image",
						"mime_type": "image/png"
					}
				]
			}
		]
	}`))

	have := feed.Items[0].MediaLinks
	want := []MediaLink{
		MediaLink{
			URL:         "http://example.org/image",
			Type:        "image",
			Description: "",
		},
	}

	if !reflect.DeepEqual(want, have) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fatal("invalid json")
	}
}

func TestJSONFeedItemLinkIsImage(t *testing.T) {
	feed, _ := Parse(strings.NewReader(`{
		"items": [
			{
				"url": "http://example.org/image.jpg"
			}
		]
	}`))

	have := feed.Items[0].MediaLinks
	want := []MediaLink{
		MediaLink{
			URL:         "http://example.org/image.jpg",
			Type:        "image",
			Description: "",
		},
	}

	if !reflect.DeepEqual(want, have) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fatal("invalid json")
	}
}

func TestJSONFeedItemContentHasImage(t *testing.T) {
	feed, _ := Parse(strings.NewReader(`{
		"items": [
			{
				"content_html": "<p>foobar</p> <img src=\"http://example.org/image\" />"
			}
		]
	}`))

	have := feed.Items[0].MediaLinks
	want := []MediaLink{
		MediaLink{
			URL:         "http://example.org/image",
			Type:        "image",
			Description: "",
		},
	}

	if !reflect.DeepEqual(want, have) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fatal("invalid json")
	}
}
