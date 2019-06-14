package cmdhandler

import (
	"fmt"
	"strings"
	"time"

	"github.com/gsmcwhirter/discord-bot-lib/v6/jsonapi"
	"github.com/gsmcwhirter/discord-bot-lib/v6/snowflake"
)

const maxLen = 1024
const maxEmbedLen = 5900
const ctn = "\n\n(continued...)"

// Response is the interface that should be returned from a command handler
type Response interface {
	SetColor(int)
	IncludeError(err error)
	HasErrors() bool
	ToString() string
	ToMessage() JSONMarshaler
	Channel() snowflake.Snowflake
	Split() []Response
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
func (r *SimpleResponse) ToMessage() JSONMarshaler {
	return jsonapi.Message{
		Content: r.ToString(),
		Tts:     false,
	}
}

// Split separates the current response into possibly-several to account for response length limits
func (r *SimpleResponse) Split() []Response {
	if len(r.ToString()) < maxLen {
		return []Response{r}
	}

	split := textSplit(r.ToString(), maxLen-len(ctn)-len(r.To)-4, "\n")

	resps := make([]Response, 0, len(split))
	for i, s := range split {
		if i < len(split)-1 {
			s += ctn
		}

		resps = append(resps, &SimpleResponse{
			To:        r.To,
			Content:   s,
			ToChannel: r.ToChannel,
		})
	}

	return resps
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
func (r *SimpleEmbedResponse) ToMessage() JSONMarshaler {
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

// Split separates the current response into possibly-several to account for response length limits
func (r *SimpleEmbedResponse) Split() []Response {
	if len(r.ToString()) < maxLen {
		return []Response{r}
	}

	split := textSplit(r.ToString(), maxLen-len(ctn)-len(r.To)-4, "\n")

	resps := make([]Response, 0, len(split))
	for i, s := range split {
		title := ""
		if i == 0 {
			title = r.Title
		}

		if i < len(split)-1 {
			s += ctn
		}

		resps = append(resps, &SimpleEmbedResponse{
			To:          r.To,
			Title:       title,
			Description: s,
			Color:       r.Color,
			FooterText:  r.FooterText,
			ToChannel:   r.ToChannel,
		})
	}

	return resps
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
func (r *EmbedResponse) ToMessage() JSONMarshaler {
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

// Split separates the current response into possibly-several to account for response length limits
func (r *EmbedResponse) Split() []Response {

	if len(r.ToString()) < maxLen {
		return []Response{r}
	}

	resps := make([]Response, 0, 2)

	// prepare messages for just the description
	descSplit := textSplit(r.Description, maxLen-len(ctn)-len(r.To)-4, "\n")
	// fmt.Printf("descSplit: %#v\n", descSplit)

	title := r.Title

	if len(descSplit) > 1 {
		for i, ds := range descSplit[:len(descSplit)-1] {
			if i > 0 {
				title = ""
			}

			resps = append(resps, &SimpleEmbedResponse{
				To:          r.To,
				Title:       title,
				Description: ds + ctn,
				Color:       r.Color,
				FooterText:  r.FooterText,
				ToChannel:   r.ToChannel,
			})
		}
	}

	// this is the first message that will contain field content
	desc := descSplit[len(descSplit)-1]
	resp := &EmbedResponse{
		To:          r.To,
		Title:       title,
		Description: desc,
		Color:       r.Color,
		FooterText:  r.FooterText,
		ToChannel:   r.ToChannel,
	}
	descRespLen := len(desc) + len(r.FooterText)

	// find out how many fields will fit on the last description message
	nextField, nextFieldSplits := r.fillResp(resp, 0, descRespLen)

	// last of the description is prepared
	resps = append(resps, resp)

	// do we need to continue an unfinished field?
	resps, resp = r.fillResps(resps, nextField, nextFieldSplits)
	nextField++

	// TODO: Go through all the remaining fields
	for nextField < len(r.Fields) {
		// find out how many fields will fit on the last description message
		nextField, nextFieldSplits = r.fillResp(resp, nextField, 0)

		// last of the description is prepared
		resps = append(resps, resp)

		// do we need to continue an unfinished field?
		resps, resp = r.fillResps(resps, nextField, nextFieldSplits)
		nextField++
	}

	return resps
}

func (r *EmbedResponse) fillResp(resp *EmbedResponse, startField, existingLen int) (int, []string) {
	var nextField int
	var nextFieldSplits []string

	// find out how many fields will fit on the last description message
	for i, f := range r.Fields[startField:] {
		nextField = i + startField

		split := textSplit(f.Val, maxLen-len(ctn), "\n")
		sliceStart := 0

		if len(split[0])+existingLen < maxEmbedLen {
			sliceStart = 1

			if len(split) > 1 {
				split[0] += ctn
			}
			resp.Fields = append(resp.Fields, EmbedField{
				Name: f.Name,
				Val:  split[0],
			})
		}

		nextFieldSplits = split[sliceStart:]

		if sliceStart == 0 || len(split) > 1 {
			break
		}
	}

	return nextField, nextFieldSplits
}

func (r *EmbedResponse) fillResps(resps []Response, nextField int, nextFieldSplits []string) ([]Response, *EmbedResponse) {
	newResps := resps
	var resp *EmbedResponse

	// do we need to continue an unfinished field?
	for _, s := range nextFieldSplits {
		resp = &EmbedResponse{
			To:          r.To,
			Title:       "",
			Description: "",
			Color:       r.Color,
			FooterText:  r.FooterText,
			ToChannel:   r.ToChannel,
		}

		resp.Fields = append(resp.Fields, EmbedField{
			Name: r.Fields[nextField].Name,
			Val:  s,
		})

		newResps = append(newResps, resp)
	}

	return newResps, resp
}

// func (r *EmbedResponse) fillResp(resps []Response, resp *EmbedResponse, nextField int, nextFieldSplits []string, existingLen int) (resps []Response, nextResp *EmbedResponse, newNext int, newNextSplits []string) {
// 	// do we need to continue an unfinished field?
// 	for i, s := range nextFieldSplits {
// 		resp.Fields = append(resp.Fields, EmbedField{
// 			Name: r.Fields[nextField].Name,
// 			Val:  s,
// 		})

// 		resps = append(resps, resp)

// 		if i < len(nextFieldSplits)-1 {
// 			resp = &EmbedResponse{
// 				To:          r.To,
// 				Title:       "",
// 				Description: "",
// 				Color:       r.Color,
// 				FooterText:  r.FooterText,
// 				ToChannel:   r.ToChannel,
// 			}
// 		}
// 	}

// 	for i, f := range r.Fields[nextField:] {
// 		newNext = i + nextField

// 		split := textSplit(f.Val, maxLen-len(ctn), "\n")
// 		sliceStart := 0

// 		if len(split[0])+existingLen < maxEmbedLen {
// 			sliceStart = 1

// 			if len(split) > 1 {
// 				split[0] += ctn
// 			}
// 			resp.Fields = append(resp.Fields, EmbedField{
// 				Name: f.Name,
// 				Val:  split[0],
// 			})
// 		}

// 		newNextSplits = split[sliceStart:]

// 		if sliceStart == 0 || len(split) > 1 {
// 			break
// 		}
// 	}

// 	return
// }
