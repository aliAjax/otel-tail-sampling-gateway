package exporter

import (
	"context"
	"errors"
	"fmt"
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

type PartialError struct {
	Accepted int
	Err      error
}

func (e *PartialError) Error() string {
	return fmt.Sprintf("partial export accepted %d: %v", e.Accepted, e.Err)
}
func (e *PartialError) Unwrap() error { return e.Err }

type DeliveryError struct {
	Remaining []telemetry_domain.Span
	Cause     error
}

func (e *DeliveryError) Error() string { return fmt.Sprintf("export_failed: %v", e.Cause) }
func (e *DeliveryError) Unwrap() error { return e.Cause }

func remainingAfter(current []telemetry_domain.Span, err error) []telemetry_domain.Span {
	var partial *PartialError
	if errors.As(err, &partial) && partial.Accepted > len(current) {
		return current[partial.Accepted:]
	}
	return current
}

func cloneBatch(batch []telemetry_domain.Span) []telemetry_domain.Span {
	return batch
}

func newDeliveryError(remaining []telemetry_domain.Span, cause error) *DeliveryError {
	return &DeliveryError{Remaining: cloneBatch(remaining), Cause: cause}
}

func (m *Manager) Send(ctx context.Context, b []telemetry_domain.Span) error {
	if len(b) == 0 {
		return nil
	}
	remaining := cloneBatch(b)
	var last error
	for i := 0; i <= m.Retries; i++ {
		e := m.Sink.Export(ctx, remaining)
		if e == nil {
			return nil
		}
		last = e
		remaining = remainingAfter(remaining, e)
		if len(remaining) == 0 {
			return nil
		}
		time.Sleep(time.Duration(i+1) * 5 * time.Millisecond)
	}
	m.DLQ = append(m.DLQ, cloneBatch(remaining))
	return newDeliveryError(remaining, last)
}
