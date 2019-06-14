package payloads

import (
	"github.com/pkg/errors"

	"github.com/gsmcwhirter/discord-bot-lib/v6/discordapi"
	"github.com/gsmcwhirter/discord-bot-lib/v6/etfapi"
)

// HeartbeatPayload is the specialized payload for sending "heartbeat" events to the discord gateway websocket
type HeartbeatPayload struct {
	Sequence int
}

// Payload converts the specialized payload to a generic etfapi.Payload
func (hp *HeartbeatPayload) Payload() (etfapi.Payload, error) {
	var err error
	p := etfapi.Payload{
		OpCode: discordapi.Heartbeat,
		Data:   map[string]etfapi.Element{},
	}

	if hp.Sequence < 0 {
		p.Data["d"], err = etfapi.NewStringElement("nil")
		if err != nil {
			return p, errors.Wrap(err, "could not create an Element for nil lastSeq")
		}
	} else {
		p.Data["d"], err = etfapi.NewInt32Element(hp.Sequence)
		if err != nil {
			return p, errors.Wrap(err, "could not create an Element for non-nil lastSeq")
		}
	}

	return p, nil
}
