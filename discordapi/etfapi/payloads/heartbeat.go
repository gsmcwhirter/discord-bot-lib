package payloads

import (
	"github.com/gsmcwhirter/discord-bot-lib/discordapi/constants"
	"github.com/gsmcwhirter/discord-bot-lib/discordapi/etfapi"
	"github.com/pkg/errors"
)

// HeartbeatPayload is the specialized payload for sending "heartbeat" events to the discord gateway websocket
type HeartbeatPayload struct {
	Sequence int
}

// Payload converts the specialized payload to a generic etfapi.Payload
func (hp *HeartbeatPayload) Payload() (p etfapi.Payload, err error) {
	p.OpCode = constants.Heartbeat
	p.Data = map[string]etfapi.Element{}

	if hp.Sequence < 0 {
		p.Data["d"], err = etfapi.NewStringElement("nil")
		if err != nil {
			err = errors.Wrap(err, "could not create an Element for nil lastSeq")
			return
		}
	} else {
		p.Data["d"], err = etfapi.NewInt32Element(hp.Sequence)
		if err != nil {
			err = errors.Wrap(err, "could not create an Element for non-nil lastSeq")
			return
		}
	}

	return
}
