package wsclient

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

	"github.com/gsmcwhirter/discord-bot-lib/logging"
	"github.com/gsmcwhirter/discord-bot-lib/util"
)

// WSClient is the api for a client that maintains an active websocket connection and hands
// off messages to be processed.
type WSClient interface {
	SetGateway(string)
	SetHandler(MessageHandler)
	Connect(string) error
	Close()
	HandleRequests(context.Context) error
	SendMessage(msg WSMessage)
}

type dependencies interface {
	Logger() log.Logger
}

type wsClient struct {
	deps dependencies

	gatewayURL string
	dialer     *websocket.Dialer
	conn       *websocket.Conn
	handler    MessageHandler

	responses chan WSMessage

	pool       *sync.WaitGroup
	poolTokens chan struct{}

	closeLock *sync.Mutex
	isClosed  bool
}

// Options enables setting up a WSClient with the desired connection settings
type Options struct {
	GatewayURL            string
	Dialer                *websocket.Dialer
	MaxConcurrentHandlers int
}

// NewWSClient creates a new WSClient
func NewWSClient(deps dependencies, options Options) WSClient {
	c := &wsClient{
		deps:       deps,
		gatewayURL: options.GatewayURL,
		closeLock:  &sync.Mutex{},
	}

	if options.Dialer != nil {
		c.dialer = options.Dialer
	} else {
		c.dialer = websocket.DefaultDialer
	}

	c.pool = &sync.WaitGroup{}
	if options.MaxConcurrentHandlers <= 0 {
		c.poolTokens = make(chan struct{}, 20)
		c.responses = make(chan WSMessage, 20)
	} else {
		c.poolTokens = make(chan struct{}, options.MaxConcurrentHandlers)
		c.responses = make(chan WSMessage, options.MaxConcurrentHandlers)
	}

	return c
}

func (c *wsClient) SetGateway(url string) {
	c.gatewayURL = url
}

func (c *wsClient) SetHandler(handler MessageHandler) {
	c.handler = handler
}

func (c *wsClient) Connect(token string) (err error) {
	ctx := util.NewRequestContext()
	logger := logging.WithContext(ctx, c.deps.Logger())

	dialHeader := http.Header{
		"Authorization": []string{fmt.Sprintf("Bot %s", token)},
	}

	var dialResp *http.Response

	_ = level.Debug(logger).Log(
		"message", "ws client dial start",
		"url", c.gatewayURL,
	)

	start := time.Now()
	c.conn, dialResp, err = c.dialer.Dial(c.gatewayURL, dialHeader)

	_ = level.Debug(logger).Log(
		"message", "ws client dial complete",
		"duration_ns", time.Since(start).Nanoseconds(),
		"status_code", dialResp.StatusCode,
	)

	if err != nil {
		return err
	}
	defer dialResp.Body.Close() // nolint: errcheck

	_ = level.Info(logger).Log("message", "ws connected")

	return
}

func (c *wsClient) Close() {
	c.pool.Wait()
	if c.conn != nil {
		c.conn.Close() // nolint: errcheck
	}
}

func (c *wsClient) gracefulClose() {
	c.closeLock.Lock()
	defer c.closeLock.Unlock()

	if c.isClosed {
		return
	}

	c.isClosed = true

	_ = c.conn.SetReadDeadline(time.Now())
}

func (c *wsClient) HandleRequests(ctx context.Context) error {
	controls, ctx := errgroup.WithContext(ctx)

	controls.Go(func() error {
		_ = level.Debug(c.deps.Logger()).Log("message", "starting response handler")
		return c.handleResponses(ctx)
	})

	controls.Go(func() error {
		_ = level.Debug(c.deps.Logger()).Log("message", "starting message reader")
		return c.readMessages(ctx)
	})

	_ = level.Info(c.deps.Logger()).Log("message", "connected and listening")

	err := controls.Wait()
	_ = level.Info(c.deps.Logger()).Log("message", "shutting down")
	return err
}

