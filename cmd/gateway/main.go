package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/example/otel-tail-sampling-gateway/internal/audit"
	"github.com/example/otel-tail-sampling-gateway/internal/buffer_store"
	"github.com/example/otel-tail-sampling-gateway/internal/exporter"
	"github.com/example/otel-tail-sampling-gateway/internal/metric_engine"
	"github.com/example/otel-tail-sampling-gateway/internal/observability"
	"github.com/example/otel-tail-sampling-gateway/internal/otlp_codec"
	"github.com/example/otel-tail-sampling-gateway/internal/quota"
	"github.com/example/otel-tail-sampling-gateway/internal/redaction"
	"github.com/example/otel-tail-sampling-gateway/internal/repository"
	"github.com/example/otel-tail-sampling-gateway/internal/sampling_engine"
	"github.com/example/otel-tail-sampling-gateway/internal/telemetry_domain"
	"github.com/example/otel-tail-sampling-gateway/internal/trace_aggregator"
	"github.com/example/otel-tail-sampling-gateway/internal/transport"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type app struct {
	repo    *repository.Repo
	agg     *trace_aggregator.Aggregator
	buf     *buffer_store.Store
	metrics *metric_engine.Engine
	quota   *quota.Limiter
	sink    *exporter.Simulator
	exp     *exporter.Manager
	audit   *audit.Log
	stats   observability.Metrics
}

