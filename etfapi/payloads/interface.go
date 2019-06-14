package payloads

import (
	"context"

	"github.com/gsmcwhirter/go-util/v3/errors"

	"github.com/gsmcwhirter/discord-bot-lib/v8/etfapi"
	"github.com/gsmcwhirter/discord-bot-lib/v8/wsclient"
)

// ETFPayload is the interface that a specialized etf api payload conforms to
type ETFPayload interface {
	Payload() (etfapi.Payload, error)
}

// ETFPayloadToMessage converts a specialized etf payload to a websocket message
func ETFPayloadToMessage(ctx context.Context, ep ETFPayload) (wsclient.WSMessage, error) {
	var m wsclient.WSMessage

	p, err := ep.Payload()
	if err != nil {
		return m, errors.Wrap(err, "could not construct Payload")
	}

	m.Ctx = ctx
	m.MessageType = wsclient.Binary
	m.MessageContents, err = p.Marshal()
	return m, errors.Wrap(err, "could not marshal payload")
}
