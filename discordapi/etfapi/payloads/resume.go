package payloads

import (
	"github.com/gsmcwhirter/discord-bot-lib/discordapi/constants"
	"github.com/gsmcwhirter/discord-bot-lib/discordapi/etfapi"
	"github.com/pkg/errors"
)

// ResumePayload TODOC
type ResumePayload struct {
	Token     string
	SessionID string
	SeqNum    int
}

// Payload TODOC
func (rp ResumePayload) Payload() (p etfapi.Payload, err error) {
	p.OpCode = constants.Resume
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
