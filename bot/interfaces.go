package bot

import (
	"context"
	"io"
	"net/http"

	"github.com/gsmcwhirter/discord-bot-lib/v19/discordapi/etf"
	"github.com/gsmcwhirter/discord-bot-lib/v19/snowflake"
	"github.com/gsmcwhirter/discord-bot-lib/v19/wsapi"
)

// JSONMarshaler is the interface implemented by types that
// can marshal themselves into valid JSON.
type JSONMarshaler interface {
	MarshalJSON() ([]byte, error)
}

// HTTPClient is the interface of an http client
type HTTPClient interface {
	SetHeaders(http.Header)
	Get(context.Context, string, *http.Header) (*http.Response, error)
	GetBody(context.Context, string, *http.Header) (*http.Response, []byte, error)
	GetJSON(context.Context, string, *http.Header, interface{}) (*http.Response, error)
	Post(context.Context, string, *http.Header, io.Reader) (*http.Response, error)
	PostBody(context.Context, string, *http.Header, io.Reader) (*http.Response, []byte, error)
	PostJSON(context.Context, string, *http.Header, io.Reader, interface{}) (*http.Response, error)
	Put(context.Context, string, *http.Header, io.Reader) (*http.Response, error)
	PutBody(context.Context, string, *http.Header, io.Reader) (*http.Response, []byte, error)
	PutJSON(context.Context, string, *http.Header, io.Reader, interface{}) (*http.Response, error)
}

type Payload interface {
	EventName() string
	Contents() map[string]etf.Element
}

// DispatcherFunc is the api that a bot expects a handler function to have
type DispatchHandlerFunc func(Payload, wsapi.WSMessage, chan<- wsapi.WSMessage) snowflake.Snowflake

// Dispatcher is the api that a bot expects a handler manager to have
type Dispatcher interface {
	ConnectToBot(*DiscordBot)
	GenerateHeartbeat(context.Context, int) (wsapi.WSMessage, error)
	AddHandler(string, DispatchHandlerFunc)
	HandleRequest(wsapi.WSMessage, chan<- wsapi.WSMessage) snowflake.Snowflake
}
