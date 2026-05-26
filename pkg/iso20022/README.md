# pkg/iso20022 — ExchangeOS ISO 20022 Toolkit

Covers 32 FX-specific schemas across **fxtr / admi / camt / reda / head** namespaces, including CLS Bank (`CLSBUS33`) and CFETS (PTPP) submitter variants.

## Layout

```
iso20022/
├── doc.go              # Package-level overview + the 32-schema list
├── README.md           # This file
├── registry/           # Version Registry + Organisation Router (+ unit tests)
│   ├── registry.go     # Descriptor, Registry, lookup helpers
│   ├── router.go       # OrganisationRouter (CLSBUS33 / CFETS / ISO fallback)
│   ├── sources.go      # Default() — registers the 32 schemas pinned by version
│   ├── errors.go       # Sentinel errors
│   └── registry_test.go
├── marshaller/         # Envelope (head.001 BAH) + canonical XML marshal/unmarshal
└── validator/          # Lightweight XSD checks + RN_FX_* BusinessRule runner
```

## Schema Catalog (32 total)

| Namespace | Count | Examples |
|-----------|-------|----------|
| `fxtr` (CLS) | 7 | 008, 013, 014, 015, 016, 017, 030 |
| `fxtr` (CFETS) | 8 | 031–038 |
| `admi` (CLS) | 6 + 1 | 002, 004, 009, 010, 011, 017, 024 |
| `camt` (CLS) | 4 | 061, 062, 063, 088 |
| `reda` (CLS) | 2 + 2 | 060, 061, 066, 067 |
| `head` (ISO) | 2 | 001, 002 |

`fxti` and `fxmt` namespaces **do not exist** in ISO 20022. Quote and Amendment are internal Revenu gRPC services that translate to `fxtr.014/015/016` (CLS) or `fxtr.031/035/036` (CFETS) on the boundary.

## Usage

```go
reg := registry.Default()                 // 32 schemas registered
router := registry.NewOrganisationRouter(
    []string{"DEUTDEFF","CHASUS33"},      // CLS member BICs
    "CFETS",                              // CFETS BIC prefix
)

org, err := router.RouteParty("DEUTDEFF", "DE") // → OrgCLS
desc, _ := reg.LookupByURN("urn:iso:std:iso:20022:tech:xsd:fxtr.014.001.05")
```

## XSDs

Pinned URLs in `registry/sources.go`. Download via `scripts/download-xsd.sh` (MS-023a follow-up). XSDs are NEVER embedded in the binary; CI verifies pinned versions are reachable.

## Validation Layers

1. **`validator.XSDValidator`** — pure-Go well-formedness + URN-in-registry check.
2. **Optional `xmllibxml` build tag** — libxml2-backed full XSD validation (out-of-scope core).
3. **`validator.BusinessRuleValidator`** — RN_FX_001..050 runtime checks; rules are cited in `modules/<bc>/domain/specifications/` and SHACL shapes (`.base/aasc/ontology/compliance/bacen-cambio-shapes.ttl`).

## Conventions

- `decimal.Decimal` for any money/rate (NEVER `float`).
- All struct types embed their own canonical XML namespace tag.
- Version bumps for an existing schema follow the descriptor `Version` field; old descriptors are kept until the new one ships in production.
