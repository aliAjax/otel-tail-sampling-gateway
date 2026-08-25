package telemetry_domain

import "time"

type Attributes map[string]string
type Resource struct {
	ServiceName    string     `json:"service_name"`
	ServiceVersion string     `json:"service_version,omitempty"`
	Attributes     Attributes `json:"attributes,omitempty"`
}
type Span struct {
	TenantID   string     `json:"tenant_id"`
	TraceID    string     `json:"trace_id"`
	SpanID     string     `json:"span_id"`
	ParentID   string     `json:"parent_id,omitempty"`
	Name       string     `json:"name"`
	Start      time.Time  `json:"start"`
	End        time.Time  `json:"end"`
	DurationMS int64      `json:"duration_ms"`
	Status     string     `json:"status"`
	HTTPStatus int        `json:"http_status,omitempty"`
	Attributes Attributes `json:"attributes,omitempty"`
	Resource   Resource   `json:"resource,omitempty"`
}

func (s Span) Duration() time.Duration {
	if s.DurationMS > 0 {
		return time.Duration(s.DurationMS) * time.Millisecond
	}
	if !s.End.IsZero() && !s.Start.IsZero() {
		return s.End.Sub(s.Start)
	}
	return 0
}

type Trace struct {
	ID        string    `json:"trace_id"`
	TenantID  string    `json:"tenant_id"`
	Spans     []Span    `json:"spans"`
	FirstSeen time.Time `json:"first_seen"`
	LastSeen  time.Time `json:"last_seen"`
	State     string    `json:"state"`
}
type SamplingDecision struct {
	TraceID   string    `json:"trace_id"`
	TenantID  string    `json:"tenant_id"`
	Keep      bool      `json:"keep"`
	Reason    string    `json:"reason"`
	Score     float64   `json:"score"`
	Explained []string  `json:"explained"`
	At        time.Time `json:"at"`
}
type MetricPoint struct {
	TenantID, Name, Temporality string
	Value                       float64
	Timestamp                   time.Time
	Attributes                  Attributes
	ExemplarTraceID             string
}
type LogRecord struct {
	TenantID, Body, Severity, TraceID string
	Timestamp                         time.Time
	Attributes                        Attributes
}
type DropReason struct {
	TenantID, Reason string
	Count            uint64
	Updated          time.Time
}
type Tenant struct {
	ID, Name  string
	Quota     uint64
	CreatedAt time.Time
}
type SamplingPolicy struct {
	ID, TenantID  string
	Enabled       bool
	Probability   float64
	ErrorOnly     bool
	MinDurationMS int64
	RarePath      string
	Priority      int
}
type AttributePolicy struct {
	ID, TenantID string
	Deny         []string
	Hash         []string
	Regex        map[string]string
	Salt         string
}
type ExporterConfig struct {
	ID, TenantID, Kind, Endpoint string
	Enabled                      bool
	BatchSize                    int
	MaxRetries                   int
}
