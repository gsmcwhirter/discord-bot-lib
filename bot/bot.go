package bot

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/gsmcwhirter/go-util/v8/errors"
	"github.com/gsmcwhirter/go-util/v8/logging/level"
	"github.com/gsmcwhirter/go-util/v8/request"
	"github.com/gsmcwhirter/go-util/v8/telemetry"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"

	"github.com/gsmcwhirter/discord-bot-lib/v22/bot/session"
	"github.com/gsmcwhirter/discord-bot-lib/v22/discordapi/entity"
	"github.com/gsmcwhirter/discord-bot-lib/v22/discordapi/jsonapi"
	"github.com/gsmcwhirter/discord-bot-lib/v22/errreport"
	"github.com/gsmcwhirter/discord-bot-lib/v22/logging"
	"github.com/gsmcwhirter/discord-bot-lib/v22/snowflake"
	"github.com/gsmcwhirter/discord-bot-lib/v22/stats"
	"github.com/gsmcwhirter/discord-bot-lib/v22/wsapi"
)

type dependencies interface {
	Logger() Logger
	DiscordJSONClient() *jsonapi.DiscordJSONClient
	WSClient() wsapi.WSClient
	MessageRateLimiter() *rate.Limiter
	ConnectRateLimiter() *rate.Limiter
	BotSession() *session.Session
	Dispatcher() Dispatcher
	ErrReporter() errreport.Reporter
	Census() *telemetry.Census
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

	UseSlashCommands    bool
	GlobalSlashCommands []entity.ApplicationCommand
	GuildSlashCommands  map[snowflake.Snowflake][]entity.ApplicationCommand
}

// HBReconfig
type hbReconfig struct {
	ctx      context.Context
	interval int
}

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

func (d *DiscordBot) SetDebug(val bool) {
	d.debug = val
}

func (d *DiscordBot) Intents() int {
	return d.intents
}

func (d *DiscordBot) Dispatcher() Dispatcher {
	return d.deps.Dispatcher()
}

func (d *DiscordBot) AuthenticateAndConnect() error {
	ctx := request.NewRequestContext()
	logger := logging.WithContext(ctx, d.deps.Logger())

	if d.config.UseSlashCommands {
		if err := d.DiffAndRegisterSlashCommands(ctx); err != nil {
			return errors.Wrap(err, "could not DiffAndRegisterSlashCommands")
		}
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

	scope := "bot"
	if d.config.UseSlashCommands {
		scope = "applications.commands%20bot"
	}
	fmt.Printf("\nTo add to a guild, go to: https://discordapp.com/api/oauth2/authorize?client_id=%s&scope=%s&permissions=%d\n\n", d.config.ClientID, scope, d.permissions)

	return nil
}

var ErrDuplicateCommand = errors.New("duplicate command")

func (d *DiscordBot) DiffAndRegisterSlashCommands(ctx context.Context) error {
	ctx, span := d.deps.Census().StartSpan(ctx, "DiscordBot.DiffAndRegisterSlashCommands")
	defer span.End()

	logger := logging.WithContext(ctx, d.deps.Logger())

	c := d.deps.DiscordJSONClient()

	level.Debug(logger).Message("starting global command registration")
	if _, err := c.BulkOverwriteGlobalCommands(ctx, d.config.ClientID, d.config.GlobalSlashCommands); err != nil {
		return errors.Wrap(err, "could not BulkOverwriteGlobalCommands")
	}

	for gid, cmds := range d.config.GuildSlashCommands {
		level.Debug(logger).Message("starting guild command registration", "gid", gid, "cmds", fmt.Sprintf("%#v", cmds))
		if _, err := c.BulkOverwriteGuildCommands(ctx, d.config.ClientID, gid, cmds); err != nil {
			return errors.Wrap(err, "could not BulkOverwriteGuildCommands", "gid", gid.ToString())
		}
	}

	return nil

	// level.Info(logger).Message("diffing global commands")
	// cmds, err := c.GetGlobalCommands(ctx, d.config.ClientID)
	// if err != nil {
	// 	return errors.Wrap(err, "could not GetGlobalCommands")
	// }

	// existing := map[jsonapi.ApplicationCommandType]map[string]snowflake.Snowflake{}
	// want := map[jsonapi.ApplicationCommandType]map[string]int{}

	// for i, v := range d.config.GlobalSlashCommands {
	// 	tcmds, ok := want[v.Type]
	// 	if !ok {
	// 		tcmds = map[string]int{
	// 			v.Name: i,
	// 		}
	// 	} else {
	// 		if _, ok := tcmds[v.Name]; ok {
	// 			return errors.WithDetails(ErrDuplicateCommand, "name", v.Name, "type", v.Type)
	// 		}

	// 		tcmds[v.Name] = i
	// 	}

	// 	want[v.Type] = tcmds
	// }

	// for _, v := range cmds {
	// 	tcmds, ok := existing[v.Type]
	// 	if !ok {
	// 		tcmds = map[string]snowflake.Snowflake{
	// 			v.Name: v.IDSnowflake,
	// 		}
	// 	} else {
	// 		// don't need to check dupes here b/c discord handles it
	// 		tcmds[v.Name] = v.IDSnowflake
	// 	}

	// 	existing[v.Type] = tcmds
	// }

	// I think bulk overwrite should work, so try that instead of manual diff

}

func (d *DiscordBot) ReconfigureHeartbeat(ctx context.Context, interval int) {
	ctx, span := d.deps.Census().StartSpan(ctx, "DiscordBot.ReconfigureHeartbeat")
	defer span.End()

	d.heartbeats <- hbReconfig{
		ctx:      ctx,
		interval: interval,
	}
}

func (d *DiscordBot) Config() Config {
	return d.config
}

func (d *DiscordBot) Disconnect() error {
	d.deps.WSClient().Close()
	return nil
}

func (d *DiscordBot) Run(ctx context.Context) error {
	if err := stats.Register(); err != nil {
		return errors.Wrap(err, "could not register stats")
	}

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

func (d *DiscordBot) LastSequence() int {
	d.seqLock.Lock()
	defer d.seqLock.Unlock()

	return d.lastSequence
}

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

			reqCtx := req.ctx
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

func (d *DiscordBot) sendHeartbeat(reqCtx context.Context) error {
	reqCtx, span := d.deps.Census().StartSpan(reqCtx, "discordBot.sendHeartbeat")
	defer span.End()

	m, err := d.deps.Dispatcher().GenerateHeartbeat(reqCtx, d.lastSequence)
	if err != nil {
		level.Error(logging.WithContext(reqCtx, d.deps.Logger())).Err("error generating heartbeat", err)
		return errors.Wrap(err, "error generating heartbeat")
	}

	err = d.deps.MessageRateLimiter().Wait(m.Ctx)
	if err != nil {
		level.Error(logging.WithContext(reqCtx, d.deps.Logger())).Err("error rate limiting", err)
		return errors.Wrap(err, "error rate limiting")
	}
	d.deps.WSClient().SendMessage(m)

	return nil
}

func (d *DiscordBot) API() *jsonapi.DiscordJSONClient {
	return d.deps.DiscordJSONClient()
}
