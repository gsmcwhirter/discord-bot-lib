package cmdhandler

import (
	"context"

	"github.com/gsmcwhirter/discord-bot-lib/v22/discordapi/entity"
	"github.com/gsmcwhirter/discord-bot-lib/v22/logging"
	"github.com/gsmcwhirter/discord-bot-lib/v22/snowflake"
)

type Interaction struct {
	entity.Interaction

	Ctx context.Context
}

var _ logging.Message = (*Interaction)(nil)

func (ix *Interaction) Context() context.Context {
	return ix.Ctx
}

func (ix *Interaction) ChannelID() snowflake.Snowflake {
	return ix.ChannelIDSnowflake
}

func (ix *Interaction) GuildID() snowflake.Snowflake {
	return ix.GuildIDSnowflake
}

func (ix *Interaction) MessageID() snowflake.Snowflake {
	return ix.IDSnowflake
}

func (ix *Interaction) UserID() snowflake.Snowflake {
	if ix.User != nil {
		return ix.User.IDSnowflake
	}

	if ix.Member != nil && ix.Member.User != nil {
		return ix.Member.User.IDSnowflake
	}

	return 0
}
