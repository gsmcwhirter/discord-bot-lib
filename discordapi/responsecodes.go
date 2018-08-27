package discordapi

import "fmt"

// ResponseCode is a type alias for websocket error codes
type ResponseCode int

// ResponseCode names
const (
	UnknownError         ResponseCode = 4000
	UnknownOpcode                     = 4001
	DecodeError                       = 4002
	NotAuthenticated                  = 4003
	AuthenticationFailed              = 4004
	AlreadyAuthenticated              = 4005
	InvalidSequence                   = 4007
	RateLimited                       = 4008
	SessionTimeout                    = 4009
	InvalidShard                      = 4010
	ShardingRequired                  = 4011
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
