// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package jsonapi

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

func easyjson7375d8c0DecodeGithubComGsmcwhirterDiscordBotLibV10Jsonapi(in *jlexer.Lexer, out *GatewayResponse) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "url":
			out.URL = string(in.String())
		case "shards":
			out.Shards = int(in.Int())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson7375d8c0EncodeGithubComGsmcwhirterDiscordBotLibV10Jsonapi(out *jwriter.Writer, in GatewayResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"url\":"
		out.RawString(prefix[1:])
		out.String(string(in.URL))
	}
	{
		const prefix string = ",\"shards\":"
		out.RawString(prefix)
		out.Int(int(in.Shards))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v GatewayResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson7375d8c0EncodeGithubComGsmcwhirterDiscordBotLibV10Jsonapi(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v GatewayResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson7375d8c0EncodeGithubComGsmcwhirterDiscordBotLibV10Jsonapi(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *GatewayResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson7375d8c0DecodeGithubComGsmcwhirterDiscordBotLibV10Jsonapi(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *GatewayResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson7375d8c0DecodeGithubComGsmcwhirterDiscordBotLibV10Jsonapi(l, v)
}
