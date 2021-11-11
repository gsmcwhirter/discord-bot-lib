package cmdhandler

import (
	"context"
	"reflect"
	"testing"

	"github.com/gsmcwhirter/discord-bot-lib/v22/snowflake"
)

func TestNewSimpleMessage(t *testing.T) {
	type args struct {
		ctx       context.Context
		userID    snowflake.Snowflake
		guildID   snowflake.Snowflake
		channelID snowflake.Snowflake
		messageID snowflake.Snowflake
		contents  string
	}
	tests := []struct {
		name string
		args args
		want Message
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSimpleMessage(tt.args.ctx, tt.args.userID, tt.args.guildID, tt.args.channelID, tt.args.messageID, tt.args.contents); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSimpleMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewWithContents(t *testing.T) {
	type args struct {
		m        Message
		contents string
	}
	tests := []struct {
		name string
		args args
		want Message
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewWithContents(tt.args.m, tt.args.contents); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWithContents() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewWithTokens(t *testing.T) {
	type args struct {
		m          Message
		tokens     []string
		contentErr error
	}
	tests := []struct {
		name string
		args args
		want Message
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewWithTokens(tt.args.m, tt.args.tokens, tt.args.contentErr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWithTokens() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewWithContext(t *testing.T) {
	type args struct {
		ctx context.Context
		m   Message
	}
	tests := []struct {
		name string
		args args
		want Message
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewWithContext(tt.args.ctx, tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWithContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_simpleMessage_Context(t *testing.T) {
	type fields struct {
		ctx        context.Context
		userID     snowflake.Snowflake
		guildID    snowflake.Snowflake
		channelID  snowflake.Snowflake
		messageID  snowflake.Snowflake
		contents   []string
		contentErr error
	}
	tests := []struct {
		name   string
		fields fields
		want   context.Context
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &simpleMessage{
				ctx:        tt.fields.ctx,
				userID:     tt.fields.userID,
				guildID:    tt.fields.guildID,
				channelID:  tt.fields.channelID,
				messageID:  tt.fields.messageID,
				contents:   tt.fields.contents,
				contentErr: tt.fields.contentErr,
			}
			if got := m.Context(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("simpleMessage.Context() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_simpleMessage_UserID(t *testing.T) {
	type fields struct {
		ctx        context.Context
		userID     snowflake.Snowflake
		guildID    snowflake.Snowflake
		channelID  snowflake.Snowflake
		messageID  snowflake.Snowflake
		contents   []string
		contentErr error
	}
	tests := []struct {
		name   string
		fields fields
		want   snowflake.Snowflake
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &simpleMessage{
				ctx:        tt.fields.ctx,
				userID:     tt.fields.userID,
				guildID:    tt.fields.guildID,
				channelID:  tt.fields.channelID,
				messageID:  tt.fields.messageID,
				contents:   tt.fields.contents,
				contentErr: tt.fields.contentErr,
			}
			if got := m.UserID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("simpleMessage.UserID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_simpleMessage_GuildID(t *testing.T) {
	type fields struct {
		ctx        context.Context
		userID     snowflake.Snowflake
		guildID    snowflake.Snowflake
		channelID  snowflake.Snowflake
		messageID  snowflake.Snowflake
		contents   []string
		contentErr error
	}
	tests := []struct {
		name   string
		fields fields
		want   snowflake.Snowflake
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &simpleMessage{
				ctx:        tt.fields.ctx,
				userID:     tt.fields.userID,
				guildID:    tt.fields.guildID,
				channelID:  tt.fields.channelID,
				messageID:  tt.fields.messageID,
				contents:   tt.fields.contents,
				contentErr: tt.fields.contentErr,
			}
			if got := m.GuildID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("simpleMessage.GuildID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_simpleMessage_ChannelID(t *testing.T) {
	type fields struct {
		ctx        context.Context
		userID     snowflake.Snowflake
		guildID    snowflake.Snowflake
		channelID  snowflake.Snowflake
		messageID  snowflake.Snowflake
		contents   []string
		contentErr error
	}
	tests := []struct {
		name   string
		fields fields
		want   snowflake.Snowflake
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &simpleMessage{
				ctx:        tt.fields.ctx,
				userID:     tt.fields.userID,
				guildID:    tt.fields.guildID,
				channelID:  tt.fields.channelID,
				messageID:  tt.fields.messageID,
				contents:   tt.fields.contents,
				contentErr: tt.fields.contentErr,
			}
			if got := m.ChannelID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("simpleMessage.ChannelID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_simpleMessage_MessageID(t *testing.T) {
	type fields struct {
		ctx        context.Context
		userID     snowflake.Snowflake
		guildID    snowflake.Snowflake
		channelID  snowflake.Snowflake
		messageID  snowflake.Snowflake
		contents   []string
		contentErr error
	}
	tests := []struct {
		name   string
		fields fields
		want   snowflake.Snowflake
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &simpleMessage{
				ctx:        tt.fields.ctx,
				userID:     tt.fields.userID,
				guildID:    tt.fields.guildID,
				channelID:  tt.fields.channelID,
				messageID:  tt.fields.messageID,
				contents:   tt.fields.contents,
				contentErr: tt.fields.contentErr,
			}
			if got := m.MessageID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("simpleMessage.MessageID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_simpleMessage_Contents(t *testing.T) {
	type fields struct {
		ctx        context.Context
		userID     snowflake.Snowflake
		guildID    snowflake.Snowflake
		channelID  snowflake.Snowflake
		messageID  snowflake.Snowflake
		contents   []string
		contentErr error
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &simpleMessage{
				ctx:        tt.fields.ctx,
				userID:     tt.fields.userID,
				guildID:    tt.fields.guildID,
				channelID:  tt.fields.channelID,
				messageID:  tt.fields.messageID,
				contents:   tt.fields.contents,
				contentErr: tt.fields.contentErr,
			}
			if got := m.Contents(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("simpleMessage.Contents() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_simpleMessage_ContentErr(t *testing.T) {
	type fields struct {
		ctx        context.Context
		userID     snowflake.Snowflake
		guildID    snowflake.Snowflake
		channelID  snowflake.Snowflake
		messageID  snowflake.Snowflake
		contents   []string
		contentErr error
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &simpleMessage{
				ctx:        tt.fields.ctx,
				userID:     tt.fields.userID,
				guildID:    tt.fields.guildID,
				channelID:  tt.fields.channelID,
				messageID:  tt.fields.messageID,
				contents:   tt.fields.contents,
				contentErr: tt.fields.contentErr,
			}
			if err := m.ContentErr(); (err != nil) != tt.wantErr {
				t.Errorf("simpleMessage.ContentErr() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
