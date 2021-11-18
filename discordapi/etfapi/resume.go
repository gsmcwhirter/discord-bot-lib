package etfapi

import (
	"github.com/gsmcwhirter/go-util/v8/errors"

	"github.com/gsmcwhirter/discord-bot-lib/v23/discordapi"
)

// ResumePayload is the specialized payload for sending "Resume" events to the discord gateway websocket
type ResumePayload struct {
	Token     string
	SessionID string
	SeqNum    int
}

// Payload converts the specialized payload to a generic Payload
func (rp *ResumePayload) Payload() (Payload, error) {
	p := Payload{
		OpCode: discordapi.Resume,
		Data:   map[string]Element{},
	}

	var err error

	// TOKEN
	p.Data["token"], err = NewStringElement(rp.Token)
	if err != nil {
		return p, errors.Wrap(err, "could not create Element for token")
	}

	p.Data["session_id"], err = NewStringElement(rp.SessionID)
	if err != nil {
		return p, errors.Wrap(err, "could not create Element for session_id")
	}

	p.Data["seq"], err = NewInt32Element(rp.SeqNum)
	if err != nil {
		return p, errors.Wrap(err, "could not create Element for seq")
	}

	return p, nil
}
