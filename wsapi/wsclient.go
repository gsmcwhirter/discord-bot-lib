package wsapi

import (
	"context"
)

// WSClient is the api for a client that maintains an active websocket connection and hands
// off messages to be processed.
type WSClient interface {
	Connect(string, string) error
	Close()
	HandleRequests(context.Context, MessageHandler) error
	SendMessage(msg WSMessage)
}
