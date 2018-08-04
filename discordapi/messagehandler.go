package discordapi

import (
	"fmt"
	"sync"

	"github.com/go-kit/kit/log/level"

	"github.com/gsmcwhirter/discord-bot-lib/discordapi/constants"
	"github.com/gsmcwhirter/discord-bot-lib/discordapi/etfapi"
	"github.com/gsmcwhirter/discord-bot-lib/discordapi/etfapi/payloads"
	"github.com/gsmcwhirter/discord-bot-lib/logging"
	"github.com/gsmcwhirter/discord-bot-lib/wsclient"
)

// DiscordMessageHandlerFunc TODOC
type DiscordMessageHandlerFunc func(*etfapi.Payload, wsclient.WSMessage, chan<- wsclient.WSMessage)

type discordMessageHandler struct {
	bot            *discordBot
	opCodeDispatch map[constants.OpCode]DiscordMessageHandlerFunc

	dispatcherLock *sync.Mutex
	eventDispatch  map[string][]DiscordMessageHandlerFunc
}

func noop(p *etfapi.Payload, req wsclient.WSMessage, resp chan<- wsclient.WSMessage) {
}

// NewDiscordMessageHandler TODOC
func newDiscordMessageHandler(bot *discordBot) *discordMessageHandler {
	c := discordMessageHandler{
		bot:            bot,
		dispatcherLock: &sync.Mutex{},
	}

	c.opCodeDispatch = map[constants.OpCode]DiscordMessageHandlerFunc{
		constants.Hello:           c.handleHello,
		constants.Heartbeat:       c.handleHeartbeat,
		constants.HeartbeatAck:    noop,
		constants.InvalidSession:  nil,
		constants.InvalidSequence: nil,
		constants.Reconnect:       nil,
		constants.Dispatch:        c.handleDispatch,
	}

	c.eventDispatch = map[string][]DiscordMessageHandlerFunc{
		"READY":        []DiscordMessageHandlerFunc{c.handleReady},
		"GUILD_CREATE": []DiscordMessageHandlerFunc{c.handleGuildCreate},
	}

	return &c
}

func (c *discordMessageHandler) addHandler(event string, handler DiscordMessageHandlerFunc) {
	c.dispatcherLock.Lock()
	defer c.dispatcherLock.Unlock()

	handlers := c.eventDispatch[event]
	c.eventDispatch[event] = append(handlers, handler)
}

func (c *discordMessageHandler) HandleRequest(req wsclient.WSMessage, resp chan<- wsclient.WSMessage) {
	logger := logging.WithContext(req.Ctx, c.bot.deps.Logger())
	_ = level.Debug(logger).Log("message", "discordapi dispatching request")

	select {
	case <-req.Ctx.Done():
		_ = level.Debug(logger).Log("message", "discordapi already done. skipping request")
		return
	default:
	}

	_ = level.Debug(logger).Log("message", "processing server message", "ws_msg", fmt.Sprintf("%v", req.MessageContents))

	p, err := etfapi.Unmarshal(req.MessageContents)
	if err != nil {
		_ = level.Error(logger).Log("message", "error unmarshaling payload", "error", err, "ws_msg", fmt.Sprintf("%v", req.MessageContents))
		return
	}

	if p.SeqNum != nil {
		c.bot.updateSequence(*p.SeqNum)
	}

	_ = level.Info(logger).Log("message", "received payload", "payload", p)

	opHandler, ok := c.opCodeDispatch[p.OpCode]
	if !ok {
		_ = level.Error(logger).Log("message", "unrecognized OpCode", "op_code", p.OpCode)
		return
	}

	if opHandler == nil {
		_ = level.Error(logger).Log("message", "no handler for OpCode", "op_code", p.OpCode)
		return
	}

	opHandler(p, req, resp)
}

func (c *discordMessageHandler) handleError(req wsclient.WSMessage, resp chan<- wsclient.WSMessage) {
	select {
	case <-req.Ctx.Done():
		return
	default:
	}

	logger := logging.WithContext(req.Ctx, c.bot.deps.Logger())
	_ = level.Error(logger).Log("message", "error code received from websocket", "ws_msg", req)
}

