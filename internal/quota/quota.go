package quota

import (
	"sync"
	"time"
)

type Limiter struct {
	mu    sync.Mutex
	limit uint64
	used  map[string]uint64
	reset time.Time
}

func New(limit uint64) *Limiter {
	return &Limiter{limit: limit, used: map[string]uint64{}, reset: time.Now().Add(time.Minute)}
}
func (l *Limiter) Allow(tenant string, n uint64) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if time.Now().After(l.reset) {
		l.used = map[string]uint64{}
		l.reset = time.Now().Add(time.Minute)
	}
	if l.used[tenant]+n > l.limit {
		return false
	}
	l.used[tenant] += n
	return true
}
func (l *Limiter) Usage(tenant string) uint64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.used[tenant]
}
