package metric_engine

import (
	"github.com/example/otel-tail-sampling-gateway/internal/telemetry_domain"
	"sync"
	"time"
)

type Engine struct {
	mu         sync.Mutex
	points     map[string]telemetry_domain.MetricPoint
	duplicates uint64
}

func New() *Engine { return &Engine{points: map[string]telemetry_domain.MetricPoint{}} }

func cloneAttrs(a telemetry_domain.Attributes) telemetry_domain.Attributes {
	if a == nil {
		return nil
	}
	out := make(telemetry_domain.Attributes, len(a))
	for k, v := range a {
		out[k] = v
	}
	return out
}

func (e *Engine) Add(p telemetry_domain.MetricPoint) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	k := p.TenantID + "/" + p.Name + "/" + p.Timestamp.UTC().Format(time.RFC3339Nano)
	if _, ok := e.points[k]; ok {
		e.duplicates++
		return false
	}
	p.Attributes = cloneAttrs(p.Attributes)
	e.points[k] = p
	return true
}
func (e *Engine) Len() int           { e.mu.Lock(); defer e.mu.Unlock(); return len(e.points) }
func (e *Engine) Duplicates() uint64 { e.mu.Lock(); defer e.mu.Unlock(); return e.duplicates }
func (e *Engine) Snapshot() []telemetry_domain.MetricPoint {
	e.mu.Lock()
	defer e.mu.Unlock()
	out := make([]telemetry_domain.MetricPoint, 0, len(e.points))
	for _, point := range e.points {
		point.Attributes = cloneAttrs(point.Attributes)
		out = append(out, point)
	}
	return out
}
