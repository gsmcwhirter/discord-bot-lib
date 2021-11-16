package jsonapi

import (
	stdjson "encoding/json" //nolint:depguard // we need this for RawMessage

	"github.com/gsmcwhirter/go-util/v8/json"
)

type MessageReference struct {
	MessageID string `json:"message_id,omitempty"`
	ChannelID string `json:"channel_id,omitempty"`
	GuildID   string `json:"guild_id,omitempty"`
}

// Message is the json object that is sent to the discord api
// to post a plain-text message to a server
type Message struct {
	Content string            `json:"content"`
	Tts     bool              `json:"tts"`
	ReplyTo *MessageReference `json:"message_reference,omitempty"`
}

func (m Message) MarshalToJSON() ([]byte, error) {
	return json.MarshalToBuffer(m)
}

// MessageWithEmbed is the json object that is sent to the discord api
// to post an embed message to a server
type MessageWithEmbed struct {
	Content string            `json:"content"`
	Tts     bool              `json:"tts"`
	Embeds  []Embed           `json:"embeds"`
	ReplyTo *MessageReference `json:"message_reference,omitempty"`
}

func (m MessageWithEmbed) MarshalToJSON() ([]byte, error) {
	return json.MarshalToBuffer(m)
}

// Embed is a json object that represents an embed in a MessageWithEmbed
type Embed struct {
	Title       string       `json:"title,omitempty"`
	Description string       `json:"description,omitempty"`
	URL         string       `json:"url,omitempty"`
	Timestamp   string       `json:"timestamp,omitempty"`
	Color       int          `json:"color,omitempty"`
	Fields      []EmbedField `json:"fields,omitempty"`
	Footer      EmbedFooter  `json:"footer,omitempty"`
}

// EmbedField is a json object that represents a field in an Embed
type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

// EmbedFooter is a json object that represents the footer of an embed
type EmbedFooter struct {
	Text string `json:"text"`
}

type InteractionCallbackType int

const (
	CallbackTypePong                   InteractionCallbackType = 1
	CallbackTypeChannelMessage         InteractionCallbackType = 4
	CallbackTypeDeferredChannelMessage InteractionCallbackType = 5
	CallbackTypeDeferredUpdate         InteractionCallbackType = 6
	CallbackTypeUpdate                 InteractionCallbackType = 7
	CallbackTypeAutocomplete           InteractionCallbackType = 8
)

type InteractionCallbackMessage struct {
	Type InteractionCallbackType `json:"type"`
	Data stdjson.RawMessage      `json:"data"`
}
