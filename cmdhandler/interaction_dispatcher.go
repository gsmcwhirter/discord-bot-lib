package cmdhandler

import (
	"github.com/gsmcwhirter/discord-bot-lib/v22/discordapi/entity"
	"github.com/gsmcwhirter/discord-bot-lib/v22/snowflake"
	"github.com/gsmcwhirter/go-util/v8/errors"
)

var ErrMalformedInteraction = errors.New("malformed interaction payload")

type InteractionDispatcher struct {
	globals []InteractionCommandHandler
	guilds  map[snowflake.Snowflake][]InteractionCommandHandler

	dispatch map[string]map[snowflake.Snowflake]InteractionHandler
}

type InteractionCommandHandler struct {
	Command entity.ApplicationCommand
	Handler InteractionHandler
}

func NewInteractionDispatcher(globals []InteractionCommandHandler) (*InteractionDispatcher, error) {
	ix := &InteractionDispatcher{
		globals:  make([]InteractionCommandHandler, 0, len(globals)),
		guilds:   make(map[snowflake.Snowflake][]InteractionCommandHandler),
		dispatch: make(map[string]map[snowflake.Snowflake]InteractionHandler),
	}

	if err := ix.LearnGlobalCommands(globals); err != nil {
		return ix, errors.Wrap(err, "could not load dispatch table")
	}

	return ix, nil
}

func (i *InteractionDispatcher) GlobalCommands() []entity.ApplicationCommand {
	cmds := make([]entity.ApplicationCommand, 0, len(i.globals))
	for _, ich := range i.globals {
		cmds = append(cmds, ich.Command)
	}

	return cmds
}

func (i *InteractionDispatcher) GuildCommands() map[snowflake.Snowflake][]entity.ApplicationCommand {
	cmds := make(map[snowflake.Snowflake][]entity.ApplicationCommand, len(i.guilds))

	for gid, gcmds := range i.guilds {
		cmds[gid] = make([]entity.ApplicationCommand, 0, len(gcmds))
		for _, ich := range gcmds {
			cmds[gid] = append(cmds[gid], ich.Command)
		}
	}

	return cmds
}

func (i *InteractionDispatcher) LearnGlobalCommands(cmds []InteractionCommandHandler) error {
	tmp := make(map[string]InteractionCommandHandler, len(i.globals)+len(cmds))

	for _, ich := range i.globals {
		tmp[ich.Command.Name] = ich
	}

	for _, ich := range cmds {
		tmp[ich.Command.Name] = ich
	}

	i.globals = make([]InteractionCommandHandler, 0, len(tmp))
	for _, ich := range tmp {
		i.globals = append(i.globals, ich)
	}

	return i.learnCommands(0, cmds)
}

func (i *InteractionDispatcher) LearnGuildCommands(gid snowflake.Snowflake, cmds []InteractionCommandHandler) error {
	tmp := make(map[string]InteractionCommandHandler, len(i.guilds[gid])+len(cmds))

	for _, ich := range i.guilds[gid] {
		tmp[ich.Command.Name] = ich
	}

	for _, ich := range cmds {
		tmp[ich.Command.Name] = ich
	}

	i.guilds[gid] = make([]InteractionCommandHandler, 0, len(tmp))
	for _, ich := range tmp {
		i.guilds[gid] = append(i.guilds[gid], ich)
	}

	return i.learnCommands(gid, cmds)
}

func (i *InteractionDispatcher) learnCommands(gid snowflake.Snowflake, cmds []InteractionCommandHandler) error {
	for _, ich := range cmds {
		handlers, ok := i.dispatch[ich.Command.Name]
		if !ok {
			handlers = make(map[snowflake.Snowflake]InteractionHandler)
		}

		i.dispatch[ich.Command.Name] = handlers
	}

	return nil
}

func (i *InteractionDispatcher) Dispatch(ix Interaction) (Response, error) {
	if ix.Data == nil {
		return nil, errors.WithDetails(ErrMalformedInteraction, "reason", "nil Data")
	}

	handlers, ok := i.dispatch[ix.Data.Name]
	if !ok {
		return nil, ErrMissingHandler
	}

	handler, ok := handlers[ix.GuildIDSnowflake]
	if !ok {
		handler, ok = handlers[0]
		if !ok {
			return nil, ErrMissingHandler
		}
	}

	return handler.HandleInteraction(ix)
}
