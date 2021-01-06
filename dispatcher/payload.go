package dispatcher

import (
	"github.com/gsmcwhirter/discord-bot-lib/v19/bot"
	"github.com/gsmcwhirter/discord-bot-lib/v19/discordapi/etf"
)

type Payload struct {
	data      map[string]etf.Element
	eventName string
}

var _ bot.Payload = (*Payload)(nil)

func (p *Payload) Contents() map[string]etf.Element { return p.data }
func (p *Payload) EventName() string                { return p.eventName }

func NewPayload(p *etf.Payload) *Payload {
	return &Payload{
		data:      p.Data,
		eventName: p.EventName,
	}
}
