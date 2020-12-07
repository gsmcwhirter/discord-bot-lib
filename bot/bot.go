package bot

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gsmcwhirter/go-util/v7/errors"
	"github.com/gsmcwhirter/go-util/v7/logging/level"
	"github.com/gsmcwhirter/go-util/v7/request"
	"github.com/gsmcwhirter/go-util/v7/telemetry"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"

	"github.com/gsmcwhirter/discord-bot-lib/v17/errreport"
	"github.com/gsmcwhirter/discord-bot-lib/v17/etfapi"
	"github.com/gsmcwhirter/discord-bot-lib/v17/etfapi/payloads"
	"github.com/gsmcwhirter/discord-bot-lib/v17/httpclient"
	"github.com/gsmcwhirter/discord-bot-lib/v17/jsonapi"
	"github.com/gsmcwhirter/discord-bot-lib/v17/logging"
	"github.com/gsmcwhirter/discord-bot-lib/v17/snowflake"
	"github.com/gsmcwhirter/discord-bot-lib/v17/stats"
	"github.com/gsmcwhirter/discord-bot-lib/v17/wsclient"
)

type dependencies interface {
	Logger() logging.Logger
	HTTPClient() httpclient.HTTPClient
	WSClient() wsclient.WSClient
	MessageRateLimiter() *rate.Limiter
	ConnectRateLimiter() *rate.Limiter
	BotSession() *etfapi.Session
	DiscordMessageHandler() DiscordMessageHandler
	ErrReporter() errreport.Reporter
	Census() *telemetry.Census
}

// DiscordMessageHandlerFunc is the api that a bot expects a handler function to have
type DiscordMessageHandlerFunc func(*etfapi.Payload, wsclient.WSMessage, chan<- wsclient.WSMessage) snowflake.Snowflake

// DiscordMessageHandler is the api that a bot expects a handler manager to have
type DiscordMessageHandler interface {
	ConnectToBot(DiscordBot)
	AddHandler(string, DiscordMessageHandlerFunc)
	HandleRequest(wsclient.WSMessage, chan<- wsclient.WSMessage) snowflake.Snowflake
}

// DiscordBot is the api for a discord bot object
type DiscordBot interface {
	AuthenticateAndConnect() error
	Disconnect() error
	Run(context.Context) error
	AddMessageHandler(event string, handler DiscordMessageHandlerFunc)
	SendMessage(context.Context, snowflake.Snowflake, JSONMarshaler) (jsonapi.MessageResponse, error)
	GetMessage(context.Context, snowflake.Snowflake, snowflake.Snowflake) (*http.Response, []byte, error)
	CreateReaction(context.Context, snowflake.Snowflake, snowflake.Snowflake, string) (*http.Response, error)
	GetGuildMember(context.Context, snowflake.Snowflake, snowflake.Snowflake) (jsonapi.GuildMemberResponse, error)
	UpdateSequence(int) bool
	ReconfigureHeartbeat(context.Context, int)
	LastSequence() int
	Config() Config
	Intents() int
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

type discordBot struct {
	config Config
	deps   dependencies

	permissions int
	intents     int

	heartbeat  *time.Ticker
	heartbeats chan hbReconfig

	seqLock      *sync.Mutex
	lastSequence int
}

var _ DiscordBot = (*discordBot)(nil)

// NewDiscordBot creates a new DiscordBot
func NewDiscordBot(deps dependencies, conf Config, permissions, intents int) DiscordBot {
	d := &discordBot{
		config: conf,
		deps:   deps,

		permissions: permissions,
		intents:     intents,

		heartbeats: make(chan hbReconfig),

		seqLock:      &sync.Mutex{},
		lastSequence: -1,
	}

	d.deps.DiscordMessageHandler().ConnectToBot(d)

	return d
}

func (d *discordBot) Intents() int {
	return d.intents
}

func (d *discordBot) AddMessageHandler(event string, handler DiscordMessageHandlerFunc) {
	d.deps.DiscordMessageHandler().AddHandler(event, handler)
}

func (d *discordBot) AuthenticateAndConnect() error {
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

	d.deps.WSClient().SetGateway(connectURL.String())
	d.deps.WSClient().SetHandler(d.deps.DiscordMessageHandler())

	err = d.deps.WSClient().Connect(d.config.BotToken)
	if err != nil {
		return errors.Wrap(err, "could not WSClient().Connect()")
	}

	fmt.Printf("\nTo add to a guild, go to: https://discordapp.com/api/oauth2/authorize?client_id=%s&scope=bot&permissions=%d\n\n", d.config.ClientID, d.permissions)

	return nil
}

func (d *discordBot) ReconfigureHeartbeat(ctx context.Context, interval int) {
	ctx, span := d.deps.Census().StartSpan(ctx, "discordBot.ReconfigureHeartbeat")
	defer span.End()

	d.heartbeats <- hbReconfig{
		ctx:      ctx,
		interval: interval,
	}
}

func (d *discordBot) Config() Config {
	return d.config
}

func (d *discordBot) Disconnect() error {
	d.deps.WSClient().Close()
	return nil
}

func (d *discordBot) Run(ctx context.Context) error {
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
		return d.deps.WSClient().HandleRequests(ctx)
	})

	return g.Wait()
}

