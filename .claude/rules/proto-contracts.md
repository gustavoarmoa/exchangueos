---
glob: "proto/**/*.proto"
---

# Rule: proto/**/*.proto

## Conventions
- Proto3 syntax obrigatorio
- Package versioning: `exchangeos.v1` (NUNCA breaking changes em v1)
- buf lint + buf breaking enforce em CI
- protoc-gen-validate para validation
- 1 file proto per service

## Naming
- RPCs: verb + noun (`BookTrade`, `GetTrade`, `ListTrades`)
- Messages: `<RPC>Request` + `<RPC>Response`
- Enums: PascalCase + suffix `_UNSPECIFIED = 0` obrigatorio

## Cite Source
Reference proto evolution policy em `.base/plans/01-architecture/patterns/220-fx-grpc-patterns.md`
