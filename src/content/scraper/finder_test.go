package scraper

import (
	"reflect"
	"testing"
)

const base = "http://example.com"

func TestFindFeedsInvalidHTML(t *testing.T) {
	x := `some nonsense`
	r := FindFeeds(x, base)
	if len(r) != 0 {
		t.Fatal("not expecting results")
	}
}

func TestFindFeedsLinks(t *testing.T) {
	x := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<title></title>
			<link rel="alternate" href="/feed.xml" type="application/rss+xml" title="rss with title">
			<link rel="alternate" href="/atom.xml" type="application/atom+xml">
			<link rel="alternate" href="/feed.json" type="application/json">
		</head>
		<body>
			<a href="/feed.xml">rss</a>
		</body>
		</html>
	`
	have := FindFeeds(x, base)

	want := []FeedLink{
		{URL: base + "/atom.xml", Title: ""},
		{URL: base + "/feed.json", Title: ""},
		{URL: base + "/feed.xml", Title: "rss with title"},
	}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fatal("invalid result")
	}
}

func TestFindFeedsGuess(t *testing.T) {
	body := `
		<!DOCTYPE html>
		<html lang="en">
		<body>
			<!-- negative -->
			<a href="/about">what is rss?</a>
			<a href="/feed/cows">moo</a>

			<!-- positive -->
			<a href="/feed.xml">subscribe</a>
			<a href="/news">rss</a>
		</body>
		</html>
	`
	have := FindFeeds(body, base)
	want := []FeedLink{
		{URL: base + "/feed.xml", Title: ""},
		{URL: base + "/news", Title: ""},
	}
	if !reflect.DeepEqual(want, have) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fatal("invalid result")
	}
}

func TestFindFeedsYouTubeOGTitle(t *testing.T) {
	body := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta property="og:title" content="My Channel">
			<link rel="alternate" href="https://www.youtube.com/feeds/videos.xml?channel_id=UCabc123" type="application/rss+xml" title="YouTube Channel">
		</head>
		<body></body>
		</html>
	`
	have := FindFeeds(body, base)

	youtubeURL := "https://www.youtube.com/feeds/videos.xml?playlist_id="
	want := []FeedLink{
		{URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCabc123", Title: "YouTube Channel - All"},
		{URL: youtubeURL + "UULVabc123", Title: "YouTube Channel - Live Streams", TitleOverride: "My Channel - Live Streams"},
		{URL: youtubeURL + "UUSHabc123", Title: "YouTube Channel - Short videos", TitleOverride: "My Channel - Short videos"},
		{URL: youtubeURL + "UULFabc123", Title: "YouTube Channel - Videos", TitleOverride: "My Channel - Videos"},
	}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fatal("invalid result")
	}
}

func TestFindFeedsYouTubeNoOGTitle(t *testing.T) {
	body := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<link rel="alternate" href="https://www.youtube.com/feeds/videos.xml?channel_id=UCxyz789" type="application/rss+xml" title="Channel Name">
		</head>
		<body></body>
		</html>
	`
	have := FindFeeds(body, base)

	youtubeURL := "https://www.youtube.com/feeds/videos.xml?playlist_id="
	want := []FeedLink{
		{URL: "https://www.youtube.com/feeds/videos.xml?channel_id=UCxyz789", Title: "Channel Name - All"},
		{URL: youtubeURL + "UULVxyz789", Title: "Channel Name - Live Streams", TitleOverride: "Channel Name - Live Streams"},
		{URL: youtubeURL + "UUSHxyz789", Title: "Channel Name - Short videos", TitleOverride: "Channel Name - Short videos"},
		{URL: youtubeURL + "UULFxyz789", Title: "Channel Name - Videos", TitleOverride: "Channel Name - Videos"},
	}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fatal("invalid result")
	}
}

func TestFindFeedsYouTubeNoChannelID(t *testing.T) {
	body := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<link rel="alternate" href="https://www.youtube.com/feeds/videos.xml?channel_id=invalid" type="application/rss+xml" title="Youtube">
		</head>
		<body></body>
		</html>
	`
	have := FindFeeds(body, base)

	want := []FeedLink{
		{URL: "https://www.youtube.com/feeds/videos.xml?channel_id=invalid", Title: "Youtube"},
	}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fatal("invalid result")
	}
}

func TestFindFeedsNonYouTubeNoTitleOverride(t *testing.T) {
	body := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<link rel="alternate" href="/blog.xml" type="application/rss+xml" title="Blog">
		</head>
		<body></body>
		</html>
	`
	have := FindFeeds(body, base)

	want := []FeedLink{
		{URL: base + "/blog.xml", Title: "Blog"},
	}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fatal("invalid result")
	}
	if have[0].TitleOverride != "" {
		t.Fatal("expected empty TitleOverride for non-YouTube feed")
	}
}

func TestFindIcons(t *testing.T) {
	body := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<title></title>
			<link rel="icon favicon" href="/favicon.ico">
			<link rel="icon macicon" href="path/to/favicon.png">
		</head>
		<body>
			
		</body>
		</html>
	`
	have := FindIcons(body, base)
	want := []string{base + "/favicon.ico", base + "/path/to/favicon.png"}
	if !reflect.DeepEqual(have, want) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fatal("invalid result")
	}
}
