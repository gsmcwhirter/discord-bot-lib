package etfapi

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/gsmcwhirter/discord-bot-lib/snowflake"
)

// MessageType represents the type of message received in a discord channel
type MessageType int

// These are the known message types
const (
	DefaultMessage              MessageType = 0
	RecipientAddMessage                     = 1
	RecipientRemoveMessage                  = 2
	CallMessage                             = 3
	ChannelNameChangeMessage                = 4
	ChannelIconChangeMessage                = 5
	ChannelPinnedMessageMessage             = 6
	GuildMemberJoinMessage                  = 7
)

// MessageTypeFromElement generates a MessageType representation from the given
// message-type Element
func MessageTypeFromElement(e Element) (t MessageType, err error) {
	temp, err := e.ToInt()
	if err != nil {
		err = errors.Wrap(err, "could not unmarshal MessageType")
	}
	t = MessageType(temp)
	return
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

// Message represents the data about a message in a discord channel
type Message struct {
	id          snowflake.Snowflake
	channelID   snowflake.Snowflake
	messageType MessageType
	author      User
	content     string
}

// ID returns the ID of the message
func (m *Message) ID() snowflake.Snowflake {
	return m.id
}

// ChannelID returns the ID of the channel the message was sent to
func (m *Message) ChannelID() snowflake.Snowflake {
	return m.channelID
}

// MessageType returns the MessageType of the message
func (m *Message) MessageType() MessageType {
	return m.messageType
}

// AuthorID returns the ID of the author of the message
func (m *Message) AuthorID() snowflake.Snowflake {
	return m.author.id
}

// ContentString returns the content of the message
func (m *Message) ContentString() string {
	return m.content
}

// MessageFromElementMap generates a new Message object from the given data
func MessageFromElementMap(eMap map[string]Element) (m Message, err error) {
	var ok bool
	var e2 Element

	e2, ok = eMap["channel_id"]
	if ok && !e2.IsNil() {
		m.channelID, err = SnowflakeFromElement(e2)
		if err != nil {
			err = errors.Wrap(err, "could not get channel_id snowflake.Snowflake")
			return
		}
	}

	m.messageType, err = MessageTypeFromElement(eMap["type"])
	if err != nil {
		err = errors.Wrap(err, "could not get messageType")
		return
	}

	e2, ok = eMap["content"]
	if ok {
		m.content, err = e2.ToString()
		if err != nil {
			err = errors.Wrap(err, "could not get content")
			return
		}
	}

	m.author, err = UserFromElement(eMap["author"])
	if err != nil {
		err = errors.Wrap(err, "could not inflate message author")
		return
	}

	return
}

// MessageFromElement generates a new Message object from the given Element
func MessageFromElement(e Element) (Message, error) {
	eMap, id, err := MapAndIDFromElement(e)
	if err != nil {
		return Message{}, err
	}

	m, err := MessageFromElementMap(eMap)
	m.id = id

	return m, err
}
