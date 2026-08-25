package otlp_codec

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/example/otel-tail-sampling-gateway/internal/telemetry_domain"
	"io"
	"net/http"
	"strings"
)

type Envelope struct {
	Traces  []telemetry_domain.Span        `json:"traces"`
	Metrics []telemetry_domain.MetricPoint `json:"metrics"`
	Logs    []telemetry_domain.LogRecord   `json:"logs"`
}

func Decode(r *http.Request, limit int64, dst any) error {
	if r.ContentLength > limit {
		return fmt.Errorf("payload_too_large")
	}
	var rd io.Reader = io.LimitReader(r.Body, limit+1)
	if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		z, e := gzip.NewReader(rd)
		if e != nil {
			return fmt.Errorf("gzip_invalid: %w", e)
		}
		defer z.Close()
		rd = z
	}
	dec := json.NewDecoder(rd)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return fmt.Errorf("otlp_decode: %w", err)
	}
	var extra json.RawMessage
	if err := dec.Decode(&extra); !errors.Is(err, io.EOF) {
		return fmt.Errorf("otlp_decode: trailing_data")
	}
	return nil
}
func Encode(v any) ([]byte, error) { return json.Marshal(v) }
