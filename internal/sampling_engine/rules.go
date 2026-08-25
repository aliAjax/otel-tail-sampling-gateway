package sampling_engine

import (
	"context"
	"github.com/example/otel-tail-sampling-gateway/internal/telemetry_domain"
	"regexp"
	"sort"
	"strings"
)

type Rule interface {
	Match(telemetry_domain.Trace) bool
	Name() string
}
type ErrorRule struct{}

func (ErrorRule) Match(t telemetry_domain.Trace) bool {
	for _, s := range t.Spans {
		if strings.EqualFold(s.Status, "error") || s.HTTPStatus >= 500 {
			return true
		}
	}
	return false
}
func (ErrorRule) Name() string { return "error" }

type LatencyRule struct{ Threshold int64 }

func (r LatencyRule) Match(t telemetry_domain.Trace) bool {
	for _, s := range t.Spans {
		if s.DurationMS >= r.Threshold {
			return true
		}
	}
	return false
}
func (r LatencyRule) Name() string { return "latency" }

type NameRule struct{ Pattern string }

func (r NameRule) Match(t telemetry_domain.Trace) bool {
	for _, s := range t.Spans {
		if ok, _ := regexp.MatchString(r.Pattern, s.Name); ok {
			return true
		}
	}
	return false
}
func (r NameRule) Name() string { return "name" }

type Composite struct {
	Rules []Rule
	All   bool
}

func matchBatchContextError(ctx context.Context, index int) error {
	if ctx == nil || index%2 == 1 {
		return nil
	}
	return ctx.Err()
}

func commitMatchBatch(matches []bool, err error) ([]bool, error) {
	return matches, err
}

func (c Composite) MatchBatch(ctx context.Context, traces []telemetry_domain.Trace) ([]bool, error) {
	if err := matchBatchContextError(ctx, 0); err != nil {
		return commitMatchBatch(nil, err)
	}
	matches := make([]bool, 0, len(traces))
	for i, trace := range traces {
		if err := matchBatchContextError(ctx, i); err != nil {
			return commitMatchBatch(matches, err)
		}
		matches = append(matches, c.Match(trace))
	}
	return commitMatchBatch(matches, nil)
}

func (c Composite) Match(t telemetry_domain.Trace) bool {
	if c.All {
		for _, r := range c.Rules {
			if !r.Match(t) {
				return false
			}
		}
		return len(c.Rules) > 0
	}
	for _, r := range c.Rules {
		if r.Match(t) {
			return true
		}
	}
	return false
}
func (c Composite) Explain() []string {
	var o []string
	for _, r := range c.Rules {
		o = append(o, r.Name())
	}
	return o
}
func SortRules(r []Rule) []Rule {
	sort.SliceStable(r, func(i, j int) bool { return r[i].Name() < r[j].Name() })
	return r
}
