# Projection Discipline — T-541 (Phase 1 of the projection line)

**Status:** SPEC — executable behaviors are the spec. This doc is a pointer, not a
duplicate. **Authoritative spec: `features/projection_discipline.feature`** (BDD) plus
`pkg/catalog/discipline.go` + `pkg/catalog/discipline_test.go` (the implementation and
its unit proof). When doc and feature disagree, the feature wins.

Date: 2026-08-27
Branch: `refactor/semantic-db-contracts` on `FlowFeel/observatory-dbtools`
Author: Flow (with Ed Phil)

---

## The axiom

- **Truth lives in the PHP semantic domain** — the manifests, `PropertyType`,
  `AllowedValues`, `EntityDefinition`, the compiler. Nothing else asserts.
- **catalog.json is a projection surface, never a record of truth.** It is
  generatable, regenerable, and designed for serialization across the PHP→Go
  boundary and as a test fixture. Any disagreement between PHP and JSON is
  settled by regenerating, never by editing JSON.
- **JSON is the wire, not the negotiation.** Negotiation happens in Go code.
  The artifact carries *facts*; it never carries *decisions*. Every field we
  add to the projection is a fork in the negotiation — so the schema is
  designed to stay boring: flat, shallow, regenerable, no derived fields.
- **Provenance rides on the wire** (Phase 2) so a projection can be
  retranslated back to the truth-state it was projected from.

## The projection line (drawn in sand, Phase 1)

1. **Flatness** — entries are flat records. Objects appear only at the artifact
   root and at the entry level; arrays only as direct children of the artifact
   (`properties`, `entities`) or of an entry (`allowed`, `aliases`,
   `properties`). No nested records, no nested lists. (JSON-hell defense.)
2. **Exclusion** — derived fields (routing tables, SMW internal codes, canonical
   title forms, normalized values) must never be serialized. Consumers derive
   them locally from declared facts. The wire never decides.
3. **OWA** — unknown non-derived fields are tolerated and surfaced as warnings,
   never as hard failures.

## What Phase 1 shipped

| Artifact | What it proves |
|----------|----------------|
| `features/projection_discipline.feature` | 5 executable scenarios: clean artifact passes; smuggled derived field flagged; OWA tolerance; nested structures rejected; routing derived from type, never the wire |
| `pkg/catalog/discipline.go` | `catalog.Inspect` — structural walk independent of typed parse; reports `derived_field_smuggled`, `nesting_depth_exceeded`, `unknown_field` (OWA warning) |
| `pkg/catalog/discipline_test.go` | 11 unit tests against the committed fixture + hostile artifacts. **Run without Docker.** |
| `pkg/audit/derivation_test.go` | Derivation totality: every fixture type routable; `ExpectedTable` is type-driven; wire field cannot influence derivation |
| `test/bdd_test.go` | BDD step defs wiring the 5 scenarios into the CI suite |

**Baseline debt registered (measured starting point):** the Go derivation tables
(`pkg/audit/validators.go` type→table, `_txt/_dat/...` codes, `catalog.SMWTitle`)
are the negotiation logic — correct placement, but previously unproven and
silently divergent. **RESOLVED by T-546:** `RoutingTotality()` now proves them total
and free of truth-divergent entries.

**Divergence surfaced and fixed (T-546):** Go's routing maps carried a `Code`
type that the PHP `PropertyType` enum never declares. That is exactly the
silent-divergence the totality contract exists to catch — Go knew a type truth
doesn't. Fixed by removing the dead `Code` routing; the load-time cross-check
now rejects any artifact emitting it until PHP truth declares it.

## Phase roadmap — Go first, PHP after (agile, each phase lands as executable behaviors)

**Phase 8 (GO) — DONE: Prove Go owns derivation (T-546, T-547, T-548)**
- T-546: totality contract — `RoutingTotality()`, `KnownTypes` registry, title
  collisions. BDD: `features/derivation_totality.feature`.
- T-547: load-time cross-check — `catalog.Parse` rejects unknown types loudly.
  BDD: `features/load_time_cross_check.feature`.
- T-548: exclusion contract — `IsDerivedField`, OWA-warn vs violation; producer
  reject is PHP-side (Phase 7). BDD: `features/exclusion_contract.feature`.
- Isolation: `TestBDDDiscipline` runs all of it without Docker.

**Phase 7 (PHP) — Provenance + determinism on the producer**
- T-543: provenance block in artifact — content hash over compiled manifests +
  compiler identity. Content-hash *only*; no timestamp (regeneration must stay
  byte-identical).
- T-544: regenerability contract (PHPUnit) — same source → identical bytes;
  manifest edit → hash changes.
- T-545: schema-shape contract (PHPUnit) — exporter emits exactly the allowed v1
  keys; nesting ≤2; flat lists; no key outside the set.

**Phase 4 — Provenance-driven retranslation (Go + CI)**
- T-549: Go loader verifies provenance (hash present, matches pinned fixture).
- T-550: flip the T-539 freshness job — CI regenerates catalog.json from PHP and
  diffs vs committed artifact. "Regeneration proves truth," not "fixture is truth."
- T-551: retranslation proof — Go rebuilds audit expectations from projection +
  provenance alone.

**Phase 5 — Cross the line (flat completeness)**
- T-552: entity-level requiredness as a *declared PHP fact* (flat `required:[...]`
  on entity), consumed by edit gate + auditor from one source. Closes the
  completeness gap without nesting or bloat.

## Open negotiation (Phase 3 fork)

Where does type→table belong? Current position: it is SMW persistence knowledge,
therefore *derived* → stays in Go, made total + load-time-verified (T-546/547).
The alternative — declaring it in PHP — bloats the projection and makes it
SMW-specific, contradicting "dbtools is a DB ops toolkit, not an SMW tool."
Default: Go-derived. The tests are the negotiation table.
