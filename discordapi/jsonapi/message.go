package jsonapi

//go:generate easyjson -all

// Message TODOC
//easyjson:json
type Message struct {
	Content string `json:"content"`
	Tts     bool   `json:"tts"`
}

// MessageWithEmbed TODOC
//easyjson:json
type MessageWithEmbed struct {
	Content string `json:"content"`
	Tts     bool   `json:"tts"`
	Embed   Embed  `json:"embed"`
}

// Embed TODOC
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

// EmbedField TODOC
//easyjson:json
type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

// EmbedFooter TODOC
//easyjson:json
type EmbedFooter struct {
	Text string `json:"text"`
}
