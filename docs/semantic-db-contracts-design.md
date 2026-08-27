# SMW Database Architecture — Semantic Contracts Design

## Status: DRAFT — Two-Pipe Design (not yet implemented)
Branch: `refactor/semantic-db-contracts` on `FlowFeel/observatory-dbtools`
Date: 2026-08-27

---

## 1. The Problem

dbtools manages the MySQL storage layer (schema, baselines, drift, curation).
The magazine now has a pure PHP semantics domain that declares 83 sovereign
properties with types, allowed values, and schema.org mappings, compiled by a
hermeneutic compiler into a sealed `PropertyCatalog`.

These two systems share a boundary — the SMW relational tables — but have no
contract between them. dbtools validates that tables *exist*; it does not
validate that stored *values* conform to declared *semantics*. The `curate`
package hardcodes property lists. The `drift` package hardcodes property IDs.
No contract verifies that the database state matches the catalog declarations.

This document is a complete accounting of the database design from the
perspective of the semantic redesign, and a proposal for the contract layer.

---

## 2. Complete Database Accounting

### 2.1 Table Categories (120 tables total)

| Category | Count | Purpose | Semantic Relation |
|----------|-------|---------|-------------------|
| MW Core | ~40 | page, revision, user, comment, slots, content | None — content storage |
| MW Links | ~10 | pagelinks, templatelinks, categorylinks, externallinks | None — structural links |
| MW Search | ~5 | searchindex, querycache, querycachetwo | None — search infra |
| MW Object Cache | ~3 | objectcache, l10n_cache, updatelog | None — cache |
| SMW Object IDs | 1 | smw_object_ids | **Direct**/subject→ID binding** |
| SMW Data Item (DI) | 7 | smw_di_blob, _bool, _coords, _number, _time, _uri, _wikipage | **Direct — stores property values by type** |
| SMW Fixed Property (FPT) | 20 | smw_fpt_type, _pval, _mdat, _inst, _subp, _impo, ... | **Direct — stores property metadata** |
|(adsbygoogle = window.adsbygoogle || []).push({});
| SMW Stats | 3 | smw_prop_stats, smw_query_links, smw_object_aux | **Indirect — usage counts** |
| SMW Concept | 1 | smw_concept_cache | None — concept cache |
| Extension tables | ~30 | echo, discussiontools, ajaxpoll, etc. | None — extension-specific |

### 2.2 SMW Tables — The Semantic Boundary

#### smw_object_ids — The Identity Table

Every entity (page, property, value) gets an integer ID here. This is the
rigid designator at the storage level — `smw_id` is the Bedeutung (reference),
`smw_title` is the Sinn (sense/name).

| Column | Role | Maps To |
|--------|------|---------|
| smw_id | Integer identity (PK, auto-increment) | PropertyName/EntityName → int |
| smw_namespace | MW namespace | Entity archetype (NS_MAIN=Article, NS_EVENT, etc.) |
| smw_title | Canonical name | PropertyName::toString() |
| smw_iw | Interwiki prefix | — |
| smw_subobject | Subobject name | Participant subobject refs |
| smw_sortkey | Sort key | — |

**Contract:** Every `Property:` page in `smw_object_ids` must correspond to a
declared property in the catalog. Every catalog property must have a
corresponding `smw_object_ids` entry (smw_namespace=102, smw_title=property name).

#### smw_di_* — Data Item Tables (Value Storage by Type)

Each table stores values for properties of a specific SMW type. The `p_id`
column references `smw_object_ids.smw_id` for the property. The `s_id`
references the subject entity.

| Table | SMW Type | Catalog PropertyType | Value Column |
|-------|----------|---------------------|-------------|
| smw_di_blob | Text, Code | TEXT | o_hash (short) /5, o_blob (long) |
| smw_di_bool | Boolean | BOOLEAN | o_value (tinyint) |
| smw_di_coords | Geographic coordinate | — (not in catalog) | o_blob |
| smw_di_number | Number, Quantity | NUMBER | o_blob, o_hash |
| smw_di_time | Date, Time | DATE | o_serialized, o_sortkey |
| smw_di_uri | URL, Email, Annotation URI | URL, EMAIL | o_serialized |
| smw_di_wikipage | Page | PAGE | o_id (→ smw_object_ids) |

**Contract:** For every row in `smw_di_*`, the `p_id` must resolve to a
declared property in the catalog whose `PropertyType` matches the table's type.
A `smw_di_time` row with a `p_id` pointing to a TEXT property is a type
violation — the value is stored in the wrong table.

