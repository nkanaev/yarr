package htmlutil

import (
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

var nodeNameRegex = regexp.MustCompile(`\w+|\*`)

func FindNodes(node *html.Node, match func(*html.Node) bool) []*html.Node {
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

func Query(node *html.Node, sel string) []*html.Node {
	matcher := NewMatcher(sel)
	return FindNodes(node, matcher.Match)
}

func NewMatcher(sel string) Matcher {
	multi := MultiMatch{}
	parts := strings.Split(sel, ",")
	for _, part := range parts {
		part := strings.TrimSpace(part)
		if nodeNameRegex.MatchString(part) {
			multi.Add(ElementMatch{Name: part})
		} else {
			panic("unsupported selector: " + part)
		}
	}
	return multi
}

type Matcher interface {
	Match(*html.Node) bool
}

type ElementMatch struct {
	Name string
}

func (m ElementMatch) Match(n *html.Node) bool {
	return n.Type == html.ElementNode && (n.Data == m.Name || m.Name == "*")
}

type MultiMatch struct {
	matchers []Matcher
}

func (m *MultiMatch) Add(matcher Matcher) {
	m.matchers = append(m.matchers, matcher)
}

func (m MultiMatch) Match(n *html.Node) bool {
	for _, matcher := range m.matchers {
		if matcher.Match(n) {
			return true
		}
	}
	return false
}
