package cmdhandler

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gsmcwhirter/go-util/v5/errors"
	"github.com/gsmcwhirter/go-util/v5/parser"
)

// ErrMissingHandler is the error thrown when an event handler cannot be found
var ErrMissingHandler = errors.New("missing handler for command")

// Options provides a way to specify configurable values when creating a CommandHandler
//
// - Placeholder is the string to be used to represent the "command"
// - PreCommand is a string representing the state of commands prior to this one
// - NoHelpOnUnknownCommands can be set to true to NOT display a help message when a command isn't known
// - HelpOnEmptyCommands can be set to true to display a help message when no command is provided
// - CaseSensitive can se set to true to make command recognition case-sensitive
type Options struct {
	Placeholder             string
	PreCommand              string
	NoHelpOnUnknownCommands bool
	HelpOnEmptyCommands     bool
	CaseSensitive           bool
}

// CommandHandler is a dispatcher for commands
type CommandHandler struct {
	parser                parser.Parser
	commands              map[string]MessageHandler
	helpCmd               string
	placeholder           string
	preCommand            string
	helpOnUnknownCommands bool
	caseSensitive         bool
}

// NewCommandHandler creates a new CommandHandler from the given parser
//
// NOTE: the parser's settings must match the Options.CaseSensitive value
func NewCommandHandler(p parser.Parser, opts Options) (*CommandHandler, error) {
	ch := CommandHandler{
		commands:              map[string]MessageHandler{},
		preCommand:            opts.PreCommand,
		helpOnUnknownCommands: !opts.NoHelpOnUnknownCommands,
		caseSensitive:         opts.CaseSensitive,
	}

	err := ch.SetParser(p)
	if err != nil {
		return nil, err
	}

	if opts.Placeholder != "" {
		ch.placeholder = opts.Placeholder
	} else {
		ch.placeholder = "command"
	}

	if opts.HelpOnEmptyCommands {
		ch.SetHandler("", NewMessageHandler(ch.showHelp))
	}
	ch.SetHandler("help", NewMessageHandler(ch.showHelp))
	return &ch, nil
}

// CommandIndicator returns the string prefix required for commands
func (ch *CommandHandler) CommandIndicator() string {
	return ch.parser.LeadChar()
}

// SetParser sets the parser for the command handler
func (ch *CommandHandler) SetParser(p parser.Parser) error {
	if p.IsCaseSensitive() != ch.caseSensitive {
		return errors.New("case sensitive mismatch")
	}

	ch.parser = p
	ch.calculateHelpCmd()
	for cmd := range ch.commands {
		ch.parser.LearnCommand(cmd)
	}

	return nil
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

// SetHandler adds a handler function for the given command, overwriting any
// previously set ones
func (ch *CommandHandler) SetHandler(cmd string, handler MessageHandler) {
	ch.parser.LearnCommand(cmd)

	if ch.caseSensitive {
		ch.commands[cmd] = handler
	}
	ch.commands[strings.ToLower(cmd)] = handler
}

func (ch *CommandHandler) getHandler(cmd string) (MessageHandler, bool) {
	if ch.caseSensitive {
		h, ok := ch.commands[cmd]
		return h, ok
	}

	h, ok := ch.commands[strings.ToLower(cmd)]
	return h, ok
}

// HandleMessage dispatches a Message to the relevant handler
func (ch *CommandHandler) HandleMessage(msg Message) (Response, error) {
	r := &SimpleEmbedResponse{
		To: UserMentionString(msg.UserID()),
	}

	var cmd string
	var err error
	var rest []string

	if len(msg.Contents()) == 0 {
		err = parser.ErrUnknownCommand
	} else {
		cmd, err = ch.parser.ParseCommand(msg.Contents()[0])
		rest = msg.Contents()[1:]
	}

	if err == parser.ErrUnknownCommand {
		if ch.helpOnUnknownCommands {
			var cmd2 string
			cmd2, err = ch.parser.ParseCommand(ch.helpCmd)
			if err != nil {
				r.Description = fmt.Sprintf("Unknown command '%s'", cmd)
				return r, err
			}

			subHandler, cmdExists := ch.getHandler(cmd2)
			if !cmdExists {
				return r, ErrMissingHandler
			}

			var s Response

			s, err = subHandler.HandleMessage(NewWithTokens(msg, rest, msg.ContentErr()))
			s.IncludeError(parser.ErrUnknownCommand)
			return s, err
		}

		return r, err
	}

	subHandler, cmdExists := ch.getHandler(cmd)

	if (err == nil || err == parser.ErrNotACommand) && cmd == "" && cmdExists {
		return subHandler.HandleMessage(NewWithTokens(msg, rest, msg.ContentErr()))
	}

	if err != nil {
		return r, errors.Wrap(err, fmt.Sprintf("%#v", msg.Contents()))
	}

	if !cmdExists {
		if cmd == "" {
			return r, parser.ErrNotACommand
		}
		return r, ErrMissingHandler
	}

	return subHandler.HandleMessage(NewWithTokens(msg, rest, msg.ContentErr()))
}
