package parser

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"io"
	"regexp"
	"strings"

	"golang.org/x/net/html/charset"
)

func firstNonEmpty(vals ...string) string {
	for _, val := range vals {
		valTrimmed := strings.TrimSpace(val)
		if len(valTrimmed) > 0 {
			return valTrimmed
		}
	}
	return ""
}

var linkRe = regexp.MustCompile(`(https?:\/\/\S+)`)

func plain2html(text string) string {
	text = linkRe.ReplaceAllString(text, `<a href="$1">$1</a>`)
	text = strings.ReplaceAll(text, "\n", "<br>")
	return text
}

func xmlDecoder(r io.Reader) *xml.Decoder {
	decoder := xml.NewDecoder(r)
	decoder.Strict = false
	decoder.CharsetReader = func(cs string, input io.Reader) (io.Reader, error) {
		r, err := charset.NewReaderLabel(cs, input)
		if err == nil {
			r = NewSafeXMLReader(r)
		}
		return r, err
	}
	return decoder
}

// XML token reader that strips the default namespace.
// It's primary purpose is to support namespaced legacy UserLand RSS feeds.
// NOTE: token readers cannot populate ",innerxml"-tagged struct fields,
// see https://github.com/golang/go/issues/39645
type rssTokenReader struct {
	Decoder   *xml.Decoder
	defaultNS string
}

func (r *rssTokenReader) Token() (xml.Token, error) {
	tok, err := r.Decoder.Token()
	if err != nil {
		return nil, err
	}

	switch t := tok.(type) {
	case xml.StartElement:
		// extract default namespace: <rss xmlns="<defaultNS>">
		if t.Name.Local == "rss" {
			for _, attr := range t.Attr {
				if attr.Name.Space == "" && attr.Name.Local == "xmlns" && attr.Value != "" {
					r.defaultNS = attr.Value
					break
				}
			}
		}

		if r.defaultNS != "" {
			// Rewrite element namespace
			if t.Name.Space == r.defaultNS {
				t.Name.Space = r.Decoder.DefaultSpace
			}
			// Rewrite attribute namespaces
			attrs := t.Attr[:0]
			for _, a := range t.Attr {
				if a.Name.Space == r.defaultNS {
					a.Name.Space = r.Decoder.DefaultSpace
				}
				attrs = append(attrs, a)
			}
			t.Attr = attrs
		}
		return t, nil
	case xml.EndElement:
		if r.defaultNS != "" && t.Name.Space == r.defaultNS {
			t.Name.Space = r.Decoder.DefaultSpace
		}
		return t, nil
	default:
		return tok, nil
	}
}

type safexmlreader struct {
	reader *bufio.Reader
	buffer *bytes.Buffer
}

func NewSafeXMLReader(r io.Reader) io.Reader {
	return &safexmlreader{
		reader: bufio.NewReader(r),
		buffer: bytes.NewBuffer(make([]byte, 0, 4096)),
	}
}

func (xr *safexmlreader) Read(p []byte) (int, error) {
	for xr.buffer.Len() < cap(p) {
		r, _, err := xr.reader.ReadRune()
		if err == io.EOF {
			if xr.buffer.Len() == 0 {
				return 0, io.EOF
			}
			break
		}
		if err != nil {
			return 0, err
		}
		if isInCharacterRange(r) {
			xr.buffer.WriteRune(r)
		}
	}
	return xr.buffer.Read(p)
}

// NOTE: copied from "encoding/xml" package
// Decide whether the given rune is in the XML Character Range, per
// the Char production of https://www.xml.com/axml/testaxml.htm,
// Section 2.2 Characters.
func isInCharacterRange(r rune) (inrange bool) {
	return r == 0x09 ||
		r == 0x0A ||
		r == 0x0D ||
		r >= 0x20 && r <= 0xD7FF ||
		r >= 0xE000 && r <= 0xFFFD ||
		r >= 0x10000 && r <= 0x10FFFF
}

// NOTE: copied from "encoding/xml" package
// procInst parses the `param="..."` or `param='...'`
// value out of the provided string, returning "" if not found.
func procInst(param, s string) string {
	// TODO: this parsing is somewhat lame and not exact.
	// It works for all actual cases, though.
	param = param + "="
	idx := strings.Index(s, param)
	if idx == -1 {
		return ""
	}
	v := s[idx+len(param):]
	if v == "" {
		return ""
	}
	if v[0] != '\'' && v[0] != '"' {
		return ""
	}
	idx = strings.IndexRune(v[1:], rune(v[0]))
	if idx == -1 {
		return ""
	}
	return v[1 : idx+1]
}