#### smw_fpt_* — Fixed Property Tables (Metadata Storage)

These tables store property *metadata* — type declarations, allowed values,
subproperty relations, imports (equivalence mappings), etc. Each FPT table
corresponds to a built-in SMW property.

| Table | SMW Built-in | Catalog Field | What It Stores |
|-------|-------------|---------------|----------------|
| smw_fpt_type | Has type | PropertyType | The type declaration for a property |
| smw_fpt_pval | Allows value | AllowedValues | Range constraints for a property |
| smw_fpt_subp | Subproperty of | — (not yet in catalog) | Subproperty hierarchy |
| smw_fpt_impo | Equivalent property | equivalentProperty | External ontology mapping (schema.org) |
| smw_fpt_inst | Instance of | — (category membership) | Entity → category |
| smw_fpt_mdat | Modification date | — (internal) | _MDAT routing (drift check target) |
| smw_fpt_dtitle | Display title | — (Display title of) | Display title override |
| smw_fpt_text | Has property | — (text storage) | Text property values |
| smw_fpt_uri | Has URI | — | URI property values |
| smw_fpt_sobj | Has subobject | — | Subobject existence |
| smw_fpt_redi | Redirect | — | Property redirects |
| smw_fpt_prec | Precision | — | Display precision |
| smw_fpt_pplb | Preferred property label | — | Human-readable label |
| smw_fpt_list | List | — | List formatting |
| smw_fpt_lcode | Language code | — | Monolingual text |
| smw_fpt_conv | Conversion | — | Unit conversion |
| smw_fpt_conc | Concept | — | Concept cache |
| smw_fpt_cdat | Creation date | — | Internal |
| smw_fpt_ask* | Ask query | — | Query links (7 tables) |
| smw_fpt_ledt | Local embedding | — | — |
| smw_fpt_serv | Service | — | — |
| smw_fpt_unit | Unit | — | Display unit |

**Contract:** For every property in the catalog:
1. `smw_fpt_type` must contain a row with `s_id` = property's `smw_id` and
   `o_serialized` matching the declared `PropertyType`.
2. If the catalog declares `AllowedValues`, `smw_fpt_pval` must contain rows
   for each allowed value with `s_id` = property's `smw_id`.
3. If the catalog declares `equivalentProperty`, `smw_fpt_impo` must contain
   a row with the external IRI.

#### smw_prop_stats — Usage Statistics

| Column | Role |
|--------|------|
| p_id | Property ID (→ smw_object_ids) |
| usage_count | Number of entities using this property |
| null_count | Number of null values |

**Contract:** Every `p_id` in `smw_prop_stats` must resolve to a catalog
property. Every catalog property with usage > 0 must have a `smw_prop_stats`
entry.

---

## 3. The Contract Layer — Three Contracts

### Contract 1: Declaration-Storage Consistency

**Direction:** Catalog → Database

Every property declared in the catalog must be present in the database with
matching metadata:
- `smw_object_ids` entry exists (smw_namespace=102, smw_title=property name)
- `smw_fpt_type` entry matches declared `PropertyType`
- `smw_fpt_pval` entries match declared `AllowedValues` (when non-empty)
- `smw_fpt_impo` entry matches declared `equivalentProperty` (when non-null)

**Violation example:** Catalog declares `Event type` as Text with allowed
values [In-Person, Virtual, Hybrid], but `smw_fpt_pval` has no rows for the
property's `s_id`. The property page was never created or was deleted.

### Contract 2: Value-Type Consistency

**Direction:** Database → Catalog

Every stored value must be in the correct data-item table for its property's
declared type:
- A property declared as DATE must have values in `smw_di_time`, not `smw_di_blob`
- A property declared as URL must have values in `smw_di_uri`, not `smw_di_blob`
- A property declared as PAGE must have values in `smw_di_wikipage`

**Violation example:** `Event start date` is declared as DATE in the catalog,
but a value "2026-08-27" is stored in `smw_di_blob` (Text table) instead of
`smw_di_time`. This happens when SMW's type routing is corrupted or when a
property's type was changed after values were already stored.

### Contract 3: Value-Range Consistency

**Direction:** Database → Catalog

Every stored value must satisfy the declared allowed-values constraint:
- If `AllowedValues` is non-empty, the stored value must be in the set
- If the type is DATE, the value must be a valid date
- If the type is URL, the value must be a valid URL

**Violation example:** `Event type` allows [In-Person, Virtual, Hybrid], but
`smw_di_blob` contains "Webinar" for some entity. This is the historical drift
the edit gate prevents going forward — this contract catches what already
slipped through.

