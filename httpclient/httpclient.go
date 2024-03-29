package httpclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gsmcwhirter/go-util/v10/errors"
	"github.com/gsmcwhirter/go-util/v10/json"
	"github.com/gsmcwhirter/go-util/v10/logging/level"
	"github.com/gsmcwhirter/go-util/v10/telemetry"

	"github.com/gsmcwhirter/discord-bot-lib/v24/logging"
)

type dependencies interface {
	Logger() Logger
	Telemetry() *telemetry.Telemeter
	HTTPDoer() Doer
}

// Logger is the interface expected for logging
type Logger = interface {
	Log(keyvals ...interface{}) error
	Message(string, ...interface{})
	Err(string, error, ...interface{})
	Printf(string, ...interface{})
}

// ErrResponse is the error that is wrapped and returned when there is a non-200 api response
var ErrResponse = errors.New("error response")

// HTTPClient is a wrapped http client for interacting with discord
type HTTPClient struct {
	deps    dependencies
	headers *http.Header

	debug bool
}

// NewHTTPClient creates a new http client
func NewHTTPClient(deps dependencies) *HTTPClient {
	return &HTTPClient{
		deps:    deps,
		headers: &http.Header{},
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

// SetDebug turns on/off debug mode
func (c *HTTPClient) SetDebug(val bool) {
	c.debug = val
}

func (c *HTTPClient) doRequest(ctx context.Context, logger Logger, method, url string, headers *http.Header, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	addHeaders(&req.Header, *c.headers)
	if headers != nil {
		addHeaders(&req.Header, *headers)
	}

	if c.debug {
		level.Debug(logger).Message("http request start",
			"method", method,
			"url", url,
			"headers", fmt.Sprintf("%+v", NonSensitiveHeaders(req.Header)),
		)
	}
	start := time.Now()
	resp, err := c.deps.HTTPDoer().Do(req)
	if err != nil {
		return nil, err
	}

	if c.debug || resp.StatusCode < 200 || resp.StatusCode >= 400 {
		level.Info(logger).Message("http request complete",
			"elapsed_ns", time.Since(start).Nanoseconds(),
			"status_code", resp.StatusCode,
		)
	}

	return resp, nil
}

// SetHeaders adds the provided headers to the set to be included in all requests
func (c *HTTPClient) SetHeaders(h http.Header) {
	addHeaders(c.headers, h)
}

// Get performs an http GET
func (c *HTTPClient) Get(ctx context.Context, url string, headers *http.Header) (*http.Response, error) {
	ctx, span := c.deps.Telemetry().StartSpan(ctx, "httpclient", "Get")
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

// GetBody performs an http GET an returns the response body
func (c *HTTPClient) GetBody(ctx context.Context, url string, headers *http.Header) (*http.Response, []byte, error) {
	ctx, span := c.deps.Telemetry().StartSpan(ctx, "httpclient", "GetBody")
	defer span.End()

	logger := logging.WithContext(ctx, c.deps.Logger())

	resp, err := c.doRequest(ctx, logger, "GET", url, headers, nil)
	if err != nil {
		return nil, nil, err
	}

	if resp.Body != nil {
		defer resp.Body.Close() //nolint:errcheck // not a real issue here
	}

	body, err := io.ReadAll(resp.Body)

	return resp, body, err
}

// GetJSON performs an http GET and unmarshals the response body into the provided target
func (c *HTTPClient) GetJSON(ctx context.Context, url string, headers *http.Header, t interface{}) (*http.Response, error) {
	ctx, span := c.deps.Telemetry().StartSpan(ctx, "httpclient", "GetJSON")
	defer span.End()

	logger := logging.WithContext(ctx, c.deps.Logger())

	resp, err := c.doRequest(ctx, logger, "GET", url, headers, nil)
	if err != nil {
		return nil, err
	}

	if resp.Body != nil {
		defer resp.Body.Close() //nolint:errcheck // not a real issue here
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= 300 {
		var body []byte
		var err error
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			err = errors.WithDetails(ErrResponse, "read_error", err.Error())
		} else {
			err = ErrResponse
		}
		return resp, errors.Wrap(err, "non-200 response", "status_code", resp.StatusCode, "response_body", string(body))
	}

	err = json.UnmarshalFromReader(resp.Body, t)
	return resp, errors.Wrap(err, "could not unmarshal json")
}

// Post performs an http POST
func (c *HTTPClient) Post(ctx context.Context, url string, headers *http.Header, body io.Reader) (*http.Response, error) {
	ctx, span := c.deps.Telemetry().StartSpan(ctx, "httpclient", "Post")
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

// PostBody performs an http POST and returns the response body
func (c *HTTPClient) PostBody(ctx context.Context, url string, headers *http.Header, body io.Reader) (*http.Response, []byte, error) {
	ctx, span := c.deps.Telemetry().StartSpan(ctx, "httpclient", "PostBody")
	defer span.End()

	logger := logging.WithContext(ctx, c.deps.Logger())

	resp, err := c.doRequest(ctx, logger, "POST", url, headers, body)
	if err != nil {
		return nil, nil, err
	}

	if resp.Body != nil {
		defer resp.Body.Close() //nolint:errcheck // not a real issue here
	}

	respBody, err := io.ReadAll(resp.Body)

	return resp, respBody, err
}

// PostJSON performs an http POST and unmarshals the response body into the provided target
func (c *HTTPClient) PostJSON(ctx context.Context, url string, headers *http.Header, body io.Reader, t interface{}) (*http.Response, error) {
	ctx, span := c.deps.Telemetry().StartSpan(ctx, "httpclient", "PostJSON")
	defer span.End()

	logger := logging.WithContext(ctx, c.deps.Logger())

	if headers == nil {
		headers = &http.Header{}
	}

	if headers.Get("content-type") == "" {
		headers.Set("content-type", "application/json")
	}

	resp, err := c.doRequest(ctx, logger, "POST", url, headers, body)
	if err != nil {
		return nil, err
	}

	if resp.Body != nil {
		defer resp.Body.Close() //nolint:errcheck // not a real issue here
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= 300 {
		var body []byte
		var err error
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			err = errors.WithDetails(ErrResponse, "read_error", err.Error())
		} else {
			err = ErrResponse
		}
		return resp, errors.Wrap(err, "non-200 response", "status_code", resp.StatusCode, "response_body", string(body))
	}

	err = json.UnmarshalFromReader(resp.Body, t)
	return resp, errors.Wrap(err, "could not unmarshal json")
}

// Put performs an http PUT
func (c *HTTPClient) Put(ctx context.Context, url string, headers *http.Header, body io.Reader) (*http.Response, error) {
	ctx, span := c.deps.Telemetry().StartSpan(ctx, "httpclient", "Put")
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

// PutBody performs an http PUT and returns the response body
func (c *HTTPClient) PutBody(ctx context.Context, url string, headers *http.Header, body io.Reader) (*http.Response, []byte, error) {
	ctx, span := c.deps.Telemetry().StartSpan(ctx, "httpclient", "PutBody")
	defer span.End()

	logger := logging.WithContext(ctx, c.deps.Logger())

	resp, err := c.doRequest(ctx, logger, "PUT", url, headers, body)
	if err != nil {
		return nil, nil, err
	}

	if resp.Body != nil {
		defer resp.Body.Close() //nolint:errcheck // not a real issue here
	}

	respBody, err := io.ReadAll(resp.Body)

	return resp, respBody, err
}

// PutJSON performs an http PUT and unmarshals the response body into the provided target
func (c *HTTPClient) PutJSON(ctx context.Context, url string, headers *http.Header, body io.Reader, t interface{}) (*http.Response, error) {
	ctx, span := c.deps.Telemetry().StartSpan(ctx, "httpclient", "PutJSON")
	defer span.End()

	logger := logging.WithContext(ctx, c.deps.Logger())

	if headers == nil {
		headers = &http.Header{}
	}

	if headers.Get("content-type") == "" {
		headers.Set("content-type", "application/json")
	}

	resp, err := c.doRequest(ctx, logger, "PUT", url, headers, body)
	if err != nil {
		return nil, err
	}

	if resp.Body != nil {
		defer resp.Body.Close() //nolint:errcheck // not a real issue here
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= 300 {
		var body []byte
		var err error
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			err = errors.WithDetails(ErrResponse, "read_error", err.Error())
		} else {
			err = ErrResponse
		}
		return resp, errors.Wrap(err, "non-200 response", "status_code", resp.StatusCode, "response_body", string(body))
	}

	err = json.UnmarshalFromReader(resp.Body, t)
	return resp, errors.Wrap(err, "could not unmarshal json")
}
