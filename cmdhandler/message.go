package cmdhandler

import (
	"context"

	"github.com/gsmcwhirter/go-util/v7/parser"

	"github.com/gsmcwhirter/discord-bot-lib/v18/snowflake"
)

var quotes = []rune{'"', '“', '”', '«', '»', '„'}

// Message is the api for a message that a command handler will respond to
type Message interface {
	Context() context.Context
	UserID() snowflake.Snowflake
	GuildID() snowflake.Snowflake
	ChannelID() snowflake.Snowflake
	MessageID() snowflake.Snowflake
	Contents() []string
	ContentErr() error
}

type simpleMessage struct {
	ctx        context.Context
	userID     snowflake.Snowflake
	guildID    snowflake.Snowflake
	channelID  snowflake.Snowflake
	messageID  snowflake.Snowflake
	contents   []string
	contentErr error
}

// NewSimpleMessage creates a new Message object
func NewSimpleMessage(ctx context.Context, userID, guildID, channelID, messageID snowflake.Snowflake, contents string) Message {
	tokens, err := parser.Tokenize(contents, ' ', '\\', quotes)

	return &simpleMessage{
		ctx:        ctx,
		userID:     userID,
		guildID:    guildID,
		channelID:  channelID,
		messageID:  messageID,
		contents:   tokens,
		contentErr: err,
	}
}

// NewWithContents clones a given message object but substitutes the Contents() with the provided string
func NewWithContents(m Message, contents string) Message {
	tokens, err := parser.Tokenize(contents, ' ', '\\', quotes)
	return NewWithTokens(m, tokens, err)
}

// NewWithTokens clones a given message object but substitutes the Contents() and ContentErr()
func NewWithTokens(m Message, tokens []string, contentErr error) Message {
	return &simpleMessage{
		ctx:        m.Context(),
		userID:     m.UserID(),
		guildID:    m.GuildID(),
		channelID:  m.ChannelID(),
		messageID:  m.MessageID(),
		contents:   tokens,
		contentErr: contentErr,
	}
}

func NewWithContext(ctx context.Context, m Message) Message {
	return &simpleMessage{
		ctx:        ctx,
		userID:     m.UserID(),
		guildID:    m.GuildID(),
		channelID:  m.ChannelID(),
		messageID:  m.MessageID(),
		contents:   m.Contents(),
		contentErr: m.ContentErr(),
	}
}

func (m *simpleMessage) Context() context.Context {
	return m.ctx
}

func (m *simpleMessage) UserID() snowflake.Snowflake {
	return m.userID
}

func (m *simpleMessage) GuildID() snowflake.Snowflake {
	return m.guildID
}

func (m *simpleMessage) ChannelID() snowflake.Snowflake {
	return m.channelID
}

func (m *simpleMessage) MessageID() snowflake.Snowflake {
	return m.messageID
}

func (m *simpleMessage) Contents() []string {
	return m.contents
}

func (m *simpleMessage) ContentErr() error {
	return m.contentErr
}
