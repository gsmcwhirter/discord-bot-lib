package main

import (
	"context"
	"fmt"

	_ "net/http/pprof"

	"github.com/gsmcwhirter/go-util/v5/deferutil"
	"github.com/gsmcwhirter/go-util/v5/pprofsidecar"

	"github.com/gsmcwhirter/discord-bot-lib/v11/bot"
)

type config struct {
	PProfHostPort string `mapstructure:"pprof_hostport"`
}

func start(c config) error {
	fmt.Printf("%+v\n", c)

	conf := bot.Config{
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

	b := bot.NewDiscordBot(deps, conf)
	err = b.AuthenticateAndConnect()
	if err != nil {
		return err
	}
	defer deferutil.CheckDefer(b.Disconnect)

	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return pprofsidecar.Run(ctx, c.PProfHostPort, nil, b.Run)
}
