package cmdhandler

import (
	"context"
	"reflect"
	"testing"

	"github.com/gsmcwhirter/go-util/v8/parser"
	"github.com/stretchr/testify/assert"
)

func TestNewCommandHandler(t *testing.T) {
	t.Parallel()

	p := parser.NewParser(parser.Options{
		CmdIndicator:  "!",
		KnownCommands: []string{"foo", "bar"},
		CaseSensitive: false,
	})

	type args struct {
		p    parser.Parser
		opts Options
	}
	tests := []struct {
		name    string
		args    args
		want    *CommandHandler
		wantErr bool
	}{
		{
			name: "basic",
			args: args{
				p: p,
				opts: Options{
					Placeholder:         "ph",
					PreCommand:          "pc",
					HelpOnEmptyCommands: true,
				},
			},
			want: &CommandHandler{
				parser: p,
				commands: map[string]MessageHandler{
					"":     nil,
					"help": nil,
				},
				helpCmd:               "!help",
				placeholder:           "ph",
				preCommand:            "pc",
				helpOnUnknownCommands: true,
				caseSensitive:         false,
			},
			wantErr: false,
		},
		{
			name: "auto-placeholder",
			args: args{
				p: p,
				opts: Options{
					PreCommand:          "pc",
					HelpOnEmptyCommands: true,
				},
			},
			want: &CommandHandler{
				parser: p,
				commands: map[string]MessageHandler{
					"":     nil,
					"help": nil,
				},
				helpCmd:               "!help",
				placeholder:           "command",
				preCommand:            "pc",
				helpOnUnknownCommands: true,
				caseSensitive:         false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := NewCommandHandler(tt.args.p, tt.args.opts)

			// hack around storing a function/handler dynamically
			tt.want.commands[""] = got.showHelpHandler()
			tt.want.commands["help"] = got.showHelpHandler()
			tt.want.helpHandler = got.showHelpHandler()
			// end hack

			if (tt.wantErr && !assert.Error(t, err)) || (!tt.wantErr && !assert.NoError(t, err)) {
				return
			}
			if err != nil {
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCommandHandler_CommandIndicator(t *testing.T) {
	t.Parallel()

	defP := parser.NewParser(parser.Options{
		CmdIndicator:  "!",
		KnownCommands: []string{"foo", "bar"},
		CaseSensitive: false,
	})

	myP := parser.NewParser(parser.Options{
		CmdIndicator:  "%",
		KnownCommands: []string{"foo", "bar"},
		CaseSensitive: false,
	})

	tests := []struct {
		name string
		ch   *CommandHandler
		want string
	}{
		{
			name: "default",
			ch:   mustNewCommandHandler(defP, Options{}),
			want: "!",
		},
		{
			name: "override",
			ch:   mustNewCommandHandler(myP, Options{}),
			want: "%",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.ch.CommandIndicator()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCommandHandler_SetParser(t *testing.T) {
	t.Parallel()

	type args struct {
		p parser.Parser
	}
	tests := []struct {
		name    string
		ch      *CommandHandler
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			ch: mustNewCommandHandler(parser.NewParser(parser.Options{
				CmdIndicator:  "!",
				KnownCommands: []string{"foo", "bar"},
				CaseSensitive: false,
			}), Options{}),
			args: args{
				p: parser.NewParser(parser.Options{
					CmdIndicator:  "%",
					KnownCommands: []string{"foo", "bar"},
					CaseSensitive: false,
				}),
			},
			wantErr: false,
		},
		{
			name: "case mismatch 1",
			ch: mustNewCommandHandler(parser.NewParser(parser.Options{
				CmdIndicator:  "!",
				KnownCommands: []string{"foo", "bar"},
				CaseSensitive: false,
			}), Options{}),
			args: args{
				p: parser.NewParser(parser.Options{
					CmdIndicator:  "%",
					KnownCommands: []string{"foo", "bar"},
					CaseSensitive: true,
				}),
			},
			wantErr: true,
		},
		{
			name: "case mismatch 2",
			ch: mustNewCommandHandler(parser.NewParser(parser.Options{
				CmdIndicator:  "!",
				KnownCommands: []string{"foo", "bar"},
				CaseSensitive: true,
			}), Options{
				CaseSensitive: true,
			}),
			args: args{
				p: parser.NewParser(parser.Options{
					CmdIndicator:  "%",
					KnownCommands: []string{"foo", "bar"},
					CaseSensitive: false,
				}),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, "!", tt.ch.CommandIndicator())

			err := tt.ch.SetParser(tt.args.p)

			if (tt.wantErr && !assert.Error(t, err)) || (!tt.wantErr && !assert.NoError(t, err)) {
				return
			}
			if err != nil {
				return
			}

			assert.Equal(t, "%", tt.ch.CommandIndicator())
		})
	}
}

func TestCommandHandler_showHelp(t *testing.T) {
	t.Parallel()
	type args struct {
		msg Message
	}
	tests := []struct {
		name     string
		ch       *CommandHandler
		args     args
		wantText string
		wantErr  bool
	}{
		{
			name: "no pre",
			ch: mustNewCommandHandler(parser.NewParser(parser.Options{
				CmdIndicator:  "!",
				KnownCommands: []string{"foo", "bar"},
			}), Options{}),
			args: args{
				msg: NewSimpleMessage(context.Background(), 1, 2, 3, 4, ""),
			},
			wantText: "\n\n*Available Commands*:\n```\nhelp\n```\n\n",
			wantErr:  false,
		},
		{
			name: "with pre",
			ch: mustNewCommandHandler(parser.NewParser(parser.Options{
				CmdIndicator:  "!",
				KnownCommands: []string{"foo", "bar"},
			}), Options{
				Placeholder: "foo",
				PreCommand:  "pre",
			}),
			args: args{
				msg: NewSimpleMessage(context.Background(), 1, 2, 3, 4, ""),
			},
			wantText: "\n\nUsage: pre [foo]\n\n\n*Available Commands*:\n```\nhelp\n```\n\n",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.ch.showHelp(tt.args.msg)

			if (tt.wantErr && !assert.Error(t, err)) || (!tt.wantErr && !assert.NoError(t, err)) {
				return
			}
			if err != nil {
				return
			}

			assert.Equal(t, tt.wantText, got.ToString())
		})
	}
}

func TestCommandHandler_SetHandler(t *testing.T) {
	t.Parallel()
	type fields struct {
		parser                parser.Parser
		commands              map[string]MessageHandler
		helpCmd               string
		placeholder           string
		preCommand            string
		helpOnUnknownCommands bool
		caseSensitive         bool
	}
	type args struct {
		cmd     string
		handler MessageHandler
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ch := &CommandHandler{
				parser:                tt.fields.parser,
				commands:              tt.fields.commands,
				helpCmd:               tt.fields.helpCmd,
				placeholder:           tt.fields.placeholder,
				preCommand:            tt.fields.preCommand,
				helpOnUnknownCommands: tt.fields.helpOnUnknownCommands,
				caseSensitive:         tt.fields.caseSensitive,
			}
			ch.SetHandler(tt.args.cmd, tt.args.handler)
		})
	}
}

func TestCommandHandler_getHandler(t *testing.T) {
	t.Parallel()
	type fields struct {
		parser                parser.Parser
		commands              map[string]MessageHandler
		helpCmd               string
		placeholder           string
		preCommand            string
		helpOnUnknownCommands bool
		caseSensitive         bool
	}
	type args struct {
		cmd string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   MessageHandler
		want1  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ch := &CommandHandler{
				parser:                tt.fields.parser,
				commands:              tt.fields.commands,
				helpCmd:               tt.fields.helpCmd,
				placeholder:           tt.fields.placeholder,
				preCommand:            tt.fields.preCommand,
				helpOnUnknownCommands: tt.fields.helpOnUnknownCommands,
				caseSensitive:         tt.fields.caseSensitive,
			}
			got, got1 := ch.getHandler(tt.args.cmd)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CommandHandler.getHandler() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CommandHandler.getHandler() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestCommandHandler_HandleMessage(t *testing.T) {
	t.Parallel()
	type fields struct {
		parser                parser.Parser
		commands              map[string]MessageHandler
		helpCmd               string
		placeholder           string
		preCommand            string
		helpOnUnknownCommands bool
		caseSensitive         bool
	}
	type args struct {
		msg Message
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Response
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ch := &CommandHandler{
				parser:                tt.fields.parser,
				commands:              tt.fields.commands,
				helpCmd:               tt.fields.helpCmd,
				placeholder:           tt.fields.placeholder,
				preCommand:            tt.fields.preCommand,
				helpOnUnknownCommands: tt.fields.helpOnUnknownCommands,
				caseSensitive:         tt.fields.caseSensitive,
			}
			got, err := ch.HandleMessage(tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("CommandHandler.HandleMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CommandHandler.HandleMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
