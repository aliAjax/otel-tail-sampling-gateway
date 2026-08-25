package telemetry_domain

import "testing"

func TestSpanCloneDeepCopiesNestedState(t *testing.T) {
	original := Span{
		Attributes: Attributes{"route": "/v1/traces"},
		Resource:   Resource{Attributes: Attributes{"zone": "a"}},
	}
	clone := original.Clone()
	clone.Attributes["route"] = "/changed"
	clone.Resource.Attributes["zone"] = "b"
	if got := original.Attributes["route"]; got != "/v1/traces" {
		t.Fatalf("original span changed to %q", got)
	}
	if got := original.Resource.Attributes["zone"]; got != "a" {
		t.Fatalf("original resource changed to %q", got)
	}
}
