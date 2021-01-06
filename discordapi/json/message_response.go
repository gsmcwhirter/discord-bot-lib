package json

import (
	"github.com/gsmcwhirter/go-util/v8/errors"

	"github.com/gsmcwhirter/discord-bot-lib/v19/snowflake"
)

// MessageResponse is the data that is received back from the discord api
type MessageResponse struct {
	ID              string                   `json:"id"`
	ChannelID       string                   `json:"channel_id"`
	GuildID         string                   `json:"guild_id"`
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
	WebhookID       string                   `json:"webhook_id"`
	Type            int                      `json:"type"`
	Flags           int                      `json:"flags"`

	// Nonce is skipped
	// Activity is skipped
	// Application is skipped
	// MessageReference is skipped

	IDSnowflake        snowflake.Snowflake
	ChannelIDSnowflake snowflake.Snowflake
	GuildIDSnowflake   snowflake.Snowflake
	WebhookIDSnowflake snowflake.Snowflake
}

func (mr *MessageResponse) Snowflakify() error {
	var err error

	if mr.IDSnowflake, err = snowflake.FromString(mr.ID); err != nil {
		return errors.Wrap(err, "could not snowflakify ID")
	}

	if mr.ChannelID != "" {
		if mr.ChannelIDSnowflake, err = snowflake.FromString(mr.ChannelID); err != nil {
			return errors.Wrap(err, "could not snowflakify ChannelID")
		}
	}

	if mr.GuildID != "" {
		if mr.GuildIDSnowflake, err = snowflake.FromString(mr.GuildID); err != nil {
			return errors.Wrap(err, "could not snowflakify GuildID")
		}
	}

	if mr.WebhookID != "" {
		if mr.WebhookIDSnowflake, err = snowflake.FromString(mr.WebhookID); err != nil {
			return errors.Wrap(err, "could not snowflakify WebhookID")
		}
	}

	if err = mr.Author.Snowflakify(); err != nil {
		return errors.Wrap(err, "could not snowflakify Author")
	}

	if err = mr.Member.Snowflakify(); err != nil {
		return errors.Wrap(err, "could not snowflakify Member")
	}

	for i := range mr.Mentions {
		m := mr.Mentions[i]
		if err = m.Snowflakify(); err != nil {
			return errors.Wrap(err, "could not snowflakify Mentions")
		}
		mr.Mentions[i] = m
	}

	for i := range mr.MentionRoles {
		m := mr.MentionRoles[i]
		if err = m.Snowflakify(); err != nil {
			return errors.Wrap(err, "could not snowflakify MentionRoles")
		}
		mr.MentionRoles[i] = m
	}

	for i := range mr.MentionChannels {
		m := mr.MentionChannels[i]
		if err = m.Snowflakify(); err != nil {
			return errors.Wrap(err, "could not snowflakify MentionChannels")
		}
		mr.MentionChannels[i] = m
	}

	for i := range mr.Attachments {
		m := mr.Attachments[i]
		if err = m.Snowflakify(); err != nil {
			return errors.Wrap(err, "could not snowflakify Attachments")
		}
		mr.Attachments[i] = m
	}

	// for i := range mr.Embeds {
	// 	m := mr.Embeds[i]
	// 	if err = m.Snowflakify(); err != nil {
	// 		return errors.Wrap(err, "could not snowflakify Embeds")
	// 	}
	// 	mr.Embeds[i] = m
	// }

	for i := range mr.Reactions {
		m := mr.Reactions[i]
		if err = m.Snowflakify(); err != nil {
			return errors.Wrap(err, "could not snowflakify Reactions")
		}
		mr.Reactions[i] = m
	}

	return nil
}

// RoleResponse is the data about a role recevied from the json api
type RoleResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Color       int    `json:"color"`
	Hoist       bool   `json:"hoist"`
	Position    int    `json:"position"`
	Permissions int    `json:"permissions"`
	Managed     bool   `json:"managed"`
	Mentionable bool   `json:"mentionable"`

	IDSnowflake snowflake.Snowflake
}

