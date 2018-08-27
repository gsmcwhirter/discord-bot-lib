package main

import (
	"github.com/gsmcwhirter/discord-bot-lib/discordapi/constants"
	"github.com/gsmcwhirter/discord-bot-lib/discordapi/etfapi"
	"github.com/gsmcwhirter/discord-bot-lib/snowflake"
)

func guildCreate(id snowflake.Snowflake, name string) ([]byte, error) {
	eMap := map[string]etfapi.Element{}

	idE, err := etfapi.NewSmallBigElement(int64(id))
	if err != nil {
		return nil, err
	}

	eMap["id"] = idE

	p := etfapi.Payload{
		OpCode:    constants.Dispatch,
		EventName: "GUILD_CREATE",
		Data:      eMap,
	}

	return p.Marshal()
}
