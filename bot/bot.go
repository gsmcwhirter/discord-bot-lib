package bot

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gsmcwhirter/go-util/v8/errors"
	"github.com/gsmcwhirter/go-util/v8/logging/level"
	"github.com/gsmcwhirter/go-util/v8/request"
	"github.com/gsmcwhirter/go-util/v8/telemetry"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"

	"github.com/gsmcwhirter/discord-bot-lib/v19/bot/session"
	"github.com/gsmcwhirter/discord-bot-lib/v19/discordapi/json"
	"github.com/gsmcwhirter/discord-bot-lib/v19/errreport"
	"github.com/gsmcwhirter/discord-bot-lib/v19/logging"
	"github.com/gsmcwhirter/discord-bot-lib/v19/snowflake"
	"github.com/gsmcwhirter/discord-bot-lib/v19/stats"
	"github.com/gsmcwhirter/discord-bot-lib/v19/wsapi"
)

type dependencies interface {
	Logger() logging.Logger
	HTTPClient() HTTPClient
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

func (d *DiscordBot) Intents() int {
	return d.intents
}

func (d *DiscordBot) AddMessageHandler(event string, handler DispatchHandlerFunc) {
	d.deps.Dispatcher().AddHandler(event, handler)
}

func (d *DiscordBot) AuthenticateAndConnect() error {
	ctx := request.NewRequestContext()
	logger := logging.WithContext(ctx, d.deps.Logger())

	err := d.deps.ConnectRateLimiter().Wait(ctx)
	if err != nil {
		return errors.Wrap(err, "connection rate limit error")
	}

	respData, err := d.GetGateway(ctx)
	if err != nil {
		return errors.Wrap(err, "could not get gateway information")
	}

	connectURL, err := url.Parse(respData.URL)
	if err != nil {
		return errors.Wrap(err, "could not parse connection url")
	}

	q := connectURL.Query()
	q.Add("v", "6")
	q.Add("encoding", "etf")
	connectURL.RawQuery = q.Encode()

	level.Info(logger).Message("connecting to gateway",
		"gateway_url", connectURL.String(),
	)

	err = d.deps.WSClient().Connect(connectURL.String(), d.config.BotToken)
	if err != nil {
		return errors.Wrap(err, "could not WSClient().Connect()")
	}

	fmt.Printf("\nTo add to a guild, go to: https://discordapp.com/api/oauth2/authorize?client_id=%s&scope=bot&permissions=%d\n\n", d.config.ClientID, d.permissions)

	return nil
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
			level.Debug(d.deps.Logger()).Message("bum-bum")
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

// ErrResponse is the error that is wrapped and returned when there is a non-200 api response
var ErrResponse = errors.New("error response")

func (d *DiscordBot) GetGateway(ctx context.Context) (json.GatewayResponse, error) {
	ctx, span := d.deps.Census().StartSpan(ctx, "DiscordBot.GetGateway")
	defer span.End()

	logger := logging.WithContext(ctx, d.deps.Logger())

	respData := json.GatewayResponse{}

	_, err := d.deps.HTTPClient().GetJSON(ctx, fmt.Sprintf("%s/gateway/bot", d.config.APIURL), nil, &respData)
	if err != nil {
		return respData, errors.Wrap(err, "could not get gateway information")
	}

	level.Debug(logger).Message("gateway response",
		"gateway_url", respData.URL,
		"gateway_shards", respData.Shards,
	)

	level.Info(logger).Message("acquired gateway url")

	return respData, nil
}

func (d *DiscordBot) SendMessage(ctx context.Context, cid snowflake.Snowflake, m JSONMarshaler) (respData json.MessageResponse, err error) {
	ctx, span := d.deps.Census().StartSpan(ctx, "DiscordBot.SendMessage")
	defer span.End()

	logger := logging.WithContext(ctx, d.deps.Logger())

	level.Info(logger).Message("sending message to channel")

	var b []byte

	b, err = m.MarshalJSON()
	if err != nil {
		return respData, errors.Wrap(err, "could not marshal message as json")
	}

	level.Info(logger).Message("sending message", "payload", string(b))
	r := bytes.NewReader(b)

	err = d.deps.MessageRateLimiter().Wait(ctx)
	if err != nil {
		return respData, errors.Wrap(err, "error waiting for rate limiter")
	}

	header := http.Header{}
	header.Add("Content-Type", "application/json")
	resp, err := d.deps.HTTPClient().PostJSON(ctx, fmt.Sprintf("%s/channels/%d/messages", d.config.APIURL, cid), &header, r, &respData)
	if err != nil {
		return respData, errors.Wrap(err, "could not complete the message send")
	}

	if err := d.deps.Census().Record(ctx, []telemetry.Measurement{stats.MessagesPostedCount.M(1)}, telemetry.Tag{Key: stats.TagStatus, Val: fmt.Sprintf("%d", resp.StatusCode)}); err != nil {
		level.Error(logger).Err("could not record stat", err)
	}

	err = respData.Snowflakify()
	if err != nil {
		return respData, errors.Wrap(err, "could not snowflakify message response information")
	}

	return respData, err
}

func (d *DiscordBot) GetMessage(ctx context.Context, cid, mid snowflake.Snowflake) (respData json.MessageResponse, err error) {
	ctx, span := d.deps.Census().StartSpan(ctx, "DiscordBot.GetMessage")
	defer span.End()

	logger := logging.WithContext(ctx, d.deps.Logger())

	level.Info(logger).Message("getting message details")

	err = d.deps.MessageRateLimiter().Wait(ctx)
	if err != nil {
		return respData, errors.Wrap(err, "error waiting for rate limiter")
	}

	_, err = d.deps.HTTPClient().GetJSON(ctx, fmt.Sprintf("%s/channels/%d/messages/%d", d.config.APIURL, cid, mid), nil, &respData)
	if err != nil {
		return respData, errors.Wrap(err, "could not complete the message get")
	}

	err = respData.Snowflakify()
	if err != nil {
		return respData, errors.Wrap(err, "could not snowflakify message information")
	}

	return respData, nil
}

func (d *DiscordBot) CreateReaction(ctx context.Context, cid, mid snowflake.Snowflake, emoji string) (resp *http.Response, err error) {
	ctx, span := d.deps.Census().StartSpan(ctx, "DiscordBot.GetMessage")
	defer span.End()

	logger := logging.WithContext(ctx, d.deps.Logger())

	level.Info(logger).Message("creating reaction")

	emoji = strings.TrimSuffix(emoji, ">")

	err = d.deps.MessageRateLimiter().Wait(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error waiting for rate limiter")
	}

	emoji = url.QueryEscape(emoji)
	resp, body, err := d.deps.HTTPClient().PutBody(ctx, fmt.Sprintf("%s/channels/%d/messages/%d/reactions/%s/@me", d.config.APIURL, cid, mid, emoji), nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not complete the reaction create")
	}

	if resp.StatusCode != http.StatusNoContent {
		err = errors.Wrap(ErrResponse, "non-204 response", "status_code", resp.StatusCode, "emoji", emoji, "response_body", string(body))
	}

	return resp, err
}

func (d *DiscordBot) GetGuildMember(ctx context.Context, gid, uid snowflake.Snowflake) (respData json.GuildMemberResponse, err error) {
	ctx, span := d.deps.Census().StartSpan(ctx, "DiscordBot.GetGuildMember")
	defer span.End()

	logger := logging.WithContext(ctx, d.deps.Logger())

	level.Info(logger).Message("getting guild member data")

	err = d.deps.MessageRateLimiter().Wait(ctx)
	if err != nil {
		return respData, errors.Wrap(err, "error waiting for rate limiter")
	}

	_, err = d.deps.HTTPClient().GetJSON(ctx, fmt.Sprintf("%s/guilds/%d/members/%d", d.config.APIURL, gid, uid), nil, &respData)
	if err != nil {
		return respData, errors.Wrap(err, "could not complete the guild member get")
	}

	err = respData.Snowflakify()
	if err != nil {
		return respData, errors.Wrap(err, "could not snowflakify guild member information")
	}

	return respData, nil
}
