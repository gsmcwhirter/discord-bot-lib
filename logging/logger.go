package logging

import (
	"context"

	log "github.com/gsmcwhirter/go-util/v5/logging"

	"github.com/gsmcwhirter/discord-bot-lib/v10/cmdhandler"
	"github.com/gsmcwhirter/discord-bot-lib/v10/request"
)

// WithContext wraps a logger with fields from a context
func WithContext(ctx context.Context, logger log.Logger) log.Logger {
	logger = log.WithContext(ctx, logger)

	if gid, ok := request.GetGuildID(ctx); ok {
		logger = log.With(logger, "guild_id", gid.ToString())
	}

	return logger
}

// WithMessage wraps a logger with fields from a cmdhandler.Message
func WithMessage(msg cmdhandler.Message, logger log.Logger) log.Logger {
	logger = log.WithContext(msg.Context(), logger)
	logger = log.With(logger, "user_id", msg.UserID().ToString(), "channel_id", msg.ChannelID().ToString(), "guild_id", msg.GuildID().ToString(), "message_id", msg.MessageID().ToString())
	return logger
}