---

## 4. Proposed Architecture

###> dbtools is a Go service. The catalog is PHP. The contracts must be
> verifiable without a PHP runtime. Solution: the catalog manifests are
> pure data (PHP files returning arrays/DTOs). A codegen step exports
> the compiled catalog to a portable format (JSON) that Go consumes.

### 4.1 Catalog Export — The Shared Artifact

A new `bin/export-catalog.php` script in the magazine repo:
- Loads the property manifests and entity compositions
- Compiles via `CatalogBuilder::build()`
- Exports to `catalog.json` — a flat JSON file with all property declarations

```json
{
  "properties": [
    {
      "name": "Event type",
      "type": "Text",
      "allowed": ["In-Person", "Virtual", "Hybrid"],
      "equivalent": "https://schema.org/eventAttendanceMode",
      "aliases": ["Type", "Attendance mode"]
    },
    ...
  ],
  "entities": [
    {
      "name": "Event",
      "gloss": "A temporal occurrence...",
      "properties": ["Event type", "Event start date", ...]
    },
    ...
  ]
}
```

This file is the contract surface — generated, not hand-edited. Go reads it,
PHP writes it. The codegen runs in CI before dbtools tests.

### 4.2 Go Packages — Redesigned

#### `pkg/catalog` (NEW) — Catalog Loader

Loads `catalog.json` into typed Go structs. Pure parsing, no DB access.

```go
type PropertyDeclaration struct {
    Name       string   `json:"name"`
    Type       string   `json:"type"`
    Allowed    []string `json:"allowed"`
    Equivalent string   `json:"equivalent"`
    Aliases    []string `json:"aliases"`
}

type EntityDefinition struct {
    Name       string   `json:"name"`
    Gloss      string   `json:"gloss"`
    Properties []string `json:"properties"`
}

type Catalog struct {
    Properties []PropertyDeclaration `json:"properties"`
    Entities   []EntityDefinition   `json:"entities"`
}
```

#### `pkg/audit` (NEW) — Contract Verifier

Three audit functions, one per contract:

```go
// Contract 1: Declaration-Storage Consistency
func AuditDeclarations(db *sql.DB, catalog *Catalog) (*AuditReport, error)

// Contract 2: Value-Type Consistency
func AuditValueTypes(db *sql.DB, catalog *Catalog) (*AuditReport, error)

// Contract 3: Value-Range Consistency
func AuditValueRanges(db *sql.DB, catalog *Catalog) (*AuditReport, error)
```

Each returns an `AuditReport` with violations (table, entity, property, value,
diagnostic). Pure logic against the DB query results + catalog — no side effects.

#### `pkg/curate` (REDESIGNED) — Catalog-Driven Curation

`RequiredProperties` is replaced by `catalog.RequiredPropertyPages()` —
generated from the loaded `catalog.json`. No more hardcoded lists.

#### `pkg/drift` (GENERALIZED) — Drift Registry

Instead of hardcoded `p_id=29`, a drift registry maps property names to their
FPT/DI table pairs:

```go
type DriftTarget struct {
    PropertyName string
    FptTable     string
    DiTable      string
    Pid          int&bsp;int
}
```

`_MDAT` is one entry. Future drift targets (e.g., `_INST` routing) are added
to the registry without new code.

---

## 5. BDD Contract Features

```gherkin
Feature: Semantic Property Declaration-Storage Consistency
  As a database engineer
  I want to verify that every catalog property is present in the SMW store
  So that the database reflects the declared semantic schema

  Scenario: All catalog properties exist in smw_object_ids
    Given a compiled property catalog
    And a database with SMW tables loaded
    Then every catalog property must have an smw_object_ids entry
    And its smw_fpt_type must match the declared type
    And its smw_fpt_pval must match the declared allowed values

  Scenario: Property missing from database
    Given a catalog/compiled property catalog with property "Event type"
    And a database where "Event type" has no smw_object_ids entry
    Then an audit should report a declaration violation
```

```gherkin
Feature: Semantic Value-Range Consistency
  As a database engineer
  I want to verify that stored property values satisfy declared constraints
  So that historical drift is detected and reported

  Scenario: Value outside allowed range
    Given a property "Event type" with allowed values [In-Person, Virtual, Hybrid]
    And a database with smw_di_blob containing "Webinar" for that property
    Then an audit should report a range violation
    And the violation diagnostic should include the allowed values
```

---

## 6. Implementation Sequencing

1. **Catalog export** (magazine side) — `bin/export-catalog.php` generates
   `catalog.json`. CI artifact. This is the shared contract surface.

