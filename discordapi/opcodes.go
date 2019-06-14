package discordapi

import "fmt"

// OpCode is a type alias for discord OpCode values
type OpCode int

// OpCode names
const (
	Heartbeat      OpCode = 1
	HeartbeatAck   OpCode = 11
	Identify       OpCode = 2
	InvalidSession OpCode = 9
	Resume         OpCode = 6
	Dispatch       OpCode = 0
	StatusUpdate   OpCode = 3
	Reconnect      OpCode = 7
	Hello          OpCode = 10
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
