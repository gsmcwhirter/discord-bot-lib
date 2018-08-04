package jsonapi

//go:generate easyjson -all

// Message TODOC
//easyjson:json
type Message struct {
	Content string
	Tts     bool
}

// MessageWithEmbed TODOC
//easyjson:json
type MessageWithEmbed struct {
	Content string
	Tts     bool
	Embed   Embed
}

// Embed TODOC
//easyjson:json
type Embed struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	URL         string `json:"url"`
	Timestamp   string `json:"timestamp,omitempty"`
	Color       int    `json:"color,omitempty"`
	Fields      []EmbedField
	Footer      EmbedFooter `json:"footer,omitempty"`
}

// EmbedField TODOC
//easyjson:json
type EmbedField struct {
	Name   string
	Value  string
	Inline bool
}

// EmbedFooter TODOC
//easyjson:json
type EmbedFooter struct {
	Text string
}