2. **`pkg/catalog`** (dbtools) — Go loader for `catalog.json`. Unit tests
   with a test fixture JSON file.

3. **`pkg/audit`** (dbtools) — Three contract auditors. Integration tests
   with testcontainers: load seed, inject known violations, verify detection.

4. **`pkg/curate` redesign** — Replace hardcoded lists with catalog-driven
   generation. Existing tests updated.

5. **BDD features** — `features/semantic_contracts.feature`. Godog step
   definitions in `test/bdd_test.go`.

6. **`pkg/drift` generalization** — Drift registry. Lower priority — only
   needed for multi-property drift tracking.

---

## 7. Patterns & Anti-Patterns (Design Grounding)

### Patterns Applied

| Pattern | Source | How |
|---------|--------|-----|
| Logic/IO Separation | `:phos.arch/fka` (PHOS-ARCH-001) | Pure PHP domain, pure Go domain, serialized artifact as boundary |
| Hexagonal (Ports & Adapters) | Phosphene architecture | `catalog.json` is the port; PHP export and Go loader are opposing adapters |
| Consumer-Driven Contracts | Pact pattern | Catalog (consumer) declares expectations; `pkg/audit` verifies provider (DB) satisfies them |
| Decomposition Before Code | War story: template-provider-contracts | Three contracts independently verifiable; design reviewed before implementation |
| Test Pyramid | dbtools existing | Unit (pkg/catalog), Integration (pkg/audit via testcontainers), BDD (Godog contract features) |
| JSON as Boundary Only | dbtools existing principle | `catalog.json` is serialization artifact; Go structs and PHP DTOs have no cross-dependency |
| Positive Predication | `:phos.arch/positive-predication` (PHOS-ARCH-009) | Contracts state what the DB *must* satisfy, not what it must avoid |
| Graduated Discovery | `:phos.arch/discovery` (PHOS-ARCH-008) | Audit reports progress from summary → per-contract → per-violation detail |

### Anti-Patterns Eliminated

| Anti-Pattern | Where | Fix |
|--------------|-------|-----|
| Hardcoded Domain Knowledge | `curate.RequiredProperties`, `drift p_id=29` | Catalog-driven generation, drift registry |
| Implicit Contract | No verification between catalog and DB | Three explicit, auditable contracts |
| Unidirectional Assumption | DB-only source of truth | Bidirectional: catalog→DB and DB→catalog |

### Anti-Patterns to Guard Against in Implementation

**A. Brittle Serialization (`catalog.json` schema drift).**
If the export format changes (new fields, restructured entities), the Go loader
breaks silently. `catalog.json` MUST include a `"version": 1` field. The Go
loader validates the version before parsing. Format changes are semver — major
version bump for breaking changes, minor for additive.

**B. Unbounded Query (audit scope explosion).**
Contract 3 (value-range) queries every stored value for every constrained
property. On a production database with 800K+ SMW object IDs, a naive
`SELECT * FROM smw_di_blob` is catastrophic. The audit MUST stream results
via cursor-based pagination (`WHERE s_id > last_id LIMIT 1000`). The audit
report accumulates violations without holding the full result set in memory.

**C. Oversimplified Test Fixture.**
The existing `bdd_test.go` creates synthetic minimal tables. Contract tests
need the actual baseline schema with realistic property declarations — at
minimum, the 5 entity archetypes with their full property sets from the
catalog. Test fixtures that don't match production shape produce false greens.
This is the lesson from the template-provider-contracts war story: "pushing
code to CI to discover what your code does is the hairball reflex."

**D. Coupled Release Cycles.**
If the magazine catalog and dbtools must ship simultaneously, we've created a
distributed monolith. The `catalog.json` artifact must be backwards-compatible
across minor versions. dbtools should degrade gracefully on missing fields
(OWA compliance — `:phos.arch/owa`, PHOS-ARCH-002): unknown properties in
the catalog are preserved and reported as warnings, not failures.

---

## 8. What Doesn't Change

- `pkg/connect` — connection management, no semantic relation
- `pkg/baseline` — schema import/export, infrastructure
- `pkg/snapshot` — schema introspection, infrastructure
- `pkg/compare` — schema diffing, infrastructure
- `pkg/migrate` — migration execution, infrastructure
- `pkg/report` — serialization, boundary
- `pkg/asset` — content boundary enforcement, orthogonal
- `pkg/search` — Elasticsearch audit, orthogonal

These packages manage database *structure*. The new packages manage database
*semantics*. The boundary is clean.
