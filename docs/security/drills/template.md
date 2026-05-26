# Tabletop Drill — <SCENARIO NAME> — <YYYY-Qn>

> Copy this file to `docs/security/drills/YYYY-Qn-<short-scenario>.md` for each drill run.

## Metadata

| Field | Value |
|-------|-------|
| Date | YYYY-MM-DD |
| Scenario | One of S-1..S-5 from incident-response.md OR bespoke |
| Facilitator | Name |
| IC under test | Name (rotates from on-call) |
| Participants | Names + roles |
| Duration | XX minutes (planned) |

## Scenario brief

> 2-3 sentence injection statement read to the room at T+0. Should be ambiguous
> enough to test diagnosis, specific enough to be actionable.

Example: "At 14:32 UTC the api SLO dashboard shows http_5xx_rate climbing through
3% in us-east1. Pod count appears stable. CRDB dashboards green. What do you do?"

## Timeline

| T+min | Actor | Action | Decision rationale |
|-------|-------|--------|--------------------|
| 0     | Facilitator | Read scenario | — |
| 2     | IC | … | … |
| …     | …  | … | … |

## Decisions made

1. …
2. …

## Observations

### What worked well
- …

### What was slow
- …

### What was wrong / missing
- …

## Action items

| Action | Owner | Target date |
|--------|-------|-------------|
| | | |

## Updates to artefacts

- [ ] `docs/security/incident-response.md` — add/update scenario steps
- [ ] `docs/security/threat-model-stride.md` — add new threat row if surfaced
- [ ] `docs/operations/canary-runbook.md` — refine mitigation steps
- [ ] `docs/security/iso27001-gap-tracker.md` — update if a control gap exposed

## Sign-off

| Role | Name | Date |
|------|------|------|
| Facilitator | | |
| Security Lead | | |
| Platform Lead | | |
