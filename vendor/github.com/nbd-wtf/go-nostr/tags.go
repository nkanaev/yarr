package nostr

import (
	"encoding/json"
	"errors"
	"iter"
	"slices"
	"strings"
)

type Tag []string

// StartsWith checks if a tag contains a prefix.
// for example,
//
//	["p", "abcdef...", "wss://relay.com"]
//
// would match against
//
//	["p", "abcdef..."]
//
// or even
//
//	["p", "abcdef...", "wss://"]
func (tag Tag) StartsWith(prefix []string) bool {
	prefixLen := len(prefix)

	if prefixLen > len(tag) {
		return false
	}
	// check initial elements for equality
	for i := 0; i < prefixLen-1; i++ {
		if prefix[i] != tag[i] {
			return false
		}
	}
	// check last element just for a prefix
	return strings.HasPrefix(tag[prefixLen-1], prefix[prefixLen-1])
}

func (tag Tag) Key() string {
	if len(tag) > 0 {
		return tag[0]
	}
	return ""
}

func (tag Tag) Value() string {
	if len(tag) > 1 {
		return tag[1]
	}
	return ""
}

func (tag Tag) Relay() string {
	if len(tag) > 2 && (tag[0] == "e" || tag[0] == "p") {
		return NormalizeURL(tag[2])
	}
	return ""
}

type Tags []Tag

// GetD gets the first "d" tag (for parameterized replaceable events) value or ""
func (tags Tags) GetD() string {
	for _, v := range tags {
		if v.StartsWith([]string{"d", ""}) {
			return v[1]
		}
	}
	return ""
}

// GetFirst gets the first tag in tags that matches the prefix, see [Tag.StartsWith]
func (tags Tags) GetFirst(tagPrefix []string) *Tag {
	for _, v := range tags {
		if v.StartsWith(tagPrefix) {
			return &v
		}
	}
	return nil
}

// GetLast gets the last tag in tags that matches the prefix, see [Tag.StartsWith]
func (tags Tags) GetLast(tagPrefix []string) *Tag {
	for i := len(tags) - 1; i >= 0; i-- {
		v := tags[i]
		if v.StartsWith(tagPrefix) {
			return &v
		}
	}
	return nil
}

// GetAll gets all the tags that match the prefix, see [Tag.StartsWith]
func (tags Tags) GetAll(tagPrefix []string) Tags {
	result := make(Tags, 0, len(tags))
	for _, v := range tags {
		if v.StartsWith(tagPrefix) {
			result = append(result, v)
		}
	}
	return result
}

// All returns an iterator for all the tags that match the prefix, see [Tag.StartsWith]
func (tags Tags) All(tagPrefix []string) iter.Seq2[int, Tag] {
	return func(yield func(int, Tag) bool) {
		for i, v := range tags {
			if v.StartsWith(tagPrefix) {
				if !yield(i, v) {
					break
				}
			}
		}
	}
}

// FilterOut returns a new slice with only the elements that match the prefix, see [Tag.StartsWith]
func (tags Tags) FilterOut(tagPrefix []string) Tags {
	filtered := make(Tags, 0, len(tags))
	for _, v := range tags {
		if !v.StartsWith(tagPrefix) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

// FilterOutInPlace removes all tags that match the prefix, but potentially reorders the tags in unpredictable ways, see [Tag.StartsWith]
func (tags *Tags) FilterOutInPlace(tagPrefix []string) {
	for i := 0; i < len(*tags); i++ {
		tag := (*tags)[i]
		if tag.StartsWith(tagPrefix) {
			// remove this by swapping the last tag into this place
			last := len(*tags) - 1
			(*tags)[i] = (*tags)[last]
			*tags = (*tags)[0:last]
			i-- // this is so we can match this just swapped item in the next iteration
		}
	}
}

// AppendUnique appends a tag if it doesn't exist yet, otherwise does nothing.
// the uniqueness comparison is done based only on the first 2 elements of the tag.
func (tags Tags) AppendUnique(tag Tag) Tags {
	n := len(tag)
	if n > 2 {
		n = 2
	}

	if tags.GetFirst(tag[:n]) == nil {
		return append(tags, tag)
	}
	return tags
}

func (t *Tags) Scan(src any) error {
	var jtags []byte

	switch v := src.(type) {
	case []byte:
		jtags = v
	case string:
		jtags = []byte(v)
	default:
		return errors.New("couldn't scan tags, it's not a json string")
	}

	json.Unmarshal(jtags, &t)
	return nil
}

func (tags Tags) ContainsAny(tagName string, values []string) bool {
	for _, tag := range tags {
		if len(tag) < 2 {
			continue
		}

		if tag[0] != tagName {
			continue
		}

		if slices.Contains(values, tag[1]) {
			return true
		}
	}

	return false
}

// Marshal Tag. Used for Serialization so string escaping should be as in RFC8259.
func (tag Tag) marshalTo(dst []byte) []byte {
	dst = append(dst, '[')
	for i, s := range tag {
		if i > 0 {
			dst = append(dst, ',')
		}
		dst = escapeString(dst, s)
	}
	dst = append(dst, ']')
	return dst
}

// MarshalTo appends the JSON encoded byte of Tags as [][]string to dst.
// String escaping is as described in RFC8259.
func (tags Tags) marshalTo(dst []byte) []byte {
	dst = append(dst, '[')
	for i, tag := range tags {
		if i > 0 {
			dst = append(dst, ',')
		}
		dst = tag.marshalTo(dst)
	}
	dst = append(dst, ']')
	return dst
}
