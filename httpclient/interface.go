package httpclient

import "net/http"

// Doer is the interface for the underlying client that this package wraps around
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}
