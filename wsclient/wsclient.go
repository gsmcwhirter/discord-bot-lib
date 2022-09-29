package wsclient

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/gsmcwhirter/go-util/v10/errors"
	"github.com/gsmcwhirter/go-util/v10/logging/level"
	"github.com/gsmcwhirter/go-util/v10/request"
	"github.com/gsmcwhirter/go-util/v10/telemetry"
	"golang.org/x/sync/errgroup"

	"github.com/gsmcwhirter/discord-bot-lib/v24/errreport"
	"github.com/gsmcwhirter/discord-bot-lib/v24/logging"
	"github.com/gsmcwhirter/discord-bot-lib/v24/snowflake"
	"github.com/gsmcwhirter/discord-bot-lib/v24/stats"
	"github.com/gsmcwhirter/discord-bot-lib/v24/wsapi"
)

type dependencies interface {
	Logger() Logger
	WSDialer() Dialer
	ErrReporter() errreport.Reporter
	Telemetry() *telemetry.Telemeter
}

// Logger is the interface expected for logging
type Logger = interface {
	Log(keyvals ...interface{}) error
	Message(string, ...interface{})
	Err(string, error, ...interface{})
	Printf(string, ...interface{})
}

// WSClient is a websocket client. It should be instantiated by NewWSClient
type WSClient struct {
	deps dependencies

	conn    Conn
	handler wsapi.MessageHandler

	responses chan wsapi.WSMessage

	pool       *sync.WaitGroup
	poolTokens chan struct{}

	closeLock *sync.Mutex
	isClosed  bool

	debug bool
}

var _ wsapi.WSClient = (*WSClient)(nil)

// Options enables setting up a WSClient with the desired connection settings
type Options struct {
	MaxConcurrentHandlers int
}

// NewWSClient creates a new WSClient
func NewWSClient(deps dependencies, options Options) *WSClient {
	c := &WSClient{
		deps:      deps,
		closeLock: &sync.Mutex{},
	}

	c.pool = &sync.WaitGroup{}
	if options.MaxConcurrentHandlers <= 0 {
		c.poolTokens = make(chan struct{}, 20)
		c.responses = make(chan wsapi.WSMessage, 20)
	} else {
		c.poolTokens = make(chan struct{}, options.MaxConcurrentHandlers)
		c.responses = make(chan wsapi.WSMessage, options.MaxConcurrentHandlers)
	}

	return c
}

// SetDebug turns on debug mode for the websocket client
func (c *WSClient) SetDebug(val bool) {
	c.debug = val
}

// Connect creates a connection to the gateway
func (c *WSClient) Connect(gatewayURL, token string) error {
	var err error
	ctx := request.NewRequestContext()
	logger := logging.WithContext(ctx, c.deps.Logger())

	dialHeader := http.Header{
		"Authorization": []string{fmt.Sprintf("Bot %s", token)},
	}

	var dialResp *http.Response

	if c.debug {
		level.Debug(logger).Message("ws client dial start",
			"url", gatewayURL,
		)
	}

	start := time.Now()
	c.conn, dialResp, err = c.deps.WSDialer().Dial(gatewayURL, dialHeader)

	level.Info(logger).Message("ws client dial complete",
		"elapsed_ns", time.Since(start).Nanoseconds(),
		"status_code", dialResp.StatusCode,
		"url", gatewayURL,
	)

	if err != nil {
		return err
	}
	if dialResp.Body != nil {
		defer dialResp.Body.Close() // nolint:errcheck // not a real issue here
	}

	// level.Info(logger).Message("ws connected")

	return nil
}

// Close closes the client
func (c *WSClient) Close() {
	c.pool.Wait()
	if c.conn != nil {
		_ = c.conn.Close()
	}
}

func (c *WSClient) gracefulClose() {
	c.closeLock.Lock()
	defer c.closeLock.Unlock()

	if c.isClosed {
		return
	}

	c.isClosed = true

	_ = c.conn.SetReadDeadline(time.Now())
}

// HandleRequests starts various goroutines to read and write to the websocket
func (c *WSClient) HandleRequests(ctx context.Context, handler wsapi.MessageHandler) error {
	controls, ctx := errgroup.WithContext(ctx)

	c.handler = handler

	controls.Go(func() error {
		defer c.deps.ErrReporter().AutoNotify(ctx)
		level.Info(c.deps.Logger()).Message("starting response handler")
		return c.handleResponses(ctx)
	})

	controls.Go(func() error {
		defer c.deps.ErrReporter().AutoNotify(ctx)
		level.Info(c.deps.Logger()).Message("starting message reader")
		return c.readMessages(ctx)
	})

	level.Info(c.deps.Logger()).Message("connected and listening")

	err := controls.Wait()
	level.Info(c.deps.Logger()).Message("shutting down")
	return err
}

func (c *WSClient) readMessages(ctx context.Context) error {
	defer level.Info(c.deps.Logger()).Message("readMessages shutdown complete")

	reader, ctx := errgroup.WithContext(ctx)
	reader.Go(func() error {
		defer c.deps.ErrReporter().AutoNotify(ctx)
		return c.doReads(ctx)
	})

	// watches for close message
	reader.Go(func() error {
		defer c.deps.ErrReporter().AutoNotify(ctx)
		<-ctx.Done()
		level.Info(c.deps.Logger()).Message("readMessages shutting down")
		c.gracefulClose()
		return ctx.Err()
	})

	return reader.Wait()
}

