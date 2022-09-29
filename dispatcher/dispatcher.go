package dispatcher

import (
	"context"
	"fmt"
	"sync"

	"github.com/gsmcwhirter/go-util/v10/errors"
	"github.com/gsmcwhirter/go-util/v10/logging/level"
	"github.com/gsmcwhirter/go-util/v10/telemetry"
	"golang.org/x/time/rate"

	"github.com/gsmcwhirter/discord-bot-lib/v24/bot"
	"github.com/gsmcwhirter/discord-bot-lib/v24/bot/session"
	"github.com/gsmcwhirter/discord-bot-lib/v24/discordapi"
	"github.com/gsmcwhirter/discord-bot-lib/v24/discordapi/etfapi"
	"github.com/gsmcwhirter/discord-bot-lib/v24/logging"
	"github.com/gsmcwhirter/discord-bot-lib/v24/snowflake"
	"github.com/gsmcwhirter/discord-bot-lib/v24/stats"
	"github.com/gsmcwhirter/discord-bot-lib/v24/wsapi"
)

type dependencies interface {
	Logger() Logger
	BotSession() *session.Session
	MessageRateLimiter() *rate.Limiter
	Telemetry() *telemetry.Telemeter
}

// Logger is the interface expected for logging
type Logger = interface {
	Log(keyvals ...interface{}) error
	Message(string, ...interface{})
	Err(string, error, ...interface{})
	Printf(string, ...interface{})
}

// Dispatcher is the op-code dispatcher
type Dispatcher struct {
	deps           dependencies
	bot            *bot.DiscordBot
	opCodeDispatch map[discordapi.OpCode]DispatchHandlerFunc

	dispatcherLock *sync.Mutex
	eventDispatch  map[string][]DispatchHandlerFunc

	debug bool
}

var _ bot.Dispatcher = (*Dispatcher)(nil)

func noop(p Payload, req wsapi.WSMessage, resp chan<- wsapi.WSMessage) snowflake.Snowflake {
	return 0
}

