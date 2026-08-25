package exporter

import (
	"context"
	"github.com/example/otel-tail-sampling-gateway/internal/telemetry_domain"
	"sync"
	"time"
)

type Pipeline struct {
	mu                sync.Mutex
	manager           *Manager
	queue             []telemetry_domain.Span
	capacity          int
	closed            bool
	accepted, dropped uint64
}

func NewPipeline(m *Manager, capacity int) *Pipeline {
	if capacity < 1 {
		capacity = 1
	}
	return &Pipeline{manager: m, capacity: capacity}
}
func (p *Pipeline) Enqueue(s telemetry_domain.Span) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed || len(p.queue) >= p.capacity {
		p.dropped++
		return false
	}
	p.queue = append(p.queue, s)
	p.accepted++
	return true
}
func (p *Pipeline) Flush(ctx context.Context) error {
	p.mu.Lock()
	b := append([]telemetry_domain.Span(nil), p.queue...)
	p.queue = nil
	p.mu.Unlock()
	if len(b) == 0 {
		return nil
	}
	return p.manager.Send(ctx, b)
}
func (p *Pipeline) Run(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				p.Flush(context.Background())
				return
			case <-ticker.C:
				p.Flush(ctx)
			}
		}
	}()
}
func (p *Pipeline) Close() { p.mu.Lock(); p.closed = true; p.mu.Unlock() }
func (p *Pipeline) Stats() (uint64, uint64, int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.accepted, p.dropped, len(p.queue)
}
