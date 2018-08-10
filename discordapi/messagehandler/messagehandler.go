package messagehandler

import (
	"fmt"
	"sync"

	"golang.org/x/time/rate"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/gsmcwhirter/discord-bot-lib/discordapi"
	"github.com/gsmcwhirter/discord-bot-lib/discordapi/constants"
	"github.com/gsmcwhirter/discord-bot-lib/discordapi/etfapi"
	"github.com/gsmcwhirter/discord-bot-lib/discordapi/etfapi/payloads"
	"github.com/gsmcwhirter/discord-bot-lib/discordapi/session"
	"github.com/gsmcwhirter/discord-bot-lib/logging"
	"github.com/gsmcwhirter/discord-bot-lib/wsclient"
)

type dependencies interface {
	Logger() log.Logger
	BotSession() *session.Session
	MessageRateLimiter() *rate.Limiter
}

// DiscordMessageHandler TODOC
type discordMessageHandler struct {
	deps           dependencies
	bot            discordapi.DiscordBot
	opCodeDispatch map[constants.OpCode]discordapi.DiscordMessageHandlerFunc

	dispatcherLock *sync.Mutex
	eventDispatch  map[string][]discordapi.DiscordMessageHandlerFunc
}

func noop(p *etfapi.Payload, req wsclient.WSMessage, resp chan<- wsclient.WSMessage) {
}

// NewDiscordMessageHandler TODOC
func NewDiscordMessageHandler(deps dependencies) discordapi.DiscordMessageHandler {
	c := discordMessageHandler{
		deps:           deps,
		dispatcherLock: &sync.Mutex{},
	}

	c.opCodeDispatch = map[constants.OpCode]discordapi.DiscordMessageHandlerFunc{
		constants.Hello:           c.handleHello,
		constants.Heartbeat:       c.handleHeartbeat,
		constants.HeartbeatAck:    noop,
		constants.InvalidSession:  nil,
		constants.InvalidSequence: nil,
		constants.Reconnect:       nil,
		constants.Dispatch:        c.handleDispatch,
	}

	c.eventDispatch = map[string][]discordapi.DiscordMessageHandlerFunc{
		"READY":        []discordapi.DiscordMessageHandlerFunc{c.handleReady},
		"GUILD_CREATE": []discordapi.DiscordMessageHandlerFunc{c.handleGuildCreate},
	}

	return &c
}

// ConnectToBot TODOC
func (c *discordMessageHandler) ConnectToBot(bot discordapi.DiscordBot) {
	c.bot = bot
}

// AddHandler TODOC
func (c *discordMessageHandler) AddHandler(event string, handler discordapi.DiscordMessageHandlerFunc) {
	c.dispatcherLock.Lock()
	defer c.dispatcherLock.Unlock()

	handlers := c.eventDispatch[event]
	c.eventDispatch[event] = append(handlers, handler)
}

// HandleRequest
func (c *discordMessageHandler) HandleRequest(req wsclient.WSMessage, resp chan<- wsclient.WSMessage) {
	logger := logging.WithContext(req.Ctx, c.deps.Logger())
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
		c.bot.UpdateSequence(*p.SeqNum)
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

	logger := logging.WithContext(req.Ctx, c.deps.Logger())
	_ = level.Error(logger).Log("message", "error code received from websocket", "ws_msg", req)
}

func (c *discordMessageHandler) handleHello(p *etfapi.Payload, req wsclient.WSMessage, resp chan<- wsclient.WSMessage) {
	select {
	case <-req.Ctx.Done():
		return
	default:
	}

	logger := logging.WithContext(req.Ctx, c.deps.Logger())
	rawInterval, ok := p.Data["heartbeat_interval"]

	if ok {
		// set heartbeat stuff
		interval, err := rawInterval.ToInt()
		if err != nil {
			_ = level.Error(logger).Log("message", "error handling hello heartbeat config", "err", err)
			return
		}

		_ = level.Info(logger).Log("message", "configuring heartbeat", "interval", interval)
		c.bot.ReconfigureHeartbeat(req.Ctx, interval)
		_ = level.Debug(logger).Log("message", "configuring heartbeat done")
	}

	// send identify
	var m wsclient.WSMessage
	var err error

	sessID := c.deps.BotSession().ID()
	if sessID != "" {
		_ = level.Info(logger).Log("message", "generating resume payload")
		rp := payloads.ResumePayload{
			Token:     c.bot.Config().BotToken,
			SessionID: sessID,
			SeqNum:    c.bot.LastSequence(),
		}

		m, err = payloads.ETFPayloadToMessage(req.Ctx, rp)
	} else {
		_ = level.Info(logger).Log("message", "generating identify payload")
		ip := payloads.IdentifyPayload{
			Token: c.bot.Config().BotToken,
			Properties: payloads.IdentifyPayloadProperties{
				OS:      c.bot.Config().OS,
				Browser: c.bot.Config().BotName,
				Device:  c.bot.Config().BotName,
			},
			LargeThreshold: 250,
			Shard: payloads.IdentifyPayloadShard{
				ID:    0,
				MaxID: 0,
			},
			Presence: payloads.IdentifyPayloadPresence{
				Game: payloads.IdentifyPayloadGame{
					Name: c.bot.Config().BotPresence,
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

	err = c.deps.MessageRateLimiter().Wait(req.Ctx)
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

	logger := logging.WithContext(req.Ctx, c.deps.Logger())
	_ = level.Info(logger).Log("message", "requesting manual heartbeat")
	c.bot.ReconfigureHeartbeat(req.Ctx, 0)
	_ = level.Debug(logger).Log("message", "manual heartbeat done")
}

func (c *discordMessageHandler) handleDispatch(p *etfapi.Payload, req wsclient.WSMessage, resp chan<- wsclient.WSMessage) {
	select {
	case <-req.Ctx.Done():
		return
	default:
	}

	logger := logging.WithContext(req.Ctx, c.deps.Logger())

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

	logger := logging.WithContext(req.Ctx, c.deps.Logger())

	err := c.deps.BotSession().UpdateFromReady(p)
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

	logger := logging.WithContext(req.Ctx, c.deps.Logger())
	_ = level.Info(logger).Log("message", "upserting guild", "pdata", fmt.Sprintf("%+v", p.Data), "tag", "GUILD_CREATE")
	err := c.deps.BotSession().UpsertGuildFromElementMap(p.Data)
	if err != nil {
		_ = level.Error(logger).Log("message", "error processing guild create", "err", err)
		return
	}
}
