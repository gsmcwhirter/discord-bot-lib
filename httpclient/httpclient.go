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

	"github.com/gsmcwhirter/discord-bot-lib/logging"
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
}

type httpClient struct {
	deps    dependencies
	client  *http.Client
	headers http.Header
}

// NewHTTPClient creates a new http client
func NewHTTPClient(deps dependencies) HTTPClient {
	return &httpClient{
		deps:    deps,
		client:  &http.Client{},
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

func (c *httpClient) Get(ctx context.Context, url string, headers *http.Header) (resp *http.Response, err error) {
	logger := logging.WithContext(ctx, c.deps.Logger())

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
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
	resp, err = c.client.Do(req)
	_ = level.Debug(logger).Log(
		"message", "http get complete",
		"duration_ns", time.Since(start).Nanoseconds(),
		"status_code", resp.StatusCode,
	)
	if err != nil {
		return
	}
	defer resp.Body.Close() //nolint: errcheck

	return
}

func (c *httpClient) GetBody(ctx context.Context, url string, headers *http.Header) (resp *http.Response, body []byte, err error) {
	logger := logging.WithContext(ctx, c.deps.Logger())

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
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
	resp, err = c.client.Do(req)
	_ = level.Debug(logger).Log(
		"message", "http get complete",
		"duration_ns", time.Since(start).Nanoseconds(),
		"status_code", resp.StatusCode,
	)
	if err != nil {
		return
	}
	defer resp.Body.Close() //nolint: errcheck
	body, err = ioutil.ReadAll(resp.Body)

	return
}

func (c *httpClient) Post(ctx context.Context, url string, headers *http.Header, body io.Reader) (resp *http.Response, err error) {
	logger := logging.WithContext(ctx, c.deps.Logger())

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return
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
	resp, err = c.client.Do(req)
	_ = level.Debug(logger).Log(
		"message", "http post complete",
		"duration_ns", time.Since(start).Nanoseconds(),
		"status_code", resp.StatusCode,
	)
	if err != nil {
		return
	}
	defer resp.Body.Close() //nolint: errcheck

	return
}

func (c *httpClient) PostBody(ctx context.Context, url string, headers *http.Header, body io.Reader) (resp *http.Response, respBody []byte, err error) {
	logger := logging.WithContext(ctx, c.deps.Logger())

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return
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
	resp, err = c.client.Do(req)
	_ = level.Debug(logger).Log(
		"message", "http post complete",
		"duration_ns", time.Since(start).Nanoseconds(),
		"status_code", resp.StatusCode,
	)
	if err != nil {
		return
	}
	defer resp.Body.Close() //nolint: errcheck
	respBody, err = ioutil.ReadAll(resp.Body)

	return
}
