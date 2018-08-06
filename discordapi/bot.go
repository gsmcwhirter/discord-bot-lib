package discordapi

import (
	"bytes"
	"context"
	"encoding/json"
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

	"github.com/gsmcwhirter/discord-bot-lib/discordapi/etfapi/payloads"
	"github.com/gsmcwhirter/discord-bot-lib/discordapi/jsonapi"
	"github.com/gsmcwhirter/discord-bot-lib/httpclient"
	"github.com/gsmcwhirter/discord-bot-lib/logging"
	"github.com/gsmcwhirter/discord-bot-lib/snowflake"
	"github.com/gsmcwhirter/discord-bot-lib/util"
	"github.com/gsmcwhirter/discord-bot-lib/wsclient"
)

// ErrResponse TODOC
var ErrResponse = errors.New("error response")

type dependencies interface {
	Logger() log.Logger
	HTTPClient() httpclient.HTTPClient
	WSClient() wsclient.WSClient
	MessageRateLimiter() *rate.Limiter
	ConnectRateLimiter() *rate.Limiter
}

// DiscordBot TODOC
type DiscordBot interface {
	AuthenticateAndConnect() error
	Disconnect() error
	Run(context.Context) error
	AddMessageHandler(event string, handler DiscordMessageHandlerFunc)
	SendMessage(context.Context, snowflake.Snowflake, json.Marshaler) (*http.Response, []byte, error)

	GuildOfChannel(snowflake.Snowflake) (snowflake.Snowflake, bool)
	IsGuildAdmin(snowflake.Snowflake, snowflake.Snowflake) bool

	ChannelName(snowflake.Snowflake) (string, bool)
}

// BotConfig TODOC
type BotConfig struct {
	ClientID     string
	ClientSecret string
	BotToken     string
	APIURL       string
	NumWorkers   int

	OS          string
	BotName     string
	BotPresence string
}

type hbReconfig struct {
	ctx      context.Context
	interval int
}

type discordBot struct {
	config         BotConfig
	deps           dependencies
	messageHandler *discordMessageHandler

	session *Session

	heartbeat  *time.Ticker
	heartbeats chan hbReconfig

	seqLock      *sync.Mutex
	lastSequence int
}

// NewDiscordBot TODOC
func NewDiscordBot(deps dependencies, conf BotConfig) DiscordBot {
	d := &discordBot{
		config: conf,
		deps:   deps,

		session:    NewSession(),
		heartbeats: make(chan hbReconfig),

		seqLock:      &sync.Mutex{},
		lastSequence: -1,
	}

	d.messageHandler = newDiscordMessageHandler(d)

	return d
}

func (d *discordBot) AddMessageHandler(event string, handler DiscordMessageHandlerFunc) {
	d.messageHandler.addHandler(event, handler)
}

func (d *discordBot) AuthenticateAndConnect() error {
	ctx := util.NewRequestContext()
	logger := logging.WithContext(ctx, d.deps.Logger())

	err := d.deps.ConnectRateLimiter().Wait(ctx)
	if err != nil {
		return err
	}

	resp, body, err := d.deps.HTTPClient().GetBody(ctx, fmt.Sprintf("%s/gateway/bot", d.config.APIURL), nil)
	if err != nil {
		return err
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
		return err
	}

	_ = level.Debug(logger).Log(
		"gateway_url", respData.URL,
		"gateway_shards", respData.Shards,
	)

	_ = level.Info(logger).Log("message", "acquired gateway url")

	connectURL, err := url.Parse(respData.URL)
	if err != nil {
		return err
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
	d.deps.WSClient().SetHandler(d.messageHandler)

	err = d.deps.WSClient().Connect(d.config.BotToken)
	if err != nil {
		return err
	}

	// See https://discordapp.com/developers/docs/topics/permissions#permissions-bitwise-permission-flags
	botPermissions := 0x00000040 // add reactions
	botPermissions |= 0x00000400 // view channel (including read messages)
	botPermissions |= 0x00000800 // send messages

	fmt.Printf("\nTo add to a guild, go to: https://discordapp.com/api/oauth2/authorize?client_id=%s&scope=bot&permissions=%d\n\n", d.config.ClientID, botPermissions)

	return nil
}

func (d *discordBot) SendMessage(ctx context.Context, cid snowflake.Snowflake, m json.Marshaler) (resp *http.Response, body []byte, err error) {
	logger := logging.WithContext(ctx, d.deps.Logger())

	_ = level.Info(logger).Log("message", "sending message to channel")

	b, err := m.MarshalJSON()
	if err != nil {
		return
	}
	_ = level.Info(logger).Log("message", "sending message", "payload", string(b))
	r := bytes.NewReader(b)

	err = d.deps.MessageRateLimiter().Wait(ctx)
	if err != nil {
		return nil, nil, err
	}

	header := http.Header{}
	header.Add("Content-Type", "application/json")
	resp, body, err = d.deps.HTTPClient().PostBody(ctx, fmt.Sprintf("%s/channels/%d/messages", d.config.APIURL, cid), &header, r)
	if err != nil {
		err = errors.Wrap(err, "could not complete the message send")
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = errors.Wrap(ErrResponse, "non-200 response")
	}

	return
}

func (d *discordBot) Disconnect() error {
	d.deps.WSClient().Close()
	return nil
}

func (d *discordBot) Run(ctx context.Context) error {
	_ = level.Debug(d.deps.Logger()).Log("message", "setting bot signal watcher")

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

func (d *discordBot) updateSequence(seq int) bool {
	d.seqLock.Lock()
	defer d.seqLock.Unlock()

	if seq < d.lastSequence {
		return false
	}
	d.lastSequence = seq
	return true
}

func (d *discordBot) heartbeatHandler(ctx context.Context) error {
	_ = level.Debug(d.deps.Logger()).Log("message", "waiting for heartbeat config")

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
			_ = level.Debug(logging.WithContext(reqCtx, d.deps.Logger())).Log("message", "manual heartbeat requested")

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
	m, err := payloads.ETFPayloadToMessage(reqCtx, payloads.HeartbeatPayload{
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

func (d *discordBot) GuildOfChannel(cid snowflake.Snowflake) (snowflake.Snowflake, bool) {
	return d.session.GuildOfChannel(cid)
}

func (d *discordBot) IsGuildAdmin(gid snowflake.Snowflake, uid snowflake.Snowflake) bool {
	g, err := d.session.Guild(gid)
	if err != nil {
		return false
	}

	return g.IsAdmin(uid)
}

func (d *discordBot) ChannelName(cid snowflake.Snowflake) (string, bool) {
	return d.session.ChannelName(cid)
}
