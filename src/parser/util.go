package parser

import (
	"encoding/xml"
	"io"
	"regexp"
	"strings"
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
	decoder.CharsetReader = func(charset string, reader io.Reader) (io.Reader, error) { return reader, nil }
	return decoder
}
