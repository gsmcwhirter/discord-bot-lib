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
	"github.com/gsmcwhirter/discord-bot-lib/logging"
	"github.com/gsmcwhirter/discord-bot-lib/wsclient"
)

type dependencies interface {
	Logger() log.Logger
	BotSession() *etfapi.Session
	MessageRateLimiter() *rate.Limiter
}

type discordMessageHandler struct {
	deps           dependencies
	bot            discordapi.DiscordBot
	opCodeDispatch map[constants.OpCode]discordapi.DiscordMessageHandlerFunc

	dispatcherLock *sync.Mutex
	eventDispatch  map[string][]discordapi.DiscordMessageHandlerFunc
}

func noop(p *etfapi.Payload, req wsclient.WSMessage, resp chan<- wsclient.WSMessage) {
}

// NewDiscordMessageHandler creates a new DiscordMessageHandler object with default state and
// session management handlers installed
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
		"READY":                []discordapi.DiscordMessageHandlerFunc{c.handleReady},
		"GUILD_CREATE":         []discordapi.DiscordMessageHandlerFunc{c.handleGuildCreate},
		"GUILD_UPDATE":         []discordapi.DiscordMessageHandlerFunc{c.handleGuildUpdate},
		"GUILD_DELETE":         []discordapi.DiscordMessageHandlerFunc{c.handleGuildDelete},
		"CHANNEL_CREATE":       []discordapi.DiscordMessageHandlerFunc{c.handleChannelCreate},
		"CHANNEL_UPDATE":       []discordapi.DiscordMessageHandlerFunc{c.handleChannelUpdate},
		"CHANNEL_DELETE":       []discordapi.DiscordMessageHandlerFunc{c.handleChannelDelete},
		"GUILD_MEMBER_ADD":     []discordapi.DiscordMessageHandlerFunc{c.handleGuildMemberCreate},
		"GUILD_MEMEBER_UPDATE": []discordapi.DiscordMessageHandlerFunc{c.handleGuildMemberUpdate},
		"GUILD_MEMBER_REMOVE":  []discordapi.DiscordMessageHandlerFunc{c.handleGuildMemberDelete},
		"GUILD_ROLE_CREATE":    []discordapi.DiscordMessageHandlerFunc{c.handleGuildRoleCreate},
		"GUILD_ROLE_UPDATE":    []discordapi.DiscordMessageHandlerFunc{c.handleGuildRoleUpdate},
		"GUILD_ROLE_DELETE":    []discordapi.DiscordMessageHandlerFunc{c.handleGuildRoleDelete},
	}

	return &c
}

func (c *discordMessageHandler) ConnectToBot(bot discordapi.DiscordBot) {
	c.bot = bot
}

func (c *discordMessageHandler) AddHandler(event string, handler discordapi.DiscordMessageHandlerFunc) {
	c.dispatcherLock.Lock()
	defer c.dispatcherLock.Unlock()

	handlers := c.eventDispatch[event]
	c.eventDispatch[event] = append(handlers, handler)
}