func (c *WSClient) doReads(ctx context.Context) error {
	defer level.Info(c.deps.Logger()).Message("websocket reader done")

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		msgType, msg, err := c.conn.ReadMessage()
		if err != nil {
			level.Error(c.deps.Logger()).Err("read error", err,
				"ws_msg_type", msgType,
				"ws_content", msg,
			)
			return errors.Wrap(err, "read error")
		}

		c.pool.Add(1)
		go c.handleMessageRead(ctx, msgType, msg)
	}
}

func (c *WSClient) handleMessageRead(ctx context.Context, msgType int, msg []byte) {
	defer c.deps.ErrReporter().Recover(ctx)
	defer c.pool.Done()

	ctx, span := c.deps.Telemetry().StartSpan(ctx, "wsclient", "handleMessageRead")
	defer span.End()

	reqCtx := request.NewRequestContextFrom(ctx)

	select {
	case <-ctx.Done():
		return
	default:
	}

	logger := logging.WithContext(reqCtx, c.deps.Logger())
	if err := stats.IncCounter(ctx, c.deps.Telemetry(), "wsclient", stats.RawMessageCount, 1); err != nil {
		level.Error(logger).Err("could not record stat", err)
	}

	mT := wsapi.MessageType(msgType)
	mC := make([]byte, len(msg))
	copy(mC, msg)

	wsMsg := wsapi.WSMessage{Ctx: reqCtx, MessageType: mT, MessageContents: mC}

	if c.debug {
		level.Debug(logger).Message("received message",
			"ws_msg_type", mT,
			"ws_msg_len", len(mC),
		)
		level.Debug(logger).Message("waiting for worker token")
	}
	c.poolTokens <- struct{}{}
	if c.debug {
		level.Debug(logger).Message("worker token acquired")
	}

	gid := c.handleRequest(wsMsg)
	span.SetAttributes(telemetry.KVString("gid", gid.ToString()))
}

func (c *WSClient) handleRequest(req wsapi.WSMessage) snowflake.Snowflake {
	ctx, span := c.deps.Telemetry().StartSpan(req.Ctx, "wsclient", "handleRequest")
	defer span.End()
	req.Ctx = ctx

	logger := logging.WithContext(req.Ctx, c.deps.Logger())

	defer func() {
		<-c.poolTokens
		if c.debug {
			level.Debug(logger).Message("released worker token")
		}
	}()

	select {
	case <-req.Ctx.Done():
		level.Info(logger).Message("handleRequest received interrupt -- shutting down")
		return 0
	default:
		if c.debug {
			level.Info(logger).Message("handleRequest dispatching request")
		}
		gid := c.handler.HandleRequest(req, c.responses)
		span.SetAttributes(telemetry.KVString("gid", gid.ToString()))
		return gid
	}
}

func (c *WSClient) handleResponses(ctx context.Context) error {
	defer func() {
		level.Info(c.deps.Logger()).Message("handleResponses shutdown complete")
	}()

	for {
		select {
		case <-ctx.Done(): // time to stop
			level.Info(c.deps.Logger()).Message("handleResponses shutting down")

			defer func() { //nolint:gocritic // not a leak
				level.Info(c.deps.Logger()).Message("gracefully closing the socket")
				err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					level.Error(c.deps.Logger()).Err("Unable to write websocket close message", err)
					return
				}
				level.Info(c.deps.Logger()).Message("close message sent")
			}()

			// drain the remaining response queue
			deadline := time.After(5 * time.Second)

		DRAIN_LOOP:
			for {
				select {
				case _, ok := <-c.responses:
					if !ok {
						close(c.responses)
						break DRAIN_LOOP
					}
				case <-deadline:
					break DRAIN_LOOP
				}
			}

			return ctx.Err()

		case resp := <-c.responses: // handle pending responses
			c.processResponse(resp) //nolint:contextcheck // context comes from the response
		}
	}
}

func (c *WSClient) processResponse(resp wsapi.WSMessage) {
	ctx, span := c.deps.Telemetry().StartSpan(resp.Ctx, "wsclient", "processResponse")
	defer span.End()
	resp.Ctx = ctx

	logger := logging.WithContext(resp.Ctx, c.deps.Logger())

	if c.debug {
		level.Debug(logger).Message("starting sending message",
			"ws_msg_type", resp.MessageType,
			"ws_msg_len", len(resp.MessageContents),
		)
	}

	if err := stats.IncCounter(ctx, c.deps.Telemetry(), "wsclient", stats.RawMessagesSentCount, 1); err != nil {
		level.Error(logger).Err("could not record stat", err)
	}

	start := time.Now()
	err := c.conn.WriteMessage(int(resp.MessageType), resp.MessageContents)

	if c.debug {
		level.Info(logger).Message("done sending message",
			"elapsed_ns", time.Since(start).Nanoseconds(),
			"ws_msg_type", resp.MessageType,
			"ws_msg_len", len(resp.MessageContents),
		)
	}

	if err != nil {
		level.Error(logger).Err("error sending message", err)
	}
}

// SendMessage queues a message to be sent to the websocket
func (c *WSClient) SendMessage(msg wsapi.WSMessage) {
	ctx, span := c.deps.Telemetry().StartSpan(msg.Ctx, "wsclient", "SendMessage")
	defer span.End()
	msg.Ctx = ctx

	if c.debug {
		logger := logging.WithContext(msg.Ctx, c.deps.Logger())
		level.Debug(logger).Message("adding message to response queue",
			"ws_msg_type", msg.MessageType,
			"ws_msg_len", len(msg.MessageContents),
		)
	}

	c.responses <- msg
}
