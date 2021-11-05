package jsonapi

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gsmcwhirter/go-util/v8/errors"
	"github.com/gsmcwhirter/go-util/v8/json"
	"github.com/gsmcwhirter/go-util/v8/logging/level"
	"github.com/gsmcwhirter/go-util/v8/telemetry"
	"golang.org/x/time/rate"

	"github.com/gsmcwhirter/discord-bot-lib/v21/logging"
	"github.com/gsmcwhirter/discord-bot-lib/v21/snowflake"
	"github.com/gsmcwhirter/discord-bot-lib/v21/stats"
)

type dependencies interface {
	Logger() Logger
	Census() *telemetry.Census
	MessageRateLimiter() *rate.Limiter
	HTTPClient() HTTPClient
}

type Logger = interface {
	Log(keyvals ...interface{}) error
	Message(string, ...interface{})
	Err(string, error, ...interface{})
	Printf(string, ...interface{})
}

// HTTPClient is the interface of an http client
type HTTPClient interface {
	SetHeaders(http.Header)
	Get(context.Context, string, *http.Header) (*http.Response, error)
	GetBody(context.Context, string, *http.Header) (*http.Response, []byte, error)
	GetJSON(context.Context, string, *http.Header, interface{}) (*http.Response, error)
	Post(context.Context, string, *http.Header, io.Reader) (*http.Response, error)
	PostBody(context.Context, string, *http.Header, io.Reader) (*http.Response, []byte, error)
	PostJSON(context.Context, string, *http.Header, io.Reader, interface{}) (*http.Response, error)
	Put(context.Context, string, *http.Header, io.Reader) (*http.Response, error)
	PutBody(context.Context, string, *http.Header, io.Reader) (*http.Response, []byte, error)
	PutJSON(context.Context, string, *http.Header, io.Reader, interface{}) (*http.Response, error)
}

// marshaler is the interface implemented by types that
// can marshal themselves into valid JSON.
type marshaler interface {
	MarshalToJSON() ([]byte, error)
}

type DiscordJSONClient struct {
	deps   dependencies
	apiURL string

	debug bool
}

func NewDiscordJSONClient(deps dependencies, apiURL string) *DiscordJSONClient {
	return &DiscordJSONClient{
		deps:   deps,
		apiURL: apiURL,
	}
}

func (d *DiscordJSONClient) SetDebug(val bool) {
	d.debug = val
}

func (d *DiscordJSONClient) GetGuildMember(ctx context.Context, gid, uid snowflake.Snowflake) (respData GuildMemberResponse, err error) {
	ctx, span := d.deps.Census().StartSpan(ctx, "DiscordBot.GetGuildMember")
	defer span.End()

	// logger := logging.WithContext(ctx, d.deps.Logger())
	// level.Info(logger).Message("getting guild member data")

	err = d.deps.MessageRateLimiter().Wait(ctx)
	if err != nil {
		return respData, errors.Wrap(err, "error waiting for rate limiter")
	}

	_, err = d.deps.HTTPClient().GetJSON(ctx, fmt.Sprintf("%s/guilds/%d/members/%d", d.apiURL, gid, uid), nil, &respData)
	if err != nil {
		return respData, errors.Wrap(err, "could not complete the guild member get")
	}

	err = respData.Snowflakify()
	if err != nil {
		return respData, errors.Wrap(err, "could not snowflakify guild member information")
	}

	return respData, nil
}

// ErrResponse is the error that is wrapped and returned when there is a non-200 api response
var ErrResponse = errors.New("error response")

func (d *DiscordJSONClient) GetGateway(ctx context.Context) (GatewayResponse, error) {
	ctx, span := d.deps.Census().StartSpan(ctx, "DiscordBot.GetGateway")
	defer span.End()

	logger := logging.WithContext(ctx, d.deps.Logger())

	respData := GatewayResponse{}

	_, err := d.deps.HTTPClient().GetJSON(ctx, fmt.Sprintf("%s/gateway/bot", d.apiURL), nil, &respData)
	if err != nil {
		return respData, errors.Wrap(err, "could not get gateway information")
	}

	if d.debug {
		level.Debug(logger).Message("gateway response",
			"gateway_url", respData.URL,
			"gateway_shards", respData.Shards,
		)
	}

	level.Info(logger).Message("acquired gateway url")

	return respData, nil
}

