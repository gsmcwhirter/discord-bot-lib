package logging

import (
	log "github.com/gsmcwhirter/go-util/v3/logging"

	"github.com/gsmcwhirter/discord-bot-lib/v8/cmdhandler"
)

var WithContext = log.WithContext

// WithMessage wraps a logger with fields from a cmdhandler.Message
func WithMessage(msg cmdhandler.Message, logger log.Logger) log.Logger {
	logger = log.WithContext(msg.Context(), logger)
	logger = log.With(logger, "user_id", msg.UserID().ToString(), "channel_id", msg.ChannelID().ToString(), "guild_id", msg.GuildID().ToString(), "message_id", msg.MessageID().ToString())
	return logger
}
