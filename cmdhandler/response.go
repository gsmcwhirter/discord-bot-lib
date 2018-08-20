package cmdhandler

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gsmcwhirter/discord-bot-lib/discordapi/jsonapi"
	"github.com/gsmcwhirter/discord-bot-lib/snowflake"
)

// Response is the interface that should be returned from a command handler
type Response interface {
	SetColor(int)
	IncludeError(err error)
	HasErrors() bool
	ToString() string
	ToMessage() json.Marshaler
	Channel() snowflake.Snowflake
}

// SimpleResponse is a Response that is intended to present plain text
type SimpleResponse struct {
	To        string
	Content   string
	ToChannel snowflake.Snowflake

	errors []error
}

// SetColor is included for the Response API but is a no-op
func (r *SimpleResponse) SetColor(color int) {}

// Channel returns the ToChannel value
func (r *SimpleResponse) Channel() snowflake.Snowflake {
	return r.ToChannel
}

// IncludeError adds an error into the response
func (r *SimpleResponse) IncludeError(err error) {
	if err == nil {
		return
	}

	r.errors = append(r.errors, err)
}

// HasErrors returns whether or not the response includes errors
func (r *SimpleResponse) HasErrors() bool {
	return len(r.errors) > 0
}

// ToString generates a plain-text representation of the response
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

// ToMessage generates an object that can be marshaled as json and sent to
// the discord http API
func (r *SimpleResponse) ToMessage() json.Marshaler {
	return jsonapi.Message{
		Content: r.ToString(),
		Tts:     false,
	}
}

// SimpleEmbedResponse is a Response that is intended to present
// text in an discord embed box but not include any embed fields
type SimpleEmbedResponse struct {
	To          string
	Title       string
	Description string
	Color       int
	FooterText  string
	ToChannel   snowflake.Snowflake

	errors []error
}

// SetColor sets the side color of the embed box
func (r *SimpleEmbedResponse) SetColor(color int) {
	r.Color = color
}

// Channel returns the ToChannel value
func (r *SimpleEmbedResponse) Channel() snowflake.Snowflake {
	return r.ToChannel
}

// IncludeError adds an error into the response
func (r *SimpleEmbedResponse) IncludeError(err error) {
	if err == nil {
		return
	}

	r.errors = append(r.errors, err)
}

// HasErrors returns whether or not the response already includes errors
func (r *SimpleEmbedResponse) HasErrors() bool {
	return len(r.errors) > 0
}

// ToString generates a plain-text representation of the response
func (r *SimpleEmbedResponse) ToString() string {
	b := strings.Builder{}

	_, _ = b.WriteString(fmt.Sprintf("%s\n\n", r.To))

	if r.Title != "" {
		_, _ = b.WriteString(fmt.Sprintf("__**%s**__\n", r.Title))
	}

	if r.Description != "" {
		_, _ = b.WriteString(fmt.Sprintf("%s\n", r.Description))
	}

	if len(r.errors) > 0 {
		for _, err := range r.errors {
			_, _ = b.WriteString(fmt.Sprintf("\nError: %v", err))
		}
		_, _ = b.WriteString("\n")
	}

	return b.String()
}

// ToMessage generates an object that can be marshaled as json and sent to
// the discord http API
func (r *SimpleEmbedResponse) ToMessage() json.Marshaler {
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

	if len(r.errors) > 0 {
		for _, err := range r.errors {
			m.Embed.Description += fmt.Sprintf("\nError: %v", err)
		}
		m.Embed.Description += "\n"
	}

	return m
}

// EmbedField is part of an EmbedResponse that represents
// an embed field
type EmbedField struct {
	Name string
	Val  string
}

// EmbedResponse is a Response that is intended to present
// text in an discord embed box, including embed fields
type EmbedResponse struct {
	To          string
	Title       string
	Description string
	Color       int
	Fields      []EmbedField
	FooterText  string
	ToChannel   snowflake.Snowflake

	errors []error
}

// SetColor sets the side color of the embed box
func (r *EmbedResponse) SetColor(color int) {
	r.Color = color
}

// Channel returns the ToChannel value
func (r *EmbedResponse) Channel() snowflake.Snowflake {
	return r.ToChannel
}

// IncludeError adds an error into the response
func (r *EmbedResponse) IncludeError(err error) {
	if err == nil {
		return
	}

	r.errors = append(r.errors, err)
}

// HasErrors returns whether or not the response already includes errors
func (r *EmbedResponse) HasErrors() bool {
	return len(r.errors) > 0
}

// ToString generates a plain-text representation of the response
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
		_, _ = b.WriteString(fmt.Sprintf("%s:\n%s\n", ef.Name, ef.Val))
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

// ToMessage generates an object that can be marshaled as json and sent to
// the discord http API
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
			m.Embed.Description += fmt.Sprintf("\nError: %v", err)
		}
		m.Embed.Description += "\n"
	}

	return m
}
