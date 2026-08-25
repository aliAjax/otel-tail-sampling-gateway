package quota

import (
	"sync"
	"time"
)

type Limiter struct {
	mu         sync.Mutex
	limit      uint64
	used       map[string]uint64
	reset      time.Time
	window     time.Duration
	now        func() time.Time
	generation uint64
}

func New(limit uint64) *Limiter {
	return NewWindow(limit, time.Minute, time.Now)
}

func NewWindow(limit uint64, window time.Duration, now func() time.Time) *Limiter {
	if window <= 0 {
		window = time.Minute
	}
	if now == nil {
		now = time.Now
	}
	start := now()
	return &Limiter{limit: limit, used: map[string]uint64{}, reset: start.Add(window), window: window, now: now, generation: 1}
}

type Decision struct {
	Allowed bool
	Window  uint64
}

func (l *Limiter) advance(now time.Time) {
	if !now.Before(l.reset) {
		l.used = map[string]uint64{}
		l.reset = now.Add(l.window)
	}
}

func (l *Limiter) AllowWithWindow(tenant string, n uint64) Decision {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.advance(l.now())
	if l.used[tenant]+n > l.limit {
		return Decision{Window: l.generation}
	}
	l.used[tenant] += n
	return Decision{Allowed: true, Window: l.generation}
}

func (l *Limiter) Allow(tenant string, n uint64) bool {
	return l.AllowWithWindow(tenant, n).Allowed
}

func (l *Limiter) Usage(tenant string) uint64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.used[tenant]
}

func (l *Limiter) Refund(window uint64, tenant string, n uint64) bool {
	return false
}
