package extension

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type NostrConfig struct {
}

type NostrOption interface {
	parser.Option
}

type nostrParser struct {
	NostrConfig
}

func NewNostrParser(opts ...NostrOption) parser.InlineParser {
	p := &nostrParser{
		NostrConfig: NostrConfig{},
	}
	return p
}

func (s *nostrParser) Trigger() []byte {
	// ' ' indicates any white spaces and a line head
	return []byte{' ', '*', '_', '~', '('}
}

var (
	protoNostr = []byte("http:")
)

func (s *nostrParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	if pc.IsInLinkLabel() {
		return nil
	}

	line, segment := block.PeekLine()
	consumes := 0
	start := segment.Start
	c := line[0]

	// Skip whitespace and any leading characters
	if c == ' ' || c == '*' || c == '_' || c == '~' || c == '(' {
		consumes++
		start++
		line = line[1:]
	}

	// Check if the line starts with "nostr:"
	if len(line) < 6 || !bytes.HasPrefix(line, []byte("nostr:")) {
		return nil
	}

	// Find the end of the nostr identifier, which is typically alphanumeric
	end := 6
	for end < len(line) && util.IsAlphaNumeric(line[end]) {
		end++
	}

	if end == 6 {
		return nil
	}

	// Extract the nostr identifier
	nostrID := line[6:end]

	// Create a new node for the "nostr:" link
	linkText := ast.NewTextSegment(text.NewSegment(start+6, start+end))
	link := ast.NewLink()
	link.Destination = []byte("nostr:" + string(nostrID))
	link.AppendChild(link, linkText)

	// Create the HTML anchor element

	/*
		htmlLink := &ast.HTMLTag{
			TagName: []byte("a"),
			Attributes: []ast.HTMLTagAttribute{
				{
					Key:   []byte("href"),
					Value: []byte("nostr:" + string(nostrID)),
				},
			},
			Inner: []ast.Node{linkText},
		}
	*/

	// Advance the reader position by the length of the processed string
	block.Advance(consumes + end)

	return link
}

func (s *nostrParser) CloseBlock(parent ast.Node, pc parser.Context) {
	// nothing to do
}

type nostr struct {
	options []NostrOption
}

// Linkify is an extension that allow you to parse text that seems like a URL.
var Nostr = &nostr{}

// NewLinkify creates a new [goldmark.Extender] that
// allow you to parse text that seems like a URL.
func NewNostr(opts ...NostrOption) goldmark.Extender {
	return &nostr{
		options: opts,
	}
}

func (e *nostr) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(NewNostrParser(e.options...), 999),
		),
	)
}
