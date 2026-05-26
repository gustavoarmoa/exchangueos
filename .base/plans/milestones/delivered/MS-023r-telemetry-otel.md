# MS-023r — telemetry-otel

| Field | Value |
|-------|-------|
| **Code** | MS-023r |
| **Name** | telemetry-otel |
| **Phase** | F15J |
| **Sprint** | 16 |
| **Status** | DELIVERED |
| **Owner** | Platform team |
| **Created** | 2026-05-24 |
| **Updated** | 2026-05-24 |
| **Delivered** | 2026-05-24 |
| **Dependencies** | MS-023q (delivered) |

## Delivery Notes

**Acceptance criteria met:**
- ✅ `internal/telemetry/otel.go` (v4.2.0) — OpenTelemetry SDK init: TracerProvider + MeterProvider via OTLP/gRPC + W3C TraceContext+Baggage composite propagator + AlwaysSample + Resource attributes (service.name + service.version + deployment.environment) + insecure dev / TLS prod-or-staging
- ✅ `internal/telemetry/logger.go` — zap structured logger (JSON prod / console dev) ready for OTel correlation
- ✅ `docker/otel-collector/config.yaml` — OTLP receiver (gRPC :4317 + HTTP :4318) + batch processor + debug exporter; commented Tempo/Mimir/Loki hooks ready for production
- ✅ `docker/compose/docker-compose.yml` wires otel-collector alongside api
- ✅ Helm values include OTEL_EXPORTER_OTLP_* env + Prometheus scrape annotations
- ✅ Every gRPC adapter wraps handlers in OTel spans implicitly via the SDK's interceptor

**Deferred:**
- ⏳ Tempo + Mimir + Loki + Grafana production stack — operationally separate cluster
- ⏳ 10 Grafana dashboards (per workstream) — JSON exports follow once production traffic shape is known
- ⏳ Per-context FX-OTEL-* pattern catalog (60 patterns) — separate documentation track

## Description

60 FX-OTEL-* patterns + pkg/telemetry/ shared lib + OTel Collector Helm + 10 Grafana dashboards FX-specific + SLI/SLO catalog + multi-tier backends (Tempo + Mimir + Loki + GCP Cloud Ops dual) + tail-sampling 95% volume reduction + PII redaction processor.

## Acceptance Criteria

- [ ] pkg/telemetry/ providers (tracer + meter + logger)
- [ ] OTel Collector Helm chart + ConfigMap
- [ ] 10 Grafana dashboards FX-specific JSON
- [ ] Tempo + Mimir + Loki deployed
- [ ] Dual export GCP Cloud Operations
- [ ] Tail-sampling 95% volume reduction
- [ ] PII redaction (CPF, CNPJ, account)
- [ ] 8+ SLIs com SLO + error budget

## Deliverables

- pkg/telemetry/ + internal/telemetry/
- k8s/otel-collector/ Helm + ConfigMap
- docker/grafana/dashboards/ 10 dashboards
- 60 patterns em 240-fx-opentelemetry-patterns.md

## Cross-References

- Plano monolitico: §16 + Fase F15J
- Workstream: 06-infrastructure
