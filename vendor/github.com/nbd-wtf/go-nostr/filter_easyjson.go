package nostr

import (
	json "encoding/json"

	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson4d398eaaDecodeGithubComNbdWtfGoNostr(in *jlexer.Lexer, out *Filter) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	out.Tags = make(TagMap)
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "ids":
			if in.IsNull() {
				in.Skip()
				out.IDs = nil
			} else {
				in.Delim('[')
				if out.IDs == nil {
					if !in.IsDelim(']') {
						out.IDs = make([]string, 0, 20)
					} else {
						out.IDs = []string{}
					}
				} else {
					out.IDs = (out.IDs)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.IDs = append(out.IDs, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "kinds":
			if in.IsNull() {
				in.Skip()
				out.Kinds = nil
			} else {
				in.Delim('[')
				if out.Kinds == nil {
					if !in.IsDelim(']') {
						out.Kinds = make([]int, 0, 8)
					} else {
						out.Kinds = []int{}
					}
				} else {
					out.Kinds = (out.Kinds)[:0]
				}
				for !in.IsDelim(']') {
					var v2 int
					v2 = int(in.Int())
					out.Kinds = append(out.Kinds, v2)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "authors":
			if in.IsNull() {
				in.Skip()
				out.Authors = nil
			} else {
				in.Delim('[')
				if out.Authors == nil {
					if !in.IsDelim(']') {
						out.Authors = make([]string, 0, 40)
					} else {
						out.Authors = []string{}
					}
				} else {
					out.Authors = (out.Authors)[:0]
				}
				for !in.IsDelim(']') {
					var v3 string
					v3 = string(in.String())
					out.Authors = append(out.Authors, v3)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "since":
			if in.IsNull() {
				in.Skip()
				out.Since = nil
			} else {
				if out.Since == nil {
					out.Since = new(Timestamp)
				}
				*out.Since = Timestamp(in.Int64())
			}
		case "until":
			if in.IsNull() {
				in.Skip()
				out.Until = nil
			} else {
				if out.Until == nil {
					out.Until = new(Timestamp)
				}
				*out.Until = Timestamp(in.Int64())
			}
		case "limit":
			out.Limit = int(in.Int())
			if out.Limit == 0 {
				out.LimitZero = true
			}
		case "search":
			out.Search = string(in.String())
		default:
			if len(key) > 1 && key[0] == '#' {
				tagValues := make([]string, 0, 40)
				if !in.IsNull() {
					in.Delim('[')
					if out.Authors == nil {
						if !in.IsDelim(']') {
							tagValues = make([]string, 0, 4)
						} else {
							tagValues = []string{}
						}
					} else {
						tagValues = (tagValues)[:0]
					}
					for !in.IsDelim(']') {
						var v3 string
						v3 = string(in.String())
						tagValues = append(tagValues, v3)
						in.WantComma()
					}
					in.Delim(']')
				}
				out.Tags[key[1:]] = tagValues
			} else {
				in.SkipRecursive()
			}
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}

func easyjson4d398eaaEncodeGithubComNbdWtfGoNostr(out *jwriter.Writer, in Filter) {
	out.RawByte('{')
	first := true
	_ = first
	if len(in.IDs) != 0 {
		const prefix string = ",\"ids\":"
		first = false
		out.RawString(prefix[1:])
		{
			out.RawByte('[')
			for v4, v5 := range in.IDs {
				if v4 > 0 {
					out.RawByte(',')
				}
				out.String(string(v5))
			}
			out.RawByte(']')
		}
	}
	if len(in.Kinds) != 0 {
		const prefix string = ",\"kinds\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v6, v7 := range in.Kinds {
				if v6 > 0 {
					out.RawByte(',')
				}
				out.Int(int(v7))
			}
			out.RawByte(']')
		}
	}
	if len(in.Authors) != 0 {
		const prefix string = ",\"authors\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v8, v9 := range in.Authors {
				if v8 > 0 {
					out.RawByte(',')
				}
				out.String(string(v9))
			}
			out.RawByte(']')
		}
	}
	if in.Since != nil {
		const prefix string = ",\"since\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(*in.Since))
	}
	if in.Until != nil {
		const prefix string = ",\"until\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(*in.Until))
	}
	if in.Limit != 0 || in.LimitZero {
		const prefix string = ",\"limit\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Limit))
	}
	if in.Search != "" {
		const prefix string = ",\"search\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Search))
	}
	for tag, values := range in.Tags {
		const prefix string = ",\"authors\":"
		if first {
			first = false
			out.RawString("\"#" + tag + "\":")
		} else {
			out.RawString(",\"#" + tag + "\":")
		}
		{
			out.RawByte('[')
			for i, v := range values {
				if i > 0 {
					out.RawByte(',')
				}
				out.String(string(v))
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Filter) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{NoEscapeHTML: true}
	easyjson4d398eaaEncodeGithubComNbdWtfGoNostr(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Filter) MarshalEasyJSON(w *jwriter.Writer) {
	w.NoEscapeHTML = true
	easyjson4d398eaaEncodeGithubComNbdWtfGoNostr(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Filter) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson4d398eaaDecodeGithubComNbdWtfGoNostr(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Filter) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson4d398eaaDecodeGithubComNbdWtfGoNostr(l, v)
}
