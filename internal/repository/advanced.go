package repository

import (
	"errors"
	"github.com/example/otel-tail-sampling-gateway/internal/telemetry_domain"
	"sort"
	"strings"
	"time"
)

func (r *Repo) DeleteTenant(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.Tenants[id]; !ok {
		return errors.New("tenant_not_found")
	}
	delete(r.Tenants, id)
	return nil
}
func (r *Repo) GetTenant(id string) (telemetry_domain.Tenant, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.Tenants[id]
	return t, ok
}
func (r *Repo) PutAttributePolicy(p telemetry_domain.AttributePolicy) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Attr[p.ID] = p
}
func (r *Repo) GetAttributePolicy(tenant string) telemetry_domain.AttributePolicy {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.Attr {
		if p.TenantID == tenant {
			return p
		}
	}
	return telemetry_domain.AttributePolicy{TenantID: tenant}
}
func (r *Repo) ListPolicies(tenant string) []telemetry_domain.SamplingPolicy {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var o []telemetry_domain.SamplingPolicy
	for _, p := range r.Policies {
		if tenant == "" || p.TenantID == tenant {
			o = append(o, p)
		}
	}
	sort.Slice(o, func(i, j int) bool { return o[i].Priority > o[j].Priority })
	return o
}
func (r *Repo) PutExporter(e telemetry_domain.ExporterConfig) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Exporters[e.ID] = e
}
func (r *Repo) ListExporters(tenant string) []telemetry_domain.ExporterConfig {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var o []telemetry_domain.ExporterConfig
	for _, e := range r.Exporters {
		if tenant == "" || e.TenantID == tenant {
			o = append(o, e)
		}
	}
	return o
}
func (r *Repo) SaveDecision(d telemetry_domain.SamplingDecision) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Decisions[d.TraceID] = d
}
func (r *Repo) GetDecision(id string) (telemetry_domain.SamplingDecision, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.Decisions[id]
	return d, ok
}
func (r *Repo) PruneDecisions(before time.Time) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	n := 0
	for id, d := range r.Decisions {
		if d.At.Before(before) {
			delete(r.Decisions, id)
			n++
		}
	}
	return n
}
func NormalizeID(v string) string { return strings.ToLower(strings.TrimSpace(v)) }
