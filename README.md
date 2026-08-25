# OpenTelemetry Tail Sampling Gateway

纯 Go 的多租户遥测接入、聚合、尾采样、脱敏和可靠导出网关。项目默认只使用标准库，外部 Kafka/OTLP 适配器不可用时使用可重复的内存/file simulator，明确返回 `adapter_unavailable`，不会伪造持久化成功。

## 启动

```bash
go run ./cmd/gateway
curl http://127.0.0.1:8080/healthz
curl -X POST http://127.0.0.1:8080/v1/otlp/traces -H 'Content-Type: application/json' -H 'X-Tenant-ID: demo' -d '{"trace_id":"abc","span_id":"s1","name":"checkout","duration_ms":12,"status":"ok"}'
```

## 模块

`telemetry_domain` 领域对象，`otlp_codec` 受限 JSON/protobuf frame 解码，`trace_aggregator` trace FSM，`sampling_engine` 策略 AST 和解释树，`metric_engine` temporality/去重，`redaction` 脱敏，`buffer_store` WAL/spill/checkpoint，`exporter` 批处理/重试/DLQ，`quota` 配额，`repository` 内存仓储，`transport` HTTP/gRPC-like 入口，`workers` 后台 worker，`observability` 指标日志，`audit` 审计。

## 验证

```bash
gofmt -w .
go vet ./...
go test ./...
go test -race ./...
go build ./...
```

非测试 Go 行数：`find . -name '*.go' ! -name '*_test.go' -print0 | xargs -0 wc -l`。`smoke/smoke.sh` 覆盖 healthz、tenant、策略、OTLP ingest、decision、buffer/drop、drain/replay。

## API

HTTP: `/healthz`, `/readyz`, `/v1/otlp/{traces,metrics,logs}`, `/api/v1/tenants`, `/sources`, `/sampling-policies`, `/attribute-policies`, `/exporters`, `/api/v1/traces/{id}/decision`, `/buffers`, `/exporter-health`, `/drop-reasons`, `/replay-jobs`, `/buffers/{id}/drain`, `/policies/{id}/validate`。每次请求产生 `request_id`，统一错误结构并限制 body 4 MiB。
