package stats

import (
	"context"
	"math"
	"sync"
	"time"
)

type ActivityRecorder struct {
	sync.RWMutex
	avg      float64
	incr     float64
	lastTick time.Time

	ticker    *time.Ticker
	timescale float64
}

func movingExpAvg(value, oldValue, deltaT, timescale float64) float64 {
	alpha := 1.0 - math.Exp(-deltaT/timescale)
	r := alpha*value + (1.0-alpha)*oldValue
	return r
}

func NewActivityRecorder(timescale float64) *ActivityRecorder {
	return &ActivityRecorder{
		RWMutex:   sync.RWMutex{},
		avg:       0,
		incr:      0,
		timescale: timescale,
		ticker:    time.NewTicker(1 * time.Second),
		lastTick:  time.Now(),
	}
}

func (r *ActivityRecorder) tick() {
	r.Lock()
	defer r.Unlock()

	now := time.Now()
	deltaT := now.Sub(r.lastTick)
	r.lastTick = now
	r.avg = movingExpAvg(r.incr, r.avg, deltaT.Seconds(), r.timescale)
	r.incr = 0
}

func (r *ActivityRecorder) Run(ctx context.Context) error {
	for {
		select {
		case <-r.ticker.C:
			r.tick()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (r *ActivityRecorder) Incr(v float64) {
	r.Lock()
	defer r.Unlock()
	r.incr += v
}

func (r *ActivityRecorder) Timescale() float64 {
	return r.timescale
}

func (r *ActivityRecorder) Avg() float64 {
	r.RLock()
	defer r.RUnlock()

	return r.avg
}