func newApp() *app {
	r := repository.New()
	r.PutTenant(telemetry_domain.Tenant{ID: "demo", Name: "Demo", Quota: 10000})
	r.PutPolicy(telemetry_domain.SamplingPolicy{ID: "default", TenantID: "demo", Enabled: true, Probability: 1})
	s := &exporter.Simulator{}
	return &app{repo: r, agg: trace_aggregator.New(10000, 30*time.Second), buf: buffer_store.New("./data"), metrics: metric_engine.New(), quota: quota.New(10000), sink: s, exp: &exporter.Manager{Sink: s, Batch: 100, Retries: 2}, audit: &audit.Log{}}
}
func reqID(r *http.Request) string {
	if x := r.Header.Get("X-Request-ID"); x != "" {
		return x
	}
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}
func (a *app) write(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
func (a *app) handler(w http.ResponseWriter, r *http.Request) {
	id := reqID(r)
	ctx, cancel := context.WithTimeout(r.Context(), transport.Deadline(r))
	defer cancel()
	r = r.WithContext(ctx)
	w.Header().Set("X-Request-ID", id)
	if r.Method == "GET" && (r.URL.Path == "/healthz" || r.URL.Path == "/readyz") {
		a.write(w, 200, map[string]any{"status": "ok", "request_id": id})
		return
	}
	if r.URL.Path == "/v1/otlp/traces" && r.Method == "POST" {
		a.ingestTrace(w, r, id)
		return
	}
	if r.URL.Path == "/v1/otlp/metrics" && r.Method == "POST" {
		a.ingestMetric(w, r, id)
		return
	}
	if r.URL.Path == "/v1/otlp/logs" && r.Method == "POST" {
		a.ingestLog(w, r, id)
		return
	}
	if r.URL.Path == "/api/v1/tenants" {
		if r.Method == "GET" {
			a.write(w, 200, a.repo.ListTenants())
			return
		}
		if r.Method == "POST" {
			var t telemetry_domain.Tenant
			if e := otlp_codec.Decode(r, 4<<20, &t); e != nil {
				a.write(w, 400, map[string]any{"error": "invalid_request", "detail": e.Error()})
				return
			}
			if t.ID == "" {
				a.write(w, 400, map[string]string{"error": "tenant_id_required"})
				return
			}
			a.repo.PutTenant(t)
			a.write(w, 201, t)
			return
		}
	}
	if strings.HasPrefix(r.URL.Path, "/api/v1/traces/") && strings.HasSuffix(r.URL.Path, "/decision") {
		traceID := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/v1/traces/"), "/decision")
		if d, ok := a.repo.Decisions[traceID]; ok {
			a.write(w, 200, d)
		} else {
			a.write(w, 404, map[string]string{"error": "not_found"})
		}
		return
	}
	if r.URL.Path == "/api/v1/buffers" {
		a.write(w, 200, map[string]any{"queued": a.buf.Len(), "drops": a.buf.Drops()})
		return
	}
	if r.URL.Path == "/api/v1/drop-reasons" {
		a.write(w, 200, a.buf.Drops())
		return
	}
	if r.URL.Path == "/api/v1/exporter-health" {
		a.write(w, 200, map[string]any{"status": "healthy", "exported": a.sink.Exported})
		return
	}
	if r.URL.Path == "/api/v1/sampling-policies" && r.Method == "POST" {
		var p telemetry_domain.SamplingPolicy
		if e := otlp_codec.Decode(r, 4<<20, &p); e != nil {
			a.write(w, 400, map[string]string{"error": e.Error()})
			return
		}
		if p.ID == "" {
			p.ID = "policy-" + id
		}
		a.repo.PutPolicy(p)
		a.write(w, 201, p)
		return
	}
	if r.URL.Path == "/api/v1/policies/validate" {
		a.write(w, 200, map[string]any{"valid": true, "explain": sampling_engine.Explain(a.repo.GetPolicy("demo"))})
		return
	}
	if r.URL.Path == "/api/v1/buffers/drain" && r.Method == "POST" {
		n := a.buf.Len()
		b := a.buf.Drain(n)
		e := a.exp.Send(r.Context(), b)
		if e != nil {
			a.write(w, 503, map[string]string{"error": e.Error()})
			return
		}
		a.write(w, 200, map[string]any{"drained": len(b)})
		return
	}
	a.write(w, 404, map[string]string{"error": "not_found", "request_id": id})
}
func (a *app) ingestTrace(w http.ResponseWriter, r *http.Request, id string) {
	var s telemetry_domain.Span
	if e := otlp_codec.Decode(r, 4<<20, &s); e != nil {
		a.write(w, 400, map[string]string{"error": "invalid_otlp", "detail": e.Error()})
		return
	}
	s.TenantID = r.Header.Get("X-Tenant-ID")
	if s.TenantID == "" {
		s.TenantID = "demo"
	}
	if !a.quota.Allow(s.TenantID, 1) {
		a.buf.Drop("quota")
		a.stats.Dropped.Add(1)
		a.write(w, 429, map[string]string{"error": "quota_exceeded"})
		return
	}
	s = redaction.Apply(s, telemetry_domain.AttributePolicy{Salt: "gateway", Hash: []string{"email", "token"}})
	t, ok := a.agg.Add(s)
	if !ok {
		a.buf.Drop("memory_budget")
		a.write(w, 202, map[string]any{"accepted": false, "request_id": id})
		return
	}
	a.stats.Ingested.Add(1)
	p := a.repo.GetPolicy(s.TenantID)
	d := (sampling_engine.Engine{}).Decide(t, p)
	if d.Keep {
		a.repo.Decisions[t.ID] = d
		a.buf.Append(s)
		a.stats.Kept.Add(1)
	} else {
		a.stats.Dropped.Add(1)
		a.buf.Drop(d.Reason)
		a.repo.Decisions[t.ID] = d
	}
	a.audit.Add(audit.Event{RequestID: id, Tenant: s.TenantID, Action: "ingest_trace", Result: d.Reason})
	a.write(w, 202, map[string]any{"accepted": true, "trace_id": s.TraceID, "decision": a.repo.Decisions[t.ID], "request_id": id})
}
func (a *app) ingestMetric(w http.ResponseWriter, r *http.Request, id string) {
	var p telemetry_domain.MetricPoint
	if e := otlp_codec.Decode(r, 4<<20, &p); e != nil {
		a.write(w, 400, map[string]string{"error": e.Error()})
		return
	}
	p.TenantID = r.Header.Get("X-Tenant-ID")
	if a.metrics.Add(p) {
		a.write(w, 202, map[string]any{"accepted": true, "request_id": id})
	} else {
		a.write(w, 202, map[string]any{"accepted": false, "reason": "duplicate", "request_id": id})
	}
}
func (a *app) ingestLog(w http.ResponseWriter, r *http.Request, id string) {
	var l telemetry_domain.LogRecord
	if e := otlp_codec.Decode(r, 4<<20, &l); e != nil {
		a.write(w, 400, map[string]string{"error": e.Error()})
		return
	}
	l.Body = redaction.SanitizeLog(l.Body)
	a.write(w, 202, map[string]any{"accepted": true, "trace_id": l.TraceID, "request_id": id})
}
func main() {
	a := newApp()
	srv := &http.Server{Addr: env("HTTP_ADDR", ":8080"), Handler: http.HandlerFunc(a.handler), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second}
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			for _, t := range a.agg.FlushExpired(time.Now()) {
				if d, ok := a.repo.Decisions[t.ID]; ok && d.Keep {
					var b []telemetry_domain.Span
					b = append(b, t.Spans...)
					a.exp.Send(context.Background(), b)
				}
			}
		}
	}()
	go srv.ListenAndServe()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

var _ = fmt.Sprintf
