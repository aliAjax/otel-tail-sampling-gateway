package observability

import (
	"sync"
	"testing"
)

func TestCounterSnapshotDoesNotMutateStore(t *testing.T) {
	c := NewCounterSet()
	c.Inc("accepted")
	s := c.Snapshot()
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); <-start; s["accepted"] = 99 }()
	go func() { defer wg.Done(); <-start; c.Inc("accepted") }()
	close(start)
	wg.Wait()
	if got := c.Get("accepted"); got != 2 {
		t.Fatalf("counter = %d, want 2", got)
	}
}
