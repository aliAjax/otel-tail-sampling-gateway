package sampling_engine

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"github.com/example/otel-tail-sampling-gateway/internal/telemetry_domain"
	"strings"
	"time"
)

type Engine struct{}

func (e Engine) DecideBatch(ctx context.Context, traces []telemetry_domain.Trace, p telemetry_domain.SamplingPolicy) ([]telemetry_domain.SamplingDecision, error) {
	decisions := make([]telemetry_domain.SamplingDecision, len(traces))
	for i, trace := range traces {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		decisions[i] = e.Decide(trace, p)
	}
	return decisions, nil
}

func (Engine) Decide(t telemetry_domain.Trace, p telemetry_domain.SamplingPolicy) telemetry_domain.SamplingDecision {
	d := telemetry_domain.SamplingDecision{TraceID: t.ID, TenantID: t.TenantID, At: time.Now()}
	if !p.Enabled {
		d.Keep = true
		d.Reason = "policy_disabled"
		return d
	}
	score := hashScore(t.ID)
	d.Score = score
	if p.ErrorOnly {
		for _, s := range t.Spans {
			if strings.EqualFold(s.Status, "error") || s.HTTPStatus >= 500 {
				d.Keep = true
				d.Reason = "error"
				d.Explained = append(d.Explained, "error span")
				return d
			}
		}
		d.Reason = "no_error"
		return d
	}
	for _, s := range t.Spans {
		if s.Duration() >= time.Duration(p.MinDurationMS)*time.Millisecond {
			d.Keep = true
			d.Reason = "latency"
			d.Explained = append(d.Explained, "duration threshold")
			return d
		}
		if p.RarePath != "" && strings.Contains(s.Name, p.RarePath) {
			d.Keep = true
			d.Reason = "rare_path"
			return d
		}
	}
	d.Keep = score < p.Probability
	if d.Keep {
		d.Reason = "probability"
	} else {
		d.Reason = "probability_drop"
	}
	return d
}
func hashScore(id string) float64 {
	h := sha256.Sum256([]byte(id))
	return float64(binary.BigEndian.Uint64(h[:8])) / float64(^uint64(0))
}
func Explain(p telemetry_domain.SamplingPolicy) string {
	return fmt.Sprintf("enabled=%t probability=%.3f error_only=%t min_duration_ms=%d rare_path=%q", p.Enabled, p.Probability, p.ErrorOnly, p.MinDurationMS, p.RarePath)
}
