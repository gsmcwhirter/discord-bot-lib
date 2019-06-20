package payloads

import (
	"github.com/gsmcwhirter/go-util/v5/errors"

	"github.com/gsmcwhirter/discord-bot-lib/v10/discordapi"
	"github.com/gsmcwhirter/discord-bot-lib/v10/etfapi"
)

// ResumePayload is the specialized payload for sending "Resume" events to the discord gateway websocket
type ResumePayload struct {
	Token     string
	SessionID string
	SeqNum    int
}

// Payload converts the specialized payload to a generic etfapi.Payload
func (rp *ResumePayload) Payload() (etfapi.Payload, error) {
	p := etfapi.Payload{
		OpCode: discordapi.Resume,
		Data:   map[string]etfapi.Element{},
	}

	var err error

	// TOKEN
	p.Data["token"], err = etfapi.NewStringElement(rp.Token)
	if err != nil {
		return p, errors.Wrap(err, "could not create Element for token")
	}

	p.Data["session_id"], err = etfapi.NewStringElement(rp.SessionID)
	if err != nil {
		return p, errors.Wrap(err, "could not create Element for session_id")
	}

	p.Data["seq"], err = etfapi.NewInt32Element(rp.SeqNum)
	if err != nil {
		return p, errors.Wrap(err, "could not create Element for seq")
	}

	return p, nil
}
