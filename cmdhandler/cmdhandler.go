package cmdhandler

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gsmcwhirter/go-util/parser"
	"github.com/pkg/errors"
)

// ErrMissingHandler TODOC
var ErrMissingHandler = errors.New("missing handler for command")

// Options TODOC
type Options struct {
	Placeholder             string
	PreCommand              string
	NoHelpOnUnknownCommands bool
	HelpOnEmptyCommands     bool
	CaseSensitive           bool
}

// CommandHandler TODOC
type CommandHandler struct {
	parser                parser.Parser
	commands              map[string]MessageHandler
	helpCmd               string
	placeholder           string
	preCommand            string
	helpOnUnknownCommands bool
	caseSensitive         bool
}

// NewCommandHandler TODOC
func NewCommandHandler(parser parser.Parser, opts Options) (*CommandHandler, error) {
	ch := CommandHandler{
		commands:              map[string]MessageHandler{},
		preCommand:            opts.PreCommand,
		helpOnUnknownCommands: !opts.NoHelpOnUnknownCommands,
		caseSensitive:         opts.CaseSensitive,
	}

	err := ch.SetParser(parser)
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

// CommandIndicator TODOC
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

// SetHandler TODOC
func (ch *CommandHandler) SetHandler(cmd string, handler MessageHandler) {
	ch.parser.LearnCommand(cmd)

	if ch.caseSensitive {
		ch.commands[cmd] = handler
	}
	ch.commands[strings.ToLower(cmd)] = handler
}

func (ch *CommandHandler) getHandler(cmd string) (h MessageHandler, ok bool) {
	if ch.caseSensitive {
		h, ok = ch.commands[cmd]
		return
	}

	h, ok = ch.commands[strings.ToLower(cmd)]
	return
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

			subHandler, cmdExists := ch.getHandler(cmd2)
			if !cmdExists {
				return r, ErrMissingHandler
			}

			s, err2 := subHandler.HandleLine(NewWithContents(msg, rest))
			s.IncludeError(parser.ErrUnknownCommand)
			return s, err2
		}

		return r, err
	}

	subHandler, cmdExists := ch.getHandler(cmd)

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
