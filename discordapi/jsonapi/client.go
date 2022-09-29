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

	"github.com/gsmcwhirter/discord-bot-lib/v24/discordapi/entity"
	"github.com/gsmcwhirter/discord-bot-lib/v24/httpclient"
	"github.com/gsmcwhirter/discord-bot-lib/v24/logging"
	"github.com/gsmcwhirter/discord-bot-lib/v24/snowflake"
	"github.com/gsmcwhirter/discord-bot-lib/v24/stats"
)

type dependencies interface {
	Logger() Logger
	Census() *telemetry.Census
	MessageRateLimiter() *rate.Limiter
	CommandRegistrationRateLimiter() *rate.Limiter
	HTTPClient() HTTPClient
}

// Logger is the interface expected for logging
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

// DiscordJSONClient is a json client for interacting wth discord
type DiscordJSONClient struct {
	deps   dependencies
	apiURL string

	debug bool
}

// NewDiscordJSONClient creates a new DiscordJSONClient
func NewDiscordJSONClient(deps dependencies, apiURL string) *DiscordJSONClient {
	return &DiscordJSONClient{
		deps:   deps,
		apiURL: apiURL,
	}
}

// SetDebug turns on or off debugging for the client
func (d *DiscordJSONClient) SetDebug(val bool) {
	d.debug = val
}

