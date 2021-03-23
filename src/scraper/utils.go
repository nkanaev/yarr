package scraper

import (
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func any(els []string, el string, match func(string, string) bool) bool {
	for _, x := range els {
		if match(x, el) {
			return true
		}
	}
	return false
}

func getAttr(node *html.Node, key string) string {
	for _, a := range node.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

func getText(node *html.Node) string {
	text := make([]string, 0)
	isTextNode := func(n *html.Node) bool {
		return n.Type == html.TextNode
	}
	for _, n := range getNodes(node, isTextNode) {
		text = append(text, strings.TrimSpace(n.Data))
	}
	return strings.Join(text, " ")
}

func getNodes(node *html.Node, match func(*html.Node) bool) []*html.Node {
	nodes := make([]*html.Node, 0)

	queue := make([]*html.Node, 0)
	queue = append(queue, node)
	for len(queue) > 0 {
		var n *html.Node
		n, queue = queue[0], queue[1:]
		if match(n) {
			nodes = append(nodes, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			queue = append(queue, c)
		}
	}
	return nodes
}

func absoluteUrl(href, base string) string {
	baseUrl, err := url.Parse(base)
	if err != nil {
		return ""
	}
	hrefUrl, err := url.Parse(href)
	if err != nil {
		return ""
	}
	return baseUrl.ResolveReference(hrefUrl).String()
}
