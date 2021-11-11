package bot

import (
	"context"

	"github.com/gsmcwhirter/discord-bot-lib/v22/discordapi/etfapi"
	"github.com/gsmcwhirter/discord-bot-lib/v22/snowflake"
	"github.com/gsmcwhirter/discord-bot-lib/v22/wsapi"
)

type Logger = interface {
	Log(keyvals ...interface{}) error
	Message(string, ...interface{})
	Err(string, error, ...interface{})
	Printf(string, ...interface{})
}

type Payload = interface {
	EventName() string
	Contents() map[string]etfapi.Element
}

// DispatcherFunc is the api that a bot expects a handler function to have
type DispatchHandlerFunc = func(Payload, wsapi.WSMessage, chan<- wsapi.WSMessage) snowflake.Snowflake

// Dispatcher is the api that a bot expects a handler manager to have
type Dispatcher interface {
	ConnectToBot(*DiscordBot)
	GenerateHeartbeat(context.Context, int) (wsapi.WSMessage, error)
	AddHandler(string, DispatchHandlerFunc)
	HandleRequest(wsapi.WSMessage, chan<- wsapi.WSMessage) snowflake.Snowflake
}
