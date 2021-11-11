package etfapi

import (
	"github.com/gsmcwhirter/go-util/v8/errors"

	"github.com/gsmcwhirter/discord-bot-lib/v22/discordapi"
)

// HeartbeatPayload is the specialized payload for sending "heartbeat" events to the discord gateway websocket
type HeartbeatPayload struct {
	Sequence int
}

// Payload converts the specialized payload to a generic etfapi.Payload
func (hp *HeartbeatPayload) Payload() (Payload, error) {
	var err error
	p := Payload{
		OpCode: discordapi.Heartbeat,
		Data:   map[string]Element{},
	}

	if hp.Sequence < 0 {
		p.Data["d"], err = NewStringElement("nil")
		if err != nil {
			return p, errors.Wrap(err, "could not create an Element for nil lastSeq")
		}
	} else {
		p.Data["d"], err = NewInt32Element(hp.Sequence)
		if err != nil {
			return p, errors.Wrap(err, "could not create an Element for non-nil lastSeq")
		}
	}

	return p, nil
}
