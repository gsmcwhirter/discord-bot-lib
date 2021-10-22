package dispatcher

import (
	"github.com/gsmcwhirter/discord-bot-lib/v20/discordapi/etf"
	"github.com/gsmcwhirter/discord-bot-lib/v20/snowflake"
	"github.com/gsmcwhirter/discord-bot-lib/v20/wsapi"
)

type Payload = interface {
	EventName() string
	Contents() map[string]etf.Element
}

// DispatcherFunc is the api that a bot expects a handler function to have
type DispatchHandlerFunc = func(Payload, wsapi.WSMessage, chan<- wsapi.WSMessage) snowflake.Snowflake
