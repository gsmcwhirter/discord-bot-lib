package httpclient

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/gsmcwhirter/discord-bot-lib/v6/logging"
)

// HTTPClient is the interface of an http client
type HTTPClient interface {
	SetHeaders(http.Header)
	Get(context.Context, string, *http.Header) (*http.Response, error)
	GetBody(context.Context, string, *http.Header) (*http.Response, []byte, error)
	Post(context.Context, string, *http.Header, io.Reader) (*http.Response, error)
	PostBody(context.Context, string, *http.Header, io.Reader) (*http.Response, []byte, error)
}

type dependencies interface {
	Logger() log.Logger
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

func (c *httpClient) SetHeaders(h http.Header) {
	addHeaders(&c.headers, h)
}

func (c *httpClient) Get(ctx context.Context, url string, headers *http.Header) (*http.Response, error) {
	logger := logging.WithContext(ctx, c.deps.Logger())

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	addHeaders(&req.Header, c.headers)
	if headers != nil {
		addHeaders(&req.Header, *headers)
	}

	_ = level.Debug(logger).Log(
		"message", "http get start",
		"url", url,
		"headers", fmt.Sprintf("%+v", NonSensitiveHeaders(req.Header)),
	)
	start := time.Now()
	resp, err := c.deps.HTTPDoer().Do(req)
	_ = level.Info(logger).Log(
		"message", "http get complete",
		"duration_ns", time.Since(start).Nanoseconds(),
		"status_code", resp.StatusCode,
	)
	if err != nil || resp.Body == nil {
		return resp, err
	}

	_ = resp.Body.Close()

	return resp, nil
}

func (c *httpClient) GetBody(ctx context.Context, url string, headers *http.Header) (*http.Response, []byte, error) {
	logger := logging.WithContext(ctx, c.deps.Logger())

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	addHeaders(&req.Header, c.headers)
	if headers != nil {
		addHeaders(&req.Header, *headers)
	}

	_ = level.Debug(logger).Log(
		"message", "http get start",
		"url", url,
		"headers", fmt.Sprintf("%+v", NonSensitiveHeaders(req.Header)),
	)
	start := time.Now()
	resp, err := c.deps.HTTPDoer().Do(req)
	_ = level.Info(logger).Log(
		"message", "http get complete",
		"duration_ns", time.Since(start).Nanoseconds(),
		"status_code", resp.StatusCode,
	)
	if err != nil || resp.Body == nil {
		return resp, nil, err
	}
	defer resp.Body.Close() //nolint: errcheck
	body, err := ioutil.ReadAll(resp.Body)

	return resp, body, err
}

func (c *httpClient) Post(ctx context.Context, url string, headers *http.Header, body io.Reader) (*http.Response, error) {
	logger := logging.WithContext(ctx, c.deps.Logger())

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	addHeaders(&req.Header, c.headers)
	if headers != nil {
		addHeaders(&req.Header, *headers)
	}

	_ = level.Debug(logger).Log(
		"message", "http post start",
		"url", url,
		"headers", fmt.Sprintf("%+v", NonSensitiveHeaders(req.Header)),
	)

	start := time.Now()
	resp, err := c.deps.HTTPDoer().Do(req)
	_ = level.Info(logger).Log(
		"message", "http post complete",
		"duration_ns", time.Since(start).Nanoseconds(),
		"status_code", resp.StatusCode,
	)
	if err != nil || resp.Body == nil {
		return resp, err
	}
	_ = resp.Body.Close() //nolint: errcheck

	return resp, nil
}

func (c *httpClient) PostBody(ctx context.Context, url string, headers *http.Header, body io.Reader) (*http.Response, []byte, error) {
	logger := logging.WithContext(ctx, c.deps.Logger())

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	addHeaders(&req.Header, c.headers)
	if headers != nil {
		addHeaders(&req.Header, *headers)
	}

	_ = level.Debug(logger).Log(
		"message", "http post start",
		"url", url,
		"headers", fmt.Sprintf("%+v", NonSensitiveHeaders(req.Header)),
	)

	start := time.Now()
	resp, err := c.deps.HTTPDoer().Do(req)
	_ = level.Info(logger).Log(
		"message", "http post complete",
		"duration_ns", time.Since(start).Nanoseconds(),
		"status_code", resp.StatusCode,
	)
	if err != nil || resp.Body == nil {
		return resp, nil, err
	}

	defer resp.Body.Close() //nolint: errcheck
	respBody, err := ioutil.ReadAll(resp.Body)

	return resp, respBody, err
}
