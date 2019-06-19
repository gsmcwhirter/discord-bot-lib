package main

import (
	"github.com/gsmcwhirter/discord-bot-lib/v9/discordapi"
	"github.com/gsmcwhirter/discord-bot-lib/v9/etfapi"
	"github.com/gsmcwhirter/discord-bot-lib/v9/snowflake"
)

func guildCreate(id snowflake.Snowflake, name string) ([]byte, error) {
	eMap := map[string]etfapi.Element{}

	idE, err := etfapi.NewSmallBigElement(int64(id))
	if err != nil {
		return nil, err
	}

	eMap["id"] = idE

	nameE, err := etfapi.NewStringElement(name)
	if err != nil {
		return nil, err
	}
	eMap["name"] = nameE

	p := etfapi.Payload{
		OpCode:    discordapi.Dispatch,
		EventName: "GUILD_CREATE",
		Data:      eMap,
	}

	return p.Marshal()
}