// NewDispatcher creates a new Dispatcher object with default state and
// session management handlers installed
func NewDispatcher(deps dependencies) *Dispatcher {
	c := &Dispatcher{
		deps:           deps,
		dispatcherLock: &sync.Mutex{},
	}

	c.opCodeDispatch = map[discordapi.OpCode]DispatchHandlerFunc{
		discordapi.Hello:          c.handleHello,
		discordapi.Heartbeat:      c.handleHeartbeat,
		discordapi.HeartbeatAck:   noop,
		discordapi.InvalidSession: nil,
		discordapi.Reconnect:      nil,
		discordapi.Dispatch:       c.handleDispatch,
	}

	c.eventDispatch = map[string][]DispatchHandlerFunc{
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

	return c
}

// SetDebug turns on/off debug mode
func (c *Dispatcher) SetDebug(val bool) {
	c.debug = val
}

// ConnectToBot attaches this dispatcher to a bot instance
func (c *Dispatcher) ConnectToBot(b *bot.DiscordBot) {
	c.bot = b
}

// GenerateHeartbeat prepares a heartbeat message to be sent
func (c *Dispatcher) GenerateHeartbeat(ctx context.Context, seqNum int) (wsapi.WSMessage, error) {
	ctx, span := c.deps.Telemetry().StartSpan(ctx, "dispatcher", "GenerateHeartbeat")
	defer span.End()

	var m wsapi.WSMessage

	m, err := ETFPayloadToMessage(ctx, &etfapi.HeartbeatPayload{
		Sequence: seqNum,
	})
	if err != nil {
		level.Error(logging.WithContext(ctx, c.deps.Logger())).Err("error formatting heartbeat", err)
		return m, errors.Wrap(err, "error formatting heartbeat")
	}

	err = c.deps.MessageRateLimiter().Wait(ctx)
	if err != nil {
		level.Error(logging.WithContext(ctx, c.deps.Logger())).Err("error rate limiting", err)
		return m, errors.Wrap(err, "error rate limiting")
	}

	return m, nil
}

// AddHandler adds a new event handler to the dispatch table
func (c *Dispatcher) AddHandler(event string, handler DispatchHandlerFunc) {
	c.dispatcherLock.Lock()
	defer c.dispatcherLock.Unlock()

	handlers := c.eventDispatch[event]
	c.eventDispatch[event] = append(handlers, handler)
}

// HandleRequest dispatches a message and queues a response, if there is one
func (c *Dispatcher) HandleRequest(req wsapi.WSMessage, resp chan<- wsapi.WSMessage) snowflake.Snowflake {
	ctx, span := c.deps.Telemetry().StartSpan(req.Ctx, "dispatcher", "HandleRequest")
	defer span.End()
	req.Ctx = ctx

	logger := logging.WithContext(req.Ctx, c.deps.Logger())
	// level.Info(logger).Message("discordapi dispatching request")

	select {
	case <-req.Ctx.Done():
		level.Info(logger).Message("discordapi already done. skipping request")
		return 0
	default:
	}

	if c.debug {
		level.Debug(logger).Message("processing server message", "ws_msg", fmt.Sprintf("%v", req.MessageContents))
	}

	p, err := etfapi.Unmarshal(req.MessageContents)
	if err != nil {
		level.Error(logger).Err("error unmarshaling payload", err, "ws_msg", fmt.Sprintf("%v", req.MessageContents))
		return 0
	}

	if err := stats.IncCounter(ctx, c.deps.Telemetry(), "dispatcher", stats.OpCodesCount, 1, telemetry.KVString(stats.TagOpCode, p.OpCode.String())); err != nil {
		level.Error(logger).Err("could not record stat", err)
	}

	if p.SeqNum != nil {
		c.bot.UpdateSequence(*p.SeqNum)
	}

	if c.debug {
		level.Debug(logger).Message("received payload", "payload", p)
	}

	opHandler, ok := c.opCodeDispatch[p.OpCode]
	if !ok {
		level.Error(logger).Message("unrecognized OpCode", "op_code", p.OpCode)
		return 0
	}

	if opHandler == nil {
		level.Error(logger).Message("no handler for OpCode", "op_code", p.OpCode)
		return 0
	}

	if c.debug {
		level.Debug(logger).Message("sending to opHandler", "op_code", p.OpCode)
	}
	return opHandler(p, req, resp)
}

func (c *Dispatcher) handleHello(p Payload, req wsapi.WSMessage, resp chan<- wsapi.WSMessage) snowflake.Snowflake {
	ctx, span := c.deps.Telemetry().StartSpan(req.Ctx, "dispatcher", "handleHello")
	defer span.End()
	req.Ctx = ctx

	select {
	case <-req.Ctx.Done():
		return 0
	default:
	}

	data := p.Contents()
	logger := logging.WithContext(req.Ctx, c.deps.Logger())
	rawInterval, ok := data["heartbeat_interval"]

	if ok {
		// set heartbeat stuff
		interval, err := rawInterval.ToInt()
		if err != nil {
			level.Error(logger).Err("error handling hello heartbeat config", err)
			return 0
		}

		level.Info(logger).Message("configuring heartbeat", "interval", interval)
		c.bot.ReconfigureHeartbeat(req.Ctx, interval)

		if c.debug {
			level.Debug(logger).Message("configuring heartbeat done")
		}
	}

	// send identify
	var m wsapi.WSMessage
	var err error

	sessID := c.deps.BotSession().ID()
	if sessID != "" {
		level.Info(logger).Message("generating resume payload")
		rp := &etfapi.ResumePayload{
			Token:     c.bot.Config().BotToken,
			SessionID: sessID,
			SeqNum:    c.bot.LastSequence(),
		}

		m, err = ETFPayloadToMessage(req.Ctx, rp)
	} else {
		level.Info(logger).Message("generating identify payload")
		ip := &etfapi.IdentifyPayload{
			Token:   c.bot.Config().BotToken,
			Intents: c.bot.Intents(),
			Properties: etfapi.IdentifyPayloadProperties{
				OS:      c.bot.Config().OS,
				Browser: c.bot.Config().BotName,
				Device:  c.bot.Config().BotName,
			},
			LargeThreshold: 250,
			Shard: etfapi.IdentifyPayloadShard{
				ID:    0,
				MaxID: 0,
			},
			Presence: etfapi.IdentifyPayloadPresence{
				Game: etfapi.IdentifyPayloadGame{
					Name: c.bot.Config().BotPresence,
					Type: 0,
				},
				Status: "online",
				Since:  0,
				AFK:    false,
			},
		}

		m, err = ETFPayloadToMessage(req.Ctx, ip)
	}

	if err != nil {
		level.Error(logger).Err("error generating identify/resume payload", err)
		return 0
	}

	err = c.deps.MessageRateLimiter().Wait(req.Ctx)
	if err != nil {
		level.Error(logger).Err("error ratelimiting", err)
		return 0
	}

	level.Info(logger).Message("sending identify/resume to channel")

	if c.debug {
		level.Debug(logger).Message("sending response to channel", "resp_message", m, "msg_len", len(m.MessageContents))
	}
	resp <- m

	return 0
}

func (c *Dispatcher) handleHeartbeat(p Payload, req wsapi.WSMessage, resp chan<- wsapi.WSMessage) snowflake.Snowflake {
	ctx, span := c.deps.Telemetry().StartSpan(req.Ctx, "dispatcher", "handleHeartbeat")
	defer span.End()
	req.Ctx = ctx

	select {
	case <-req.Ctx.Done():
		return 0
	default:
	}

	logger := logging.WithContext(req.Ctx, c.deps.Logger())
	level.Info(logger).Message("requesting manual heartbeat")
	c.bot.ReconfigureHeartbeat(req.Ctx, 0)
	if c.debug {
		level.Debug(logger).Message("manual heartbeat done")
	}

	return 0
}

func (c *Dispatcher) handleDispatch(p Payload, req wsapi.WSMessage, resp chan<- wsapi.WSMessage) snowflake.Snowflake {
	ctx, span := c.deps.Telemetry().StartSpan(req.Ctx, "dispatcher", "handleDispatch")
	defer span.End()
	req.Ctx = ctx

	select {
	case <-req.Ctx.Done():
		return 0
	default:
	}

	logger := logging.WithContext(req.Ctx, c.deps.Logger())

	if err := stats.IncCounter(ctx, c.deps.Telemetry(), "dispatcher", stats.RawEventsCount, 1, telemetry.KVString(stats.TagEventName, p.EventName())); err != nil {
		level.Error(logger).Err("could not record stat", err)
	}

	if c.debug {
		level.Debug(logger).Message("looking up event dispatch handler", "event_name", p.EventName())
	}

	c.dispatcherLock.Lock()
	eventHandlers, ok := c.eventDispatch[p.EventName()]
	c.dispatcherLock.Unlock()

	if !ok {
		if c.debug {
			level.Debug(logger).Message("no event dispatch handler found", "event_name", p.EventName())
		}
		return 0
	}

	var guildID snowflake.Snowflake

	level.Info(logger).Message("processing event", "event_name", p.EventName())
	for _, eventHandler := range eventHandlers {
		if gid := eventHandler(p, req, resp); gid != 0 {
			span.SetAttributes(telemetry.KVString("gid", gid.ToString()))
			guildID = gid
		}
	}

	return guildID
}

func (c *Dispatcher) handleReady(p Payload, req wsapi.WSMessage, resp chan<- wsapi.WSMessage) snowflake.Snowflake {
	ctx, span := c.deps.Telemetry().StartSpan(req.Ctx, "dispatcher", "handleReady")
	defer span.End()
	req.Ctx = ctx

	select {
	case <-req.Ctx.Done():
		return 0
	default:
	}

	logger := logging.WithContext(req.Ctx, c.deps.Logger())

	err := c.deps.BotSession().UpdateFromReady(p.Contents())
	if err != nil {
		level.Error(logger).Err("error setting up session", err)
	}

	return 0
}

func (c *Dispatcher) handleGuildCreate(p Payload, req wsapi.WSMessage, resp chan<- wsapi.WSMessage) snowflake.Snowflake {
	ctx, span := c.deps.Telemetry().StartSpan(req.Ctx, "dispatcher", "handleGuildCreate")
	defer span.End()
	req.Ctx = ctx

	select {
	case <-req.Ctx.Done():
		return 0
	default:
	}

	logger := logging.WithContext(req.Ctx, c.deps.Logger())

	if c.debug {
		level.Debug(logger).Message("upserting guild debug", "pdata", fmt.Sprintf("%+v", p.Contents()), "event_name", "GUILD_CREATE")
	}

	data := p.Contents()
	gid, err := c.deps.BotSession().UpsertGuildFromElementMap(data)
	level.Info(logger).Message("upserting guild", "event_name", "GUILD_CREATE", "guild_id_elem", fmt.Sprintf("%+v", data["id"]), "guild_id", gid)
	if err != nil {
		level.Error(logger).Err("error processing guild create", err)
	}
	span.SetAttributes(telemetry.KVString("gid", gid.ToString()))

	return gid
}

func (c *Dispatcher) handleGuildUpdate(p Payload, req wsapi.WSMessage, resp chan<- wsapi.WSMessage) snowflake.Snowflake {
	ctx, span := c.deps.Telemetry().StartSpan(req.Ctx, "dispatcher", "handleGuildUpdate")
	defer span.End()
	req.Ctx = ctx

	select {
	case <-req.Ctx.Done():
		return 0
	default:
	}

	logger := logging.WithContext(req.Ctx, c.deps.Logger())
	data := p.Contents()
	if c.debug {
		level.Debug(logger).Message("upserting guild debug", "pdata", fmt.Sprintf("%+v", data), "event_name", "GUILD_UPDATE")
	}
	gid, err := c.deps.BotSession().UpsertGuildFromElementMap(data)
	level.Info(logger).Message("upserting guild", "event_name", "GUILD_UPDATE", "guild_id_elem", fmt.Sprintf("%+v", data["id"]), "guild_id", gid)
	if err != nil {
		level.Error(logger).Err("error processing guild update", err)
	}
	span.SetAttributes(telemetry.KVString("gid", gid.ToString()))

	return gid
}

func (c *Dispatcher) handleGuildDelete(p Payload, req wsapi.WSMessage, resp chan<- wsapi.WSMessage) snowflake.Snowflake {
	ctx, span := c.deps.Telemetry().StartSpan(req.Ctx, "dispatcher", "handleGuildDelete")
	defer span.End()
	req.Ctx = ctx

	select {
	case <-req.Ctx.Done():
		return 0
	default:
	}

	logger := logging.WithContext(req.Ctx, c.deps.Logger())
	data := p.Contents()
	if c.debug {
		level.Debug(logger).Message("deleting guild debug", "pdata", fmt.Sprintf("%+v", data), "event_name", "GUILD_DELETE")
	}
	gid, err := c.deps.BotSession().UpsertGuildFromElementMap(data)
	level.Info(logger).Message("upserting guild", "event_name", "GUILD_DELETE", "guild_id_elem", fmt.Sprintf("%+v", data["id"]), "guild_id", gid)
	if err != nil {
		level.Error(logger).Err("error processing guild delete", err)
	}
	span.SetAttributes(telemetry.KVString("gid", gid.ToString()))

	return gid
}

func (c *Dispatcher) handleChannelCreate(p Payload, req wsapi.WSMessage, resp chan<- wsapi.WSMessage) snowflake.Snowflake {
	ctx, span := c.deps.Telemetry().StartSpan(req.Ctx, "dispatcher", "handleChannelCreate")
	defer span.End()
	req.Ctx = ctx

	select {
	case <-req.Ctx.Done():
		return 0
	default:
	}

	logger := logging.WithContext(req.Ctx, c.deps.Logger())
	data := p.Contents()
	if c.debug {
		level.Debug(logger).Message("upserting channel debug", "pdata", fmt.Sprintf("%+v", data), "event_name", "CHANNEL_CREATE")
	}
	gid, cid, err := c.deps.BotSession().UpsertChannelFromElementMap(data)
	level.Info(logger).Message("upserting channel", "event_name", "CHANNEL_CREATE", "channel_id_elem", fmt.Sprintf("%+v", data["id"]), "guild_id", gid, "channel_id", cid)
	if err != nil {
		level.Error(logger).Err("error processing channel create", err)
	}
	span.SetAttributes(telemetry.KVString("gid", gid.ToString()), telemetry.KVString("cid", cid.ToString()))

	return gid
}

func (c *Dispatcher) handleChannelUpdate(p Payload, req wsapi.WSMessage, resp chan<- wsapi.WSMessage) snowflake.Snowflake {
	ctx, span := c.deps.Telemetry().StartSpan(req.Ctx, "dispatcher", "handleChannelUpdate")
	defer span.End()
	req.Ctx = ctx

	select {
	case <-req.Ctx.Done():
		return 0
	default:
	}

	logger := logging.WithContext(req.Ctx, c.deps.Logger())
	data := p.Contents()
	if c.debug {
		level.Debug(logger).Message("upserting channel debug", "pdata", fmt.Sprintf("%+v", data), "event_name", "CHANNEL_UPDATE")
	}
	gid, cid, err := c.deps.BotSession().UpsertChannelFromElementMap(data)
	level.Info(logger).Message("upserting channel", "event_name", "CHANNEL_UPDATE", "channel_id_elem", fmt.Sprintf("%+v", data["id"]), "guild_id", gid, "channel_id", cid)
	if err != nil {
		level.Error(logger).Err("error processing channel update", err)
	}
	span.SetAttributes(telemetry.KVString("gid", gid.ToString()), telemetry.KVString("cid", cid.ToString()))

	return gid
}

func (c *Dispatcher) handleChannelDelete(p Payload, req wsapi.WSMessage, resp chan<- wsapi.WSMessage) snowflake.Snowflake {
	ctx, span := c.deps.Telemetry().StartSpan(req.Ctx, "dispatcher", "handleChannelDelete")
	defer span.End()
	req.Ctx = ctx

	select {
	case <-req.Ctx.Done():
		return 0
	default:
	}

	logger := logging.WithContext(req.Ctx, c.deps.Logger())
	data := p.Contents()
	if c.debug {
		level.Debug(logger).Message("deleting channel debug", "pdata", fmt.Sprintf("%+v", data), "event_name", "CHANNEL_DELETE")
	}
	gid, cid, err := c.deps.BotSession().UpsertChannelFromElementMap(data)
	level.Info(logger).Message("upserting channel", "event_name", "CHANNEL_DELETE", "channel_id_elem", fmt.Sprintf("%+v", data["id"]), "guild_id", gid, "channel_id", cid)
	if err != nil {
		level.Error(logger).Err("error processing channel delete", err)
	}
	span.SetAttributes(telemetry.KVString("gid", gid.ToString()), telemetry.KVString("cid", cid.ToString()))

	return gid
}

func (c *Dispatcher) handleGuildMemberCreate(p Payload, req wsapi.WSMessage, resp chan<- wsapi.WSMessage) snowflake.Snowflake {
	ctx, span := c.deps.Telemetry().StartSpan(req.Ctx, "dispatcher", "handleGuildMemberCreate")
	defer span.End()
	req.Ctx = ctx

	select {
	case <-req.Ctx.Done():
		return 0
	default:
	}

	logger := logging.WithContext(req.Ctx, c.deps.Logger())
	data := p.Contents()
	if c.debug {
		level.Debug(logger).Message("upserting guild member debug", "pdata", fmt.Sprintf("%+v", data), "event_name", "GUILD_MEMBER_ADD")
	}
	gid, err := c.deps.BotSession().UpsertGuildMemberFromElementMap(data)
	level.Info(logger).Message("upserting guild member", "event_name", "GUILD_MEMBER_ADD", "guild_id_elem", fmt.Sprintf("%+v", data["guild_id"]), "guild_id", gid)
	if err != nil {
		level.Error(logger).Err("error processing guild member create", err)
	}
	span.SetAttributes(telemetry.KVString("gid", gid.ToString()))

	return gid
}

func (c *Dispatcher) handleGuildMemberUpdate(p Payload, req wsapi.WSMessage, resp chan<- wsapi.WSMessage) snowflake.Snowflake {
	ctx, span := c.deps.Telemetry().StartSpan(req.Ctx, "dispatcher", "handleGuildMemberUpdate")
	defer span.End()
	req.Ctx = ctx

	select {
	case <-req.Ctx.Done():
		return 0
	default:
	}

	logger := logging.WithContext(req.Ctx, c.deps.Logger())
	data := p.Contents()
	if c.debug {
		level.Debug(logger).Message("upserting guild member debug", "pdata", fmt.Sprintf("%+v", data), "event_name", "GUILD_MEMBER_UPDATE")
	}
	gid, err := c.deps.BotSession().UpsertGuildMemberFromElementMap(data)
	level.Info(logger).Message("upserting guild member", "event_name", "GUILD_MEMBER_UPDATE", "guild_id_elem", fmt.Sprintf("%+v", data["guild_id"]), "guild_id", gid)
	if err != nil {
		level.Error(logger).Err("error processing guild member update", err)
	}
	span.SetAttributes(telemetry.KVString("gid", gid.ToString()))

	return gid
}

func (c *Dispatcher) handleGuildMemberDelete(p Payload, req wsapi.WSMessage, resp chan<- wsapi.WSMessage) snowflake.Snowflake {
	ctx, span := c.deps.Telemetry().StartSpan(req.Ctx, "dispatcher", "handleGuildMemberDelete")
	defer span.End()
	req.Ctx = ctx

	select {
	case <-req.Ctx.Done():
		return 0
	default:
	}

	logger := logging.WithContext(req.Ctx, c.deps.Logger())
	data := p.Contents()
	if c.debug {
		level.Debug(logger).Message("deleting guild member debug", "pdata", fmt.Sprintf("%+v", data), "event_name", "GUILD_MEMBER_REMOVE")
	}
	gid, err := c.deps.BotSession().UpsertGuildMemberFromElementMap(data)
	level.Info(logger).Message("upserting guild member", "event_name", "GUILD_MEMBER_REMOVE", "guild_id_elem", fmt.Sprintf("%+v", data["guild_id"]), "guild_id", gid)
	if err != nil {
		level.Error(logger).Err("error processing guild member delete", err)
	}
	span.SetAttributes(telemetry.KVString("gid", gid.ToString()))

	return gid
}

func (c *Dispatcher) handleGuildRoleCreate(p Payload, req wsapi.WSMessage, resp chan<- wsapi.WSMessage) snowflake.Snowflake {
	ctx, span := c.deps.Telemetry().StartSpan(req.Ctx, "dispatcher", "handleGuildRoleCreate")
	defer span.End()
	req.Ctx = ctx

	select {
	case <-req.Ctx.Done():
		return 0
	default:
	}

	logger := logging.WithContext(req.Ctx, c.deps.Logger())
	data := p.Contents()
	if c.debug {
		level.Debug(logger).Message("upserting guild role debug", "pdata", fmt.Sprintf("%+v", data), "event_name", "GUILD_ROLE_CREATE")
	}
	gid, err := c.deps.BotSession().UpsertGuildRoleFromElementMap(data)
	level.Info(logger).Message("upserting guild role", "event_name", "GUILD_ROLE_CREATE", "guild_id_elem", fmt.Sprintf("%+v", data["guild_id"]), "guild_id", gid)
	if err != nil {
		level.Error(logger).Err("error processing guild role create", err)
	}
	span.SetAttributes(telemetry.KVString("gid", gid.ToString()))

	return gid
}

func (c *Dispatcher) handleGuildRoleUpdate(p Payload, req wsapi.WSMessage, resp chan<- wsapi.WSMessage) snowflake.Snowflake {
	ctx, span := c.deps.Telemetry().StartSpan(req.Ctx, "dispatcher", "handleGuildRoleUpdate")
	defer span.End()
	req.Ctx = ctx

	select {
	case <-req.Ctx.Done():
		return 0
	default:
	}

	logger := logging.WithContext(req.Ctx, c.deps.Logger())
	data := p.Contents()
	if c.debug {
		level.Debug(logger).Message("upserting guild role debug", "pdata", fmt.Sprintf("%+v", data), "event_name", "GUILD_ROLE_UPDATE")
	}
	gid, err := c.deps.BotSession().UpsertGuildRoleFromElementMap(data)
	level.Info(logger).Message("upserting guild role", "event_name", "GUILD_ROLE_UPDATE", "guild_id_elem", fmt.Sprintf("%+v", data["guild_id"]), "guild_id", gid)
	if err != nil {
		level.Error(logger).Err("error processing guild role update", err)
	}
	span.SetAttributes(telemetry.KVString("gid", gid.ToString()))

	return gid
}

func (c *Dispatcher) handleGuildRoleDelete(p Payload, req wsapi.WSMessage, resp chan<- wsapi.WSMessage) snowflake.Snowflake {
	ctx, span := c.deps.Telemetry().StartSpan(req.Ctx, "dispatcher", "handleGuildRoleDelete")
	defer span.End()
	req.Ctx = ctx

	select {
	case <-req.Ctx.Done():
		return 0
	default:
	}

	logger := logging.WithContext(req.Ctx, c.deps.Logger())
	data := p.Contents()
	if c.debug {
		level.Debug(logger).Message("deleting guild role debug", "pdata", fmt.Sprintf("%+v", data), "event_name", "GUILD_ROLE_DELETE")
	}
	gid, err := c.deps.BotSession().UpsertGuildRoleFromElementMap(data)
	level.Info(logger).Message("upserting guild role", "event_name", "GUILD_ROLE_DELETE", "guild_id_elem", fmt.Sprintf("%+v", data["guild_id"]), "guild_id", gid)
	if err != nil {
		level.Error(logger).Err("error processing guild role delete", err)
	}
	span.SetAttributes(telemetry.KVString("gid", gid.ToString()))

	return gid
}
