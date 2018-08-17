package util

import (
	"context"

	"github.com/rs/xid"
)

// ContextKey is a wrapper type for our keys attached to a context
type ContextKey string

// NewRequestContext creates a new request context from context.Background()
func NewRequestContext() context.Context {
	return NewRequestContextFrom(context.Background())
}

// NewRequestContextFrom creates a new request context from an existing context
// regenerating a request_id value
func NewRequestContextFrom(ctx context.Context) context.Context {
	return context.WithValue(ctx, ContextKey("request_id"), GenerateRequestID())
}

// GenerateRequestID generates a new random-enough request id for a request context
func GenerateRequestID() string {
	return xid.New().String()
}
