package messagehandler

import (
	"fmt"
	"sync"

	log "github.com/gsmcwhirter/go-util/v3/logging"
	"github.com/gsmcwhirter/go-util/v3/logging/level"
	"golang.org/x/time/rate"

	"github.com/gsmcwhirter/discord-bot-lib/v8/bot"
	"github.com/gsmcwhirter/discord-bot-lib/v8/discordapi"
	"github.com/gsmcwhirter/discord-bot-lib/v8/etfapi"
	"github.com/gsmcwhirter/discord-bot-lib/v8/etfapi/payloads"
	"github.com/gsmcwhirter/discord-bot-lib/v8/logging"
	"github.com/gsmcwhirter/discord-bot-lib/v8/wsclient"
)

type dependencies interface {
	Logger() log.Logger
	BotSession() *etfapi.Session
	MessageRateLimiter() *rate.Limiter
}

type discordMessageHandler struct {
	deps           dependencies
	bot            bot.DiscordBot
	opCodeDispatch map[discordapi.OpCode]bot.DiscordMessageHandlerFunc

	dispatcherLock *sync.Mutex
	eventDispatch  map[string][]bot.DiscordMessageHandlerFunc
}

func noop(p *etfapi.Payload, req wsclient.WSMessage, resp chan<- wsclient.WSMessage) {
}

// NewDiscordMessageHandler creates a new DiscordMessageHandler object with default state and
// session management handlers installed
func NewDiscordMessageHandler(deps dependencies) bot.DiscordMessageHandler {
	c := discordMessageHandler{
		deps:           deps,
		dispatcherLock: &sync.Mutex{},
	}

	c.opCodeDispatch = map[discordapi.OpCode]bot.DiscordMessageHandlerFunc{
		discordapi.Hello:          c.handleHello,
		discordapi.Heartbeat:      c.handleHeartbeat,
		discordapi.HeartbeatAck:   noop,
		discordapi.InvalidSession: nil,
		discordapi.Reconnect:      nil,
		discordapi.Dispatch:       c.handleDispatch,
	}

	c.eventDispatch = map[string][]bot.DiscordMessageHandlerFunc{
		"READY":               {c.handleReady},
		"GUILD_CREATE":        {c.handleGuildCreate},
		"GUILD_UPDATE":        {c.handleGuildUpdate},
		"GUILD_DELETE":        {c.handleGuildDelete},
		"CHANNEL_CREATE":      {c.handleChannelCreate},
		"CHANNEL_UPDATE":      {c.handleChannelUpdate},
		"CHANNEL_DELETE":      {c.handleChannelDelete},
		"GUILD_MEMBER_ADD":    {c.handleGuildMemberCreate},
		"GUILD_MEMBER_UPDATE": {c.handleGuildMemberUpdate},
		"GUILD_MEMBER_REMOVE": {c.handleGuildMemberDelete},
		"GUILD_ROLE_CREATE":   {c.handleGuildRoleCreate},
		"GUILD_ROLE_UPDATE":   {c.handleGuildRoleUpdate},
		"GUILD_ROLE_DELETE":   {c.handleGuildRoleDelete},
	}

	return &c
}

func (c *discordMessageHandler) ConnectToBot(b bot.DiscordBot) {
	c.bot = b
}

func (c *discordMessageHandler) AddHandler(event string, handler bot.DiscordMessageHandlerFunc) {
	c.dispatcherLock.Lock()
	defer c.dispatcherLock.Unlock()

	handlers := c.eventDispatch[event]
	c.eventDispatch[event] = append(handlers, handler)
}

