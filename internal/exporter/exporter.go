package exporter

import (
	"context"
	"errors"
	"github.com/example/otel-tail-sampling-gateway/internal/telemetry_domain"
	"sync"
	"time"
)

type Sink interface {
	Export(context.Context, []telemetry_domain.Span) error
}
type Simulator struct {
	mu       sync.Mutex
	Exported uint64
	Fail     bool
}

func (s *Simulator) Export(ctx context.Context, b []telemetry_domain.Span) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Fail {
		return errors.New("adapter_unavailable")
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		s.Exported += uint64(len(b))
		return nil
	}
}

type Manager struct {
	Sink    Sink
	Batch   int
	Retries int
	DLQ     [][]telemetry_domain.Span
}

func (m *Manager) Send(ctx context.Context, b []telemetry_domain.Span) error {
	if len(b) == 0 {
		return nil
	}
	for i := 0; i <= m.Retries; i++ {
		e := m.Sink.Export(ctx, b)
		if e == nil {
			return nil
		}
		if err := sleep(ctx, time.Duration(i+1)*5*time.Millisecond); err != nil {
			return err
		}
	}
	m.DLQ = append(m.DLQ, b)
	return errors.New("export_failed")
}

func sleep(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
