package audit

import "testing"

func TestAuditAddOwnsEvent(t *testing.T) {
	log := &Log{}
	event := Event{RequestID: "req-1", Details: map[string]string{"status": "accepted"}, Labels: []string{"ingest"}}
	log.Add(event)
	event.Details["status"] = "mutated"
	event.Labels[0] = "changed"
	stored := log.List()[0]
	if stored.Details["status"] != "accepted" || stored.Labels[0] != "ingest" {
		t.Fatalf("stored event followed caller mutation: %#v", stored)
	}
}

func TestAuditListReturnsSnapshot(t *testing.T) {
	log := &Log{}
	log.Add(Event{RequestID: "req-1", Details: map[string]string{"status": "accepted"}, Labels: []string{"ingest"}})
	first := log.List()
	first[0].Details["status"] = "mutated"
	first[0].Labels[0] = "changed"
	second := log.List()[0]
	if second.Details["status"] != "accepted" || second.Labels[0] != "ingest" {
		t.Fatalf("published audit snapshot shares storage: %#v", second)
	}
}