func (d *DiscordJSONClient) SendMessage(ctx context.Context, cid snowflake.Snowflake, m marshaler) (respData MessageResponse, err error) {
	ctx, span := d.deps.Census().StartSpan(ctx, "DiscordBot.SendMessage")
	defer span.End()

	logger := logging.WithContext(ctx, d.deps.Logger())

	// level.Info(logger).Message("sending message to channel")

	var b []byte

	b, err = m.MarshalToJSON()
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
	resp, err := d.deps.HTTPClient().PostJSON(ctx, fmt.Sprintf("%s/channels/%d/messages", d.apiURL, cid), &header, r, &respData)
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

func (d *DiscordJSONClient) GetMessage(ctx context.Context, cid, mid snowflake.Snowflake) (respData MessageResponse, err error) {
	ctx, span := d.deps.Census().StartSpan(ctx, "DiscordBot.GetMessage")
	defer span.End()

	// logger := logging.WithContext(ctx, d.deps.Logger())
	// level.Info(logger).Message("getting message details")

	err = d.deps.MessageRateLimiter().Wait(ctx)
	if err != nil {
		return respData, errors.Wrap(err, "error waiting for rate limiter")
	}

	_, err = d.deps.HTTPClient().GetJSON(ctx, fmt.Sprintf("%s/channels/%d/messages/%d", d.apiURL, cid, mid), nil, &respData)
	if err != nil {
		return respData, errors.Wrap(err, "could not complete the message get")
	}

	err = respData.Snowflakify()
	if err != nil {
		return respData, errors.Wrap(err, "could not snowflakify message information")
	}

	return respData, nil
}

func (d *DiscordJSONClient) CreateReaction(ctx context.Context, cid, mid snowflake.Snowflake, emoji string) (resp *http.Response, err error) {
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
	resp, body, err := d.deps.HTTPClient().PutBody(ctx, fmt.Sprintf("%s/channels/%d/messages/%d/reactions/%s/@me", d.apiURL, cid, mid, emoji), nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not complete the reaction create")
	}

	if resp.StatusCode != http.StatusNoContent {
		err = errors.Wrap(ErrResponse, "non-204 response", "status_code", resp.StatusCode, "emoji", emoji, "response_body", string(body))
	}

	return resp, err
}

func (d *DiscordJSONClient) GetGlobalCommands(ctx context.Context, aid string) (cmds []ApplicationCommand, err error) {
	ctx, span := d.deps.Census().StartSpan(ctx, "DiscordBot.GetGlobalCommands")
	defer span.End()

	logger := logging.WithContext(ctx, d.deps.Logger())
	level.Info(logger).Message("listing global commands", "aid", aid)

	err = d.deps.MessageRateLimiter().Wait(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error waiting for rate limiter")
	}

	_, err = d.deps.HTTPClient().GetJSON(ctx, fmt.Sprintf("%s/applications/%s/commands", d.apiURL, aid), nil, &cmds)
	if err != nil {
		return nil, errors.Wrap(err, "could not get global commands", "aid", aid)
	}

	for i := range cmds {
		if err := cmds[i].Snowflakify(); err != nil {
			return nil, errors.Wrap(err, "could not snowflakify command")
		}
	}

	return cmds, err
}

func (d *DiscordJSONClient) BulkOverwriteGlobalCommands(ctx context.Context, aid string, cmds []ApplicationCommand) (resCmds []ApplicationCommand, err error) {
	ctx, span := d.deps.Census().StartSpan(ctx, "DiscordBot.BulkOverwriteGlobalCommands")
	defer span.End()

	logger := logging.WithContext(ctx, d.deps.Logger())

	var b []byte

	b, err = json.Marshal(cmds)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal commands as json")
	}

	r := bytes.NewReader(b)

	level.Info(logger).Message("overwriting global commands", "aid", aid, "num_commands", len(cmds))

	err = d.deps.MessageRateLimiter().Wait(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error waiting for rate limiter")
	}

	_, err = d.deps.HTTPClient().PutJSON(ctx, fmt.Sprintf("%s/applications/%s/commands", d.apiURL, aid), nil, r, &resCmds)
	if err != nil {
		return nil, errors.Wrap(err, "could not overwrite global commands", "aid", aid)
	}

	for i := range resCmds {
		if err := resCmds[i].Snowflakify(); err != nil {
			return nil, errors.Wrap(err, "could not snowflakify command")
		}
	}

	return resCmds, err
}

func (d *DiscordJSONClient) GetGuildCommands(ctx context.Context, aid string, gid snowflake.Snowflake) (cmds []ApplicationCommand, err error) {
	ctx, span := d.deps.Census().StartSpan(ctx, "DiscordBot.GetGuildCommands")
	defer span.End()

	logger := logging.WithContext(ctx, d.deps.Logger())
	level.Info(logger).Message("listing guild commands", "aid", aid, "gid", gid.ToString())

	err = d.deps.MessageRateLimiter().Wait(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error waiting for rate limiter")
	}

	_, err = d.deps.HTTPClient().GetJSON(ctx, fmt.Sprintf("%s/applications/%s/guilds/%d/commands", d.apiURL, aid, gid), nil, &cmds)
	if err != nil {
		return nil, errors.Wrap(err, "could not get guild commands", "aid", aid, "gid", gid.ToString())
	}

	for i := range cmds {
		if err := cmds[i].Snowflakify(); err != nil {
			return nil, errors.Wrap(err, "could not snowflakify command")
		}
	}

	return cmds, err
}

func (d *DiscordJSONClient) BulkOverwriteGuildCommands(ctx context.Context, aid string, gid snowflake.Snowflake, cmds []ApplicationCommand) (resCmds []ApplicationCommand, err error) {
	ctx, span := d.deps.Census().StartSpan(ctx, "DiscordBot.BulkOverwriteGuildCommands")
	defer span.End()

	logger := logging.WithContext(ctx, d.deps.Logger())

	var b []byte

	b, err = json.Marshal(cmds)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal commands as json")
	}

	r := bytes.NewReader(b)

	level.Info(logger).Message("overwriting global commands", "aid", aid, "gid", gid.ToString(), "num_commands", len(cmds))

	err = d.deps.MessageRateLimiter().Wait(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error waiting for rate limiter")
	}

	_, err = d.deps.HTTPClient().PutJSON(ctx, fmt.Sprintf("%s/applications/%s/guilds/%d/commands", d.apiURL, aid, gid), nil, r, &resCmds)
	if err != nil {
		return nil, errors.Wrap(err, "could not overwrite global commands", "aid", aid)
	}

	for i := range resCmds {
		if err := resCmds[i].Snowflakify(); err != nil {
			return nil, errors.Wrap(err, "could not snowflakify command")
		}
	}

	return resCmds, err
}
