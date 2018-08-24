package wsclient

import (
	"net/http"
	"time"
)

// Conn is the interface for a *websocket.Conn as used in this wrapper package
type Conn interface {
	Close() error
	SetReadDeadline(time.Time) error
	ReadMessage() (int, []byte, error)
	WriteMessage(int, []byte) error
}

// Dialer is the interface for a *websocket.Dialer as used in this wrapper package
type Dialer interface {
	Dial(string, http.Header) (Conn, *http.Response, error)
}
