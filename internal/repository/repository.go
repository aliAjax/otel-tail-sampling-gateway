package repository

import (
	"github.com/example/otel-tail-sampling-gateway/internal/telemetry_domain"
	"sync"
	"time"
)

type Repo struct {
	mu        sync.RWMutex
	Tenants   map[string]telemetry_domain.Tenant
	Policies  map[string]telemetry_domain.SamplingPolicy
	Attr      map[string]telemetry_domain.AttributePolicy
	Exporters map[string]telemetry_domain.ExporterConfig
	Decisions map[string]telemetry_domain.SamplingDecision
}

func New() *Repo {
	return &Repo{Tenants: map[string]telemetry_domain.Tenant{}, Policies: map[string]telemetry_domain.SamplingPolicy{}, Attr: map[string]telemetry_domain.AttributePolicy{}, Exporters: map[string]telemetry_domain.ExporterConfig{}, Decisions: map[string]telemetry_domain.SamplingDecision{}}
}
func (r *Repo) PutTenant(t telemetry_domain.Tenant) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now()
	}
	r.Tenants[t.ID] = t
}
func (r *Repo) ListTenants() []telemetry_domain.Tenant {
	r.mu.RLock()
	defer r.mu.RUnlock()
	o := make([]telemetry_domain.Tenant, 0, len(r.Tenants))
	for _, t := range r.Tenants {
		o = append(o, t)
	}
	return o
}
func (r *Repo) PutPolicy(p telemetry_domain.SamplingPolicy) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Policies[p.ID] = p
}
func (r *Repo) GetPolicy(tenant string) telemetry_domain.SamplingPolicy {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.Policies {
		if p.TenantID == tenant {
			return p
		}
	}
	return telemetry_domain.SamplingPolicy{TenantID: tenant, Enabled: true, Probability: 1}
}
