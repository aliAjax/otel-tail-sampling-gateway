package redaction

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/example/otel-tail-sampling-gateway/internal/telemetry_domain"
	"regexp"
	"strings"
)

func Apply(s telemetry_domain.Span, p telemetry_domain.AttributePolicy) telemetry_domain.Span {
	out := s
	out.Attributes = telemetry_domain.Attributes{}
	deny := map[string]bool{}
	for _, k := range p.Deny {
		deny[k] = true
	}
	hashes := map[string]bool{}
	for _, k := range p.Hash {
		hashes[k] = true
	}
	for k, v := range s.Attributes {
		if deny[k] {
			continue
		}
		if hashes[k] {
			h := sha256.Sum256([]byte(p.Salt + v))
			v = hex.EncodeToString(h[:])
		}
		if pat, ok := p.Regex[k]; ok {
			if re, e := regexp.Compile(pat); e == nil {
				v = re.ReplaceAllString(v, "[REDACTED]")
			}
		}
		out.Attributes[k] = v
	}
	return out
}
func SanitizeLog(body string) string {
	for _, x := range []string{"password", "token", "secret"} {
		body = strings.ReplaceAll(body, x+"=", x+"=[REDACTED]")
	}
	return body
}