// GetGuildMember retrieves information about a guild memeber
func (d *DiscordJSONClient) GetGuildMember(ctx context.Context, gid, uid snowflake.Snowflake) (respData entity.GuildMember, err error) {
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

// Err is the error that is wrapped and returned when there is a non-200 api response
var Err = errors.New("error response")

// GetGateway retrieves information about the gateway
func (d *DiscordJSONClient) GetGateway(ctx context.Context) (entity.Gateway, error) {
	ctx, span := d.deps.Census().StartSpan(ctx, "DiscordBot.GetGateway")
	defer span.End()

	logger := logging.WithContext(ctx, d.deps.Logger())

	respData := entity.Gateway{}

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

// SendMessage sends a message to a channel
func (d *DiscordJSONClient) SendMessage(ctx context.Context, cid snowflake.Snowflake, m marshaler) (respData entity.Message, err error) {
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

	header := &http.Header{}
	header.Set("Content-Type", "application/json")
	resp, err := d.deps.HTTPClient().PostJSON(ctx, fmt.Sprintf("%s/channels/%d/messages", d.apiURL, cid), header, r, &respData)
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

// SendInteractionMessage sends an interaction response message
func (d *DiscordJSONClient) SendInteractionMessage(ctx context.Context, ixID snowflake.Snowflake, ixToken string, m marshaler) error {
	err := d.sendInteractionResponse(ctx, ixID, ixToken, m, CallbackTypeChannelMessage)

	if err2 := d.deps.Census().Record(ctx, []telemetry.Measurement{stats.InteractionResponsesCount.M(1)}); err2 != nil {
		logger := logging.WithContext(ctx, d.deps.Logger())
		level.Error(logger).Err("could not record stat", err)
	}

	return err
}

// SendInteractionAutocomplete sends an interaction autocomplete response
func (d *DiscordJSONClient) SendInteractionAutocomplete(ctx context.Context, ixID snowflake.Snowflake, ixToken string, m marshaler) error {
	err := d.sendInteractionResponse(ctx, ixID, ixToken, m, CallbackTypeAutocomplete)

	if err2 := d.deps.Census().Record(ctx, []telemetry.Measurement{stats.InteractionAutocompletesCount.M(1)}); err2 != nil {
		logger := logging.WithContext(ctx, d.deps.Logger())
		level.Error(logger).Err("could not record stat", err)
	}

	return err
}

// DeferInteractionResponse sends a deferral for an interaction response
func (d *DiscordJSONClient) DeferInteractionResponse(ctx context.Context, ixID snowflake.Snowflake, ixToken string) error {
	ctx, span := d.deps.Census().StartSpan(ctx, "DiscordBot.SendMessage")
	defer span.End()

	logger := logging.WithContext(ctx, d.deps.Logger())

	var b []byte
	var err error

	im := InteractionCallbackMessage{
		Type: CallbackTypeDeferredChannelMessage,
	}

	b, err = json.Marshal(im)
	if err != nil {
		return errors.Wrap(err, "could not marshal InteractionCallbackMessage")
	}

	level.Info(logger).Message("sending message", "payload", string(b))
	r := bytes.NewReader(b)

	err = d.deps.MessageRateLimiter().Wait(ctx)
	if err != nil {
		return errors.Wrap(err, "error waiting for rate limiter")
	}

	header := &http.Header{}
	header.Set("Content-Type", "application/json")
	resp, body, err := d.deps.HTTPClient().PostBody(ctx, fmt.Sprintf("%s/interactions/%d/%s/callback", d.apiURL, ixID, ixToken), header, r)
	if err != nil {
		return errors.Wrap(err, "could not complete the message send")
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= 300 {
		return errors.Wrap(httpclient.ErrResponse, "non-200 response", "status_code", resp.StatusCode, "response_body", string(body))
	}

	if err := d.deps.Census().Record(ctx, []telemetry.Measurement{stats.InteractionDeferralsCount.M(1)}); err != nil {
		level.Error(logger).Err("could not record stat", err)
	}

	return err
}

func (d *DiscordJSONClient) sendInteractionResponse(ctx context.Context, ixID snowflake.Snowflake, ixToken string, m marshaler, typ InteractionCallbackType) (err error) {
	ctx, span := d.deps.Census().StartSpan(ctx, "DiscordBot.SendMessage")
	defer span.End()

	logger := logging.WithContext(ctx, d.deps.Logger())

	var b []byte

	b, err = m.MarshalToJSON()
	if err != nil {
		return errors.Wrap(err, "could not marshal message as json")
	}

	im := InteractionCallbackMessage{
		Type: typ,
	}
	if err := im.Data.UnmarshalJSON(b); err != nil {
		return errors.Wrap(err, "could not fill InteractionCallbackMessage Data")
	}

	b, err = json.Marshal(im)
	if err != nil {
		return errors.Wrap(err, "could not marshal InteractionCallbackMessage")
	}

	level.Info(logger).Message("sending message", "payload", string(b))
	r := bytes.NewReader(b)

	err = d.deps.MessageRateLimiter().Wait(ctx)
	if err != nil {
		return errors.Wrap(err, "error waiting for rate limiter")
	}

	header := &http.Header{}
	header.Set("Content-Type", "application/json")
	resp, body, err := d.deps.HTTPClient().PostBody(ctx, fmt.Sprintf("%s/interactions/%d/%s/callback", d.apiURL, ixID, ixToken), header, r)
	if err != nil {
		return errors.Wrap(err, "could not complete the message send")
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= 300 {
		return errors.Wrap(httpclient.ErrResponse, "non-200 response", "status_code", resp.StatusCode, "response_body", string(body))
	}

	return err
}

// GetInteractionResponse retrieves an interaction response
func (d *DiscordJSONClient) GetInteractionResponse(ctx context.Context, aid snowflake.Snowflake, ixToken string) (respData entity.Message, err error) {
	ctx, span := d.deps.Census().StartSpan(ctx, "DiscordBot.GetMessage")
	defer span.End()

	// logger := logging.WithContext(ctx, d.deps.Logger())
	// level.Info(logger).Message("getting message details")

	err = d.deps.MessageRateLimiter().Wait(ctx)
	if err != nil {
		return respData, errors.Wrap(err, "error waiting for rate limiter")
	}

	_, err = d.deps.HTTPClient().GetJSON(ctx, fmt.Sprintf("%s/webhooks/%d/%s/messages/@original", d.apiURL, aid, ixToken), nil, &respData)
	if err != nil {
		return respData, errors.Wrap(err, "could not complete the message get")
	}

	err = respData.Snowflakify()
	if err != nil {
		return respData, errors.Wrap(err, "could not snowflakify message information")
	}

	return respData, nil
}

// GetMessage retrieves information about a discord message
func (d *DiscordJSONClient) GetMessage(ctx context.Context, cid, mid snowflake.Snowflake) (respData entity.Message, err error) {
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

// CreateReaction adds a reaction to a message
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
		err = errors.Wrap(Err, "non-204 response", "status_code", resp.StatusCode, "emoji", emoji, "response_body", string(body))
	}

	return resp, err
}

// GetGlobalCommands gets the registered global commands
func (d *DiscordJSONClient) GetGlobalCommands(ctx context.Context, aid string) (cmds []entity.ApplicationCommand, err error) {
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

// BulkOverwriteGlobalCommands overwrites the global commands
func (d *DiscordJSONClient) BulkOverwriteGlobalCommands(ctx context.Context, aid string, cmds []entity.ApplicationCommand) (resCmds []entity.ApplicationCommand, err error) {
	ctx, span := d.deps.Census().StartSpan(ctx, "DiscordBot.BulkOverwriteGlobalCommands")
	defer span.End()

	logger := logging.WithContext(ctx, d.deps.Logger())

	var b []byte

	b, err = json.Marshal(cmds)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal commands as json")
	}

	level.Debug(logger).Message("body debug", "body", string(b))

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

	level.Debug(logger).Message("done overwriting globals; snowflakifying")

	for i := range resCmds {
		if err := resCmds[i].Snowflakify(); err != nil {
			return nil, errors.Wrap(err, "could not snowflakify command")
		}
	}

	return resCmds, err
}

// GetGuildCommands gets the currently registered guild commands
func (d *DiscordJSONClient) GetGuildCommands(ctx context.Context, aid string, gid snowflake.Snowflake) (cmds []entity.ApplicationCommand, err error) {
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

// BulkOverwriteGuildCommands overwrites the guild commands
func (d *DiscordJSONClient) BulkOverwriteGuildCommands(ctx context.Context, aid string, gid snowflake.Snowflake, cmds []entity.ApplicationCommand) (resCmds []entity.ApplicationCommand, err error) {
	ctx, span := d.deps.Census().StartSpan(ctx, "DiscordBot.BulkOverwriteGuildCommands")
	defer span.End()

	logger := logging.WithContext(ctx, d.deps.Logger())

	var b []byte

	level.Debug(logger).Message("starting marshal")

	b, err = json.Marshal(cmds)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal commands as json")
	}

	level.Debug(logger).Message("body debug", "body", string(b))

	r := bytes.NewReader(b)

	level.Info(logger).Message("overwriting guild commands", "aid", aid, "gid", gid, "num_commands", len(cmds))

	err = d.deps.CommandRegistrationRateLimiter().Wait(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error waiting for rate limiter")
	}

	_, err = d.deps.HTTPClient().PutJSON(ctx, fmt.Sprintf("%s/applications/%s/guilds/%d/commands", d.apiURL, aid, gid), nil, r, &resCmds)
	if err != nil {
		return nil, errors.Wrap(err, "could not overwrite guild commands", "aid", aid, "gid", gid)
	}

	for i := range resCmds {
		if err := resCmds[i].Snowflakify(); err != nil {
			return nil, errors.Wrap(err, "could not snowflakify command")
		}
	}

	return resCmds, err
}

// BulkOverwriteGuildCommandPermissions overwrites the guild command permissions
func (d *DiscordJSONClient) BulkOverwriteGuildCommandPermissions(ctx context.Context, aid string, gid snowflake.Snowflake, perms []entity.ApplicationCommandPermissions) (resPerms []entity.ApplicationCommandPermissions, err error) {
	ctx, span := d.deps.Census().StartSpan(ctx, "DiscordBot.BulkOverwriteGuildCommandPermissions")
	defer span.End()

	logger := logging.WithContext(ctx, d.deps.Logger())

	var b []byte

	level.Debug(logger).Message("starting marshal")

	b, err = json.Marshal(perms)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal permissions as json")
	}

	level.Debug(logger).Message("body debug", "body", string(b))

	r := bytes.NewReader(b)

	level.Info(logger).Message("overwriting guild command permissions", "aid", aid, "gid", gid, "num_permissions", len(perms))

	err = d.deps.MessageRateLimiter().Wait(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error waiting for rate limiter")
	}

	_, err = d.deps.HTTPClient().PutJSON(ctx, fmt.Sprintf("%s/applications/%s/guilds/%d/commands/permissions", d.apiURL, aid, gid), nil, r, &resPerms)
	if err != nil {
		return nil, errors.Wrap(err, "could not overwrite guild command permissions", "aid", aid, "gid", gid)
	}

	for i := range resPerms {
		if err := resPerms[i].Snowflakify(); err != nil {
			return nil, errors.Wrap(err, "could not snowflakify command permissions")
		}
	}

	return resPerms, err
}
