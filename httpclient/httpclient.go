package httpclient

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gsmcwhirter/go-util/v7/logging/level"
	"github.com/gsmcwhirter/go-util/v7/telemetry"

	"github.com/gsmcwhirter/discord-bot-lib/v18/logging"
)

// HTTPClient is the interface of an http client
type HTTPClient interface {
	SetHeaders(http.Header)
	Get(context.Context, string, *http.Header) (*http.Response, error)
	GetBody(context.Context, string, *http.Header) (*http.Response, []byte, error)
	Post(context.Context, string, *http.Header, io.Reader) (*http.Response, error)
	PostBody(context.Context, string, *http.Header, io.Reader) (*http.Response, []byte, error)
	Put(context.Context, string, *http.Header, io.Reader) (*http.Response, error)
	PutBody(context.Context, string, *http.Header, io.Reader) (*http.Response, []byte, error)
}

type dependencies interface {
	Logger() logging.Logger
	Census() *telemetry.Census
	HTTPDoer() Doer
}

type httpClient struct {
	deps    dependencies
	headers http.Header
}

// NewHTTPClient creates a new http client
func NewHTTPClient(deps dependencies) HTTPClient {
	return &httpClient{
		deps:    deps,
		headers: http.Header{},
	}
}

func addHeaders(to *http.Header, from http.Header) {
	for k, v := range from {
		to.Del(k)
		for _, v2 := range v {
			to.Add(k, v2)
		}
	}
}

func (c *httpClient) doRequest(ctx context.Context, logger logging.Logger, method, url string, headers *http.Header, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	addHeaders(&req.Header, c.headers)
	if headers != nil {
		addHeaders(&req.Header, *headers)
	}

	level.Debug(logger).Message("http request start",
		"method", method,
		"url", url,
		"headers", fmt.Sprintf("%+v", NonSensitiveHeaders(req.Header)),
	)
	start := time.Now()
	resp, err := c.deps.HTTPDoer().Do(req)
	if err != nil {
		return nil, err
	}

	level.Info(logger).Message("http request complete",
		"elapsed_ns", time.Since(start).Nanoseconds(),
		"status_code", resp.StatusCode,
	)

	return resp, nil
}

func (c *httpClient) SetHeaders(h http.Header) {
	addHeaders(&c.headers, h)
}

func (c *httpClient) Get(ctx context.Context, url string, headers *http.Header) (*http.Response, error) {
	ctx, span := c.deps.Census().StartSpan(ctx, "httpClient.Get")
	defer span.End()

	logger := logging.WithContext(ctx, c.deps.Logger())

	resp, err := c.doRequest(ctx, logger, "GET", url, headers, nil)
	if err != nil {
		return nil, err
	}

	if resp.Body != nil {
		_ = resp.Body.Close()
	}

	return resp, nil
}

func (c *httpClient) GetBody(ctx context.Context, url string, headers *http.Header) (*http.Response, []byte, error) {
	ctx, span := c.deps.Census().StartSpan(ctx, "httpClient.GetBody")
	defer span.End()

	logger := logging.WithContext(ctx, c.deps.Logger())

	resp, err := c.doRequest(ctx, logger, "GET", url, headers, nil)
	if err != nil {
		return nil, nil, err
	}

	if resp.Body != nil {
		defer resp.Body.Close() //nolint:errcheck // not a real issue here
	}

	body, err := ioutil.ReadAll(resp.Body)

	return resp, body, err
}

func (c *httpClient) Post(ctx context.Context, url string, headers *http.Header, body io.Reader) (*http.Response, error) {
	ctx, span := c.deps.Census().StartSpan(ctx, "httpClient.Post")
	defer span.End()

	logger := logging.WithContext(ctx, c.deps.Logger())

	resp, err := c.doRequest(ctx, logger, "POST", url, headers, body)
	if err != nil {
		return nil, err
	}

	if resp.Body != nil {
		_ = resp.Body.Close()
	}

	return resp, nil
}

func (c *httpClient) PostBody(ctx context.Context, url string, headers *http.Header, body io.Reader) (*http.Response, []byte, error) {
	ctx, span := c.deps.Census().StartSpan(ctx, "httpClient.PostBody")
	defer span.End()

	logger := logging.WithContext(ctx, c.deps.Logger())

	resp, err := c.doRequest(ctx, logger, "POST", url, headers, body)
	if err != nil {
		return nil, nil, err
	}

	if resp.Body != nil {
		defer resp.Body.Close() //nolint:errcheck // not a real issue here
	}

	respBody, err := ioutil.ReadAll(resp.Body)

	return resp, respBody, err
}

func (c *httpClient) Put(ctx context.Context, url string, headers *http.Header, body io.Reader) (*http.Response, error) {
	ctx, span := c.deps.Census().StartSpan(ctx, "httpClient.Put")
	defer span.End()

	logger := logging.WithContext(ctx, c.deps.Logger())

	resp, err := c.doRequest(ctx, logger, "PUT", url, headers, body)
	if err != nil {
		return nil, err
	}

	if resp.Body != nil {
		_ = resp.Body.Close()
	}

	return resp, nil
}

func (c *httpClient) PutBody(ctx context.Context, url string, headers *http.Header, body io.Reader) (*http.Response, []byte, error) {
	ctx, span := c.deps.Census().StartSpan(ctx, "httpClient.PutBody")
	defer span.End()

	logger := logging.WithContext(ctx, c.deps.Logger())

	resp, err := c.doRequest(ctx, logger, "PUT", url, headers, body)
	if err != nil {
		return nil, nil, err
	}

	if resp.Body != nil {
		defer resp.Body.Close() //nolint:errcheck // not a real issue here
	}

	respBody, err := ioutil.ReadAll(resp.Body)

	return resp, respBody, err
}
