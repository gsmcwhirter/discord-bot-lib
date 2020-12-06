package jsonapi

import "github.com/gsmcwhirter/discord-bot-lib/v13/snowflake"

//go:generate easyjson -all

// MessageResponse is the data that is received back from the discord api
//easyjson:json
type MessageResponse struct {
	ID              snowflake.Snowflake      `json:"id"`
	ChannelID       snowflake.Snowflake      `json:"channel_id"`
	GuildID         snowflake.Snowflake      `json:"guild_id"`
	Author          UserResponse             `json:"author"`
	Member          GuildMemberResponse      `json:"member"`
	Content         string                   `json:"content"`
	Timestamp       string                   `json:"timestamp"`        // ISO8601
	EditedTimestamp string                   `json:"edited_timestamp"` // ISO8601
	TTS             bool                     `json:"tts"`
	MentionEveryone bool                     `json:"mention_everyone"`
	Mentions        []UserResponse           `json:"mentions"`
	MentionRoles    []RoleResponse           `json:"mention_roles"`
	MentionChannels []ChannelMentionResponse `json:"mention_channels"`
	Attachments     []AttachmentResponse     `json:"attachments"`
	Embeds          []EmbedResponse          `json:"embeds"`
	Reactions       []ReactionResponse       `json:"reactions"`
	Pinned          bool                     `json:"pinned"`
	WebhookID       snowflake.Snowflake      `json:"webhook_id"`
	Type            int                      `json:"type"`
	Flags           int                      `json:"flags"`

	// Nonce is skipped
	// Activity is skipped
	// Application is skipped
	// MessageReference is skipped
}

// RoleResponse is the data about a role recevied from the json api
//easyjson:json
type RoleResponse struct {
	ID          snowflake.Snowflake `json:"id"`
	Name        string              `json:"name"`
	Color       int                 `json:"color"`
	Hoist       bool                `json:"hoist"`
	Position    int                 `json:"position"`
	Permissions int                 `json:"permissions"`
	Managed     bool                `json:"managed"`
	Mentionable bool                `json:"mentionable"`
}

// ChannelMentionResponse is the data about a channel mention recevied from the json api
//easyjson:json
type ChannelMentionResponse struct {
	ID      snowflake.Snowflake `json:"id"`
	GuildID snowflake.Snowflake `json:"guild_id"`
	Type    int                 `json:"type"`
	Name    string              `json:"name"`
}

// AttachmentResponse is the data about an attachment recevied from the json api
//easyjson:json
type AttachmentResponse struct {
	ID       snowflake.Snowflake `json:"id"`
	Filename string              `json:"filename"`
	Size     int                 `json:"size"`
	URL      string              `json:"url"`
	ProxyURL string              `json:"proxy_url"`
	Height   int                 `json:"height"`
	Width    int                 `json:"width"`
}

// ReactionResponse is the data about a reaction received from the json api
//easyjson:json
type ReactionResponse struct {
	Count int           `json:"count"`
	Me    bool          `json:"me"`
	Emoji EmojiResponse `json:"emoji"`
}

// EmojiResponse is the data about an emoji recevied from the json api
//easyjson:json
type EmojiResponse struct {
	ID            snowflake.Snowflake `json:"id"`
	Name          string              `json:"name"`
	Roles         []RoleResponse      `json:"roles"`
	User          UserResponse        `json:"user"`
	RequireColons bool                `json:"require_colons"`
	Managed       bool                `json:"managed"`
	Animated      bool                `json:"animated"`
	Available     bool                `json:"available"`
}

// EmbedResponse is the data about a message embed received from the json api
//easyjson:json
type EmbedResponse struct {
	Title       string                `json:"title"`
	Type        string                `json:"type"`
	Description string                `json:"description"`
	URL         string                `json:"url"`
	Timestamp   string                `json:"timestamp"` //ISO8601
	Color       int                   `json:"color"`
	Footer      EmbedFooterResponse   `json:"footer"`
	Image       EmbedImageResponse    `json:"image"`
	Thumbnail   EmbedImageResponse    `json:"thumbnail"`
	Video       EmbedImageResponse    `json:"video"`
	Provider    EmbedProviderResponse `json:"provider"`
	Author      EmbedAuthorResponse   `json:"author"`
	Fields      []EmbedFieldResponse  `json:"fields"`
}

// EmbedFooterResponse is the data about an embed footer recevied from the json api
//easyjson:json
type EmbedFooterResponse struct {
	Text         string `json:"text"`
	IconURL      string `json:"icon_url"`
	ProxyIconURL string `json:"proxy_icon_url"`
}

// EmbedFieldResponse is the data about an embed field received from the json api
//easyjson:json
type EmbedFieldResponse struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

// EmbedImageResponse is the data about an embed thumbnail received from the json api
//easyjson:json
type EmbedImageResponse struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`
}

// EmbedProviderResponse is the data about an embed provider recevied from the json api
//easyjson:json
type EmbedProviderResponse struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// EmbedAuthorResponse is the data about an embed author recevied from the json api
//easyjson:json
type EmbedAuthorResponse struct {
	Name         string `json:"name"`
	URL          string `json:"url"`
	IconURL      string `json:"icon_url"`
	ProxyIconURL string `json:"proxy_icon_url"`
}
