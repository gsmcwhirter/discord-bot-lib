package cmdhandler

import (
	"fmt"
	"sort"

	"github.com/pkg/errors"

	"github.com/gsmcwhirter/go-util/parser"
)

// ErrMissingHandler TODOC
var ErrMissingHandler = errors.New("missing handler for command")

// Options TODOC
type Options struct {
	Placeholder string
	PreCommand  string
}

// CommandHandler TODOC
type CommandHandler struct {
	parser      parser.Parser
	commands    map[string]LineHandler
	helpCmd     string
	placeholder string
	preCommand  string
}

// NewCommandHandler TODOC
func NewCommandHandler(parser parser.Parser, opts Options) *CommandHandler {
	ch := CommandHandler{
		commands:   map[string]LineHandler{},
		preCommand: opts.PreCommand,
	}
	ch.SetParser(parser)

	if opts.Placeholder != "" {
		ch.placeholder = opts.Placeholder
	} else {
		ch.placeholder = "command"
	}

	ch.SetHandler("", NewLineHandler(ch.showHelp))
	ch.SetHandler("help", NewLineHandler(ch.showHelp))
	return &ch
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

func (ch *CommandHandler) showHelp(user, guild string, line string) (Response, error) {
	var helpStr string
	if ch.preCommand != "" {
		helpStr = fmt.Sprintf("Usage: %s [%s]\n\n", ch.preCommand, ch.placeholder)
	}

	helpStr += fmt.Sprintf("Available %ss:\n", ch.placeholder)
	cNames := make([]string, 0, len(ch.commands))
	for cmd := range ch.commands {
		cNames = append(cNames, cmd)
	}
	sort.Strings(cNames)

	for _, cmd := range cNames {
		if cmd != "" {
			helpStr += fmt.Sprintf("  %s\n", cmd)
		}
	}

	return &SimpleResponse{
		To:      user,
		Content: helpStr,
	}, nil
}

// SetHandler TODOC
func (ch *CommandHandler) SetHandler(cmd string, handler LineHandler) {
	ch.parser.LearnCommand(cmd)
	ch.commands[cmd] = handler
}

// HandleLine TODOC
func (ch *CommandHandler) HandleLine(user, guild string, line string) (Response, error) {
	r := &SimpleResponse{
		To: user,
	}

	cmd, rest, err := ch.parser.ParseCommand(line)

	subHandler, cmdExists := ch.commands[cmd]
	if err == parser.ErrNotACommand && cmd != "" && !cmdExists {
		return r, err
	}

	if err == parser.ErrUnknownCommand {
		var cmd2 string
		cmd2, rest, err = ch.parser.ParseCommand(ch.helpCmd)

		subHandler, cmdExists = ch.commands[cmd2]
		if !cmdExists {
			return r, ErrMissingHandler
		}

		if err != nil {
			r.Content = fmt.Sprintf("Unknown command '%s'", cmd)
			return r, err
		}

		s, sherr := subHandler.HandleLine(user, guild, rest)
		if sherr == nil {
			sherr = parser.ErrUnknownCommand
		}
		return s, sherr
	}

	if err != nil && err != parser.ErrNotACommand {
		return r, err
	}

	if !cmdExists {
		return r, ErrMissingHandler
	}

	return subHandler.HandleLine(user, guild, rest)
}
