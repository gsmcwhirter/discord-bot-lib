package httpclient

import (
	"context"
	"sync"
	"time"
)

// RateLimit is a parameterized limit set that a RateLimiter uses
type RateLimit struct {
	Bucket       string
	Limit        int
	LastReported time.Time
	ResetAt      time.Time
}

// RateLimiter is a rate limiting feature that accepts outside input dynamically for limits per time
type RateLimiter struct {
	mu sync.Mutex

	buckets map[string]RateLimit
}

// Wait waits until there is sufficient limit room to continue
func (r *RateLimiter) Wait(ctx context.Context, b string) error {
	r.mu.Lock()

	bk, ok := r.buckets[b]
	if !ok {
		r.mu.Unlock()
		return nil
	}

	if bk.Limit > 0 {
		bk.Limit--
		r.buckets[b] = bk
		r.mu.Unlock()
		return nil
	}

	r.mu.Unlock()
	waitTime := time.Until(bk.ResetAt)
	select {
	case <-time.After(waitTime):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Report tunes the bucket's limit
func (r *RateLimiter) Report(b string, limit int, reset, at time.Time) {
	r.mu.Lock()
	defer r.mu.Unlock()

	bk, ok := r.buckets[b]
	if !ok {
		return
	}

	if at.Before(bk.LastReported) {
		return
	}

	bk.Limit = limit
	bk.ResetAt = reset
	bk.LastReported = at
	r.buckets[b] = bk
}
