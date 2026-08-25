package metric_engine

import (
	"sync"
	"testing"
	"time"

	"github.com/example/otel-tail-sampling-gateway/internal/telemetry_domain"
)

func TestMetricPointOwnsAttributes(t *testing.T) {
	for round := 0; round < 10; round++ {
		attrs := telemetry_domain.Attributes{"route": "ingest", "region": "east"}
		e := New()
		point := telemetry_domain.MetricPoint{TenantID: "tenant-a", Name: "spans", Timestamp: time.Unix(int64(round), 0), Attributes: attrs}
		if !e.Add(point) {
			t.Fatalf("round %d: first metric was rejected", round)
		}
		attrs["route"] = "mutated"
		delete(attrs, "region")
		stored := e.Snapshot()[0]
		if stored.Attributes["route"] != "ingest" || stored.Attributes["region"] != "east" {
			t.Fatalf("round %d: stored metric changed after caller mutation: %#v", round, stored.Attributes)
		}
	}
}

func TestMetricSnapshotOwnsAttributes(t *testing.T) {
	e := New()
	point := telemetry_domain.MetricPoint{TenantID: "tenant-b", Name: "latency", Timestamp: time.Unix(99, 0), Attributes: telemetry_domain.Attributes{"unit": "ms"}}
	if !e.Add(point) {
		t.Fatal("first metric was rejected")
	}
	first := e.Snapshot()
	first[0].Attributes["unit"] = "seconds"
	second := e.Snapshot()
	if second[0].Attributes["unit"] != "ms" {
		t.Fatalf("snapshot mutation escaped into engine: %#v", second[0].Attributes)
	}
	attrs := point.Attributes
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < 100; i++ {
			attrs["unit"] = "ms"
		}
	}()
	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < 100; i++ {
			_ = e.Snapshot()[0].Attributes["unit"]
		}
	}()
	close(start)
	wg.Wait()
}
