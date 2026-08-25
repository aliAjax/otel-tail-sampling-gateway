package transport

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDeadlineRejectsZeroAndParsesMilliseconds(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("X-Timeout-Ms", " 25 ")
	if got := Deadline(req); got != 25*time.Millisecond {
		t.Fatalf("trimmed timeout = %s, want 25ms", got)
	}
	req.Header.Set("X-Timeout-Ms", "0")
	if got := Deadline(req); got != 10*time.Second {
		t.Fatalf("zero timeout = %s, want fallback", got)
	}
}
