package bot

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gsmcwhirter/go-util/v5/errors"
	"github.com/gsmcwhirter/go-util/v5/logging/level"
	"github.com/gsmcwhirter/go-util/v5/request"
	census "github.com/gsmcwhirter/go-util/v5/stats"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"

	"github.com/gsmcwhirter/discord-bot-lib/v11/errreport"
	"github.com/gsmcwhirter/discord-bot-lib/v11/etfapi"
	"github.com/gsmcwhirter/discord-bot-lib/v11/etfapi/payloads"
	"github.com/gsmcwhirter/discord-bot-lib/v11/httpclient"
	"github.com/gsmcwhirter/discord-bot-lib/v11/jsonapi"
	"github.com/gsmcwhirter/discord-bot-lib/v11/logging"
	"github.com/gsmcwhirter/discord-bot-lib/v11/snowflake"
	"github.com/gsmcwhirter/discord-bot-lib/v11/stats"
	"github.com/gsmcwhirter/discord-bot-lib/v11/wsclient"
)

// ErrResponse is the error that is wrapped and returned when there is a non-200 api response
var ErrResponse = errors.New("error response")

type dependencies interface {
	Logger() logging.Logger
	HTTPClient() httpclient.HTTPClient
	WSClient() wsclient.WSClient
	MessageRateLimiter() *rate.Limiter
	ConnectRateLimiter() *rate.Limiter
	BotSession() *etfapi.Session
	DiscordMessageHandler() DiscordMessageHandler
	ErrReporter() errreport.Reporter
	Census() *census.Census
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
	SendMessage(context.Context, snowflake.Snowflake, JSONMarshaler) (*http.Response, []byte, error)
	UpdateSequence(int) bool
	ReconfigureHeartbeat(context.Context, int)
	LastSequence() int
	Config() Config
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

	heartbeat  *time.Ticker
	heartbeats chan hbReconfig

	seqLock      *sync.Mutex
	lastSequence int
}

// NewDiscordBot creates a new DiscordBot
func NewDiscordBot(deps dependencies, conf Config) DiscordBot {
	d := &discordBot{
		config: conf,
		deps:   deps,

		heartbeats: make(chan hbReconfig),

		seqLock:      &sync.Mutex{},
		lastSequence: -1,
	}

	d.deps.DiscordMessageHandler().ConnectToBot(d)

	return d
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

	resp, body, err := d.deps.HTTPClient().GetBody(ctx, fmt.Sprintf("%s/gateway/bot", d.config.APIURL), nil)
	if err != nil {
		return errors.Wrap(err, "could not get gateway information")
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Wrap(ErrResponse, "non-200 response")
	}

	level.Debug(logger).Message("response stats",
		"response_body", body,
		"response_bytes", len(body),
	)

	respData := jsonapi.GatewayResponse{}
	err = respData.UnmarshalJSON(body)
	if err != nil {
		return errors.Wrap(err, "could not unmarshal gateway information")
	}

	level.Debug(logger).Message("gateway response",
		"gateway_url", respData.URL,
		"gateway_shards", respData.Shards,
	)

	level.Info(logger).Message("acquired gateway url")

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

	// See https://discordapp.com/developers/docs/topics/permissions#permissions-bitwise-permission-flags
	botPermissions := 0x00000040 // add reactions
	botPermissions |= 0x00000400 // view channel (including read messages)
	botPermissions |= 0x00000800 // send messages
	botPermissions |= 0x00002000 // manage messages
	botPermissions |= 0x00010000 // read message history
	botPermissions |= 0x00020000 // mention everyone
	botPermissions |= 0x04000000 // change own nickname

	fmt.Printf("\nTo add to a guild, go to: https://discordapp.com/api/oauth2/authorize?client_id=%s&scope=bot&permissions=%d\n\n", d.config.ClientID, botPermissions)

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

func (d *discordBot) SendMessage(ctx context.Context, cid snowflake.Snowflake, m JSONMarshaler) (resp *http.Response, body []byte, err error) {
	ctx, span := d.deps.Census().StartSpan(ctx, "discordBot.SendMessage")
	defer span.End()

	logger := logging.WithContext(ctx, d.deps.Logger())

	level.Info(logger).Message("sending message to channel")

	var b []byte

	b, err = m.MarshalJSON()
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not marshal message as json")
	}

	level.Info(logger).Message("sending message", "payload", string(b))
	r := bytes.NewReader(b)

	err = d.deps.MessageRateLimiter().Wait(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error waiting for rate limiter")
	}

	header := http.Header{}
	header.Add("Content-Type", "application/json")
	resp, body, err = d.deps.HTTPClient().PostBody(ctx, fmt.Sprintf("%s/channels/%d/messages", d.config.APIURL, cid), &header, r)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not complete the message send")
	}

	if err := d.deps.Census().Record(ctx, []census.Measurement{stats.MessagesPostedCount.M(1)}, census.Tag{Key: stats.TagStatus, Val: fmt.Sprintf("%d", resp.StatusCode)}); err != nil {
		level.Error(logger).Err("could not record stat", err)
	}

	if resp.StatusCode != http.StatusOK {
		err = errors.Wrap(ErrResponse, "non-200 response")
	}

	return resp, body, err
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
