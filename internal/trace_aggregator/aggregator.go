package trace_aggregator

import (
	"github.com/example/otel-tail-sampling-gateway/internal/telemetry_domain"
	"sync"
	"time"
)

type Aggregator struct {
	mu     sync.Mutex
	traces map[string]*telemetry_domain.Trace
	max    int
	ttl    time.Duration
}

func New(max int, ttl time.Duration) *Aggregator {
	return &Aggregator{traces: map[string]*telemetry_domain.Trace{}, max: max, ttl: ttl}
}
func (a *Aggregator) Add(s telemetry_domain.Span) (telemetry_domain.Trace, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	t := a.traces[s.TraceID]
	now := time.Now()
	if t == nil {
		if len(a.traces) >= a.max {
			return telemetry_domain.Trace{}, false
		}
		t = &telemetry_domain.Trace{ID: s.TraceID, TenantID: s.TenantID, FirstSeen: now, State: "open"}
		a.traces[s.TraceID] = t
	}
	t.Spans = append(t.Spans, s)
	t.LastSeen = now
	return *t, true
}
func (a *Aggregator) FlushExpired(now time.Time) []telemetry_domain.Trace {
	a.mu.Lock()
	defer a.mu.Unlock()
	var out []telemetry_domain.Trace
	for id, t := range a.traces {
		if now.Sub(t.LastSeen) >= a.ttl {
			t.State = "flushed"
			out = append(out, *t)
			delete(a.traces, id)
		}
	}
	return out
}
func (a *Aggregator) Get(id string) (telemetry_domain.Trace, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	t, ok := a.traces[id]
	if !ok {
		return telemetry_domain.Trace{}, false
	}
	return *t, true
}
func (a *Aggregator) Len() int { a.mu.Lock(); defer a.mu.Unlock(); return len(a.traces) }
