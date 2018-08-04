package wsclient

import (
	"context"
	"fmt"
)

// WSMessage TODOC
type WSMessage struct {
	Ctx             context.Context
	MessageType     MessageType
	MessageContents []byte
}

func (m WSMessage) String() string {
	return fmt.Sprintf("WSMessage{Type=%v, Contents=%v}", m.MessageType, m.MessageContents)
}