func (rr *RoleResponse) Snowflakify() error {
	var err error

	if rr.IDSnowflake, err = snowflake.FromString(rr.ID); err != nil {
		return errors.Wrap(err, "could not snowflakify ID")
	}

	return nil
}

// ChannelMentionResponse is the data about a channel mention recevied from the json api
type ChannelMentionResponse struct {
	ID      string `json:"id"`
	GuildID string `json:"guild_id"`
	Type    int    `json:"type"`
	Name    string `json:"name"`

	IDSnowflake      snowflake.Snowflake
	GuildIDSnowflake snowflake.Snowflake
}

func (cmr *ChannelMentionResponse) Snowflakify() error {
	var err error

	if cmr.IDSnowflake, err = snowflake.FromString(cmr.ID); err != nil {
		return errors.Wrap(err, "could not snowflakify ID")
	}

	if cmr.GuildIDSnowflake, err = snowflake.FromString(cmr.GuildID); err != nil {
		return errors.Wrap(err, "could not snowflakify GuildID")
	}

	return nil
}

// AttachmentResponse is the data about an attachment recevied from the json api
type AttachmentResponse struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	Size     int    `json:"size"`
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`

	IDSnowflake snowflake.Snowflake
}

func (ar *AttachmentResponse) Snowflakify() error {
	var err error
	if ar.IDSnowflake, err = snowflake.FromString(ar.ID); err != nil {
		return errors.Wrap(err, "could not snowflakify ID")
	}

	return nil
}

// ReactionResponse is the data about a reaction received from the json api
type ReactionResponse struct {
	Count int           `json:"count"`
	Me    bool          `json:"me"`
	Emoji EmojiResponse `json:"emoji"`
}

func (rr *ReactionResponse) Snowflakify() error {
	if err := rr.Emoji.Snowflakify(); err != nil {
		return errors.Wrap(err, "could not snowflakify Emoji")
	}

	return nil
}

// EmojiResponse is the data about an emoji recevied from the json api
type EmojiResponse struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	Roles         []RoleResponse `json:"roles"`
	User          UserResponse   `json:"user"`
	RequireColons bool           `json:"require_colons"`
	Managed       bool           `json:"managed"`
	Animated      bool           `json:"animated"`
	Available     bool           `json:"available"`

	IDSnowflake snowflake.Snowflake
}

func (er *EmojiResponse) Snowflakify() error {
	var err error

	if er.ID != "" {
		if er.IDSnowflake, err = snowflake.FromString(er.ID); err != nil {
			return errors.Wrap(err, "could not snowflakify ID")
		}
	}

	if err = er.User.Snowflakify(); err != nil {
		return errors.Wrap(err, "could not snowflakify User")
	}

	for i := range er.Roles {
		m := er.Roles[i]
		if err = m.Snowflakify(); err != nil {
			return errors.Wrap(err, "could not snowflakify Roles")
		}
		er.Roles[i] = m
	}

	return nil
}

// EmbedResponse is the data about a message embed received from the json api
type EmbedResponse struct {
	Title       string                `json:"title"`
	Type        string                `json:"type"`
	Description string                `json:"description"`
	URL         string                `json:"url"`
	Timestamp   string                `json:"timestamp"` // ISO8601
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
type EmbedFooterResponse struct {
	Text         string `json:"text"`
	IconURL      string `json:"icon_url"`
	ProxyIconURL string `json:"proxy_icon_url"`
}

// EmbedFieldResponse is the data about an embed field received from the json api
type EmbedFieldResponse struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

// EmbedImageResponse is the data about an embed thumbnail received from the json api
type EmbedImageResponse struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`
}

// EmbedProviderResponse is the data about an embed provider recevied from the json api
type EmbedProviderResponse struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// EmbedAuthorResponse is the data about an embed author recevied from the json api
type EmbedAuthorResponse struct {
	Name         string `json:"name"`
	URL          string `json:"url"`
	IconURL      string `json:"icon_url"`
	ProxyIconURL string `json:"proxy_icon_url"`
}