func (c *wsClient) doReads(ctx context.Context) error {
	defer level.Info(c.deps.Logger()).Log("message", "websocket reader done") //nolint: errcheck

	for {
		msgType, msg, err := c.conn.ReadMessage()
		if err != nil {
			_ = level.Error(c.deps.Logger()).Log(
				"message", "read error",
				"error", err,
				"ws_msg_type", msgType,
				"ws_content", msg,
			)
			return errors.Wrap(err, "read error")
		}

		reqCtx := util.NewRequestContextFrom(ctx)
		mT := MessageType(msgType)
		mC := make([]byte, len(msg))
		copy(mC, msg)

		wsMsg := WSMessage{Ctx: reqCtx, MessageType: mT, MessageContents: mC}
		_ = level.Info(logging.WithContext(reqCtx, c.deps.Logger())).Log(
			"message", "received message",
			"ws_msg_type", mT,
			"ws_msg_len", len(mC),
		)

		_ = level.Debug(logging.WithContext(reqCtx, c.deps.Logger())).Log(
			"message", "waiting for worker token",
		)
		c.poolTokens <- struct{}{}
		_ = level.Info(logging.WithContext(reqCtx, c.deps.Logger())).Log(
			"message", "worker token acquired",
		)
		c.pool.Add(1)
		go c.handleRequest(wsMsg)
	}
}

func (c *wsClient) readMessages(ctx context.Context) error {
	defer level.Info(c.deps.Logger()).Log("message", "readMessages shutdown complete") //nolint: errcheck

	reader, ctx := errgroup.WithContext(ctx)
	reader.Go(func() error {
		return c.doReads(ctx)
	})

	// watches for close message
	reader.Go(func() error {
		<-ctx.Done()
		_ = level.Info(c.deps.Logger()).Log("message", "readMessages shutting down")
		c.gracefulClose()
		return ctx.Err()
	})

	return reader.Wait()

}

func (c *wsClient) handleRequest(req WSMessage) {
	defer c.pool.Done()

	logger := logging.WithContext(req.Ctx, c.deps.Logger())

	defer func() {
		<-c.poolTokens
		_ = level.Info(logger).Log("message", "released worker token")
	}()

	select {
	case <-req.Ctx.Done():
		_ = level.Info(logger).Log("message", "handleRequest received interrupt -- shutting down")
		return
	default:
		_ = level.Debug(logger).Log("message", "handleRequest dispatching request")
		c.handler.HandleRequest(req, c.responses)
	}
}

func (c *wsClient) processResponse(resp WSMessage) {
	logger := logging.WithContext(resp.Ctx, c.deps.Logger())

	_ = level.Debug(logger).Log(
		"message", "starting sending message",
		"ws_msg_type", resp.MessageType,
		"ws_msg_len", len(resp.MessageContents),
	)

	start := time.Now()
	err := c.conn.WriteMessage(int(resp.MessageType), resp.MessageContents)

	_ = level.Info(logging.WithContext(resp.Ctx, c.deps.Logger())).Log(
		"message", "done sending message",
		"elapsed_ns", time.Since(start).Nanoseconds(),
	)

	if err != nil {
		_ = level.Error(logger).Log(
			"message", "error sending message",
			"error", err,
		)
	}
}

func (c *wsClient) handleResponses(ctx context.Context) error {
	defer func() {
		_ = level.Info(c.deps.Logger()).Log("message", "handleResponses shutdown complete")
	}()

	for {
		select {
		case <-ctx.Done(): // time to stop
			_ = level.Info(c.deps.Logger()).Log("message", "handleResponses shutting down")

			defer func() {
				_ = level.Debug(c.deps.Logger()).Log("message", "gracefully closing the socket")
				err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					_ = level.Error(c.deps.Logger()).Log("message", "Unable to write websocket close message", "error", err)
					return
				}
				_ = level.Info(c.deps.Logger()).Log("message", "close message sent")
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
			c.processResponse(resp)
		}
	}
}

func (c *wsClient) SendMessage(msg WSMessage) {
	logger := logging.WithContext(msg.Ctx, c.deps.Logger())
	_ = level.Debug(logger).Log(
		"message", "adding message to response queue",
		"ws_msg_type", msg.MessageType,
		"ws_msg_len", len(msg.MessageContents),
	)

	c.responses <- msg
}
