package bot_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gsmcwhirter/go-util/v10/deferutil"
	"github.com/gsmcwhirter/go-util/v10/telemetry"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/metric/nonrecording"
	"golang.org/x/time/rate"

	"github.com/gsmcwhirter/discord-bot-lib/v24/bot"
	"github.com/gsmcwhirter/discord-bot-lib/v24/bot/session"
	"github.com/gsmcwhirter/discord-bot-lib/v24/discordapi/jsonapi"
	"github.com/gsmcwhirter/discord-bot-lib/v24/dispatcher"
	"github.com/gsmcwhirter/discord-bot-lib/v24/errreport"
	"github.com/gsmcwhirter/discord-bot-lib/v24/httpclient"
	"github.com/gsmcwhirter/discord-bot-lib/v24/wsapi"
	"github.com/gsmcwhirter/discord-bot-lib/v24/wsclient"
)

type nopLogger struct{}

func (l nopLogger) Log(kv ...interface{}) error              { return nil }
func (l nopLogger) Err(m string, e error, kv ...interface{}) {}
func (l nopLogger) Message(m string, kv ...interface{})      {}
func (l nopLogger) Printf(f string, a ...interface{})        {}

type mockHTTPDoer struct{}

func (d *mockHTTPDoer) Do(req *http.Request) (*http.Response, error) {
	if req.URL.Path == "/applications/test id/commands" {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte("[]"))),
		}, nil
	}

	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader([]byte("{}"))),
	}, nil
}

type mockWSConn struct{}

func (c *mockWSConn) Close() error                    { return nil }
func (c *mockWSConn) SetReadDeadline(time.Time) error { return nil }
func (c *mockWSConn) ReadMessage() (int, []byte, error) {
	return 0, nil, nil
}

func (c *mockWSConn) WriteMessage(int, []byte) error {
	return nil
}

type mockWSDialer struct{}

func (d *mockWSDialer) Dial(string, http.Header) (wsclient.Conn, *http.Response, error) {
	return &mockWSConn{}, &http.Response{StatusCode: 200}, nil
}

type testSpanExporter struct{}

var _ telemetry.SpanExporter = (*testSpanExporter)(nil)

func (ce *testSpanExporter) ExportSpans(ctx context.Context, spans []telemetry.ReadOnlySpan) error {
	for _, sd := range spans {
		fmt.Printf("Name: %s\nAttributes: %+v\nStatus: %v\nStartTime: %s\nEndTime: %s\n\n",
			sd.Name(), sd.Attributes(), sd.Status(), sd.StartTime(), sd.EndTime())
	}
	return nil
}

func (ce *testSpanExporter) Shutdown(ctx context.Context) error {
	return nil
}

type mockdeps struct {
	logger    bot.Logger
	doer      httpclient.Doer
	http      *httpclient.HTTPClient
	wsd       wsclient.Dialer
	ws        *wsclient.WSClient
	msgrl     *rate.Limiter
	cnxrl     *rate.Limiter
	cregrl    *rate.Limiter
	session   *session.Session
	mh        bot.Dispatcher
	rep       errreport.Reporter
	telemeter *telemetry.Telemeter
	jsc       *jsonapi.DiscordJSONClient
}

func (d *mockdeps) Logger() bot.Logger                            { return d.logger }
func (d *mockdeps) HTTPDoer() httpclient.Doer                     { return d.doer }
func (d *mockdeps) HTTPClient() jsonapi.HTTPClient                { return d.http }
func (d *mockdeps) WSDialer() wsclient.Dialer                     { return d.wsd }
func (d *mockdeps) WSClient() wsapi.WSClient                      { return d.ws }
func (d *mockdeps) MessageRateLimiter() *rate.Limiter             { return d.msgrl }
func (d *mockdeps) ConnectRateLimiter() *rate.Limiter             { return d.cnxrl }
func (d *mockdeps) CommandRegistrationRateLimiter() *rate.Limiter { return d.cregrl }
func (d *mockdeps) BotSession() *session.Session                  { return d.session }
func (d *mockdeps) Dispatcher() bot.Dispatcher                    { return d.mh }
func (d *mockdeps) ErrReporter() errreport.Reporter               { return d.rep }
func (d *mockdeps) Telemetry() *telemetry.Telemeter               { return d.telemeter }
func (d *mockdeps) DiscordJSONClient() *jsonapi.DiscordJSONClient { return d.jsc }

func TestDiscordBot(t *testing.T) {
	t.Parallel()

	conf := bot.Config{
		ClientID:     "test id",
		ClientSecret: "test secret",
		BotToken:     "test token",
		APIURL:       "http://localhost",
		NumWorkers:   10,
		OS:           "Test OS",
		BotName:      "test bot",
		BotPresence:  "test presence",
	}

	deps := &mockdeps{
		logger:  nopLogger{},
		doer:    &mockHTTPDoer{},
		wsd:     &mockWSDialer{},
		msgrl:   rate.NewLimiter(rate.Every(60*time.Second), 120),
		cnxrl:   rate.NewLimiter(rate.Every(5*time.Second), 1),
		cregrl:  rate.NewLimiter(rate.Every(1*time.Second), 2),
		session: session.NewSession(),
		rep:     errreport.NopReporter{},
	}

	deps.telemeter = telemetry.NewTelemeter("test", "test", "test", new(testSpanExporter), nonrecording.NewNoopMeterProvider(), 1.0)

	deps.ws = wsclient.NewWSClient(deps, wsclient.Options{
		MaxConcurrentHandlers: conf.NumWorkers,
	})
	deps.mh = dispatcher.NewDispatcher(deps)

	deps.http = httpclient.NewHTTPClient(deps)
	deps.jsc = jsonapi.NewDiscordJSONClient(deps, conf.APIURL)

	b := bot.NewDiscordBot(deps, conf, 0, 0)
	err := b.AuthenticateAndConnect()
	if assert.NoError(t, err) {
		defer deferutil.CheckDefer(b.Disconnect)
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = b.Run(ctx)
	}
}
