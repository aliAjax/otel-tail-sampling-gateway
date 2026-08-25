package sampling_engine

import (
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

func (c Composite) Match(t telemetry_domain.Trace) bool {
	SortRules(c.Rules)
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
