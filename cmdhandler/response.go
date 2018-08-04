package cmdhandler

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gsmcwhirter/discord-bot-lib/discordapi/jsonapi"
)

// Response TODOC
type Response interface {
	ToString() string
	ToMessage() json.Marshaler
}

// SimpleResponse TODOC
type SimpleResponse struct {
	To      string
	Content string
}

// ToString TODOC
func (r *SimpleResponse) ToString() string {
	return fmt.Sprintf(`%s

%s
`, r.To, r.Content)
}

// ToMessage TODOC
func (r *SimpleResponse) ToMessage() json.Marshaler {
	return jsonapi.Message{
		Content: r.ToString(),
		Tts:     false,
	}
}

// EmbedField TODOC
type EmbedField struct {
	Name string
	Val  string
}

// EmbedResponse TODOC
type EmbedResponse struct {
	To          string
	Title       string
	Description string
	Color       int
	Fields      []EmbedField
	FooterText  string
}

// ToString TODOC
func (r *EmbedResponse) ToString() string {
	b := strings.Builder{}

	_, _ = b.WriteString(fmt.Sprintf("%s\n\n", r.To))

	if r.Title != "" {
		_, _ = b.WriteString(fmt.Sprintf("__**%s**__\n", r.Title))
	}

	if r.Description != "" {
		_, _ = b.WriteString(fmt.Sprintf("%s\n", r.Description))
	}

	for _, ef := range r.Fields {
		_, _ = b.WriteString(fmt.Sprintf("%s:\n```\n%s\n```\n", ef.Name, ef.Val))
	}

	if r.FooterText != "" {
		_, _ = b.WriteString(fmt.Sprintf("%s\n", r.FooterText))
	}

	return b.String()
}

// ToMessage TODOC
func (r *EmbedResponse) ToMessage() json.Marshaler {
	m := jsonapi.MessageWithEmbed{
		Content: fmt.Sprintf("%s\n", r.To),
		Tts:     false,
		Embed: jsonapi.Embed{
			Timestamp: time.Now().Format(time.RFC3339),
		},
	}

	if r.Title != "" {
		m.Embed.Title = r.Title
	}

	if r.Description != "" {
		m.Embed.Description = r.Description
	}

	if r.Color != 0 {
		m.Embed.Color = r.Color
	}

	if r.FooterText != "" {
		m.Embed.Footer.Text = r.FooterText
	}

	m.Embed.Fields = make([]jsonapi.EmbedField, 0, len(r.Fields))
	for _, ef := range r.Fields {
		m.Embed.Fields = append(m.Embed.Fields, jsonapi.EmbedField{
			Name:  ef.Name,
			Value: ef.Val,
		})
	}

	return m
}