func (c *discordMessageHandler) HandleRequest(req wsclient.WSMessage, resp chan<- wsclient.WSMessage) {
	logger := logging.WithContext(req.Ctx, c.deps.Logger())
	level.Info(logger).Message("discordapi dispatching request")

	select {
	case <-req.Ctx.Done():
		level.Info(logger).Message("discordapi already done. skipping request")
		return
	default:
	}

	level.Debug(logger).Message("processing server message", "ws_msg", fmt.Sprintf("%v", req.MessageContents))

	p, err := etfapi.Unmarshal(req.MessageContents)
	if err != nil {
		level.Error(logger).Err("error unmarshaling payload", err, "ws_msg", fmt.Sprintf("%v", req.MessageContents))
		return
	}

	if p.SeqNum != nil {
		c.bot.UpdateSequence(*p.SeqNum)
	}

	level.Debug(logger).Message("received payload", "payload", p)

	opHandler, ok := c.opCodeDispatch[p.OpCode]
	if !ok {
		level.Error(logger).Message("unrecognized OpCode", "op_code", p.OpCode)
		return
	}

	if opHandler == nil {
		level.Error(logger).Message("no handler for OpCode", "op_code", p.OpCode)
		return
	}

	level.Info(logger).Message("sending to opHandler", "op_code", p.OpCode)
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
			level.Error(logger).Err("error handling hello heartbeat config", err)
			return
		}

		level.Info(logger).Message("configuring heartbeat", "interval", interval)
		c.bot.ReconfigureHeartbeat(req.Ctx, interval)
		level.Debug(logger).Message("configuring heartbeat done")
	}

	// send identify
	var m wsclient.WSMessage
	var err error

	sessID := c.deps.BotSession().ID()
	if sessID != "" {
		level.Info(logger).Message("generating resume payload")
		rp := &payloads.ResumePayload{
			Token:     c.bot.Config().BotToken,
			SessionID: sessID,
			SeqNum:    c.bot.LastSequence(),
		}

		m, err = payloads.ETFPayloadToMessage(req.Ctx, rp)
	} else {
		level.Info(logger).Message("generating identify payload")
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
		level.Error(logger).Err("error generating identify/resume payload", err)
		return
	}

	err = c.deps.MessageRateLimiter().Wait(req.Ctx)
	if err != nil {
		level.Error(logger).Err("error ratelimiting", err)
		return
	}

	level.Info(logger).Message("sending identify/resume to channel")
	level.Debug(logger).Message("sending response to channel", "resp_message", m, "msg_len", len(m.MessageContents))
	resp <- m
}

func (c *discordMessageHandler) handleHeartbeat(p *etfapi.Payload, req wsclient.WSMessage, resp chan<- wsclient.WSMessage) {
	select {
	case <-req.Ctx.Done():
		return
	default:
	}

	logger := logging.WithContext(req.Ctx, c.deps.Logger())
	level.Info(logger).Message("requesting manual heartbeat")
	c.bot.ReconfigureHeartbeat(req.Ctx, 0)
	level.Debug(logger).Message("manual heartbeat done")
}

