package cmdhandler

import (
	"github.com/gsmcwhirter/discord-bot-lib/v22/discordapi/entity"
	"github.com/gsmcwhirter/discord-bot-lib/v22/snowflake"
)

type InteractionDispatcher struct {
	globals []InteractionCommandHandler
	guilds  map[snowflake.Snowflake][]InteractionCommandHandler
}

type InteractionCommandHandler struct {
	Command entity.ApplicationCommand
	Handler InteractionHandler
}

func NewInteractionDispatcher(globals []InteractionCommandHandler) (*InteractionDispatcher, error) {
	return &InteractionDispatcher{
		globals: globals,
		guilds:  make(map[snowflake.Snowflake][]InteractionCommandHandler),
	}, nil
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
