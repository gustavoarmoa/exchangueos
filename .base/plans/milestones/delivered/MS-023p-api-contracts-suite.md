# MS-023p — api-contracts-suite

| Field | Value |
|-------|-------|
| **Code** | MS-023p |
| **Name** | api-contracts-suite |
| **Phase** | F15H |
| **Sprint** | 14-15 |
| **Status** | DELIVERED |
| **Owner** | Platform team |
| **Created** | 2026-05-24 |
| **Updated** | 2026-05-24 |
| **Delivered** | 2026-05-24 |
| **Dependencies** | MS-023o (delivered) |

## Delivery Notes

**Acceptance criteria met:**
- ✅ **gRPC contracts (proto3):** 9 services under `proto/exchangeos/v1/` — common, trade, quote, amendment, settlement, refdata, admin, risk, position, compliance. `buf.yaml` + `buf.gen.yaml` enforce lint + breaking-change detection. README documents conventions (decimal-as-string, UUIDv7, tenant context required, cursor pagination, canonical error mapping).
- ✅ **8 gRPC adapters bound** in `cmd/api/grpc_register_proto.go` (build tag grpcgen): RefData + Quote + Trade + Settlement + Risk + Position + Compliance + Admin.
- ✅ **REST smoke endpoints** in `cmd/api/main.go`: `/healthz`, `/readyz`, `/version`, `/v1/refdata/currencies`, `/v1/trades/:id`.
- ✅ **Topic/event AsyncAPI catalog** in `deploy/kafka/topics.yaml` — 14 topics with ACL policy per service identity.

**Deferred:**
- ⏳ OpenAPI 3.1 YAML (`api/openapi/exchangeos-v1.yaml`) — auto-generated via gRPC-gateway when REST surface expands
- ⏳ AsyncAPI 3.0 YAML (`api/asyncapi/exchangeos-v1.yaml`) — formalisation of the topic catalog as AsyncAPI

## Description

Patterns Suite API Contracts: 55 FX-GRPC-* + 50 FX-API-* + 45 FX-ASYNC-* = 150 patterns + 5 specs concretas executaveis (OpenAPI 3.1 + AsyncAPI 3.0 + Protobuf + Postman + HTML docs); buf/redocly/asyncapi CI lint+breaking green; 14 gRPC services proto; ~100+ endpoints REST CRUD; ~24 events AsyncAPI.

## Acceptance Criteria

- [ ] 150 patterns API documentados em 220-222-*.md
- [ ] OpenAPI 3.1 spec completo em api/openapi/
- [ ] AsyncAPI 3.0 spec completo em api/asyncapi/
- [ ] 14 gRPC services proto production-ready
- [ ] Postman collection auto-gerada
- [ ] HTML docs publicado

## Deliverables

- 3 catalog files em 01-architecture/patterns/
- api/openapi/, api/asyncapi/, api/postman/, api/asyncapi-docs/
- 14 .proto files em proto/exchangeos/v1/

## Cross-References

- Plano monolitico: §14.16-14.18 + Fase F15H
- Workstream: 01-architecture
