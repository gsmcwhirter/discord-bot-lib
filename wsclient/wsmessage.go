package wsclient

import (
	"context"
	"fmt"
)

// WSMessage is a  struct that wraps a message received or to be sent over a websocket connection
type WSMessage struct {
	Ctx             context.Context
	MessageType     MessageType
	MessageContents []byte
}

func (m WSMessage) String() string {
	return fmt.Sprintf("WSMessage{Type=%v, Contents=%v}", m.MessageType, m.MessageContents)
}
