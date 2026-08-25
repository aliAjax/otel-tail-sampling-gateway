package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/example/otel-tail-sampling-gateway/internal/buffer_store"
	"github.com/example/otel-tail-sampling-gateway/internal/exporter"
	"github.com/example/otel-tail-sampling-gateway/internal/telemetry_domain"
)

type cancelAwareSink struct {
	started chan struct{}
	once    sync.Once
}

func (s *cancelAwareSink) Export(ctx context.Context, _ []telemetry_domain.Span) error {
	s.once.Do(func() { close(s.started) })
	<-ctx.Done()
	return ctx.Err()
}

func TestDrainPreservesRequestCancellation(t *testing.T) {
	sink := &cancelAwareSink{started: make(chan struct{})}
	a1 := &app{buf: buffer_store.New(t.TempDir())}
	a2 := &app{buf: buffer_store.New(t.TempDir())}
	if err := a1.buf.Append(telemetry_domain.Span{TraceID: "trace-ctx-1", SpanID: "span-ctx-1"}); err != nil {
		t.Fatal(err)
	}
	if err := a2.buf.Append(telemetry_domain.Span{TraceID: "trace-ctx-2", SpanID: "span-ctx-2"}); err != nil {
		t.Fatal(err)
	}
	a1.exp = &exporter.Manager{Sink: sink, Retries: 0}
	a2.exp = &exporter.Manager{Sink: sink, Retries: 0}
	ctx1, cancel1 := context.WithCancel(context.Background())
	ctx2, cancel2 := context.WithCancel(context.Background())
	req1 := httptest.NewRequest(http.MethodPost, "/api/v1/buffers/drain", nil).WithContext(ctx1)
	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/buffers/drain", nil).WithContext(ctx2)
	req1.Header.Set("X-Timeout-Ms", "1000")
	req2.Header.Set("X-Timeout-Ms", "1000")
	resp1, resp2 := httptest.NewRecorder(), httptest.NewRecorder()
	done := make(chan struct{}, 2)
	start := make(chan struct{})
	go func() { <-start; a1.handler(resp1, req1); done <- struct{}{} }()
	go func() { <-start; a2.handler(resp2, req2); done <- struct{}{} }()
	close(start)
	select {
	case <-sink.started:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("export did not start")
	}
	cancel1()
	cancel2()
	select {
	case <-done:
		select {
		case <-done:
		case <-time.After(200 * time.Millisecond):
			t.Fatal("second drain did not stop after request cancellation")
		}
	case <-time.After(200 * time.Millisecond):
		panic("drain did not stop after request cancellation")
	}
	if resp1.Code != http.StatusServiceUnavailable || resp2.Code != http.StatusServiceUnavailable {
		t.Fatalf("statuses = %d/%d, want 503/503", resp1.Code, resp2.Code)
	}
}
