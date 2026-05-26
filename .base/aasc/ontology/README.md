# ExchangeOS Ontology Suite

OWL 2 DL ontologies for the 14 bounded contexts + bridges (FIBO, ISO 20022) + SHACL shapes for compliance.

## Layout

```
.base/aasc/ontology/
├── README.md                   # this file
├── CHANGELOG.md                # version history (semver per ontology)
├── core/                       # 14 per-BC TTL files (1 per bounded context)
│   ├── trade.ttl               ✅ representative (v1.2.0)
│   ├── quote.ttl               ⏳
│   ├── refdata.ttl             ⏳
│   ├── cls_settlement.ttl      ⏳
│   ├── ...
├── bridges/                    # 9 alignment files (FIBO, ISO 20022, CFETS)
│   ├── fibo-trade.ttl
│   ├── iso20022-fxtr.ttl
│   ├── ...
├── shapes/                     # 8 SHACL validation shapes
│   └── trade-shapes.ttl
└── compliance/                 # 5 BACEN/regulation-specific shapes
    └── bacen-cambio-shapes.ttl
```

## Versioning

Per Allenty convention, each TTL carries:

- `owl:versionIRI` — e.g. `http://exchangeos.revenu.tech/ontology/trade/1.2.0`
- `owl:versionInfo` — same semver
- `dct:title` bilingual labels (en/pt)

CHANGELOG.md tracks per-ontology bump rationale (minor: add class/property backward-compat; major: rename/remove).

## Validation

```bash
# pyshacl validation against shapes/
pyshacl -s shapes/trade-shapes.ttl core/trade.ttl

# HermiT consistency check
java -jar HermiT.jar core/trade.ttl
```

## FIBO Alignment

Target: ≥ 80% of FX-relevant classes referenced via `skos:closeMatch` or `owl:equivalentClass`. See `bridges/fibo-*.ttl` for the per-domain alignment files.

## Status

- ✅ trade.ttl (this milestone) — representative; 13 classes + 4 object props + 5 datatype props
- ⏳ Remaining 13 core/ + 9 bridges/ + 8 shapes/ + 5 compliance/ — next sprint
