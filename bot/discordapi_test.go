package bot_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"

	"github.com/gsmcwhirter/discord-bot-lib/bot"
	"github.com/gsmcwhirter/discord-bot-lib/etfapi"
	"github.com/gsmcwhirter/discord-bot-lib/httpclient"
	"github.com/gsmcwhirter/discord-bot-lib/messagehandler"
	"github.com/gsmcwhirter/discord-bot-lib/wsclient"
)

type mockHTTPDoer struct{}

func (d *mockHTTPDoer) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte("{}"))),
	}, nil
}

// type mockHTTPClient struct {
// }

// func (c *mockHTTPClient) SetHeaders(http.Header) {}
// func (c *mockHTTPClient) Get(ctx context.Context, url string, h *http.Header) (*http.Response, error) {
// 	return &http.Response{StatusCode: 200}, nil
// }
// func (c *mockHTTPClient) GetBody(ctx context.Context, url string, h *http.Header) (*http.Response, []byte, error) {
// 	return &http.Response{StatusCode: 200}, nil, nil
// }
// func (c *mockHTTPClient) Post(ctx context.Context, url string, h *http.Header, b io.Reader) (*http.Response, error) {
// 	return &http.Response{StatusCode: 200}, nil
// }
// func (c *mockHTTPClient) PostBody(ctx context.Context, url string, h *http.Header, b io.Reader) (*http.Response, []byte, error) {
// 	return &http.Response{StatusCode: 200}, nil, nil
// }

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

type mockdeps struct {
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

func (d *mockdeps) Logger() log.Logger                               { return d.logger }
func (d *mockdeps) HTTPDoer() httpclient.Doer                        { return d.doer }
func (d *mockdeps) HTTPClient() httpclient.HTTPClient                { return d.http }
func (d *mockdeps) WSDialer() wsclient.Dialer                        { return d.wsd }
func (d *mockdeps) WSClient() wsclient.WSClient                      { return d.ws }
func (d *mockdeps) MessageRateLimiter() *rate.Limiter                { return d.msgrl }
func (d *mockdeps) ConnectRateLimiter() *rate.Limiter                { return d.cnxrl }
func (d *mockdeps) BotSession() *etfapi.Session                      { return d.session }
func (d *mockdeps) DiscordMessageHandler() bot.DiscordMessageHandler { return d.mh }

func TestDiscordBot(t *testing.T) {
	conf := bot.BotConfig{
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
		logger:  log.NewNopLogger(),
		doer:    &mockHTTPDoer{},
		wsd:     &mockWSDialer{},
		msgrl:   rate.NewLimiter(rate.Every(60*time.Second), 120),
		cnxrl:   rate.NewLimiter(rate.Every(5*time.Second), 1),
		session: etfapi.NewSession(),
	}

	deps.http = httpclient.NewHTTPClient(deps)
	deps.ws = wsclient.NewWSClient(deps, wsclient.Options{
		GatewayURL:            conf.APIURL,
		MaxConcurrentHandlers: conf.NumWorkers,
	})
	deps.mh = messagehandler.NewDiscordMessageHandler(deps)

	bot := bot.NewDiscordBot(deps, conf)
	err := bot.AuthenticateAndConnect()
	if assert.Nil(t, err) {
		defer bot.Disconnect()
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		bot.Run(ctx)
	}
}
