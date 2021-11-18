package entity

import (
	"fmt"

	"github.com/gsmcwhirter/go-util/v8/errors"

	"github.com/gsmcwhirter/discord-bot-lib/v23/discordapi/etfapi"
	"github.com/gsmcwhirter/discord-bot-lib/v23/snowflake"
)

// MessageType represents the type of message received in a discord channel
type MessageType int

// These are the known message types
const (
	DefaultMessage              MessageType = 0
	RecipientAddMessage         MessageType = 1
	RecipientRemoveMessage      MessageType = 2
	CallMessage                 MessageType = 3
	ChannelNameChangeMessage    MessageType = 4
	ChannelIconChangeMessage    MessageType = 5
	ChannelPinnedMessageMessage MessageType = 6
	GuildMemberJoinMessage      MessageType = 7
)

// MessageTypeFromElement generates a MessageType representation from the given
// message-type Element
func MessageTypeFromElement(e etfapi.Element) (MessageType, error) {
	temp, err := e.ToInt()
	t := MessageType(temp)
	return t, errors.Wrap(err, "could not unmarshal MessageType")
}

func (t MessageType) String() string {
	switch t {
	case DefaultMessage:
		return "DEFAULT"
	case RecipientAddMessage:
		return "RECIPIENT_ADD"
	case RecipientRemoveMessage:
		return "RECIPIENT_REMOVE"
	case CallMessage:
		return "CALL"
	case ChannelNameChangeMessage:
		return "CHANNEL_NAME_CHANGE"
	case ChannelIconChangeMessage:
		return "CHANNEL_ICON_CHANGE"
	case ChannelPinnedMessageMessage:
		return "CHANNEL_PINNED_MESSAGE"
	case GuildMemberJoinMessage:
		return "GUILD_MEMBER_JOIN"
	default:
		return fmt.Sprintf("(unknown: %d)", int(t))
	}
}

// Message is the data that is received back from the discord api
type Message struct {
	IDString            string            `json:"id"`
	ChannelIDString     string            `json:"channel_id"`
	GuildIDString       string            `json:"guild_id"`
	Author              User              `json:"author"`
	Member              GuildMember       `json:"member"`
	Content             string            `json:"content"`
	Timestamp           string            `json:"timestamp"`        // ISO8601
	EditedTimestamp     string            `json:"edited_timestamp"` // ISO8601
	TTS                 bool              `json:"tts"`
	MentionEveryone     bool              `json:"mention_everyone"`
	Mentions            []User            `json:"mentions"`
	MentionRolesStrings []string          `json:"mention_roles"`
	MentionChannels     []Channel         `json:"mention_channels"`
	Attachments         []Attachment      `json:"attachments"`
	Embeds              []Embed           `json:"embeds"`
	Reactions           []MessageReaction `json:"reactions"`
	Pinned              bool              `json:"pinned"`
	WebhookID           string            `json:"webhook_id"`
	Type                MessageType       `json:"type"`
	Flags               int               `json:"flags"`

	// Nonce is skipped
	// Activity is skipped
	// Application is skipped
	// MessageReference is skipped

	IDSnowflake        snowflake.Snowflake `json:"-"`
	ChannelIDSnowflake snowflake.Snowflake `json:"-"`
	GuildIDSnowflake   snowflake.Snowflake `json:"-"`
	WebhookIDSnowflake snowflake.Snowflake `json:"-"`

	MentionRoles []snowflake.Snowflake `json:"-"`
}

