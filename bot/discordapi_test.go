package bot_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/gsmcwhirter/go-util/v8/deferutil"
	"github.com/gsmcwhirter/go-util/v8/telemetry"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"

	"github.com/gsmcwhirter/discord-bot-lib/v19/bot"
	"github.com/gsmcwhirter/discord-bot-lib/v19/bot/session"
	"github.com/gsmcwhirter/discord-bot-lib/v19/discordapi/json"
	"github.com/gsmcwhirter/discord-bot-lib/v19/dispatcher"
	"github.com/gsmcwhirter/discord-bot-lib/v19/errreport"
	"github.com/gsmcwhirter/discord-bot-lib/v19/httpclient"
	"github.com/gsmcwhirter/discord-bot-lib/v19/stats"
	"github.com/gsmcwhirter/discord-bot-lib/v19/wsapi"
	"github.com/gsmcwhirter/discord-bot-lib/v19/wsclient"
)

type nopLogger struct{}

func (l nopLogger) Log(kv ...interface{}) error              { return nil }
func (l nopLogger) Err(m string, e error, kv ...interface{}) {}
func (l nopLogger) Message(m string, kv ...interface{})      {}
func (l nopLogger) Printf(f string, a ...interface{})        {}

type mockHTTPDoer struct{}

func (d *mockHTTPDoer) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte("{}"))),
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

type mockWSDialer struct {
}

func (d *mockWSDialer) Dial(string, http.Header) (wsclient.Conn, *http.Response, error) {
	return &mockWSConn{}, &http.Response{StatusCode: 200}, nil
}

type customTraceExporter struct{}

func (ce *customTraceExporter) ExportSpan(sd *telemetry.SpanData) {
	fmt.Printf("Name: %s\nTraceID: %x\nSpanID: %x\nParentSpanID: %x\nStartTime: %s\nEndTime: %s\nAnnotations: %+v\n\n",
		sd.Name, sd.TraceID, sd.SpanID, sd.ParentSpanID, sd.StartTime, sd.EndTime, sd.Annotations)
}

type customMetricsExporter struct{}

func (ce *customMetricsExporter) ExportView(vd *telemetry.ViewData) {
	fmt.Printf("vd.View: %+v\n%#v\n", vd.View, vd.Rows)
	for i, row := range vd.Rows {
		fmt.Printf("\tRow: %d: %#v\n", i, row)
	}
	fmt.Printf("StartTime: %s EndTime: %s\n\n", vd.Start.Round(0), vd.End.Round(0))
}

type mockdeps struct {
	logger  bot.Logger
	doer    httpclient.Doer
	http    *httpclient.HTTPClient
	wsd     wsclient.Dialer
	ws      *wsclient.WSClient
	msgrl   *rate.Limiter
	cnxrl   *rate.Limiter
	session *session.Session
	mh      bot.Dispatcher
	rep     errreport.Reporter
	census  *telemetry.Census
	actRec  *stats.ActivityRecorder
	jsc     *json.DiscordJSONClient
}

func (d *mockdeps) Logger() bot.Logger                              { return d.logger }
func (d *mockdeps) HTTPDoer() httpclient.Doer                       { return d.doer }
func (d *mockdeps) HTTPClient() json.HTTPClient                     { return d.http }
func (d *mockdeps) WSDialer() wsclient.Dialer                       { return d.wsd }
func (d *mockdeps) WSClient() wsapi.WSClient                        { return d.ws }
func (d *mockdeps) MessageRateLimiter() *rate.Limiter               { return d.msgrl }
func (d *mockdeps) ConnectRateLimiter() *rate.Limiter               { return d.cnxrl }
func (d *mockdeps) BotSession() *session.Session                    { return d.session }
func (d *mockdeps) Dispatcher() bot.Dispatcher                      { return d.mh }
func (d *mockdeps) ErrReporter() errreport.Reporter                 { return d.rep }
func (d *mockdeps) Census() *telemetry.Census                       { return d.census }
func (d *mockdeps) MessageHandlerRecorder() *stats.ActivityRecorder { return d.actRec }
func (d *mockdeps) DiscordJSONClient() *json.DiscordJSONClient      { return d.jsc }

func TestDiscordBot(t *testing.T) {
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
		session: session.NewSession(),
		rep:     errreport.NopReporter{},
		actRec:  stats.NewActivityRecorder(30.0),
	}

	deps.census = telemetry.NewCensus(telemetry.Options{
		StatsExporter: new(customMetricsExporter),
		TraceExporter: new(customTraceExporter),
	})

	deps.ws = wsclient.NewWSClient(deps, wsclient.Options{
		MaxConcurrentHandlers: conf.NumWorkers,
	})
	deps.mh = dispatcher.NewDispatcher(deps)

	deps.http = httpclient.NewHTTPClient(deps)
	deps.jsc = json.NewDiscordJSONClient(deps, conf.APIURL)

	b := bot.NewDiscordBot(deps, conf, 0, 0)
	err := b.AuthenticateAndConnect()
	if assert.Nil(t, err) {
		defer deferutil.CheckDefer(b.Disconnect)
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = b.Run(ctx)
	}
}
