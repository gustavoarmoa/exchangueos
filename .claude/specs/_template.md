---
feature: <feature-name>
status: DRAFT | REVIEW | APPROVED | IMPLEMENTED
priority: P0 | P1 | P2
sprint: <X>
owner: <team>
related_milestone: MS-XXX
related_rn_fx: [RN_FX_NNN, ...]
---

# Spec: <Feature Name>

## Why
<Business justification>

## What
<Functional description>

## Acceptance Criteria

- [ ] Criterion 1
- [ ] Criterion 2

## API Contract

```proto
// proto/exchangeos/v1/<service>.proto
service XxxService {
  rpc XxxYyy(XxxRequest) returns (XxxResponse);
}
```

## Data Contract

```sql
-- migrations/000NNN_xxx.up.sql
CREATE TABLE IF NOT EXISTS xxx (...);
```

## Test Cases (TDD)

```go
func TestXxx_HappyPath(t *testing.T) { ... }
func TestXxx_ValidationErrors(t *testing.T) { ... }
```

## Cross-References

- Patterns: FX-DDD-NNN, FX-EDA-NNN
- Ontology: .base/aasc/ontology/core/<bc>.ttl
- Flows: RFLW.024.NNN.NN
- ERDs: .base/erds/domain/erd-<bc>-domain.md
