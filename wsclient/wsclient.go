package wsclient

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/gsmcwhirter/go-util/v3/errors"
	log "github.com/gsmcwhirter/go-util/v3/logging"
	"github.com/gsmcwhirter/go-util/v3/logging/level"
	"github.com/gsmcwhirter/go-util/v3/request"
	"golang.org/x/sync/errgroup"

	"github.com/gsmcwhirter/discord-bot-lib/v7/logging"
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
	WSDialer() Dialer
}

type wsClient struct {
	deps dependencies

	gatewayURL string
	conn       Conn
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
	MaxConcurrentHandlers int
}

// NewWSClient creates a new WSClient
func NewWSClient(deps dependencies, options Options) WSClient {
	c := &wsClient{
		deps:       deps,
		gatewayURL: options.GatewayURL,
		closeLock:  &sync.Mutex{},
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

func (c *wsClient) Connect(token string) error {
	var err error
	ctx := request.NewRequestContext()
	logger := logging.WithContext(ctx, c.deps.Logger())

	dialHeader := http.Header{
		"Authorization": []string{fmt.Sprintf("Bot %s", token)},
	}

	var dialResp *http.Response

	level.Debug(logger).Message("ws client dial start",
		"url", c.gatewayURL,
	)

	start := time.Now()
	c.conn, dialResp, err = c.deps.WSDialer().Dial(c.gatewayURL, dialHeader)

	level.Info(logger).Message("ws client dial complete",
		"duration_ns", time.Since(start).Nanoseconds(),
		"status_code", dialResp.StatusCode,
		"url", c.gatewayURL,
	)

	if err != nil {
		return err
	}
	if dialResp.Body != nil {
		defer dialResp.Body.Close() // nolint: errcheck
	}

	level.Info(logger).Message("ws connected")

	return nil
}

func (c *wsClient) Close() {
	c.pool.Wait()
	if c.conn != nil {
		_ = c.conn.Close()
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
		level.Info(c.deps.Logger()).Message("starting response handler")
		return c.handleResponses(ctx)
	})

	controls.Go(func() error {
		level.Info(c.deps.Logger()).Message("starting message reader")
		return c.readMessages(ctx)
	})

	level.Info(c.deps.Logger()).Message("connected and listening")

	err := controls.Wait()
	level.Info(c.deps.Logger()).Message("shutting down")
	return err
}

func (c *wsClient) handleMessageRead(ctx context.Context, msgType int, msg []byte) {
	defer c.pool.Done()

	select {
	case <-ctx.Done():
		return
	default:
	}

	reqCtx := request.NewRequestContextFrom(ctx)
	mT := MessageType(msgType)
	mC := make([]byte, len(msg))
	copy(mC, msg)

	logger := logging.WithContext(reqCtx, c.deps.Logger())

	wsMsg := WSMessage{Ctx: reqCtx, MessageType: mT, MessageContents: mC}
	level.Info(logger).Message("received message",
		"ws_msg_type", mT,
		"ws_msg_len", len(mC),
	)

	level.Debug(logger).Message("waiting for worker token")
	c.poolTokens <- struct{}{}
	level.Info(logger).Message("worker token acquired")

	c.handleRequest(wsMsg)
}

func (c *wsClient) doReads(ctx context.Context) error {
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

func (c *wsClient) readMessages(ctx context.Context) error {
	defer level.Info(c.deps.Logger()).Message("readMessages shutdown complete")

	reader, ctx := errgroup.WithContext(ctx)
	reader.Go(func() error {
		return c.doReads(ctx)
	})

	// watches for close message
	reader.Go(func() error {
		<-ctx.Done()
		level.Info(c.deps.Logger()).Message("readMessages shutting down")
		c.gracefulClose()
		return ctx.Err()
	})

	return reader.Wait()

}

func (c *wsClient) handleRequest(req WSMessage) {
	logger := logging.WithContext(req.Ctx, c.deps.Logger())

	defer func() {
		<-c.poolTokens
		level.Info(logger).Message("released worker token")
	}()

	select {
	case <-req.Ctx.Done():
		level.Info(logger).Message("handleRequest received interrupt -- shutting down")
		return
	default:
		level.Info(logger).Message("handleRequest dispatching request")
		c.handler.HandleRequest(req, c.responses)
	}
}

func (c *wsClient) processResponse(resp WSMessage) {
	logger := logging.WithContext(resp.Ctx, c.deps.Logger())

	level.Debug(logger).Message("starting sending message",
		"ws_msg_type", resp.MessageType,
		"ws_msg_len", len(resp.MessageContents),
	)

	start := time.Now()
	err := c.conn.WriteMessage(int(resp.MessageType), resp.MessageContents)

	level.Info(logger).Message("done sending message",
		"elapsed_ns", time.Since(start).Nanoseconds(),
		"ws_msg_type", resp.MessageType,
		"ws_msg_len", len(resp.MessageContents),
	)

	if err != nil {
		level.Error(logger).Err("error sending message", err)
	}
}

func (c *wsClient) handleResponses(ctx context.Context) error {
	defer func() {
		level.Info(c.deps.Logger()).Message("handleResponses shutdown complete")
	}()

	for {
		select {
		case <-ctx.Done(): // time to stop
			level.Info(c.deps.Logger()).Message("handleResponses shutting down")

			defer func() {
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
			c.processResponse(resp)
		}
	}
}

func (c *wsClient) SendMessage(msg WSMessage) {
	logger := logging.WithContext(msg.Ctx, c.deps.Logger())
	level.Debug(logger).Message("adding message to response queue",
		"ws_msg_type", msg.MessageType,
		"ws_msg_len", len(msg.MessageContents),
	)

	c.responses <- msg
}
