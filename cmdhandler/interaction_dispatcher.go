package cmdhandler

import (
	"github.com/gsmcwhirter/go-util/v10/errors"

	"github.com/gsmcwhirter/discord-bot-lib/v24/discordapi/entity"
	"github.com/gsmcwhirter/discord-bot-lib/v24/snowflake"
)

// ErrMalformedInteraction is returned when an interaction payload is malformed
var ErrMalformedInteraction = errors.New("malformed interaction payload")

// InteractionDispatcher is responsible for dispatching interaction requests to handlers
type InteractionDispatcher struct {
	globals map[string]InteractionCommandHandler
	guilds  map[snowflake.Snowflake]map[string]InteractionCommandHandler
}

// InteractionCommandHandler is the interface for an interaction handler
type InteractionCommandHandler interface {
	Command() entity.ApplicationCommand
	Handler() InteractionHandler
	AutocompleteHandler() AutocompleteHandler
}

type interactionCommandHandler struct {
	command      entity.ApplicationCommand
	handler      InteractionHandler
	autocomplete AutocompleteHandler
}

var _ InteractionCommandHandler = (*interactionCommandHandler)(nil)

func (ich *interactionCommandHandler) Command() entity.ApplicationCommand { return ich.command }
func (ich *interactionCommandHandler) Handler() InteractionHandler        { return ich.handler }
func (ich *interactionCommandHandler) AutocompleteHandler() AutocompleteHandler {
	return ich.autocomplete
}

// NewInteractionDispatcher creates a new interaction dispatcher
func NewInteractionDispatcher(globals []InteractionCommandHandler) (*InteractionDispatcher, error) {
	ix := &InteractionDispatcher{
		globals: make(map[string]InteractionCommandHandler, len(globals)),
		guilds:  make(map[snowflake.Snowflake]map[string]InteractionCommandHandler),
	}

	if err := ix.LearnGlobalCommands(globals); err != nil {
		return ix, errors.Wrap(err, "could not load dispatch table")
	}

	return ix, nil
}

// GlobalCommands returns the global commands for the InteractionDispatcher
func (i *InteractionDispatcher) GlobalCommands() []entity.ApplicationCommand {
	cmds := make([]entity.ApplicationCommand, 0, len(i.globals))
	for _, ich := range i.globals {
		cmds = append(cmds, ich.Command())
	}

	return cmds
}

// GuildCommands returns the guild commands for the InteractionDispatcher
func (i *InteractionDispatcher) GuildCommands() map[snowflake.Snowflake][]entity.ApplicationCommand {
	cmds := make(map[snowflake.Snowflake][]entity.ApplicationCommand, len(i.guilds))

	for gid, gcmds := range i.guilds {
		cmds[gid] = make([]entity.ApplicationCommand, 0, len(gcmds))
		for _, ich := range gcmds {
			cmds[gid] = append(cmds[gid], ich.Command())
		}
	}

	return cmds
}

// LearnGlobalCommands inserts the povided global commands into the dispatch table
func (i *InteractionDispatcher) LearnGlobalCommands(cmds []InteractionCommandHandler) error {
	for _, ich := range cmds {
		i.globals[ich.Command().Name] = ich
	}

	return nil
}

// LearnGuildCommands inserts the provided guild commands into the dispatch table
func (i *InteractionDispatcher) LearnGuildCommands(gid snowflake.Snowflake, cmds []InteractionCommandHandler) error {
	gcmds, ok := i.guilds[gid]
	if !ok {
		gcmds = make(map[string]InteractionCommandHandler, len(cmds))
	}

	for _, ich := range cmds {
		gcmds[ich.Command().Name] = ich
	}

	i.guilds[gid] = gcmds

	return nil
}

// Dispatch sends the interaction to the appropriate dispatcher
func (i *InteractionDispatcher) Dispatch(ix *Interaction) (Response, []Response, error) {
	if ix.Data == nil {
		return nil, nil, errors.WithDetails(ErrMalformedInteraction, "reason", "nil Data")
	}

	handlers := i.guilds[ix.GuildID()]
	handler, ok := handlers[ix.Data.Name]
	if !ok {
		handlers = i.globals
		handler, ok = handlers[ix.Data.Name]
		if !ok {
			return nil, nil, ErrMissingHandler
		}
	}

	return handler.Handler().HandleInteraction(ix)
}

// Autocomplete returns autocomplete information for the command
func (i *InteractionDispatcher) Autocomplete(ix *Interaction) ([]entity.ApplicationCommandOptionChoice, error) {
	if ix.Data == nil {
		return nil, errors.WithDetails(ErrMalformedInteraction, "reason", "nil Data")
	}

	handlers := i.guilds[ix.GuildID()]
	handler, ok := handlers[ix.Data.Name]
	if !ok {
		handlers = i.globals
		handler, ok = handlers[ix.Data.Name]
		if !ok {
			return nil, ErrMissingHandler
		}
	}

	return handler.AutocompleteHandler().Autocomplete(ix)
}
