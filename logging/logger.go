package logging

import (
	"context"

	log "github.com/gsmcwhirter/go-util/v8/logging"

	"github.com/gsmcwhirter/discord-bot-lib/v24/request"
	"github.com/gsmcwhirter/discord-bot-lib/v24/snowflake"
)

// Message is the expected interface of a message
type Message interface {
	Context() context.Context
	UserID() snowflake.Snowflake
	GuildID() snowflake.Snowflake
	ChannelID() snowflake.Snowflake
	MessageID() snowflake.Snowflake
}

// WithContext wraps a logger with fields from a context
func WithContext(ctx context.Context, logger log.Logger) log.Logger {
	logger = log.WithContext(ctx, logger)

	if gid, ok := request.GetGuildID(ctx); ok {
		logger = log.With(logger, "guild_id", gid.ToString())
	}

	return logger
}

// WithMessage wraps a logger with fields from a cmdhandler.Message
func WithMessage(msg Message, logger log.Logger) log.Logger {
	logger = log.WithContext(msg.Context(), logger)
	logger = log.With(logger, "user_id", msg.UserID().ToString(), "channel_id", msg.ChannelID().ToString(), "guild_id", msg.GuildID().ToString(), "message_id", msg.MessageID().ToString())
	return logger
}
