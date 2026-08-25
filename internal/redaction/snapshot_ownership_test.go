package redaction

import (
	"testing"

	"github.com/example/otel-tail-sampling-gateway/internal/telemetry_domain"
)

func TestRedactedSpanOwnsAttributes(t *testing.T) {
	input := telemetry_domain.Span{
		Attributes: telemetry_domain.Attributes{"email": "person@example.test", "route": "/checkout"},
		Resource:   telemetry_domain.Resource{Attributes: telemetry_domain.Attributes{"region": "east"}},
	}
	redacted := Apply(input, telemetry_domain.AttributePolicy{Hash: []string{"email"}, Salt: "salt"})
	input.Attributes["route"] = "/mutated"
	input.Resource.Attributes["region"] = "west"
	if got := redacted.Attributes["route"]; got != "/checkout" {
		t.Fatalf("redacted attribute changed to %q", got)
	}
	if got := redacted.Resource.Attributes["region"]; got != "east" {
		t.Fatalf("redacted resource changed to %q", got)
	}
}
