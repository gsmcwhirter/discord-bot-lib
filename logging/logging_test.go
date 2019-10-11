package logging

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gsmcwhirter/go-util/v5/request"

	"github.com/gsmcwhirter/discord-bot-lib/v12/cmdhandler"
	"github.com/gsmcwhirter/discord-bot-lib/v12/snowflake"
)

type mockLogger struct {
	calls [][]interface{}
}

func (l *mockLogger) Log(keyvals ...interface{}) error {
	l.calls = append(l.calls, keyvals)
	return nil
}

func (l *mockLogger) Err(msg string, err error, keyvals ...interface{}) {
	keyvals = append([]interface{}{"message", msg, "error", err}, keyvals...)
	_ = l.Log(keyvals...)
}

func (l *mockLogger) Message(msg string, keyvals ...interface{}) {
	keyvals = append([]interface{}{"message", msg}, keyvals...)
	_ = l.Log(keyvals...)
}

func (l *mockLogger) Printf(f string, args ...interface{}) {
	msg := fmt.Sprintf(f, args...)
	_ = l.Log("message", msg)
}

func TestWithMessage(t *testing.T) {
	mock := &mockLogger{}

	ctx := request.NewRequestContext()
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

		rid, ok := request.GetRequestID(ctx)
		if !ok {
			rid = "(unknown)"
		}
		assert.Equal(t, []string{"request_id", rid, "user_id", "1", "channel_id", "3", "guild_id", "2", "message_id", "4", "message", "foo"}, callArgs)
	}

}
