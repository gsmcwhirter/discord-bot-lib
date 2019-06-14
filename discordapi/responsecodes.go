package discordapi

import "fmt"

// ResponseCode is a type alias for websocket error codes
type ResponseCode int

// ResponseCode names
const (
	UnknownError         ResponseCode = 4000
	UnknownOpcode        ResponseCode = 4001
	DecodeError          ResponseCode = 4002
	NotAuthenticated     ResponseCode = 4003
	AuthenticationFailed ResponseCode = 4004
	AlreadyAuthenticated ResponseCode = 4005
	InvalidSequence      ResponseCode = 4007
	RateLimited          ResponseCode = 4008
	SessionTimeout       ResponseCode = 4009
	InvalidShard         ResponseCode = 4010
	ShardingRequired     ResponseCode = 4011
)

func (c ResponseCode) String() string {
	switch c {
	case UnknownError:
		return "UnknownError"
	case UnknownOpcode:
		return "UnknownOpcode"
	case DecodeError:
		return "DecodeError"
	case NotAuthenticated:
		return "NotAuthenticated"
	case AuthenticationFailed:
		return "AuthenticationFailed"
	case AlreadyAuthenticated:
		return "AlreadyAuthenticated"
	case InvalidSequence:
		return "InvalidSequence"
	case RateLimited:
		return "RateLimited"
	case SessionTimeout:
		return "SessionTimeout"
	case InvalidShard:
		return "InvalidShard"
	case ShardingRequired:
		return "ShardingRequired"
	default:
		return fmt.Sprintf("(unknown: %d)", int(c))
	}
}