func (m *Message) Snowflakify() error {
	var err error

	if m.IDSnowflake, err = snowflake.FromString(m.IDString); err != nil {
		return errors.Wrap(err, "could not snowflakify ID")
	}

	if m.ChannelIDString != "" {
		if m.ChannelIDSnowflake, err = snowflake.FromString(m.ChannelIDString); err != nil {
			return errors.Wrap(err, "could not snowflakify ChannelID")
		}
	}

	if m.GuildIDString != "" {
		if m.GuildIDSnowflake, err = snowflake.FromString(m.GuildIDString); err != nil {
			return errors.Wrap(err, "could not snowflakify GuildID")
		}
	}

	if m.WebhookID != "" {
		if m.WebhookIDSnowflake, err = snowflake.FromString(m.WebhookID); err != nil {
			return errors.Wrap(err, "could not snowflakify WebhookID")
		}
	}

	if err = m.Author.Snowflakify(); err != nil {
		return errors.Wrap(err, "could not snowflakify Author")
	}

	if err = m.Member.Snowflakify(); err != nil {
		return errors.Wrap(err, "could not snowflakify Member")
	}

	for i := range m.Mentions {
		mn := m.Mentions[i]
		if err = mn.Snowflakify(); err != nil {
			return errors.Wrap(err, "could not snowflakify Mentions")
		}
		m.Mentions[i] = mn
	}

	m.MentionRoles = make([]snowflake.Snowflake, len(m.MentionRolesStrings))
	for i := range m.MentionRolesStrings {
		m.MentionRoles[i], err = snowflake.FromString(m.MentionRolesStrings[i])
		if err != nil {
			return errors.Wrap(err, "could not snowflakify MentionRoles")
		}
	}

	for i := range m.MentionChannels {
		mc := m.MentionChannels[i]
		if err = mc.Snowflakify(); err != nil {
			return errors.Wrap(err, "could not snowflakify MentionChannels")
		}
		m.MentionChannels[i] = mc
	}

	for i := range m.Attachments {
		ma := m.Attachments[i]
		if err = ma.Snowflakify(); err != nil {
			return errors.Wrap(err, "could not snowflakify Attachments")
		}
		m.Attachments[i] = ma
	}

	// for i := range m.Embeds {
	// 	m := m.Embeds[i]
	// 	if err = m.Snowflakify(); err != nil {
	// 		return errors.Wrap(err, "could not snowflakify Embeds")
	// 	}
	// 	m.Embeds[i] = m
	// }

	for i := range m.Reactions {
		mr := m.Reactions[i]
		if err = mr.Snowflakify(); err != nil {
			return errors.Wrap(err, "could not snowflakify Reactions")
		}
		m.Reactions[i] = mr
	}

	return nil
}

// ID returns the ID of the message
func (m *Message) ID() snowflake.Snowflake {
	return m.IDSnowflake
}

// ChannelID returns the ID of the channel the message was sent to
func (m *Message) ChannelID() snowflake.Snowflake {
	return m.ChannelIDSnowflake
}

// MessageType returns the MessageType of the message
func (m *Message) MessageType() MessageType {
	return m.Type
}

// AuthorID returns the ID of the author of the message
func (m *Message) AuthorID() snowflake.Snowflake {
	return m.Author.IDSnowflake
}

// ContentString returns the content of the message
func (m *Message) ContentString() string {
	return m.Content
}

// MessageFromElementMap generates a new Message object from the given data
func MessageFromElementMap(eMap map[string]etfapi.Element) (Message, error) {
	var m Message
	var err error

	e2, ok := eMap["id"]
	if ok {
		m.IDSnowflake, err = etfapi.SnowflakeFromUnknownElement(e2)
		if err != nil {
			return m, errors.Wrap(err, "could not get message_id snowflake.Snowflake")
		}

		m.IDString = m.IDSnowflake.ToString()
	}

	e2, ok = eMap["channel_id"]
	if ok && !e2.IsNil() {
		m.ChannelIDSnowflake, err = etfapi.SnowflakeFromUnknownElement(e2)
		if err != nil {
			return m, errors.Wrap(err, "could not get channel_id snowflake.Snowflake")
		}
		m.ChannelIDString = m.ChannelIDSnowflake.ToString()
	}

	m.Type, err = MessageTypeFromElement(eMap["type"])
	if err != nil {
		return m, errors.Wrap(err, "could not get messageType")
	}

	e2, ok = eMap["content"]
	if ok {
		m.Content, err = e2.ToString()
		if err != nil {
			return m, errors.Wrap(err, "could not get content")
		}
	}

	m.Author, err = UserFromElement(eMap["author"])
	if err != nil {
		return m, errors.Wrap(err, "could not inflate message author")
	}

	return m, nil
}

// MessageFromElement generates a new Message object from the given Element
func MessageFromElement(e etfapi.Element) (Message, error) {
	eMap, _, err := etfapi.MapAndIDFromElement(e)
	if err != nil {
		return Message{}, err
	}

	m, err := MessageFromElementMap(eMap)

	return m, err
}
