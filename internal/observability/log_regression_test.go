package observability

import (
	"sync"
	"testing"
)

func TestLogDoesNotMutateInput(t *testing.T) {
	fields := map[string]any{"tenant": "north", "labels": []string{"cold"}}
	Log("ingest", fields)
	if _, ok := fields["event"]; ok {
		t.Fatal("Log mutated the caller-owned fields map")
	}
	fields["labels"].([]string)[0] = "warm"
	if got := fields["labels"].([]string)[0]; got != "warm" {
		t.Fatalf("caller labels changed unexpectedly: %s", got)
	}
}

func TestLogNilFieldsAndConcurrentCalls(t *testing.T) {
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			Log("heartbeat", nil)
		}()
	}
	close(start)
	wg.Wait()
}