func (c *discordMessageHandler) handleDispatch(p *etfapi.Payload, req wsclient.WSMessage, resp chan<- wsclient.WSMessage) {
	select {
	case <-req.Ctx.Done():
		return
	default:
	}

	logger := logging.WithContext(req.Ctx, c.deps.Logger())

	level.Info(logger).Message("looking up event dispatch handler", "event_name", p.EventName)

	c.dispatcherLock.Lock()
	eventHandlers, ok := c.eventDispatch[p.EventName]
	c.dispatcherLock.Unlock()

	if !ok {
		level.Debug(logger).Message("no event dispatch handler found", "event_name", p.EventName)
		return
	}

	level.Info(logger).Message("processing event", "event_name", p.EventName)
	for _, eventHandler := range eventHandlers {
		eventHandler(p, req, resp)
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
		level.Error(logger).Err("error setting up session", err)
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
	level.Info(logger).Message("upserting guild", "event_name", "GUILD_CREATE", "guild_id_elem", fmt.Sprintf("%+v", p.Data["id"]))
	level.Debug(logger).Message("upserting guild debug", "pdata", fmt.Sprintf("%+v", p.Data), "event_name", "GUILD_CREATE")
	err := c.deps.BotSession().UpsertGuildFromElementMap(p.Data)
	if err != nil {
		level.Error(logger).Err("error processing guild create", err)
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
	level.Info(logger).Message("upserting guild", "event_name", "GUILD_UPDATE", "guild_id_elem", fmt.Sprintf("%+v", p.Data["id"]))
	level.Debug(logger).Message("upserting guild debug", "pdata", fmt.Sprintf("%+v", p.Data), "event_name", "GUILD_UPDATE")
	err := c.deps.BotSession().UpsertGuildFromElementMap(p.Data)
	if err != nil {
		level.Error(logger).Err("error processing guild update", err)
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
	level.Info(logger).Message("upserting guild", "event_name", "GUILD_DELETE", "guild_id_elem", fmt.Sprintf("%+v", p.Data["id"]))
	level.Debug(logger).Message("upserting guild debug", "pdata", fmt.Sprintf("%+v", p.Data), "event_name", "GUILD_DELETE")
	err := c.deps.BotSession().UpsertGuildFromElementMap(p.Data)
	if err != nil {
		level.Error(logger).Err("error processing guild delete", err)
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
	level.Info(logger).Message("upserting channel", "event_name", "CHANNEL_CREATE", "channel_id_elem", fmt.Sprintf("%+v", p.Data["id"]))
	level.Debug(logger).Message("upserting channel debug", "pdata", fmt.Sprintf("%+v", p.Data), "event_name", "CHANNEL_CREATE")
	err := c.deps.BotSession().UpsertChannelFromElementMap(p.Data)
	if err != nil {
		level.Error(logger).Err("error processing channel create", err)
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
	level.Info(logger).Message("upserting channel", "event_name", "CHANNEL_UPDATE", "channel_id_elem", fmt.Sprintf("%+v", p.Data["id"]))
	level.Debug(logger).Message("upserting channel debug", "pdata", fmt.Sprintf("%+v", p.Data), "event_name", "CHANNEL_UPDATE")
	err := c.deps.BotSession().UpsertChannelFromElementMap(p.Data)
	if err != nil {
		level.Error(logger).Err("error processing channel update", err)
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
	level.Info(logger).Message("upserting channel", "event_name", "CHANNEL_DELETE", "channel_id_elem", fmt.Sprintf("%+v", p.Data["id"]))
	level.Debug(logger).Message("upserting channel debug", "pdata", fmt.Sprintf("%+v", p.Data), "event_name", "CHANNEL_DELETE")
	err := c.deps.BotSession().UpsertChannelFromElementMap(p.Data)
	if err != nil {
		level.Error(logger).Err("error processing channel delete", err)
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
	level.Info(logger).Message("upserting guild member", "event_name", "GUILD_MEMBER_ADD", "guild_id_elem", fmt.Sprintf("%+v", p.Data["guild_id"]))
	level.Debug(logger).Message("upserting guild member debug", "pdata", fmt.Sprintf("%+v", p.Data), "event_name", "GUILD_MEMBER_ADD")
	err := c.deps.BotSession().UpsertGuildMemberFromElementMap(p.Data)
	if err != nil {
		level.Error(logger).Err("error processing guild member create", err)
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
	level.Info(logger).Message("upserting guild member", "event_name", "GUILD_MEMBER_UPDATE", "guild_id_elem", fmt.Sprintf("%+v", p.Data["guild_id"]))
	level.Debug(logger).Message("upserting guild member debug", "pdata", fmt.Sprintf("%+v", p.Data), "event_name", "GUILD_MEMBER_UPDATE")
	err := c.deps.BotSession().UpsertGuildMemberFromElementMap(p.Data)
	if err != nil {
		level.Error(logger).Err("error processing guild member update", err)
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
	level.Info(logger).Message("upserting guild member", "event_name", "GUILD_MEMBER_REMOVE", "guild_id_elem", fmt.Sprintf("%+v", p.Data["guild_id"]))
	level.Debug(logger).Message("upserting guild member debug", "pdata", fmt.Sprintf("%+v", p.Data), "event_name", "GUILD_MEMBER_REMOVE")
	err := c.deps.BotSession().UpsertGuildMemberFromElementMap(p.Data)
	if err != nil {
		level.Error(logger).Err("error processing guild member delete", err)
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
	level.Info(logger).Message("upserting guild role", "event_name", "GUILD_ROLE_CREATE", "guild_id_elem", fmt.Sprintf("%+v", p.Data["guild_id"]))
	level.Debug(logger).Message("upserting guild role debug", "pdata", fmt.Sprintf("%+v", p.Data), "event_name", "GUILD_ROLE_CREATE")
	err := c.deps.BotSession().UpsertGuildRoleFromElementMap(p.Data)
	if err != nil {
		level.Error(logger).Err("error processing guild role create", err)
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
	level.Info(logger).Message("upserting guild role", "event_name", "GUILD_ROLE_UPDATE", "guild_id_elem", fmt.Sprintf("%+v", p.Data["guild_id"]))
	level.Debug(logger).Message("upserting guild role debug", "pdata", fmt.Sprintf("%+v", p.Data), "event_name", "GUILD_ROLE_UPDATE")
	err := c.deps.BotSession().UpsertGuildRoleFromElementMap(p.Data)
	if err != nil {
		level.Error(logger).Err("error processing guild role update", err)
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
	level.Info(logger).Message("upserting guild role", "event_name", "GUILD_ROLE_DELETE", "guild_id_elem", fmt.Sprintf("%+v", p.Data["guild_id"]))
	level.Debug(logger).Message("upserting guild role debug", "pdata", fmt.Sprintf("%+v", p.Data), "event_name", "GUILD_ROLE_DELETE")
	err := c.deps.BotSession().UpsertGuildRoleFromElementMap(p.Data)
	if err != nil {
		level.Error(logger).Err("error processing guild role delete", err)
		return
	}
}
