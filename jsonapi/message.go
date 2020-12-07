package jsonapi

//go:generate easyjson -all

type MessageReference struct {
	MessageID string `json:"message_id,omitempty"`
	ChannelID string `json:"channel_id,omitempty"`
	GuildID   string `json:"guild_id,omitempty"`
}

// Message is the json object that is sent to the discord api
// to post a plain-text message to a server
//easyjson:json
type Message struct {
	Content string           `json:"content"`
	Tts     bool             `json:"tts"`
	ReplyTo MessageReference `json:"message_reference,omitempty"`
}

// MessageWithEmbed is the json object that is sent to the discord api
// to post an embed message to a server
//easyjson:json
type MessageWithEmbed struct {
	Content string           `json:"content"`
	Tts     bool             `json:"tts"`
	Embed   Embed            `json:"embed"`
	ReplyTo MessageReference `json:"message_reference,omitempty"`
}

// Embed is a json object that represents an embed in a MessageWithEmbed
//easyjson:json
type Embed struct {
	Title       string       `json:"title,omitempty"`
	Description string       `json:"description,omitempty"`
	URL         string       `json:"url"`
	Timestamp   string       `json:"timestamp,omitempty"`
	Color       int          `json:"color,omitempty"`
	Fields      []EmbedField `json:"fields"`
	Footer      EmbedFooter  `json:"footer,omitempty"`
}

// EmbedField is a json object that represents a field in an Embed
//easyjson:json
type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

// EmbedFooter is a json object that represents the footer of an embed
//easyjson:json
type EmbedFooter struct {
	Text string `json:"text"`
}
