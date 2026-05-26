# 06 — Infrastructure

> **Workstream:** Infrastructure
> **Versao:** 1.0.0
> **Status:** DRAFT

## Conteudo

| Documento | Status | Descricao |
|-----------|--------|-----------|
| `local-deploy.md` | TODO | Deploy local + shared CRDB hub TLS (§17 monolitico) |
| `cross-platform-tooling.md` | TODO | Task + Makefile + PowerShell + Bash POSIX (§21 monolitico) |
| `observability.md` | TODO | OpenTelemetry stack + Tempo + Mimir + Loki + Grafana + GCP Cloud Ops dual (§16 monolitico) |
| `kubernetes.md` | TODO | GKE Autopilot 1.29+ + Helm + Argo Rollouts + GitOps ArgoCD + Falco + OPA Gatekeeper |
| `terraform-gcp.md` | TODO | Terraform 1.5+ + GCP (KMS CMEK + WIF + VPC SC + PSC + Cloud Armor + Binary Auth) |
| `docker.md` | TODO | Multi-stage distroless + non-root + SBOM + Cosign + SLSA L3 |
| `database.md` | TODO | CockroachDB Dedicated (multi-region NY+London+Sao Paulo) + CDC + backups |
| `kafka.md` | TODO | Kafka KRaft 3-broker RF=3 + mTLS + ACLs + Schema Registry |
| `flink.md` | TODO | Apache Flink K8s Operator + RocksDB state + CEP fraud detection |
| `vault.md` | TODO | HashiCorp Vault (HCP managed) + PKI 90d certs + Secret rotation 30d |
| `ibm-mq.md` | TODO | IBM MQ wrapper (reuso paymentos pattern) para SWIFT FIN bridge |
| `migration-sequence.md` | TODO | Sequence de migrations CockroachDB |
| `deployment-reference.md` | TODO | Reference deployments per environment (dev/staging/prod) |
| `first-deploy-runbook.md` | TODO | Runbook primeiro deploy production |
| `operations/` | TODO | Runbooks operacionais |
| `standards/` | TODO | Infra standards |

## Stack Cross-Platform

| Layer | Tool | macOS | Linux | Windows | WSL2 | Alpine |
|-------|------|-------|-------|---------|------|--------|
| Task runner | **Task** (taskfile.dev) | ✅ | ✅ | ✅ | ✅ | ✅ |
| Backward compat | **Makefile** (auto-gen) | ✅ | ✅ | ⚠ | ✅ | ✅ |
| Win32 fallback | **PowerShell 7+** | (cross) | (cross) | ✅ | (n/a) | (n/a) |
| Container | **Docker Compose v2** | ✅ | ✅ | ✅ | ✅ | n/a |

## Sources

- §16 (Telemetry OTel) + §17 (Deploy Local + CRDB hub) + §21 (Cross-Platform Tooling) do [snapshot monolitico v3.11.7](../_archive/allenty-v3.11.7-monolithic-plan.md)
- Pattern de referencia: [LedgerOS 06-infrastructure](../../../../ledgeros/.base/plans/06-infrastructure/)
