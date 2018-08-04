package payloads

import (
	"github.com/gsmcwhirter/discord-bot-lib/discordapi/constants"
	"github.com/gsmcwhirter/discord-bot-lib/discordapi/etfapi"
	"github.com/pkg/errors"
)

// HeartbeatPayload TODOC
type HeartbeatPayload struct {
	Sequence int
}

// Payload TODOC
func (hp HeartbeatPayload) Payload() (p etfapi.Payload, err error) {
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

// HeartbeatAckPayload TODOC
type HeartbeatAckPayload struct {
}

// Payload TODOC
func (hp HeartbeatAckPayload) Payload() (p etfapi.Payload, err error) {
	p.OpCode = constants.Heartbeat
	return
}
