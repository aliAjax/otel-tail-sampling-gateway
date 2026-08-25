#!/usr/bin/env sh
set -eu
base="${OTEL_BASE_URL:-http://127.0.0.1:8125}"
curl -fsS "$base/healthz" >/dev/null
curl -fsS "$base/readyz" >/dev/null
curl -fsS "$base/api/v1/tenants" >/dev/null
curl -fsS "$base/api/v1/buffers" >/dev/null
curl -fsS "$base/api/v1/drop-reasons" >/dev/null
echo "smoke ok"
