package dispatcher

import (
	"github.com/gsmcwhirter/discord-bot-lib/v21/discordapi/etfapi"
	"github.com/gsmcwhirter/discord-bot-lib/v21/snowflake"
	"github.com/gsmcwhirter/discord-bot-lib/v21/wsapi"
)

type Payload = interface {
	EventName() string
	Contents() map[string]etfapi.Element
}

// DispatcherFunc is the api that a bot expects a handler function to have
type DispatchHandlerFunc = func(Payload, wsapi.WSMessage, chan<- wsapi.WSMessage) snowflake.Snowflake
