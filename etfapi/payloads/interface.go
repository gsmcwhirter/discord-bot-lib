package payloads

import (
	"context"

	"github.com/gsmcwhirter/discord-bot-lib/etfapi"
	"github.com/gsmcwhirter/discord-bot-lib/wsclient"
	"github.com/pkg/errors"
)

// ETFPayload is the interface that a specialized etf api payload conforms to
type ETFPayload interface {
	Payload() (etfapi.Payload, error)
}

// ETFPayloadToMessage converts a specialized etf payload to a websocket message
func ETFPayloadToMessage(ctx context.Context, ep ETFPayload) (m wsclient.WSMessage, err error) {
	var p etfapi.Payload

	p, err = ep.Payload()
	if err != nil {
		err = errors.Wrap(err, "could not construct Payload")
		return
	}

	m.Ctx = ctx
	m.MessageType = wsclient.Binary
	m.MessageContents, err = p.Marshal()
	err = errors.Wrap(err, "could not marshal payload")
	return
}
