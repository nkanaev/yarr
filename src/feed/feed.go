package feed

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"
)

var UnknownFormat = errors.New("unknown feed format")

type processor func(r io.Reader) (*Feed, error)

func detect(lookup string) (string, processor) {
	lookup = strings.TrimSpace(lookup)
	if lookup[0] == '{' {
		return "json", ParseJSON
	}
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
	return "", nil
}

func Parse(r io.Reader) (*Feed, error) {
	var x [1024]byte
	numread, err := r.Read(x[:])
	fmt.Println(numread, err)
	if err != nil {
		return nil, fmt.Errorf("Failed to read: %s", err)
	}

	_, callback := detect(string(x[:]))
	if callback == nil {
		return nil, UnknownFormat
	}
	return callback(r)
}
