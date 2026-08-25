package sampling_engine

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"

	"github.com/example/otel-tail-sampling-gateway/internal/telemetry_domain"
)

func batchTraces() []telemetry_domain.Trace {
	return []telemetry_domain.Trace{
		{ID: "trace-first", TenantID: "tenant"},
		{ID: "trace-second", TenantID: "tenant"},
		{ID: "trace-third", TenantID: "tenant"},
	}
}

func TestDecideBatchHonorsPreCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	got, err := (Engine{}).DecideBatch(ctx, batchTraces(), telemetry_domain.SamplingPolicy{Enabled: true, Probability: 1})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected canceled context, got %v", err)
	}
	if got != nil {
		t.Fatalf("canceled batch published %d decisions", len(got))
	}
}

func TestDecideBatchPreservesInputOrder(t *testing.T) {
	got, err := (Engine{}).DecideBatch(context.Background(), batchTraces(), telemetry_domain.SamplingPolicy{Enabled: true, Probability: 1})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"trace-first", "trace-second", "trace-third"}
	for i := range want {
		if got[i].TraceID != want[i] {
			t.Fatalf("decision %d has trace %q, want %q", i, got[i].TraceID, want[i])
		}
	}
}

type cancelOnFirstRule struct {
	cancel context.CancelFunc
	calls  *atomic.Int32
}

func (r cancelOnFirstRule) Match(telemetry_domain.Trace) bool {
	if r.calls.Add(1) == 1 {
		r.cancel()
	}
	return true
}

func (cancelOnFirstRule) Name() string { return "cancel-on-first" }

func TestMatchBatchStopsAfterCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var calls atomic.Int32
	c := Composite{Rules: []Rule{cancelOnFirstRule{cancel: cancel, calls: &calls}}}
	_, err := c.MatchBatch(ctx, batchTraces())
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected canceled context, got %v", err)
	}
	if got := calls.Load(); got != 1 {
		t.Fatalf("batch evaluated %d traces after cancellation, want 1", got)
	}
}

func TestMatchBatchDoesNotPublishPartialResults(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var calls atomic.Int32
	c := Composite{Rules: []Rule{cancelOnFirstRule{cancel: cancel, calls: &calls}}}
	got, err := c.MatchBatch(ctx, batchTraces())
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected canceled context, got %v", err)
	}
	if got != nil {
		t.Fatalf("canceled batch published partial matches: %v", got)
	}
}
