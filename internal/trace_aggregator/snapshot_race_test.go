package trace_aggregator

import (
	"sync"
	"testing"
	"time"

	"github.com/example/otel-tail-sampling-gateway/internal/telemetry_domain"
)

func runParallel(left, right func()) {
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); <-start; left() }()
	go func() { defer wg.Done(); <-start; right() }()
	close(start)
	wg.Wait()
}

func testSpan() telemetry_domain.Span {
	return telemetry_domain.Span{TraceID: "trace", SpanID: "span", Attributes: telemetry_domain.Attributes{"region": "east"}}
}

func TestAddSnapshotDoesNotMutateStoredTrace(t *testing.T) {
	a := New(10, time.Minute)
	got, _ := a.Add(testSpan())
	runParallel(func() { got.Spans[0].Attributes["region"] = "west" }, func() { a.Add(testSpan()) })
	stored, _ := a.Get("trace")
	if stored.Spans[0].Attributes["region"] != "east" {
		t.Fatal("Add snapshot mutated store")
	}
}

func TestGetSnapshotDoesNotMutateStoredTrace(t *testing.T) {
	a := New(10, time.Minute)
	a.Add(testSpan())
	got, _ := a.Get("trace")
	runParallel(func() { got.Spans[0].SpanID = "changed" }, func() { a.Add(testSpan()) })
	stored, _ := a.Get("trace")
	if stored.Spans[0].SpanID != "span" {
		t.Fatal("Get snapshot mutated store")
	}
}

func TestFlushSnapshotOwnsNestedAttributes(t *testing.T) {
	a := New(10, 0)
	attrs := telemetry_domain.Attributes{"region": "east"}
	a.Add(telemetry_domain.Span{TraceID: "trace", SpanID: "one", Attributes: attrs})
	a.Add(telemetry_domain.Span{TraceID: "trace", SpanID: "two", Attributes: attrs})
	flushed := a.FlushExpired(time.Now())
	runParallel(func() { flushed[0].Spans[0].Attributes["region"] = "changed" }, func() { flushed[0].Spans[1].Attributes["other"] = "value" })
	if flushed[0].Spans[1].Attributes["region"] != "east" {
		t.Fatal("flushed spans share attributes")
	}
}
