package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"golang.org/x/time/rate"

	"github.com/gsmcwhirter/discord-bot-lib/bot"
	"github.com/gsmcwhirter/discord-bot-lib/etfapi"
	"github.com/gsmcwhirter/discord-bot-lib/httpclient"
	"github.com/gsmcwhirter/discord-bot-lib/messagehandler"
	"github.com/gsmcwhirter/discord-bot-lib/snowflake"
	"github.com/gsmcwhirter/discord-bot-lib/wsclient"
)

type dependencies struct {
	logger  log.Logger
	doer    httpclient.Doer
	http    httpclient.HTTPClient
	wsd     wsclient.Dialer
	ws      wsclient.WSClient
	msgrl   *rate.Limiter
	cnxrl   *rate.Limiter
	session *etfapi.Session
	mh      bot.DiscordMessageHandler
}

func (d *dependencies) Close()                                           {}
func (d *dependencies) Logger() log.Logger                               { return d.logger }
func (d *dependencies) HTTPDoer() httpclient.Doer                        { return d.doer }
func (d *dependencies) HTTPClient() httpclient.HTTPClient                { return d.http }
func (d *dependencies) WSDialer() wsclient.Dialer                        { return d.wsd }
func (d *dependencies) WSClient() wsclient.WSClient                      { return d.ws }
func (d *dependencies) MessageRateLimiter() *rate.Limiter                { return d.msgrl }
func (d *dependencies) ConnectRateLimiter() *rate.Limiter                { return d.cnxrl }
func (d *dependencies) BotSession() *etfapi.Session                      { return d.session }
func (d *dependencies) DiscordMessageHandler() bot.DiscordMessageHandler { return d.mh }

type mockHTTPDoer struct{}

func (d *mockHTTPDoer) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte("{}"))),
	}, nil
}

func createDependencies(c config, conf bot.BotConfig) (d *dependencies, err error) {
	gcreate, err := guildCreate(snowflake.Snowflake(12345), "Test Guild!")
	if err != nil {
		return nil, err
	}

	d = &dependencies{
		doer: &mockHTTPDoer{},
		wsd: &mockWSDialer{
			MsgType: int(wsclient.Binary),
			First: [][]byte{
				nil, // TODO: identify
			},
			Repeat: [][]byte{
				gcreate,
			},
		},
		msgrl:   rate.NewLimiter(rate.Every(60*time.Second), 120),
		cnxrl:   rate.NewLimiter(rate.Every(5*time.Second), 1),
		session: etfapi.NewSession(),
	}

	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "timestamp", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	d.logger = logger

	d.http = httpclient.NewHTTPClient(d)
	d.ws = wsclient.NewWSClient(d, wsclient.Options{
		GatewayURL:            conf.APIURL,
		MaxConcurrentHandlers: conf.NumWorkers,
	})
	d.mh = messagehandler.NewDiscordMessageHandler(d)

	return
}
