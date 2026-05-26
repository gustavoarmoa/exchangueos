# 01 — Architecture

> **Workstream:** Architecture
> **Versao:** 1.0.0
> **Status:** DRAFT

## Conteudo

| Documento | Status | Descricao |
|-----------|--------|-----------|
| `system-architecture.md` | TODO | Visao geral arquitetura ExchangeOS standalone |
| `actor-model.md` | TODO | Entity Relationship & Actor Model (CLS Bank + Settlement Members + CFETS + Counterparties) — §2.2 monolitico |
| `folder-structure.md` | TODO | Estrutura de pastas (padrao PaymentOS/AuthorityOS) — §3 monolitico |
| `flows-suite.md` | TODO | Flows Suite Master Index — `.base/flows/` (85 flows RFLW.024.NNN.NN) — §12 monolitico |
| `erds-suite.md` | TODO | ERDs Suite Master Index — `.base/erds/` (23 ERDs + 16 SQL DDL) — §13 monolitico |
| `legacy-swift-mt-bridge.md` | TODO | Legacy SWIFT MT Bridge (MT300/MT304/MT202) — §2.4 monolitico |
| `ddd-implementation-plan.md` | TODO | DDD implementation per BC |
| `context-map.md` | TODO | Bounded Context Map (14 BCs) |
| `streaming-architecture.md` | TODO | Kafka + Flink streaming architecture |
| `patterns/` | TODO | 20 catalogos FX-* patterns (850 patterns totais) |
| `product/` | TODO | Product architecture views |
| `standards/` | TODO | Architecture standards |

## Patterns Catalog (`patterns/`)

20 catalogos com **850 patterns Tier-1**:

| Code | Patterns | Foco | File |
|------|----------|------|------|
| FX-GP-* | 40 | Golang application | `patterns/200-fx-golang-patterns.md` |
| FX-DDD-* | 35 | Domain-Driven Design | `patterns/201-fx-ddd-patterns.md` |
| FX-EDA-* | 45 | Event-Driven Architecture | `patterns/202-fx-eda-patterns.md` |
| FX-CP-* | 50 | CockroachDB | `patterns/205-fx-cockroachdb-patterns.md` |
| FX-KP-* | 60 | Kafka | `patterns/206-fx-kafka-patterns.md` |
| FX-FP-* | 40 | Apache Flink | `patterns/207-fx-flink-patterns.md` |
| FX-DS-* | 50 | DevSecOps & CI-CD | `patterns/210-fx-devsecops-cicd-patterns.md` |
| FX-K8S-* | 40 | Kubernetes | `patterns/211-fx-kubernetes-patterns.md` |
| FX-IAC-* | 40 | Terraform + GCP | `patterns/212-fx-terraform-gcp-patterns.md` |
| FX-DOC-* | 20 | Docker & Container | `patterns/213-fx-docker-container-patterns.md` |
| FX-GRPC-* | 55 | gRPC | `patterns/220-fx-grpc-patterns.md` |
| FX-API-* | 50 | REST/OpenAPI/CRUD | `patterns/221-fx-api-rest-patterns.md` |
| FX-ASYNC-* | 45 | AsyncAPI 3.0 | `patterns/222-fx-asyncapi-patterns.md` |
| FX-IAM-* | 50 | IAM + RBAC + ISO 27001 | `patterns/230-fx-iam-rbac-patterns.md` |
| FX-OTEL-* | 60 | OpenTelemetry Go | `patterns/240-fx-opentelemetry-patterns.md` |
| FX-TEST-* | 40 | Testing strategy | `patterns/250-fx-testing-patterns.md` |
| FX-QA-* | 35 | TDD + E2E + Security Local | `patterns/260-fx-qa-tdd-e2e-patterns.md` |
| FX-SYNC-* | 40 | Database Sync + Cross-Module | `patterns/270-fx-sync-cross-module-patterns.md` |
| FX-INT-* | 25 | Integration Verification | `patterns/280-fx-integration-verification-patterns.md` |
| FX-XOS-* | 20 | Cross-Platform Tooling | `patterns/290-fx-cross-platform-patterns.md` |
| FX-COMMIT-* | 25 | Pre-Commit HARD Enforcement | `patterns/300-fx-precommit-enforcement-patterns.md` |

## Sources

- §2.1-2.4 (Modulo + Actor Model + CLS Cycle + SWIFT MT) + §3 (Estrutura) + §12 (Flows) + §13 (ERDs) + §14 (Patterns Suite) do [snapshot monolitico v3.11.7](../_archive/allenty-v3.11.7-monolithic-plan.md)
- Pattern de referencia: [LedgerOS 01-architecture](../../../../ledgeros/.base/plans/01-architecture/)
