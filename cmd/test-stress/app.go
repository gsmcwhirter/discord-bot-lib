package main

import (
	"context"
	"fmt"

	_ "net/http/pprof"

	"github.com/gsmcwhirter/discord-bot-lib/discordapi"
	"github.com/gsmcwhirter/go-util/pprofsidecar"
)

type config struct {
	PProfHostPort string `mapstructure:"pprof_hostport"`
}

func start(c config) error {
	fmt.Printf("%+v\n", c)

	conf := discordapi.BotConfig{
		ClientID:     "test id",
		ClientSecret: "test secret",
		BotToken:     "test token",
		APIURL:       "http://localhost",
		NumWorkers:   10,
		OS:           "Test OS",
		BotName:      "test bot",
		BotPresence:  "test presence",
	}

	deps, err := createDependencies(c, conf)
	if err != nil {
		return err
	}
	defer deps.Close()

	bot := discordapi.NewDiscordBot(deps, conf)
	err = bot.AuthenticateAndConnect()
	if err != nil {
		return err
	}
	defer bot.Disconnect()

	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return pprofsidecar.Run(ctx, c.PProfHostPort, nil, bot.Run)
}
