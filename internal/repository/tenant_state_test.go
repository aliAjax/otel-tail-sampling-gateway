package repository

import (
	"testing"

	"github.com/example/otel-tail-sampling-gateway/internal/telemetry_domain"
)

func repoWithTenantState() *Repo {
	r := New()
	r.PutTenant(telemetry_domain.Tenant{ID: "gone"})
	r.PutTenant(telemetry_domain.Tenant{ID: "keep"})
	r.PutPolicy(telemetry_domain.SamplingPolicy{ID: "p-gone", TenantID: "gone"})
	r.PutPolicy(telemetry_domain.SamplingPolicy{ID: "p-keep", TenantID: "keep"})
	r.PutAttributePolicy(telemetry_domain.AttributePolicy{ID: "a-gone", TenantID: "gone"})
	r.PutAttributePolicy(telemetry_domain.AttributePolicy{ID: "a-keep", TenantID: "keep"})
	r.PutExporter(telemetry_domain.ExporterConfig{ID: "e-gone", TenantID: "gone"})
	r.PutExporter(telemetry_domain.ExporterConfig{ID: "e-keep", TenantID: "keep"})
	r.SaveDecision(telemetry_domain.SamplingDecision{TraceID: "d-gone", TenantID: "gone"})
	r.SaveDecision(telemetry_domain.SamplingDecision{TraceID: "d-keep", TenantID: "keep"})
	return r
}

func TestDeleteTenantRemovesPolicies(t *testing.T) {
	r := repoWithTenantState()
	if err := r.DeleteTenant("gone"); err != nil {
		t.Fatal(err)
	}
	if _, ok := r.Policies["p-gone"]; ok {
		t.Fatal("deleted tenant policy survived")
	}
	if _, ok := r.Policies["p-keep"]; !ok {
		t.Fatal("other tenant policy was removed")
	}
}

func TestDeleteTenantRemovesAttributePolicies(t *testing.T) {
	r := repoWithTenantState()
	if err := r.DeleteTenant("gone"); err != nil {
		t.Fatal(err)
	}
	if _, ok := r.Attr["a-gone"]; ok {
		t.Fatal("deleted tenant attribute policy survived")
	}
	if _, ok := r.Attr["a-keep"]; !ok {
		t.Fatal("other tenant attribute policy was removed")
	}
}

func TestDeleteTenantRemovesExporters(t *testing.T) {
	r := repoWithTenantState()
	if err := r.DeleteTenant("gone"); err != nil {
		t.Fatal(err)
	}
	if _, ok := r.Exporters["e-gone"]; ok {
		t.Fatal("deleted tenant exporter survived")
	}
	if _, ok := r.Exporters["e-keep"]; !ok {
		t.Fatal("other tenant exporter was removed")
	}
}

func TestDeleteTenantRemovesDecisions(t *testing.T) {
	r := repoWithTenantState()
	if err := r.DeleteTenant("gone"); err != nil {
		t.Fatal(err)
	}
	if _, ok := r.Decisions["d-gone"]; ok {
		t.Fatal("deleted tenant decision survived")
	}
	if _, ok := r.Decisions["d-keep"]; !ok {
		t.Fatal("other tenant decision was removed")
	}
}

func TestGetPolicyReturnsHighestPriority(t *testing.T) {
	r := New()
	r.PutPolicy(telemetry_domain.SamplingPolicy{ID: "low", TenantID: "t", Priority: 1})
	r.PutPolicy(telemetry_domain.SamplingPolicy{ID: "high", TenantID: "t", Priority: 50})
	r.PutPolicy(telemetry_domain.SamplingPolicy{ID: "middle", TenantID: "t", Priority: 10})
	for i := 0; i < 20; i++ {
		if got := r.GetPolicy("t").ID; got != "high" {
			t.Fatalf("unstable policy selection: %q", got)
		}
	}
}