func (c *discordMessageHandler) HandleRequest(req wsclient.WSMessage, resp chan<- wsclient.WSMessage) {
	logger := logging.WithContext(req.Ctx, c.deps.Logger())
	_ = level.Info(logger).Log("message", "discordapi dispatching request")

	select {
	case <-req.Ctx.Done():
		_ = level.Info(logger).Log("message", "discordapi already done. skipping request")
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

	_ = level.Debug(logger).Log("message", "received payload", "payload", p)

	opHandler, ok := c.opCodeDispatch[p.OpCode]
	if !ok {
		_ = level.Error(logger).Log("message", "unrecognized OpCode", "op_code", p.OpCode)
		return
	}

	if opHandler == nil {
		_ = level.Error(logger).Log("message", "no handler for OpCode", "op_code", p.OpCode)
		return
	}

	_ = level.Info(logger).Log("message", "sending to opHandler", "op_code", p.OpCode)
	opHandler(p, req, resp)
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
		rp := &payloads.ResumePayload{
			Token:     c.bot.Config().BotToken,
			SessionID: sessID,
			SeqNum:    c.bot.LastSequence(),
		}

		m, err = payloads.ETFPayloadToMessage(req.Ctx, rp)
	} else {
		_ = level.Info(logger).Log("message", "generating identify payload")
		ip := &payloads.IdentifyPayload{
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

	_ = level.Info(logger).Log("message", "looking up event dispatch handler", "event_name", p.EventName)

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
	_ = level.Info(logger).Log("message", "upserting guild", "event_name", "GUILD_CREATE", "guild_id_elem", fmt.Sprintf("%+v", p.Data["id"]))
	_ = level.Debug(logger).Log("message", "upserting guild debug", "pdata", fmt.Sprintf("%+v", p.Data), "event_name", "GUILD_CREATE")
	err := c.deps.BotSession().UpsertGuildFromElementMap(p.Data)
	if err != nil {
		_ = level.Error(logger).Log("message", "error processing guild create", "err", err)
		return
	}
}

func (c *discordMessageHandler) handleGuildUpdate(p *etfapi.Payload, req wsclient.WSMessage, resp chan<- wsclient.WSMessage) {
	select {
	case <-req.Ctx.Done():
		return
	default:
	}

	logger := logging.WithContext(req.Ctx, c.deps.Logger())
	_ = level.Info(logger).Log("message", "upserting guild", "event_name", "GUILD_UPDATE", "guild_id_elem", fmt.Sprintf("%+v", p.Data["id"]))
	_ = level.Debug(logger).Log("message", "upserting guild debug", "pdata", fmt.Sprintf("%+v", p.Data), "event_name", "GUILD_UPDATE")
	err := c.deps.BotSession().UpsertGuildFromElementMap(p.Data)
	if err != nil {
		_ = level.Error(logger).Log("message", "error processing guild update", "err", err)
		return
	}
}

func (c *discordMessageHandler) handleGuildDelete(p *etfapi.Payload, req wsclient.WSMessage, resp chan<- wsclient.WSMessage) {
	select {
	case <-req.Ctx.Done():
		return
	default:
	}

	logger := logging.WithContext(req.Ctx, c.deps.Logger())
	_ = level.Info(logger).Log("message", "upserting guild", "event_name", "GUILD_DELETE", "guild_id_elem", fmt.Sprintf("%+v", p.Data["id"]))
	_ = level.Debug(logger).Log("message", "upserting guild debug", "pdata", fmt.Sprintf("%+v", p.Data), "event_name", "GUILD_DELETE")
	err := c.deps.BotSession().UpsertGuildFromElementMap(p.Data)
	if err != nil {
		_ = level.Error(logger).Log("message", "error processing guild delete", "err", err)
		return
	}
}

func (c *discordMessageHandler) handleChannelCreate(p *etfapi.Payload, req wsclient.WSMessage, resp chan<- wsclient.WSMessage) {
	select {
	case <-req.Ctx.Done():
		return
	default:
	}

	logger := logging.WithContext(req.Ctx, c.deps.Logger())
	_ = level.Info(logger).Log("message", "upserting channel", "event_name", "CHANNEL_CREATE", "channel_id_elem", fmt.Sprintf("%+v", p.Data["id"]))
	_ = level.Debug(logger).Log("message", "upserting channel debug", "pdata", fmt.Sprintf("%+v", p.Data), "event_name", "CHANNEL_CREATE")
	err := c.deps.BotSession().UpsertChannelFromElementMap(p.Data)
	if err != nil {
		_ = level.Error(logger).Log("message", "error processing channel create", "err", err)
		return
	}
}

func (c *discordMessageHandler) handleChannelUpdate(p *etfapi.Payload, req wsclient.WSMessage, resp chan<- wsclient.WSMessage) {
	select {
	case <-req.Ctx.Done():
		return
	default:
	}

	logger := logging.WithContext(req.Ctx, c.deps.Logger())
	_ = level.Info(logger).Log("message", "upserting channel", "event_name", "CHANNEL_UPDATE", "channel_id_elem", fmt.Sprintf("%+v", p.Data["id"]))
	_ = level.Debug(logger).Log("message", "upserting channel debug", "pdata", fmt.Sprintf("%+v", p.Data), "event_name", "CHANNEL_UPDATE")
	err := c.deps.BotSession().UpsertChannelFromElementMap(p.Data)
	if err != nil {
		_ = level.Error(logger).Log("message", "error processing channel update", "err", err)
		return
	}
}

func (c *discordMessageHandler) handleChannelDelete(p *etfapi.Payload, req wsclient.WSMessage, resp chan<- wsclient.WSMessage) {
	select {
	case <-req.Ctx.Done():
		return
	default:
	}

	logger := logging.WithContext(req.Ctx, c.deps.Logger())
	_ = level.Info(logger).Log("message", "upserting channel", "event_name", "CHANNEL_DELETE", "channel_id_elem", fmt.Sprintf("%+v", p.Data["id"]))
	_ = level.Debug(logger).Log("message", "upserting channel debug", "pdata", fmt.Sprintf("%+v", p.Data), "event_name", "CHANNEL_DELETE")
	err := c.deps.BotSession().UpsertChannelFromElementMap(p.Data)
	if err != nil {
		_ = level.Error(logger).Log("message", "error processing channel delete", "err", err)
		return
	}
}

func (c *discordMessageHandler) handleGuildMemberCreate(p *etfapi.Payload, req wsclient.WSMessage, resp chan<- wsclient.WSMessage) {
	select {
	case <-req.Ctx.Done():
		return
	default:
	}

	logger := logging.WithContext(req.Ctx, c.deps.Logger())
	_ = level.Info(logger).Log("message", "upserting guild member", "event_name", "GUILD_MEMBER_ADD", "guild_id_elem", fmt.Sprintf("%+v", p.Data["guild_id"]))
	_ = level.Debug(logger).Log("message", "upserting guild member debug", "pdata", fmt.Sprintf("%+v", p.Data), "event_name", "GUILD_MEMBER_ADD")
	err := c.deps.BotSession().UpsertGuildMemberFromElementMap(p.Data)
	if err != nil {
		_ = level.Error(logger).Log("message", "error processing guild member create", "err", err)
		return
	}
}

func (c *discordMessageHandler) handleGuildMemberUpdate(p *etfapi.Payload, req wsclient.WSMessage, resp chan<- wsclient.WSMessage) {
	select {
	case <-req.Ctx.Done():
		return
	default:
	}

	logger := logging.WithContext(req.Ctx, c.deps.Logger())
	_ = level.Info(logger).Log("message", "upserting guild member", "event_name", "GUILD_MEMBER_UPDATE", "guild_id_elem", fmt.Sprintf("%+v", p.Data["guild_id"]))
	_ = level.Debug(logger).Log("message", "upserting guild member debug", "pdata", fmt.Sprintf("%+v", p.Data), "event_name", "GUILD_MEMBER_UPDATE")
	err := c.deps.BotSession().UpsertGuildMemberFromElementMap(p.Data)
	if err != nil {
		_ = level.Error(logger).Log("message", "error processing guild member update", "err", err)
		return
	}
}

func (c *discordMessageHandler) handleGuildMemberDelete(p *etfapi.Payload, req wsclient.WSMessage, resp chan<- wsclient.WSMessage) {
	select {
	case <-req.Ctx.Done():
		return
	default:
	}

	logger := logging.WithContext(req.Ctx, c.deps.Logger())
	_ = level.Info(logger).Log("message", "upserting guild member", "event_name", "GUILD_MEMBER_REMOVE", "guild_id_elem", fmt.Sprintf("%+v", p.Data["guild_id"]))
	_ = level.Debug(logger).Log("message", "upserting guild member debug", "pdata", fmt.Sprintf("%+v", p.Data), "event_name", "GUILD_MEMBER_REMOVE")
	err := c.deps.BotSession().UpsertGuildMemberFromElementMap(p.Data)
	if err != nil {
		_ = level.Error(logger).Log("message", "error processing guild member delete", "err", err)
		return
	}
}

func (c *discordMessageHandler) handleGuildRoleCreate(p *etfapi.Payload, req wsclient.WSMessage, resp chan<- wsclient.WSMessage) {
	select {
	case <-req.Ctx.Done():
		return
	default:
	}

	logger := logging.WithContext(req.Ctx, c.deps.Logger())
	_ = level.Info(logger).Log("message", "upserting guild role", "event_name", "GUILD_ROLE_CREATE", "guild_id_elem", fmt.Sprintf("%+v", p.Data["guild_id"]))
	_ = level.Debug(logger).Log("message", "upserting guild role debug", "pdata", fmt.Sprintf("%+v", p.Data), "event_name", "GUILD_ROLE_CREATE")
	err := c.deps.BotSession().UpsertGuildRoleFromElementMap(p.Data)
	if err != nil {
		_ = level.Error(logger).Log("message", "error processing guild role create", "err", err)
		return
	}
}

func (c *discordMessageHandler) handleGuildRoleUpdate(p *etfapi.Payload, req wsclient.WSMessage, resp chan<- wsclient.WSMessage) {
	select {
	case <-req.Ctx.Done():
		return
	default:
	}

	logger := logging.WithContext(req.Ctx, c.deps.Logger())
	_ = level.Info(logger).Log("message", "upserting guild role", "event_name", "GUILD_ROLE_UPDATE", "guild_id_elem", fmt.Sprintf("%+v", p.Data["guild_id"]))
	_ = level.Debug(logger).Log("message", "upserting guild role debug", "pdata", fmt.Sprintf("%+v", p.Data), "event_name", "GUILD_ROLE_UPDATE")
	err := c.deps.BotSession().UpsertGuildRoleFromElementMap(p.Data)
	if err != nil {
		_ = level.Error(logger).Log("message", "error processing guild role update", "err", err)
		return
	}
}

func (c *discordMessageHandler) handleGuildRoleDelete(p *etfapi.Payload, req wsclient.WSMessage, resp chan<- wsclient.WSMessage) {
	select {
	case <-req.Ctx.Done():
		return
	default:
	}

	logger := logging.WithContext(req.Ctx, c.deps.Logger())
	_ = level.Info(logger).Log("message", "upserting guild role", "event_name", "GUILD_ROLE_DELETE", "guild_id_elem", fmt.Sprintf("%+v", p.Data["guild_id"]))
	_ = level.Debug(logger).Log("message", "upserting guild role debug", "pdata", fmt.Sprintf("%+v", p.Data), "event_name", "GUILD_ROLE_DELETE")
	err := c.deps.BotSession().UpsertGuildRoleFromElementMap(p.Data)
	if err != nil {
		_ = level.Error(logger).Log("message", "error processing guild role delete", "err", err)
		return
	}
}
