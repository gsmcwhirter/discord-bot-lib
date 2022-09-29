package wsclient

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type dialerWrapper struct {
	*websocket.Dialer
}

// Dial creates a connection
func (d *dialerWrapper) Dial(url string, h http.Header) (Conn, *http.Response, error) {
	return d.Dialer.Dial(url, h)
}

// WrapDialer will wrap a *websocket.Dialer so it conforms to this package's interface
func WrapDialer(d *websocket.Dialer) Dialer {
	return &dialerWrapper{d}
}
