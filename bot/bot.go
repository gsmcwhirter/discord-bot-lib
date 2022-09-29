package bot

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/gsmcwhirter/go-util/v10/errors"
	"github.com/gsmcwhirter/go-util/v10/logging/level"
	"github.com/gsmcwhirter/go-util/v10/request"
	"github.com/gsmcwhirter/go-util/v10/telemetry"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"

	"github.com/gsmcwhirter/discord-bot-lib/v24/bot/session"
	"github.com/gsmcwhirter/discord-bot-lib/v24/discordapi/entity"
	"github.com/gsmcwhirter/discord-bot-lib/v24/discordapi/jsonapi"
	"github.com/gsmcwhirter/discord-bot-lib/v24/errreport"
	"github.com/gsmcwhirter/discord-bot-lib/v24/logging"
	"github.com/gsmcwhirter/discord-bot-lib/v24/snowflake"
	"github.com/gsmcwhirter/discord-bot-lib/v24/wsapi"
)

type dependencies interface {
	Logger() Logger
	DiscordJSONClient() *jsonapi.DiscordJSONClient
	WSClient() wsapi.WSClient
	MessageRateLimiter() *rate.Limiter
	ConnectRateLimiter() *rate.Limiter
	CommandRegistrationRateLimiter() *rate.Limiter
	BotSession() *session.Session
	Dispatcher() Dispatcher
	ErrReporter() errreport.Reporter
	Telemetry() *telemetry.Telemeter
}

// Config is the set of configuration options for creating a DiscordBot with NewDiscordBot
type Config struct {
	ClientID     string
	ClientSecret string
	BotToken     string
	APIURL       string
	NumWorkers   int

	OS          string
	BotName     string
	BotPresence string

	GlobalSlashCommands []entity.ApplicationCommand
}

// HBReconfig
type hbReconfig struct {
	ctx      context.Context
	interval int
}

// DiscordBot is the actal bot
type DiscordBot struct {
	config Config
	deps   dependencies

	permissions int
	intents     int

	heartbeat  *time.Ticker
	heartbeats chan hbReconfig

	seqLock      *sync.Mutex
	lastSequence int

	debug bool
}

// NewDiscordBot creates a new DiscordBot
func NewDiscordBot(deps dependencies, conf Config, permissions, intents int) *DiscordBot {
	d := &DiscordBot{
		config: conf,
		deps:   deps,

		permissions: permissions,
		intents:     intents,

		heartbeats: make(chan hbReconfig),

		seqLock:      &sync.Mutex{},
		lastSequence: -1,
	}

	d.deps.Dispatcher().ConnectToBot(d)

	return d
}

// SetDebug turns on/off debug mode
func (d *DiscordBot) SetDebug(val bool) {
	d.debug = val
}

// Intents returns the combined discord intents
func (d *DiscordBot) Intents() int {
	return d.intents
}

// Dispatcher returns the bot dispatcher
func (d *DiscordBot) Dispatcher() Dispatcher {
	return d.deps.Dispatcher()
}

// AuthenticateAndConnect sets up the bot to run
func (d *DiscordBot) AuthenticateAndConnect() error {
	ctx := request.NewRequestContext()
	logger := logging.WithContext(ctx, d.deps.Logger())

	if err := d.RegisterGlobalCommands(ctx); err != nil {
		return errors.Wrap(err, "could not RegisterGlobalCommands")
	}

	err := d.deps.ConnectRateLimiter().Wait(ctx)
	if err != nil {
		return errors.Wrap(err, "connection rate limit error")
	}

	respData, err := d.deps.DiscordJSONClient().GetGateway(ctx)
	if err != nil {
		return errors.Wrap(err, "could not get gateway information")
	}

	connectURL, err := url.Parse(respData.URL)
	if err != nil {
		return errors.Wrap(err, "could not parse connection url")
	}

	q := connectURL.Query()
	q.Add("v", "9")
	q.Add("encoding", "etf")
	connectURL.RawQuery = q.Encode()

	level.Info(logger).Message("connecting to gateway",
		"gateway_url", connectURL.String(),
	)

	err = d.deps.WSClient().Connect(connectURL.String(), d.config.BotToken)
	if err != nil {
		return errors.Wrap(err, "could not WSClient().Connect()")
	}

	scope := "applications.commands%20bot"
	fmt.Printf("\nTo add to a guild, go to: https://discordapp.com/api/oauth2/authorize?client_id=%s&scope=%s&permissions=%d\n\n", d.config.ClientID, scope, d.permissions)

	return nil
}

// ErrDuplicateCommand represents having multiple commands with the same name
var ErrDuplicateCommand = errors.New("duplicate command")

// RegisterGlobalCommands registers the global bot commands with discord
func (d *DiscordBot) RegisterGlobalCommands(ctx context.Context) error {
	ctx, span := d.deps.Telemetry().StartSpan(ctx, "bot", "RegisterGlobalCommands")
	defer span.End()

	logger := logging.WithContext(ctx, d.deps.Logger())
	c := d.deps.DiscordJSONClient()

	level.Debug(logger).Message("starting global command registration")
	if _, err := c.BulkOverwriteGlobalCommands(ctx, d.config.ClientID, d.config.GlobalSlashCommands); err != nil {
		return errors.Wrap(err, "could not BulkOverwriteGlobalCommands")
	}

	return nil
}

