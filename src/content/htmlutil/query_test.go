package htmlutil

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestQuery(t *testing.T) {
	node, _ := html.Parse(strings.NewReader(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<title></title>
		</head>
		<body>
			<div>
				<p>test</p>
			</div>
		</body>
		</html>
	`))
	nodes := Query(node, "p")
	match := (
		len(nodes) == 1 &&
		nodes[0].Type == html.ElementNode &&
		nodes[0].Data == "p")		
	if !match {
		t.Fatalf("incorrect match: %#v", nodes)
	}
}

func TestQueryMulti(t *testing.T) {
	node, _ := html.Parse(strings.NewReader(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<title></title>
		</head>
		<body>
			<p>foo</p>
			<div>
				<p>bar</p>
				<span>baz</span>
			</div>
		</body>
		</html>
	`))
	nodes := Query(node, "p , span")
	match := (
		len(nodes) == 3 &&
		nodes[0].Type == html.ElementNode && nodes[0].Data == "p" &&
		nodes[1].Type == html.ElementNode && nodes[1].Data == "p" &&
		nodes[2].Type == html.ElementNode && nodes[2].Data == "span")	
	if !match {
		for i, n := range nodes {
			t.Logf("%d: %s", i, HTML(n))
		}
		t.Fatal("incorrect match")
	}
}

func TestClosest(t *testing.T) {
	html, _ := html.Parse(strings.NewReader(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<title></title>
		</head>
		<body>
			<div class="foo">
				<p><a class="bar" href=""></a></p>
			</div>
		</body>
		</html>
	`))
	link := Query(html, "a")
	if link == nil || Attr(link[0], "class") != "bar" {
		t.FailNow()
	}
	wrap := Closest(link[0], "div")
	if wrap == nil || Attr(wrap, "class") != "foo" {
		t.FailNow()
	}
}
