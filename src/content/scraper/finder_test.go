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

	want := map[string]string{
		base + "/feed.xml":  "rss with title",
		base + "/atom.xml":  "",
		base + "/feed.json": "",
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
	want := map[string]string{
		base + "/feed.xml": "",
		base + "/news":     "",
	}
	if !reflect.DeepEqual(want, have) {
		t.Logf("want: %#v", want)
		t.Logf("have: %#v", have)
		t.Fatal("invalid result")
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
