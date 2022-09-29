package dispatcher

import (
	"github.com/gsmcwhirter/discord-bot-lib/v24/discordapi/etfapi"
	"github.com/gsmcwhirter/discord-bot-lib/v24/snowflake"
	"github.com/gsmcwhirter/discord-bot-lib/v24/wsapi"
)

// Payload is the interface expected for an etf payload
type Payload = interface {
	EventName() string
	Contents() map[string]etfapi.Element
}

// DispatchHandlerFunc is the api that a bot expects a handler function to have
type DispatchHandlerFunc = func(Payload, wsapi.WSMessage, chan<- wsapi.WSMessage) snowflake.Snowflake
