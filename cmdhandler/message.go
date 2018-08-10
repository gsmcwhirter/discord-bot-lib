package cmdhandler

import "github.com/gsmcwhirter/discord-bot-lib/snowflake"

// Message TODOC
type Message interface {
	UserID() snowflake.Snowflake
	GuildID() snowflake.Snowflake
	ChannelID() snowflake.Snowflake
	MessageID() snowflake.Snowflake
	Contents() string
}

// SimpleMessage TODOC
type SimpleMessage struct {
	userID    snowflake.Snowflake
	guildID   snowflake.Snowflake
	channelID snowflake.Snowflake
	messageID snowflake.Snowflake
	contents  string
}

// NewSimpleMessage TODOC
func NewSimpleMessage(userID, guildID, channelID, messageID snowflake.Snowflake, contents string) *SimpleMessage {
	return &SimpleMessage{
		userID:    userID,
		guildID:   guildID,
		channelID: channelID,
		messageID: messageID,
		contents:  contents,
	}
}

// NewWithContents TODOC
func NewWithContents(m Message, contents string) *SimpleMessage {
	return &SimpleMessage{
		userID:    m.UserID(),
		guildID:   m.GuildID(),
		channelID: m.ChannelID(),
		messageID: m.MessageID(),
		contents:  contents,
	}
}

// UserID TODOC
func (m *SimpleMessage) UserID() snowflake.Snowflake {
	return m.userID
}

// GuildID TODOC
func (m *SimpleMessage) GuildID() snowflake.Snowflake {
	return m.guildID
}

// ChannelID TODOC
func (m *SimpleMessage) ChannelID() snowflake.Snowflake {
	return m.channelID
}

// MessageID TODO
func (m *SimpleMessage) MessageID() snowflake.Snowflake {
	return m.messageID
}

// Contents TODOC
func (m *SimpleMessage) Contents() string {
	return m.contents
}
