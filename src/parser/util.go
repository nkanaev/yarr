package parser

import (
	"bufio"
	"encoding/xml"
	"io"
	"regexp"
	"strings"
	"unicode/utf8"

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
	decoder := xml.NewDecoder(NewSafeXMLReader(r))
	decoder.Strict = false
	decoder.CharsetReader = charset.NewReaderLabel
	return decoder
}

type safexmlreader struct {
	reader  *bufio.Reader
	buffer  []byte
	isEOF   bool
	runebuf []byte
}

func NewSafeXMLReader(r io.Reader) io.Reader {
	return &safexmlreader{
		reader:  bufio.NewReader(r),
		runebuf: make([]byte, 6),
	}
}

func (xr *safexmlreader) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	for len(xr.buffer) < cap(p) {
		r, _, err := xr.reader.ReadRune()
		if err == io.EOF {
			xr.isEOF = true
			break
		}
		if err != nil {
			return 0, err
		}
		if isInCharacterRange(r) {
			size := utf8.EncodeRune(xr.runebuf, r)
			xr.buffer = append(xr.buffer, xr.runebuf[:size]...)
		}
	}

	if xr.isEOF && len(xr.buffer) == 0 {
		return 0, io.EOF
	}

	n := cap(p)
	if len(xr.buffer) < n {
		n = len(xr.buffer)
	}
	copy(p, xr.buffer[:n])
	xr.buffer = xr.buffer[n:]
	return n, nil
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
