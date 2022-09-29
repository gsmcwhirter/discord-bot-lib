package cmdhandler

import (
	"context"

	"github.com/gsmcwhirter/discord-bot-lib/v24/discordapi/entity"
	"github.com/gsmcwhirter/discord-bot-lib/v24/logging"
	"github.com/gsmcwhirter/discord-bot-lib/v24/snowflake"
)

// Interaction represents an interaction request and its context
type Interaction struct {
	entity.Interaction

	Ctx context.Context
}

var _ logging.Message = (*Interaction)(nil)

// Context returns the interaction request context
func (ix *Interaction) Context() context.Context {
	return ix.Ctx
}

// ChannelID returns the interaction channel ID
func (ix *Interaction) ChannelID() snowflake.Snowflake {
	return ix.ChannelIDSnowflake
}

// GuildID returns the interaction guild ID
func (ix *Interaction) GuildID() snowflake.Snowflake {
	return ix.GuildIDSnowflake
}

// MessageID returns the interaction ID
func (ix *Interaction) MessageID() snowflake.Snowflake {
	return ix.IDSnowflake
}

// UserID returns the interaction user ID
func (ix *Interaction) UserID() snowflake.Snowflake {
	if ix.User != nil {
		return ix.User.IDSnowflake
	}

	if ix.Member != nil && ix.Member.User != nil {
		return ix.Member.User.IDSnowflake
	}

	return 0
}
