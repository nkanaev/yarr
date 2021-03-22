package feed

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"
)

var UnknownFormat = errors.New("unknown feed format")

type processor func(r io.Reader) (*Feed, error)

func sniff(lookup string) (string, processor) {
	lookup = strings.TrimSpace(lookup)
	switch lookup[0] {
	case '<':
		decoder := xml.NewDecoder(strings.NewReader(lookup))
		for {
			token, _ := decoder.Token()
			if token == nil {
				break
			}
			if el, ok := token.(xml.StartElement); ok {
				switch el.Name.Local {
				case "rss":
					return "rss", ParseRSS
				case "RDF":
					return "rss", ParseRDF
				case "feed":
					return "atom", ParseAtom
				}
			}
		}
	case '{':
		return "json", ParseJSON
	}
	return "", nil
}

func Parse(r io.Reader) (*Feed, error) {
	chunk := make([]byte, 64)
	numread, err := r.Read(chunk)
	fmt.Println(numread, err)
	if err != nil {
		return nil, fmt.Errorf("Failed to read: %s", err)
	}

	_, callback := sniff(string(chunk))
	if callback == nil {
		return nil, UnknownFormat
	}

	r = io.MultiReader(bytes.NewReader(chunk), r)
	return callback(r)
}
