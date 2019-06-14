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

func easyjson4086215fDecodeGithubComGsmcwhirterDiscordBotLibV7Jsonapi(in *jlexer.Lexer, out *MessageWithEmbed) {
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
		case "content":
			out.Content = string(in.String())
		case "tts":
			out.Tts = bool(in.Bool())
		case "embed":
			(out.Embed).UnmarshalEasyJSON(in)
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
func easyjson4086215fEncodeGithubComGsmcwhirterDiscordBotLibV7Jsonapi(out *jwriter.Writer, in MessageWithEmbed) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"content\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Content))
	}
	{
		const prefix string = ",\"tts\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.Tts))
	}
	{
		const prefix string = ",\"embed\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.Embed).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v MessageWithEmbed) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson4086215fEncodeGithubComGsmcwhirterDiscordBotLibV7Jsonapi(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v MessageWithEmbed) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson4086215fEncodeGithubComGsmcwhirterDiscordBotLibV7Jsonapi(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *MessageWithEmbed) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson4086215fDecodeGithubComGsmcwhirterDiscordBotLibV7Jsonapi(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *MessageWithEmbed) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson4086215fDecodeGithubComGsmcwhirterDiscordBotLibV7Jsonapi(l, v)
}
func easyjson4086215fDecodeGithubComGsmcwhirterDiscordBotLibV7Jsonapi1(in *jlexer.Lexer, out *Message) {
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
		case "content":
			out.Content = string(in.String())
		case "tts":
			out.Tts = bool(in.Bool())
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
func easyjson4086215fEncodeGithubComGsmcwhirterDiscordBotLibV7Jsonapi1(out *jwriter.Writer, in Message) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"content\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Content))
	}
	{
		const prefix string = ",\"tts\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.Tts))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Message) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson4086215fEncodeGithubComGsmcwhirterDiscordBotLibV7Jsonapi1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Message) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson4086215fEncodeGithubComGsmcwhirterDiscordBotLibV7Jsonapi1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Message) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson4086215fDecodeGithubComGsmcwhirterDiscordBotLibV7Jsonapi1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Message) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson4086215fDecodeGithubComGsmcwhirterDiscordBotLibV7Jsonapi1(l, v)
}
func easyjson4086215fDecodeGithubComGsmcwhirterDiscordBotLibV7Jsonapi2(in *jlexer.Lexer, out *EmbedFooter) {
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
		case "text":
			out.Text = string(in.String())
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
func easyjson4086215fEncodeGithubComGsmcwhirterDiscordBotLibV7Jsonapi2(out *jwriter.Writer, in EmbedFooter) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"text\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Text))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v EmbedFooter) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson4086215fEncodeGithubComGsmcwhirterDiscordBotLibV7Jsonapi2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v EmbedFooter) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson4086215fEncodeGithubComGsmcwhirterDiscordBotLibV7Jsonapi2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *EmbedFooter) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson4086215fDecodeGithubComGsmcwhirterDiscordBotLibV7Jsonapi2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *EmbedFooter) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson4086215fDecodeGithubComGsmcwhirterDiscordBotLibV7Jsonapi2(l, v)
}
func easyjson4086215fDecodeGithubComGsmcwhirterDiscordBotLibV7Jsonapi3(in *jlexer.Lexer, out *EmbedField) {
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
		case "name":
			out.Name = string(in.String())
		case "value":
			out.Value = string(in.String())
		case "inline":
			out.Inline = bool(in.Bool())
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
func easyjson4086215fEncodeGithubComGsmcwhirterDiscordBotLibV7Jsonapi3(out *jwriter.Writer, in EmbedField) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"value\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Value))
	}
	{
		const prefix string = ",\"inline\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.Inline))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v EmbedField) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson4086215fEncodeGithubComGsmcwhirterDiscordBotLibV7Jsonapi3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v EmbedField) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson4086215fEncodeGithubComGsmcwhirterDiscordBotLibV7Jsonapi3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *EmbedField) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson4086215fDecodeGithubComGsmcwhirterDiscordBotLibV7Jsonapi3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *EmbedField) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson4086215fDecodeGithubComGsmcwhirterDiscordBotLibV7Jsonapi3(l, v)
}
func easyjson4086215fDecodeGithubComGsmcwhirterDiscordBotLibV7Jsonapi4(in *jlexer.Lexer, out *Embed) {
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
		case "title":
			out.Title = string(in.String())
		case "description":
			out.Description = string(in.String())
		case "url":
			out.URL = string(in.String())
		case "timestamp":
			out.Timestamp = string(in.String())
		case "color":
			out.Color = int(in.Int())
		case "fields":
			if in.IsNull() {
				in.Skip()
				out.Fields = nil
			} else {
				in.Delim('[')
				if out.Fields == nil {
					if !in.IsDelim(']') {
						out.Fields = make([]EmbedField, 0, 1)
					} else {
						out.Fields = []EmbedField{}
					}
				} else {
					out.Fields = (out.Fields)[:0]
				}
				for !in.IsDelim(']') {
					var v1 EmbedField
					(v1).UnmarshalEasyJSON(in)
					out.Fields = append(out.Fields, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "footer":
			(out.Footer).UnmarshalEasyJSON(in)
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
func easyjson4086215fEncodeGithubComGsmcwhirterDiscordBotLibV7Jsonapi4(out *jwriter.Writer, in Embed) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Title != "" {
		const prefix string = ",\"title\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Title))
	}
	if in.Description != "" {
		const prefix string = ",\"description\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"url\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.URL))
	}
	if in.Timestamp != "" {
		const prefix string = ",\"timestamp\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Timestamp))
	}
	if in.Color != 0 {
		const prefix string = ",\"color\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Color))
	}
	{
		const prefix string = ",\"fields\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Fields == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Fields {
				if v2 > 0 {
					out.RawByte(',')
				}
				(v3).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	if true {
		const prefix string = ",\"footer\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.Footer).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Embed) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson4086215fEncodeGithubComGsmcwhirterDiscordBotLibV7Jsonapi4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Embed) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson4086215fEncodeGithubComGsmcwhirterDiscordBotLibV7Jsonapi4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Embed) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson4086215fDecodeGithubComGsmcwhirterDiscordBotLibV7Jsonapi4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Embed) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson4086215fDecodeGithubComGsmcwhirterDiscordBotLibV7Jsonapi4(l, v)
}
