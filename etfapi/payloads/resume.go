package payloads

import (
	"github.com/gsmcwhirter/discord-bot-lib/discordapi"
	"github.com/gsmcwhirter/discord-bot-lib/etfapi"
	"github.com/pkg/errors"
)

// ResumePayload is the specialized payload for sending "Resume" events to the discord gateway websocket
type ResumePayload struct {
	Token     string
	SessionID string
	SeqNum    int
}

// Payload converts the specialized payload to a generic etfapi.Payload
func (rp *ResumePayload) Payload() (p etfapi.Payload, err error) {
	p.OpCode = discordapi.Resume
	p.Data = map[string]etfapi.Element{}

	// TOKEN
	p.Data["token"], err = etfapi.NewStringElement(rp.Token)
	if err != nil {
		err = errors.Wrap(err, "could not create Element for token")
		return
	}

	p.Data["session_id"], err = etfapi.NewStringElement(rp.SessionID)
	if err != nil {
		err = errors.Wrap(err, "could not create Element for session_id")
		return
	}

	p.Data["seq"], err = etfapi.NewInt32Element(rp.SeqNum)
	if err != nil {
		err = errors.Wrap(err, "could not create Element for seq")
		return
	}

	return
}
