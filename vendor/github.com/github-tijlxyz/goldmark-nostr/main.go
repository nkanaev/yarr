package extension

import (
	"bytes"
	"fmt"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type NostrConfig struct {
	NostrLink string
	Strict    bool
}

type NostrOption interface {
	parser.Option
	SetNostrOption(*NostrConfig)
}

type nostrParser struct {
	NostrConfig
}

func NewNostrParser(opts ...NostrOption) parser.InlineParser {
	p := &nostrParser{
		NostrConfig: NostrConfig{},
	}
	for _, o := range opts {
		o.SetNostrOption(&p.NostrConfig)
	}
	return p
}

func WithStrict() NostrOption {
	return &withStrict{
		value: true,
	}
}

func WithNostrLink(link string) NostrOption {
	return &withNostrLink{
		value: link,
	}
}

type withNostrLink struct {
	value string
}

func (o *withNostrLink) SetParserOption(c *parser.Config) {
	c.Options["NostrLink"] = o.value
}
func (o *withNostrLink) SetNostrOption(p *NostrConfig) {
	p.NostrLink = o.value
}

type withStrict struct {
	value bool
}

func (o *withStrict) SetParserOption(c *parser.Config) {
	c.Options["Strict"] = o.value
}
func (o *withStrict) SetNostrOption(p *NostrConfig) {
	p.Strict = o.value
}

func (s *nostrParser) Trigger() []byte {
	// ' ' indicates any white spaces and a line head
	return []byte{' ', '*', '_', '~', '('}
}

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

	link.Destination = []byte(fmt.Sprintf(s.NostrConfig.NostrLink, string(nostrID)))

	if s.NostrConfig.Strict {
		link.SetAttribute([]byte("rel"), "noopener noreferrer")
		link.SetAttribute([]byte("target"), "_blank")
		link.SetAttribute([]byte("referrerpolicy"), "no-referrer")
	}

	link.AppendChild(link, linkText)

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

var Nostr = &nostr{}

func New(opts ...NostrOption) goldmark.Extender {
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