func (d *discordBot) LastSequence() int {
	d.seqLock.Lock()
	defer d.seqLock.Unlock()

	return d.lastSequence
}

func (d *discordBot) UpdateSequence(seq int) bool {
	d.seqLock.Lock()
	defer d.seqLock.Unlock()

	if seq < d.lastSequence {
		return false
	}
	d.lastSequence = seq
	return true
}

func (d *discordBot) heartbeatHandler(ctx context.Context) error {
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

func (d *discordBot) sendHeartbeat(reqCtx context.Context) error {
	reqCtx, span := d.deps.Census().StartSpan(reqCtx, "discordBot.sendHeartbeat")
	defer span.End()

	m, err := payloads.ETFPayloadToMessage(reqCtx, &payloads.HeartbeatPayload{
		Sequence: d.lastSequence,
	})
	if err != nil {
		level.Error(logging.WithContext(reqCtx, d.deps.Logger())).Err("error formatting heartbeat", err)
		return errors.Wrap(err, "error formatting heartbeat")
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

func (d *discordBot) GetGateway(ctx context.Context) (jsonapi.GatewayResponse, error) {
	ctx, span := d.deps.Census().StartSpan(ctx, "jsonapiclient.GetGateway")
	defer span.End()

	logger := logging.WithContext(ctx, d.deps.Logger())

	respData := jsonapi.GatewayResponse{}

	resp, body, err := d.deps.HTTPClient().GetBody(ctx, fmt.Sprintf("%s/gateway/bot", d.config.APIURL), nil)
	if err != nil {
		return respData, errors.Wrap(err, "could not get gateway information")
	}
	if resp.StatusCode != http.StatusOK {
		return respData, errors.Wrap(ErrResponse, "non-200 response")
	}

	level.Debug(logger).Message("response stats",
		"response_body", body,
		"response_bytes", len(body),
	)

	err = respData.UnmarshalJSON(body)
	if err != nil {
		return respData, errors.Wrap(err, "could not unmarshal gateway information")
	}

	level.Debug(logger).Message("gateway response",
		"gateway_url", respData.URL,
		"gateway_shards", respData.Shards,
	)

	level.Info(logger).Message("acquired gateway url")

	return respData, nil
}

func (d *discordBot) SendMessage(ctx context.Context, cid snowflake.Snowflake, m JSONMarshaler) (respData jsonapi.MessageResponse, err error) {
	ctx, span := d.deps.Census().StartSpan(ctx, "jsonapi.SendMessage")
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
	resp, body, err := d.deps.HTTPClient().PostBody(ctx, fmt.Sprintf("%s/channels/%d/messages", d.config.APIURL, cid), &header, r)
	if err != nil {
		return respData, errors.Wrap(err, "could not complete the message send")
	}

	if err := d.deps.Census().Record(ctx, []telemetry.Measurement{stats.MessagesPostedCount.M(1)}, telemetry.Tag{Key: stats.TagStatus, Val: fmt.Sprintf("%d", resp.StatusCode)}); err != nil {
		level.Error(logger).Err("could not record stat", err)
	}

	if resp.StatusCode != http.StatusOK {
		err = errors.Wrap(ErrResponse, "non-200 response", "status_code", resp.StatusCode)
	}

	err = respData.UnmarshalJSON(body)
	if err != nil {
		return respData, errors.Wrap(err, "could not unmarshal message response information")
	}

	err = respData.Snowflakify()
	if err != nil {
		return respData, errors.Wrap(err, "could not snowflakify message response information")
	}

	return respData, err
}

func (d *discordBot) GetMessage(ctx context.Context, cid, mid snowflake.Snowflake) (resp *http.Response, body []byte, err error) {
	ctx, span := d.deps.Census().StartSpan(ctx, "jsonapi.GetMessage")
	defer span.End()

	logger := logging.WithContext(ctx, d.deps.Logger())

	level.Info(logger).Message("getting message details")

	err = d.deps.MessageRateLimiter().Wait(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error waiting for rate limiter")
	}

	resp, body, err = d.deps.HTTPClient().GetBody(ctx, fmt.Sprintf("%s/channels/%d/messages/%d", d.config.APIURL, cid, mid), nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not complete the message get")
	}

	if resp.StatusCode != http.StatusOK {
		err = errors.Wrap(ErrResponse, "non-200 response", "status_code", resp.StatusCode)
	}

	return resp, body, err
}

func (d *discordBot) CreateReaction(ctx context.Context, cid, mid snowflake.Snowflake, emoji string) (resp *http.Response, err error) {
	ctx, span := d.deps.Census().StartSpan(ctx, "jsonapi.GetMessage")
	defer span.End()

	logger := logging.WithContext(ctx, d.deps.Logger())

	level.Info(logger).Message("creating reaction")

	err = d.deps.MessageRateLimiter().Wait(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error waiting for rate limiter")
	}

	emoji = url.QueryEscape(emoji)
	resp, err = d.deps.HTTPClient().Put(ctx, fmt.Sprintf("%s/channels/%d/messages/%d/reactions/%s/@me", d.config.APIURL, cid, mid, emoji), nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not complete the reaction create")
	}

	if resp.StatusCode != http.StatusNoContent {
		err = errors.Wrap(ErrResponse, "non-204 response", "status_code", resp.StatusCode)
	}

	return resp, err
}

func (d *discordBot) GetGuildMember(ctx context.Context, gid, uid snowflake.Snowflake) (respData jsonapi.GuildMemberResponse, err error) {
	ctx, span := d.deps.Census().StartSpan(ctx, "jsonapi.GetGuildMember")
	defer span.End()

	logger := logging.WithContext(ctx, d.deps.Logger())

	level.Info(logger).Message("getting guild member data")

	err = d.deps.MessageRateLimiter().Wait(ctx)
	if err != nil {
		return respData, errors.Wrap(err, "error waiting for rate limiter")
	}

	resp, body, err := d.deps.HTTPClient().GetBody(ctx, fmt.Sprintf("%s/guilds/%d/members/%d", d.config.APIURL, gid, uid), nil)
	if err != nil {
		return respData, errors.Wrap(err, "could not complete the guild member get")
	}

	if resp.StatusCode != http.StatusOK {
		return respData, errors.Wrap(ErrResponse, "non-200 response", "status_code", resp.StatusCode)
	}

	err = respData.UnmarshalJSON(body)
	if err != nil {
		return respData, errors.Wrap(err, "could not unmarshal guild member information")
	}

	err = respData.Snowflakify()
	if err != nil {
		return respData, errors.Wrap(err, "could not snowflakify guild member information")
	}

	return respData, nil
}
