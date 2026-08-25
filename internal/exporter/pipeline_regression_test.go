package exporter

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/example/otel-tail-sampling-gateway/internal/telemetry_domain"
)

type shutdownProbe struct {
	ctxSeen chan error
	err     error
}

func (s *shutdownProbe) Export(ctx context.Context, _ []telemetry_domain.Span) error {
	s.ctxSeen <- ctx.Err()
	return s.err
}

func TestRunShutdownPropagatesCancellation(t *testing.T) {
	probe := &shutdownProbe{ctxSeen: make(chan error, 1)}
	p := NewPipeline(&Manager{Sink: probe}, 2)
	p.Enqueue(telemetry_domain.Span{TraceID: "shutdown"})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	p.Run(ctx, time.Hour)
	select {
	case got := <-probe.ctxSeen:
		if !errors.Is(got, context.Canceled) {
			t.Fatalf("shutdown flush used context %v", got)
		}
	case <-time.After(time.Second):
		t.Fatal("shutdown flush did not run")
	}
}

func TestRunShutdownRestoresFailedBatch(t *testing.T) {
	probe := &shutdownProbe{ctxSeen: make(chan error, 1), err: errors.New("downstream")}
	p := NewPipeline(&Manager{Sink: probe}, 2)
	p.Enqueue(telemetry_domain.Span{TraceID: "retain"})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	p.Run(ctx, time.Hour)
	<-probe.ctxSeen
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		_, _, queued := p.Stats()
		if queued == 1 {
			return
		}
		time.Sleep(time.Millisecond)
	}
	_, _, queued := p.Stats()
	t.Fatalf("failed shutdown flush dropped %d queued spans", 1-queued)
}
