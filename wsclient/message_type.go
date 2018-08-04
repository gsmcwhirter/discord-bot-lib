package wsclient

import "fmt"

// MessageType TODOC
type MessageType int

// Websocket Message types
const (
	Text   MessageType = 1
	Binary             = 2
)

func (t MessageType) String() string {
	switch t {
	case Text:
		return "Text"
	case Binary:
		return "Binary"
	default:
		return fmt.Sprintf("(unknown: %d)", int(t))
	}
}
