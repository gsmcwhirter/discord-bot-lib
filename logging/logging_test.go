package logging

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gsmcwhirter/discord-bot-lib/v6/cmdhandler"
	"github.com/gsmcwhirter/discord-bot-lib/v6/snowflake"
	"github.com/gsmcwhirter/discord-bot-lib/v6/util"
)

type mockLogger struct {
	calls [][]interface{}
}

func (l *mockLogger) Log(keyvals ...interface{}) error {
	l.calls = append(l.calls, keyvals)
	return nil
}

func TestWithContextOk(t *testing.T) {
	mock := &mockLogger{}

	ctx := util.NewRequestContext()
	logger := WithContext(ctx, mock)

	logger.Log("message", "foo")

	if assert.Equal(t, 1, len(mock.calls)) {
		call := mock.calls[0]
		callArgs := make([]string, 0, len(call))
		for _, arg := range call {
			argStr, ok := arg.(string)
			if ok {
				callArgs = append(callArgs, argStr)
			}
		}

		assert.Equal(t, []string{"request_id", ctx.Value(util.ContextKey("request_id")).(string), "message", "foo"}, callArgs)
	}

}

func TestWithContextMissing(t *testing.T) {
	mock := &mockLogger{}

	ctx := context.Background()
	logger := WithContext(ctx, mock)

	logger.Log("message", "foo")

	if assert.Equal(t, 1, len(mock.calls)) {
		call := mock.calls[0]
		callArgs := make([]string, 0, len(call))
		for _, arg := range call {
			argStr, ok := arg.(string)
			if ok {
				callArgs = append(callArgs, argStr)
			}
		}

		assert.Equal(t, []string{"request_id", "unknown", "message", "foo"}, callArgs)
	}

}

func TestWithMessage(t *testing.T) {
	mock := &mockLogger{}

	ctx := util.NewRequestContext()
	msg := cmdhandler.NewSimpleMessage(ctx, snowflake.Snowflake(1), snowflake.Snowflake(2), snowflake.Snowflake(3), snowflake.Snowflake(4), "test")

	logger := WithMessage(msg, mock)

	logger.Log("message", "foo")

	if assert.Equal(t, 1, len(mock.calls)) {
		call := mock.calls[0]
		callArgs := make([]string, 0, len(call))
		for _, arg := range call {
			argStr, ok := arg.(string)
			if ok {
				callArgs = append(callArgs, argStr)
			}
		}

		assert.Equal(t, []string{"request_id", ctx.Value(util.ContextKey("request_id")).(string), "user_id", "1", "channel_id", "3", "guild_id", "2", "message_id", "4", "message", "foo"}, callArgs)
	}

}
