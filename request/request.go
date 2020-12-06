package request

import (
	"context"

	"github.com/gsmcwhirter/discord-bot-lib/v15/snowflake"
)

// ContextKey is a wrapper type for our keys attached to a context
type ContextKey string

func WithGuildID(ctx context.Context, gid snowflake.Snowflake) context.Context {
	return context.WithValue(ctx, ContextKey("guild_id"), gid)
}

func GetGuildID(ctx context.Context) (snowflake.Snowflake, bool) {
	gid, ok := ctx.Value(ContextKey("guild_id")).(snowflake.Snowflake)
	return gid, ok
}
