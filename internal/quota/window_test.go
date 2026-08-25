package quota

import (
	"testing"
	"time"
)

func TestQuotaUsageResetsExpiredWindow(t *testing.T) {
	l := New(10)
	if !l.Allow("tenant-a", 4) {
		t.Fatal("initial quota was rejected")
	}
	l.mu.Lock()
	l.reset = time.Now().Add(-time.Second)
	l.mu.Unlock()
	if got := l.Usage("tenant-a"); got != 0 {
		t.Fatalf("expired usage=%d, want 0", got)
	}
	if !l.Allow("tenant-a", 10) {
		t.Fatal("new window did not accept the full quota")
	}
}
