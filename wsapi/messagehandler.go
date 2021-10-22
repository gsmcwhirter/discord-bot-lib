package wsapi

import (
	"github.com/gsmcwhirter/discord-bot-lib/v20/snowflake"
)

// MessageHandler is the api of an object that has a HandleRequest MessageHandlerFunc method
type MessageHandler interface {
	HandleRequest(WSMessage, chan<- WSMessage) snowflake.Snowflake
}

// MessageHandlerFunc is the api of a handler that processes a WSMessage
// and sends a response back over the response channe;
type MessageHandlerFunc func(m WSMessage, resp chan<- WSMessage) snowflake.Snowflake

type messageHandler struct {
	handler MessageHandlerFunc
}

// NewMessageHandler creates a MessageHandler from the given MessageHandlerFunc
func NewMessageHandler(h MessageHandlerFunc) MessageHandler {
	return &messageHandler{handler: h}
}

func (mh *messageHandler) HandleRequest(req WSMessage, resp chan<- WSMessage) snowflake.Snowflake {
	return mh.handler(req, resp)
}
