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
	SetColor(int)
	IncludeError(err error)
	ToString() string
	ToMessage() json.Marshaler
}

// SimpleResponse TODOC
type SimpleResponse struct {
	To      string
	Content string

	errors []error
}

// SetColor TODOC
func (r *SimpleResponse) SetColor(color int) {}

// IncludeError TODOC
func (r *SimpleResponse) IncludeError(err error) {
	if err == nil {
		return
	}

	r.errors = append(r.errors, err)
}

// ToString TODOC
func (r *SimpleResponse) ToString() string {
	s := fmt.Sprintf(`%s

%s
`, r.To, r.Content)

	if len(r.errors) > 0 {
		for _, err := range r.errors {
			s += fmt.Sprintf("\nError: %v", err)
		}
		s += "\n"
	}

	return s
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

	errors []error
}

// SetColor TODOC
func (r *EmbedResponse) SetColor(color int) {
	r.Color = color
}

// IncludeError TODOC
func (r *EmbedResponse) IncludeError(err error) {
	if err == nil {
		return
	}

	r.errors = append(r.errors, err)
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

	if len(r.errors) > 0 {
		for _, err := range r.errors {
			_, _ = b.WriteString(fmt.Sprintf("\nError: %v", err))
		}
		_, _ = b.WriteString("\n")
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

	if len(r.errors) > 0 {
		for _, err := range r.errors {
			m.Content += fmt.Sprintf("\nError: %v", err)
		}
		m.Content += "\n"
	}

	return m
}
