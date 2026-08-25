package transport

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type RequestMeta struct {
	ID          string
	Tenant      string
	Started     time.Time
	ContentType string
}

func Meta(r *http.Request) RequestMeta {
	return RequestMeta{ID: r.Header.Get("X-Request-ID"), Tenant: r.Header.Get("X-Tenant-ID"), Started: time.Now(), ContentType: r.Header.Get("Content-Type")}
}
func ValidateTenant(t string) error {
	if t == "" {
		return fmt.Errorf("tenant_required")
	}
	if len(t) > 128 {
		return fmt.Errorf("tenant_too_long")
	}
	for _, r := range t {
		if !(r == '-' || r == '_' || r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9') {
			return fmt.Errorf("tenant_invalid")
		}
	}
	return nil
}
func ValidateSignal(s string) error {
	switch strings.ToLower(s) {
	case "traces", "metrics", "logs":
		return nil
	default:
		return fmt.Errorf("signal_invalid")
	}
}
func StatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	if strings.Contains(err.Error(), "too_large") {
		return 413
	}
	if strings.Contains(err.Error(), "quota") {
		return 429
	}
	return 400
}
func ErrorBody(id string, e error) map[string]string {
	return map[string]string{"code": "invalid_request", "message": e.Error(), "request_id": id}
}
func ContentAllowed(r *http.Request) bool {
	v := r.Header.Get("Content-Type")
	return strings.Contains(v, "json") || strings.Contains(v, "protobuf") || v == ""
}
func Deadline(r *http.Request) time.Duration {
	if v := r.Header.Get("X-Timeout-Ms"); v != "" {
		var n int64
		fmt.Sscan(v, &n)
		if n >= 0 && n < 60000 {
			return time.Duration(n) * time.Millisecond
		}
	}
	return 10 * time.Second
}
func IsCompressed(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Content-Encoding"), "gzip")
}
func Header(r *http.Request, k string) string { return strings.TrimSpace(r.Header.Get(k)) }
func RequestAge(m RequestMeta) time.Duration  { return time.Since(m.Started) }
func EnsureRequestID(r *http.Request) string {
	if x := Header(r, "X-Request-ID"); x != "" {
		return x
	}
	return fmt.Sprintf("req-%d", time.Now().UnixNano())
}
