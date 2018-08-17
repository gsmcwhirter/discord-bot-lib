package constants

import "fmt"

// OpCode is a type alias for discord OpCode values
type OpCode int

// OpCode names
const (
	Heartbeat      OpCode = 1
	HeartbeatAck          = 11
	Identify              = 2
	InvalidSession        = 9
	Resume                = 6
	Dispatch              = 0
	StatusUpdate          = 3
	Reconnect             = 7
	Hello                 = 10
)

func (c OpCode) String() string {
	switch c {
	case Heartbeat:
		return "Heartbeat"
	case HeartbeatAck:
		return "HeartbeatAck"
	case Identify:
		return "Identify"
	case InvalidSession:
		return "InvalidSession"
	case Resume:
		return "Resume"
	case Dispatch:
		return "Dispatch"
	case StatusUpdate:
		return "StatusUpdate"
	case Reconnect:
		return "Reconnect"
	case Hello:
		return "Hello"
	default:
		return fmt.Sprintf("(unknown: %d)", int(c))
	}
}
