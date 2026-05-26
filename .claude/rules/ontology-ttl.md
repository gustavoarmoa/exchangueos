---
glob: ".base/aasc/ontology/**/*.ttl"
---

# Rule: .base/aasc/ontology/*.ttl

## Conventions
- OWL 2 DL profile compliant
- Header com `owl:versionIRI` + `owl:versionInfo "1.2.0"` + bilingual labels (en/pt)
- Namespace `http://exchangeos.revenu.tech/ontology/<modulo>#`
- Per-module file (1 modulo = 1 TTL no core/)
- Reuso LedgerOS via `owl:imports <http://ledgeros.revenu.tech/ontology/finance>`

## Validation
- pyshacl validate em CI
- HermiT consistency check
- FIBO alignment via `skos:closeMatch`
- ISO 20022 alignment via `skos:related`

## Version Bump
- Add class/property backward-compat: minor
- Rename/remove: major
- Bump policy em `.base/aasc/ontology/CHANGELOG.md`
