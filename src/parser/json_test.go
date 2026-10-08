package parser

import (
	"fmt"
	"strings"
	"testing"
)

func TestParseJSONAttachmentNumbers(t *testing.T) {
	testcases := []struct {
		field   string
		value   string
		wantErr bool
	}{
		{"size_in_bytes", "120", false},
		{"size_in_bytes", "120.0", false},
		{"size_in_bytes", "1.2e2", false},
		{"size_in_bytes", `"120"`, true},
		{"duration_in_seconds", "120", false},
		{"duration_in_seconds", "120.0", false},
		{"duration_in_seconds", "120.5", false},
		{"duration_in_seconds", "1.2e2", false},
		{"duration_in_seconds", `"120"`, true},
	}
	for _, testcase := range testcases {
		t.Run(testcase.field+"/"+testcase.value, func(t *testing.T) {
			data := fmt.Sprintf(`{
				"version": "https://jsonfeed.org/version/1.1",
				"title": "Podcast Feed",
				"items": [{
					"id": "episode-1",
					"content_text": "Episode one.",
					"attachments": [{
						"url": "https://example.org/audio.mp3",
						"mime_type": "audio/mpeg",
						"title": "Episode 1",
						%q: %s
					}]
				}]
			}`, testcase.field, testcase.value)
			feed, err := Parse(strings.NewReader(data))
			if testcase.wantErr {
				if err == nil {
					t.Fatal("expected a parse error for a non-numeric attachment field")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if len(feed.Items) != 1 || len(feed.Items[0].MediaLinks) != 1 {
				t.Fatalf("expected one item with one attachment, got %#v", feed)
			}
			want := MediaLink{
				URL:         "https://example.org/audio.mp3",
				Type:        "audio",
				Description: "Episode 1",
			}
			if got := feed.Items[0].MediaLinks[0]; got != want {
				t.Errorf("media link: got %#v, want %#v", got, want)
			}
		})
	}
}
