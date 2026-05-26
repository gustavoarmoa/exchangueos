---
name: ontology-validate
description: Run SHACL + OWL 2 DL validation completa em .base/aasc/ontology/
allowed-tools: [Bash, Read, Glob]
---

# Skill: /ontology-validate

## Trigger
`/ontology-validate [--scope <core|bridges|shapes|compliance|domains|all>]`

## Workflow (via `ontology-shacl` agent)
1. Run Apache Jena OWL 2 DL profile check
2. Run HermiT consistency reasoning
3. Run pyshacl: shapes/*.ttl vs fixtures/*.ttl
4. Validate FIBO coverage (target >= 80% relevant classes)
5. Validate ISO 20022 coverage (target 100% dos 32 schemas FX)
6. Report: triples count + validation result + drift
