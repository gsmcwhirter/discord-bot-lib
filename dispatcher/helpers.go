package dispatcher

import (
	"context"

	"github.com/gsmcwhirter/go-util/v10/errors"

	"github.com/gsmcwhirter/discord-bot-lib/v24/discordapi/etfapi"
	"github.com/gsmcwhirter/discord-bot-lib/v24/wsapi"
)

// ETFPayload is the interface that a specialized etf api payload conforms to
type ETFPayload interface {
	Payload() (etfapi.Payload, error)
}

// ETFPayloadToMessage converts a specialized etf payload to a websocket message
func ETFPayloadToMessage(ctx context.Context, ep ETFPayload) (wsapi.WSMessage, error) {
	var m wsapi.WSMessage

	p, err := ep.Payload()
	if err != nil {
		return m, errors.Wrap(err, "could not construct Payload")
	}

	m.Ctx = ctx
	m.MessageType = wsapi.Binary
	m.MessageContents, err = p.Marshal()
	return m, errors.Wrap(err, "could not marshal payload")
}
