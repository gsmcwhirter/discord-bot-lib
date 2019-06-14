package bot

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"golang.org/x/time/rate"

	"github.com/gsmcwhirter/discord-bot-lib/v6/etfapi"
	"github.com/gsmcwhirter/discord-bot-lib/v6/etfapi/payloads"
	"github.com/gsmcwhirter/discord-bot-lib/v6/httpclient"
	"github.com/gsmcwhirter/discord-bot-lib/v6/jsonapi"
	"github.com/gsmcwhirter/discord-bot-lib/v6/logging"
	"github.com/gsmcwhirter/discord-bot-lib/v6/snowflake"
	"github.com/gsmcwhirter/discord-bot-lib/v6/util"
	"github.com/gsmcwhirter/discord-bot-lib/v6/wsclient"
)

// ErrResponse is the error that is wrapped and returned when there is a non-200 api response
var ErrResponse = errors.New("error response")

type dependencies interface {
	Logger() log.Logger
	HTTPClient() httpclient.HTTPClient
	WSClient() wsclient.WSClient
	MessageRateLimiter() *rate.Limiter
	ConnectRateLimiter() *rate.Limiter
	BotSession() *etfapi.Session
	DiscordMessageHandler() DiscordMessageHandler
}

// DiscordMessageHandlerFunc is the api that a bot expects a handler function to have
type DiscordMessageHandlerFunc func(*etfapi.Payload, wsclient.WSMessage, chan<- wsclient.WSMessage)

// DiscordMessageHandler is the api that a bot expects a handler manager to have
type DiscordMessageHandler interface {
	ConnectToBot(DiscordBot)
	AddHandler(string, DiscordMessageHandlerFunc)
	HandleRequest(wsclient.WSMessage, chan<- wsclient.WSMessage)
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
	ctx := util.NewRequestContext()
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

	_ = level.Debug(logger).Log(
		"response_body", body,
		"response_bytes", len(body),
	)

	respData := jsonapi.GatewayResponse{}
	err = respData.UnmarshalJSON(body)
	if err != nil {
		return errors.Wrap(err, "could not unmarshal gateway information")
	}

	_ = level.Debug(logger).Log(
		"gateway_url", respData.URL,
		"gateway_shards", respData.Shards,
	)

	_ = level.Info(logger).Log("message", "acquired gateway url")

	connectURL, err := url.Parse(respData.URL)
	if err != nil {
		return errors.Wrap(err, "could not parse connection url")
	}

	q := connectURL.Query()
	q.Add("v", "6")
	q.Add("encoding", "etf")
	connectURL.RawQuery = q.Encode()

	_ = level.Info(logger).Log(
		"message", "connecting to gateway",
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
	d.heartbeats <- hbReconfig{
		ctx:      ctx,
		interval: interval,
	}
}

func (d *discordBot) Config() Config {
	return d.config
}

func (d *discordBot) SendMessage(ctx context.Context, cid snowflake.Snowflake, m JSONMarshaler) (resp *http.Response, body []byte, err error) {
	logger := logging.WithContext(ctx, d.deps.Logger())

	_ = level.Info(logger).Log("message", "sending message to channel")

	var b []byte

	b, err = m.MarshalJSON()
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not marshal message as json")
	}

	_ = level.Info(logger).Log("message", "sending message", "payload", string(b))
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
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return d.heartbeatHandler(ctx)
	})

	g.Go(func() error {
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
	_ = level.Info(d.deps.Logger()).Log("message", "waiting for heartbeat config")

	// wait for init
	if d.heartbeat == nil {
		select {
		case <-ctx.Done():
			_ = level.Info(d.deps.Logger()).Log("message", "heartbeat loop stopping before config")
			return ctx.Err()
		case req := <-d.heartbeats:
			if req.interval > 0 {
				d.heartbeat = time.NewTicker(time.Duration(req.interval) * time.Millisecond)
				_ = level.Info(d.deps.Logger()).Log("message", "starting heartbeat loop", "interval", req.interval)
			}
		}
	}

	// in the groove
	for {
		select {
		case <-ctx.Done(): // quit
			_ = level.Info(d.deps.Logger()).Log("message", "heartbeat quitting at request")
			d.heartbeat.Stop()
			return ctx.Err()

		case req := <-d.heartbeats: // reconfigure
			if req.interval > 0 {
				_ = level.Info(d.deps.Logger()).Log("message", "reconfiguring heartbeat loop", "interval", req.interval)
				d.heartbeat.Stop()
				d.heartbeat = time.NewTicker(time.Duration(req.interval) * time.Millisecond)
				continue
			}

			reqCtx := req.ctx
			if reqCtx == nil {
				reqCtx = util.NewRequestContextFrom(ctx)
			}
			_ = level.Info(logging.WithContext(reqCtx, d.deps.Logger())).Log("message", "manual heartbeat requested")

			err := d.sendHeartbeat(reqCtx)
			if err != nil {
				return err
			}

		case <-d.heartbeat.C: // tick
			_ = level.Debug(d.deps.Logger()).Log("message", "bum-bum")
			reqCtx := util.NewRequestContextFrom(ctx)
			err := d.sendHeartbeat(reqCtx)
			if err != nil {
				return err
			}
		}
	}
}

func (d *discordBot) sendHeartbeat(reqCtx context.Context) error {
	m, err := payloads.ETFPayloadToMessage(reqCtx, &payloads.HeartbeatPayload{
		Sequence: d.lastSequence,
	})
	if err != nil {
		_ = level.Error(logging.WithContext(reqCtx, d.deps.Logger())).Log("message", "error formatting heartbeat", "err", err)
		return errors.Wrap(err, "error formatting heartbeat")
	}

	err = d.deps.MessageRateLimiter().Wait(m.Ctx)
	if err != nil {
		_ = level.Error(logging.WithContext(reqCtx, d.deps.Logger())).Log("message", "error rate limiting", "err", err)
		return errors.Wrap(err, "error rate limiting")
	}
	d.deps.WSClient().SendMessage(m)

	return nil
}
