#!/bin/sh
set -eu
curl -fsS http://127.0.0.1:8080/healthz
curl -fsS -X POST http://127.0.0.1:8080/api/v1/tenants -H 'Content-Type: application/json' -d '{"id":"smoke","name":"Smoke","quota":100}'
curl -fsS -X POST http://127.0.0.1:8080/v1/otlp/traces -H 'Content-Type: application/json' -H 'X-Tenant-ID: demo' -d '{"trace_id":"smoke-trace","span_id":"s1","name":"smoke","duration_ms":5,"status":"ok"}'
curl -fsS http://127.0.0.1:8080/api/v1/traces/smoke-trace/decision
curl -fsS http://127.0.0.1:8080/api/v1/buffers
curl -fsS -X POST http://127.0.0.1:8080/api/v1/buffers/drain
