package metric_engine

import (
	"github.com/example/otel-tail-sampling-gateway/internal/quota"
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
func pointKey(p telemetry_domain.MetricPoint, window uint64) string {
	return p.TenantID + "/" + p.Name + "/" + p.Timestamp.UTC().Format(time.RFC3339Nano)
}
func (e *Engine) Add(p telemetry_domain.MetricPoint) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	k := pointKey(p, 0)
	if _, ok := e.points[k]; ok {
		e.duplicates++
		return false
	}
	e.points[k] = p
	return true
}

func (e *Engine) AddWithinQuota(p telemetry_domain.MetricPoint, limiter *quota.Limiter, units uint64) bool {
	decision := limiter.AllowWithWindow(p.TenantID, units)
	if !decision.Allowed {
		return false
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	k := pointKey(p, decision.Window)
	if _, ok := e.points[k]; ok {
		e.duplicates++
		return false
	}
	e.points[k] = p
	return true
}

func (e *Engine) Len() int           { e.mu.Lock(); defer e.mu.Unlock(); return len(e.points) }
func (e *Engine) Duplicates() uint64 { e.mu.Lock(); defer e.mu.Unlock(); return e.duplicates }
