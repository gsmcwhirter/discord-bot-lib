package cmdhandler

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gsmcwhirter/go-util/parser"
	"github.com/pkg/errors"

	"github.com/gsmcwhirter/discord-bot-lib/discordapi"
)

// ErrMissingHandler TODOC
var ErrMissingHandler = errors.New("missing handler for command")

// Options TODOC
type Options struct {
	Placeholder             string
	PreCommand              string
	NoHelpOnUnknownCommands bool
	HelpOnEmptyCommands     bool
}

// CommandHandler TODOC
type CommandHandler struct {
	parser                parser.Parser
	commands              map[string]MessageHandler
	helpCmd               string
	placeholder           string
	preCommand            string
	helpOnUnknownCommands bool
	bot                   discordapi.DiscordBot
}

// NewCommandHandler TODOC
func NewCommandHandler(parser parser.Parser, opts Options) *CommandHandler {
	ch := CommandHandler{
		commands:              map[string]MessageHandler{},
		preCommand:            opts.PreCommand,
		helpOnUnknownCommands: !opts.NoHelpOnUnknownCommands,
	}
	ch.SetParser(parser)

	if opts.Placeholder != "" {
		ch.placeholder = opts.Placeholder
	} else {
		ch.placeholder = "command"
	}

	if opts.HelpOnEmptyCommands {
		ch.SetHandler("", NewMessageHandler(ch.showHelp))
	}
	ch.SetHandler("help", NewMessageHandler(ch.showHelp))
	return &ch
}

// ConnectToBot TODOC
func (ch *CommandHandler) ConnectToBot(b discordapi.DiscordBot) {
	ch.bot = b
}

// DiscordBot TODOC
func (ch *CommandHandler) DiscordBot() discordapi.DiscordBot {
	return ch.bot
}

// CommandIndicator TODOC
func (ch *CommandHandler) CommandIndicator() string {
	return ch.parser.LeadChar()
}

// SetParser sets the parser for the command handler
func (ch *CommandHandler) SetParser(p parser.Parser) {
	ch.parser = p
	ch.calculateHelpCmd()
	for cmd := range ch.commands {
		ch.parser.LearnCommand(cmd)
	}
}

func (ch *CommandHandler) calculateHelpCmd() {
	ch.helpCmd = ch.parser.LeadChar() + "help"
}

func (ch *CommandHandler) showHelp(msg Message) (Response, error) {
	r := &EmbedResponse{
		To: UserMentionString(msg.UserID()),
	}

	if ch.preCommand != "" {
		r.Description = fmt.Sprintf("Usage: %s [%s]\n\n", ch.preCommand, ch.placeholder)
	}

	cNames := make([]string, 0, len(ch.commands))
	for cmd := range ch.commands {
		if cmd == "" {
			continue
		}
		cNames = append(cNames, cmd)
	}
	sort.Strings(cNames)

	r.Fields = []EmbedField{
		{
			Name: "*Available Commands*",
			Val:  fmt.Sprintf("```\n%s\n```\n", strings.Join(cNames, "\n")),
		},
	}

	return r, nil
}

// SetHandler TODOC
func (ch *CommandHandler) SetHandler(cmd string, handler MessageHandler) {
	ch.parser.LearnCommand(cmd)
	ch.commands[cmd] = handler
}

// HandleLine TODOC
func (ch *CommandHandler) HandleLine(msg Message) (Response, error) {
	r := &SimpleEmbedResponse{
		To: UserMentionString(msg.UserID()),
	}

	cmd, rest, err := ch.parser.ParseCommand(msg.Contents())
	if err == parser.ErrUnknownCommand {
		if ch.helpOnUnknownCommands {
			cmd2, rest, err2 := ch.parser.ParseCommand(ch.helpCmd)
			if err2 != nil {
				r.Description = fmt.Sprintf("Unknown command '%s'", cmd)
				return r, err2
			}

			subHandler, cmdExists := ch.commands[cmd2]
			if !cmdExists {
				return r, ErrMissingHandler
			}

			s, err2 := subHandler.HandleLine(NewWithContents(msg, rest))
			s.IncludeError(parser.ErrUnknownCommand)
			return s, err2
		}

		return r, err
	}

	subHandler, cmdExists := ch.commands[cmd]

	if (err == nil || err == parser.ErrNotACommand) && cmd == "" && cmdExists {
		return subHandler.HandleLine(NewWithContents(msg, rest))
	}

	if err != nil {
		return r, err
	}

	if !cmdExists {
		if cmd == "" {
			return r, parser.ErrNotACommand
		}
		return r, ErrMissingHandler
	}

	return subHandler.HandleLine(NewWithContents(msg, rest))
}