// RegisterGuildCommands registers the guild-specific commands for a guild with discord
func (d *DiscordBot) RegisterGuildCommands(ctx context.Context, gid snowflake.Snowflake, cmds []entity.ApplicationCommand) ([]entity.ApplicationCommand, error) {
	ctx, span := d.deps.Telemetry().StartSpan(ctx, "bot", "RegisterGuildCommands", telemetry.WithAttributes(telemetry.KVString("gid", gid.ToString())))
	defer span.End()

	logger := logging.WithContext(ctx, d.deps.Logger())
	c := d.deps.DiscordJSONClient()

	level.Debug(logger).Message("starting guild command registration", "gid", gid, "cmds", fmt.Sprintf("%#v", cmds))
	learned, err := c.BulkOverwriteGuildCommands(ctx, d.config.ClientID, gid, cmds)
	return learned, errors.Wrap(err, "could not BulkOverwriteGuildCommands", "gid", gid.ToString())
}

// ReconfigureHeartbeat re-configures the heartbeat ticker
func (d *DiscordBot) ReconfigureHeartbeat(ctx context.Context, interval int) {
	ctx, span := d.deps.Telemetry().StartSpan(ctx, "bot", "ReconfigureHeartbeat")
	defer span.End()

	d.heartbeats <- hbReconfig{
		ctx:      ctx,
		interval: interval,
	}
}

// Config returns the bot config
func (d *DiscordBot) Config() Config {
	return d.config
}

// Disconnect stops the bot
func (d *DiscordBot) Disconnect() error {
	d.deps.WSClient().Close()
	return nil
}

// Run starts handling websocket requests and heartbeats after calling AuthenticateAndConnect
func (d *DiscordBot) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		defer d.deps.ErrReporter().AutoNotify(ctx)
		return d.heartbeatHandler(ctx)
	})

	g.Go(func() error {
		defer d.deps.ErrReporter().AutoNotify(ctx)
		return d.deps.WSClient().HandleRequests(ctx, d.deps.Dispatcher())
	})

	return g.Wait()
}

// LastSequence is the last sequence number seen
func (d *DiscordBot) LastSequence() int {
	d.seqLock.Lock()
	defer d.seqLock.Unlock()

	return d.lastSequence
}

// UpdateSequence updates the sequence number if it is newer
func (d *DiscordBot) UpdateSequence(seq int) bool {
	d.seqLock.Lock()
	defer d.seqLock.Unlock()

	if seq < d.lastSequence {
		return false
	}
	d.lastSequence = seq
	return true
}

func (d *DiscordBot) heartbeatHandler(ctx context.Context) error {
	level.Info(d.deps.Logger()).Message("waiting for heartbeat config")

	// wait for init
	if d.heartbeat == nil {
		select {
		case <-ctx.Done():
			level.Info(d.deps.Logger()).Message("heartbeat loop stopping before config")
			return ctx.Err()
		case req := <-d.heartbeats:
			if req.interval > 0 {
				d.heartbeat = time.NewTicker(time.Duration(req.interval) * time.Millisecond)
				level.Info(d.deps.Logger()).Message("starting heartbeat loop", "interval", req.interval)
			}
		}
	}

	// in the groove
	for {
		select {
		case <-ctx.Done(): // quit
			level.Info(d.deps.Logger()).Message("heartbeat quitting at request")
			d.heartbeat.Stop()
			return ctx.Err()

		case req := <-d.heartbeats: // reconfigure
			if req.interval > 0 {
				level.Info(d.deps.Logger()).Message("reconfiguring heartbeat loop", "interval", req.interval)
				d.heartbeat.Stop()
				d.heartbeat = time.NewTicker(time.Duration(req.interval) * time.Millisecond)
				continue
			}

			reqCtx := req.ctx // nolint:contextcheck // not a real issue -- the function context is not per-request
			if reqCtx == nil {
				reqCtx = request.NewRequestContextFrom(ctx)
			}
			level.Info(logging.WithContext(reqCtx, d.deps.Logger())).Message("manual heartbeat requested")

			err := d.sendHeartbeat(reqCtx)
			if err != nil {
				return err
			}

		case <-d.heartbeat.C: // tick
			if d.debug {
				level.Debug(d.deps.Logger()).Message("bum-bum")
			}
			reqCtx := request.NewRequestContextFrom(ctx)
			err := d.sendHeartbeat(reqCtx)
			if err != nil {
				return err
			}
		}
	}
}

func (d *DiscordBot) sendHeartbeat(ctx context.Context) error {
	ctx, span := d.deps.Telemetry().StartSpan(ctx, "bot", "sendHeartbeat")
	defer span.End()

	m, err := d.deps.Dispatcher().GenerateHeartbeat(ctx, d.lastSequence)
	if err != nil {
		level.Error(logging.WithContext(ctx, d.deps.Logger())).Err("error generating heartbeat", err)
		return errors.Wrap(err, "error generating heartbeat")
	}

	err = d.deps.MessageRateLimiter().Wait(ctx)
	if err != nil {
		level.Error(logging.WithContext(ctx, d.deps.Logger())).Err("error rate limiting", err)
		return errors.Wrap(err, "error rate limiting")
	}
	d.deps.WSClient().SendMessage(m)

	return nil
}

// API returns the DiscordJSONClient
func (d *DiscordBot) API() *jsonapi.DiscordJSONClient {
	return d.deps.DiscordJSONClient()
}