func (c *discordMessageHandler) handleHello(p *etfapi.Payload, req wsclient.WSMessage, resp chan<- wsclient.WSMessage) {
	select {
	case <-req.Ctx.Done():
		return
	default:
	}

	logger := logging.WithContext(req.Ctx, c.bot.deps.Logger())
	rawInterval, ok := p.Data["heartbeat_interval"]

	if ok {
		// set heartbeat stuff
		interval, err := rawInterval.ToInt()
		if err != nil {
			_ = level.Error(logger).Log("message", "error handling hello heartbeat config", "err", err)
			return
		}

		_ = level.Info(logger).Log("message", "configuring heartbeat", "interval", interval)
		c.bot.heartbeats <- hbReconfig{
			ctx:      req.Ctx,
			interval: interval,
		}
		_ = level.Debug(logger).Log("message", "configuring heartbeat done")

	}

	// send identify
	var m wsclient.WSMessage
	var err error

	sessID := c.bot.session.ID()
	if sessID != "" {
		_ = level.Info(logger).Log("message", "generating resume payload")
		rp := payloads.ResumePayload{
			Token:     c.bot.config.BotToken,
			SessionID: sessID,
			SeqNum:    c.bot.LastSequence(),
		}

		m, err = payloads.ETFPayloadToMessage(req.Ctx, rp)
	} else {
		_ = level.Info(logger).Log("message", "generating identify payload")
		ip := payloads.IdentifyPayload{
			Token: c.bot.config.BotToken,
			Properties: payloads.IdentifyPayloadProperties{
				OS:      "linux",
				Browser: "eso-have-want-bot#0286",
				Device:  "eso-have-want-bot#0286",
			},
			LargeThreshold: 250,
			Shard: payloads.IdentifyPayloadShard{
				ID:    0,
				MaxID: 0,
			},
			Presence: payloads.IdentifyPayloadPresence{
				Game: payloads.IdentifyPayloadGame{
					Name: "List Manager 2018",
					Type: 0,
				},
				Status: "online",
				Since:  0,
				AFK:    false,
			},
		}

		m, err = payloads.ETFPayloadToMessage(req.Ctx, ip)
	}

	if err != nil {
		_ = level.Error(logger).Log("message", "error generating identify/resume payload", "err", err)
		return
	}

	err = c.bot.deps.MessageRateLimiter().Wait(req.Ctx)
	if err != nil {
		_ = level.Error(logger).Log("message", "error ratelimiting", "err", err)
		return
	}

	_ = level.Info(logger).Log("message", "sending identify/resume to channel")
	_ = level.Debug(logger).Log("message", "sending response to channel", "message", m, "msg_len", len(m.MessageContents))
	resp <- m
}

func (c *discordMessageHandler) handleHeartbeat(p *etfapi.Payload, req wsclient.WSMessage, resp chan<- wsclient.WSMessage) {
	select {
	case <-req.Ctx.Done():
		return
	default:
	}

	logger := logging.WithContext(req.Ctx, c.bot.deps.Logger())
	_ = level.Info(logger).Log("message", "requesting manual heartbeat")
	c.bot.heartbeats <- hbReconfig{
		ctx: req.Ctx,
	}
	_ = level.Debug(logger).Log("message", "manual heartbeat done")
}

func (c *discordMessageHandler) handleDispatch(p *etfapi.Payload, req wsclient.WSMessage, resp chan<- wsclient.WSMessage) {
	select {
	case <-req.Ctx.Done():
		return
	default:
	}

	logger := logging.WithContext(req.Ctx, c.bot.deps.Logger())

	c.dispatcherLock.Lock()
	eventHandlers, ok := c.eventDispatch[p.EventName]
	c.dispatcherLock.Unlock()

	if ok {
		_ = level.Info(logger).Log("message", "processing event", "event_name", p.EventName)
		for _, eventHandler := range eventHandlers {
			eventHandler(p, req, resp)
		}
	}
}

func (c *discordMessageHandler) handleReady(p *etfapi.Payload, req wsclient.WSMessage, resp chan<- wsclient.WSMessage) {
	select {
	case <-req.Ctx.Done():
		return
	default:
	}

	logger := logging.WithContext(req.Ctx, c.bot.deps.Logger())

	err := c.bot.session.UpdateFromReady(p)
	if err != nil {
		_ = level.Error(logger).Log("message", "error setting up session", "err", err)
		return
	}
}

func (c *discordMessageHandler) handleGuildCreate(p *etfapi.Payload, req wsclient.WSMessage, resp chan<- wsclient.WSMessage) {
	select {
	case <-req.Ctx.Done():
		return
	default:
	}

	logger := logging.WithContext(req.Ctx, c.bot.deps.Logger())
	_ = level.Info(logger).Log("message", "upserting guild", "pdata", fmt.Sprintf("%+v", p.Data), "tag", "GUILD_CREATE")
	err := c.bot.session.UpsertGuildFromElementMap(p.Data)
	if err != nil {
		_ = level.Error(logger).Log("message", "error processing guild create", "err", err)
		return
	}
}
