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

func easyjsonF642ad3eDecodeGithubComNbdWtfGoNostr(in *jlexer.Lexer, out *Event) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	out.extra = make(map[string]any)
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(true)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = in.String()
		case "pubkey":
			out.PubKey = in.String()
		case "created_at":
			out.CreatedAt = Timestamp(in.Int64())
		case "kind":
			out.Kind = in.Int()
		case "tags":
			if in.IsNull() {
				in.Skip()
				out.Tags = nil
			} else {
				in.Delim('[')
				if out.Tags == nil {
					if !in.IsDelim(']') {
						out.Tags = make(Tags, 0, 7)
					} else {
						out.Tags = Tags{}
					}
				} else {
					out.Tags = (out.Tags)[:0]
				}
				for !in.IsDelim(']') {
					var v1 Tag
					if in.IsNull() {
						in.Skip()
						v1 = nil
					} else {
						in.Delim('[')
						if !in.IsDelim(']') {
							v1 = make(Tag, 0, 5)
						} else {
							v1 = Tag{}
						}
						for !in.IsDelim(']') {
							var v2 string
							v2 = string(in.String())
							v1 = append(v1, v2)
							in.WantComma()
						}
						in.Delim(']')
					}
					out.Tags = append(out.Tags, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "content":
			out.Content = in.String()
		case "sig":
			out.Sig = in.String()
		default:
			out.extra[key] = in.Interface()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}

func easyjsonF642ad3eEncodeGithubComNbdWtfGoNostr(out *jwriter.Writer, in Event) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = "\"kind\":"
		out.RawString(prefix)
		out.Int(in.Kind)
	}
	{
		if in.ID != "" {
			const prefix string = ",\"id\":"
			out.RawString(prefix)
			out.String(in.ID)
		}
	}
	{
		if in.PubKey != "" {
			const prefix string = ",\"pubkey\":"
			out.RawString(prefix)
			out.String(in.PubKey)
		}
	}
	{
		const prefix string = ",\"created_at\":"
		out.RawString(prefix)
		out.Int64(int64(in.CreatedAt))
	}
	{
		const prefix string = ",\"tags\":"
		out.RawString(prefix)
		out.RawByte('[')
		for v3, v4 := range in.Tags {
			if v3 > 0 {
				out.RawByte(',')
			}
			out.RawByte('[')
			for v5, v6 := range v4 {
				if v5 > 0 {
					out.RawByte(',')
				}
				out.String(string(v6))
			}
			out.RawByte(']')
		}
		out.RawByte(']')
	}
	{
		const prefix string = ",\"content\":"
		out.RawString(prefix)
		out.String(in.Content)
	}
	{
		if in.Sig != "" {
			const prefix string = ",\"sig\":"
			out.RawString(prefix)
			out.String(in.Sig)
		}
	}
	{
		for key, value := range in.extra {
			out.RawString(",\"" + key + "\":")
			out.Raw(json.Marshal(value))
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Event) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{NoEscapeHTML: true}
	easyjsonF642ad3eEncodeGithubComNbdWtfGoNostr(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Event) MarshalEasyJSON(w *jwriter.Writer) {
	w.NoEscapeHTML = true
	easyjsonF642ad3eEncodeGithubComNbdWtfGoNostr(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Event) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonF642ad3eDecodeGithubComNbdWtfGoNostr(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Event) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF642ad3eDecodeGithubComNbdWtfGoNostr(l, v)
}
