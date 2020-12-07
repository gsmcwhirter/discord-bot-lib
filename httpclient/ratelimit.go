package httpclient

import (
	"context"
	"sync"
	"time"
)

type RateLimit struct {
	Bucket       string
	Limit        int
	LastReported time.Time
	ResetAt      time.Time
}

type RateLimiter struct {
	sync.Mutex

	buckets map[string]RateLimit
}

func (r *RateLimiter) Wait(ctx context.Context, b string) error {
	r.Lock()

	bk, ok := r.buckets[b]
	if !ok {
		r.Unlock()
		return nil
	}

	if bk.Limit > 0 {
		bk.Limit--
		r.buckets[b] = bk
		r.Unlock()
		return nil
	}

	r.Unlock()
	waitTime := time.Until(bk.ResetAt)
	select {
	case <-time.After(waitTime):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (r *RateLimiter) Report(b string, limit int, reset, at time.Time) {
	r.Lock()
	defer r.Unlock()

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
