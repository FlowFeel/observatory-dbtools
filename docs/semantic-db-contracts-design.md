# SMW Database Architecture — Semantic Contracts Design

## Status: SPEC — Two-Pipe Architecture (authoritative)
Branch: `refactor/semantic-db-contracts` on `FlowFeel/observatory-dbtools`
Date: 2026-08-27
Authors: Ed Phil, Flow

---

## 1. The Two-Pipe Architecture

The Two-Pipe architecture establishes an authoritative contract between
declarative domain intelligence and relational persistence. Semantic MediaWiki
represents a dual storage paradigm where abstract assertions are partitioned
across identity lookup tables (`smw_object_ids`), metadata registries
(`smw_fpt_*`), and typed data-item tables (`smw_di_*`). Decoupling the PHP
domain compiler from the Go database engine via a versioned serialization
artifact (`catalog.json`) preserves language independence while enforcing
rigorous structural assertions over the persistence layer.

```
┌──────────────────────────────────────┐
│   PHP Semantic Domain (The Producer) │
│   • Sovereign Property Manifests     │
│   • Entity Composition Manifests     │
│   • Hermeneutic Compiler             │
└──────────────────┬───────────────────┘
                   │
                   ▼  (bin/export-catalog.php)
┌──────────────────────────────────────┐
│   catalog.json (Version 1)           │  <--- Portable Serialized Contract
└──────────────────┬───────────────────┘
                   │
                   ▼  (pkg/catalog Loader)
┌──────────────────────────────────────┐
│   Go dbtools (The Auditor)           │
│   • pkg/audit  (Contracts 1, 2, 3)   │
│   • pkg/curate (Page Generation)     │
│   • pkg/drift  (Topological Registry)│
└──────────────────┬───────────────────┘
                   │
                   ▼  (Cursor-Streamed Verification)
┌──────────────────────────────────────┐
│   MySQL Semantic MediaWiki Store     │
│   • smw_object_ids, smw_fpt_*, smw_di_* │
└──────────────────────────────────────┘
```

---

## 2. Complete Database Accounting

### 2.1 Table Categories (120 tables total)

| Category | Count | Purpose | Semantic Relation |
|----------|-------|---------|-------------------|
| MW Core | ~40 | page, revision, user, comment, slots, content | None — content storage |
| MW Links | ~10 | pagelinks, templatelinks, categorylinks, externallinks | None — structural links |
| MW Search | ~5 | searchindex, querycache, querycachetwo | None — search infra |
| MW Object Cache | ~3 | objectcache, l10n_cache, updatelog | None — cache |
| SMW Object IDs | 1 | smw_object_ids | **Direct — subject→ID binding** |
| SMW Data Item (DI) | 7 | smw_di_blob, _bool, _coords, _number, _time, _uri, _wikipage | **Direct — stores property values by type** |
| SMW Fixed Property (FPT) | 20 | smw_fpt_type, _pval, _mdat, _inst, _subp, _impo, ... | **Direct — stores property metadata** |
| SMW Stats | 3 | smw_prop_stats, smw_query_links, smw_object_aux | **Indirect — usage counts** |
| SMW Concept | 1 | smw_concept_cache | None — concept cache |
| Extension tables | ~30 | echo, discussiontools, ajaxpoll, etc. | None — extension-specific |

### 2.2 SMW Tables — The Semantic Boundary

#### smw_object_ids — The Identity Table

Every entity (page, property, value) gets an integer ID. `smw_id` is the
Bedeutung (reference), `smw_title` is the Sinn (sense/name).

| Column | Role | Maps To |
|--------|------|---------|
| smw_id | Integer identity (PK, auto-increment) | PropertyName/EntityName → int |
| smw_namespace | MW namespace | Entity archetype (NS_MAIN=Article, NS_EVENT, etc.) |
| smw_title | Canonical name | PropertyName::toString() |
| smw_iw | Interwiki prefix | — |
| smw_subobject | Subobject name | Participant subobject refs |
| smw_sortkey | Sort key | — |

#### smw_di_* — Data Item Tables (Value Storage by Type)

| Table | SMW Type | Catalog PropertyType | Value Column |
|-------|----------|---------------------|-------------|
| smw_di_blob | Text, Code | TEXT | o_hash (short), o_blob (long) |
| smw_di_bool | Boolean | BOOLEAN | o_value (tinyint) |
| smw_di_coords | Geographic coordinate | — (not in catalog) | o_blob |
| smw_di_number | Number, Quantity | NUMBER | o_blob, o_hash |
| smw_di_time | Date, Time | DATE | o_serialized, o_sortkey |
| smw_di_uri | URL, Email, Annotation URI | URL, EMAIL | o_serialized |
| smw_di_wikipage | Page | PAGE | o_id (→ smw_object_ids) |

#### smw_fpt_* — Fixed Property Tables (Metadata Storage)

| Table | SMW Built-in | Catalog Field | What It Stores |
|-------|-------------|---------------|----------------|
| smw_fpt_type | Has type | PropertyType | The type declaration for a property |
| smw_fpt_pval | Allows value | AllowedValues | Range constraints for a property |
| smw_fpt_subp | Subproperty of | — (not yet in catalog) | Subproperty hierarchy |
| smw_fpt_impo | Equivalent property | equivalentProperty | External ontology mapping (schema.org) |
| smw_fpt_inst | Instance of | — (category membership) | Entity → category |
| smw_fpt_mdat | Modification date | — (internal) | _MDAT routing (drift check target) |
| smw_fpt_dtitle | Display title | — | Display title override |
| smw_fpt_text | Has property | — | Text property values |
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

#### smw_prop_stats — Usage Statistics

| Column | Role |
|--------|------|
| p_id | Property ID (→ smw_object_ids) |
| usage_count | Number of entities using this property |
| null_count | Number of null values |

---

## 3. The Three Contracts

### Contract 1: Declaration-Storage Consistency

Contract 1 enforces top-down structural realization. Natural language gloss: an
authoritative verification policy ensuring that every abstract property defined
in the catalog is materialized as a physical entity within the database.

**Identity Allocation:** Asserts that every property name exists within
`smw_object_ids` under the dedicated property namespace (`smw_namespace = 102`).
Missing entries indicate uninstantiated property definitions that cannot
participate in graph joins.

**Metadata Table Synchronization:** Validates that the property's semantic type
in `smw_fpt_type`, its range bounds in `smw_fpt_pval`, and its external
vocabulary mapping in `smw_fpt_impo` match the catalog declaration.

**Transactional State Awareness:** Detects incomplete MediaWiki job queue
executions where property pages were created in wikitext but have not yet been
parsed into the fixed property metadata tables. The contract distinguishes
between "property was never declared" (true violation) and "property declared
but not yet parsed" (transient state — report as warning, not failure).

### Contract 2: Value-Type Consistency and Table Routing

Contract 2 enforces physical type integrity from the database back to the domain
model. Natural language gloss: a relational routing assertion verifying that
stored attribute values reside exclusively within the data-item table designated
for their declared semantic type.

**Mitigating Historical Type Drift:** When an administrator alters a property's
type in wikitext (for instance, changing a property from unstructured text to an
ISO date), the semantic engine does not retroactively migrate existing rows
from `smw_di_blob` to `smw_di_time`. Contract 2 acts as a forensic diagnostic
identifying orphaned values trapped in deprecated physical tables.

**Foreign Key and Table Alignment:** Traverses every row in each `smw_di_*`
table, verifies that the foreign key property identifier (`p_id`) resolves to a
known catalog property, and confirms that the physical table aligns with the
declared `PropertyType`.

| Stored Value Table | Declared PropertyType | Storage Integrity |
|--------------------|-----------------------|--------------------|
| smw_di_blob | Text, Code | Valid Routing |
| smw_di_time | Date | Valid Routing |
| smw_di_uri | URL, Email | Valid Routing |
| smw_di_wikipage | Page | Valid Routing |
| smw_di_blob | Date (Historical Drift) | Routing Violation |

### Contract 3: Value-Range Consistency and Historical Invariants

Contract 3 audits historical values against domain invariants. Natural language
gloss: a semantic audit verifying that all persisted literal values conform to
range restrictions and syntactic formatting constraints.

**Set-Theoretic Range Verification:** Inspects discrete literal values stored in
`smw_di_blob` or `smw_di_wikipage` and asserts membership within the declared
`AllowedValues` set.

**Syntactic Scalar Validation:** Verifies that date strings in `smw_di_time`
adhere to standardized ISO temporal specifications and that URI strings in
`smw_di_uri` conform to absolute web addresses.

**Catching Legacy Data:** Identifies corrupt or non-standard assertions that
were committed before the introduction of real-time edit-filtering hooks.

---

## 4. Go Architectural Subsystems in dbtools

The Go implementation decomposes verification into isolated functional packages,
preventing runtime coupling with the PHP application.

### The Catalog Ingestion Layer (`pkg/catalog`)

Loads and parses `catalog.json` into typed, read-only Go memory structures.
Natural language gloss: a stateless parser providing an in-memory knowledge
index for database auditing routines. It validates that the serialized artifact
contains a compatible major version header before exposing the domain schema.

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
    Version   int                  `json:"version"`
    Properties []PropertyDeclaration `json:"properties"`
    Entities   []EntityDefinition   `json:"entities"`
}
```

### The Streaming Audit Engine (`pkg/audit`)

Executes validation queries across large datasets without memory exhaustion.
Natural language gloss: a batch-processing engine that queries database rows
using bounded cursor streams and evaluates them against the catalog.

**Cursor-Based Paging:** Executes queries using indexed cursor ranges
(`WHERE s_id > last_seen_id ORDER BY s_id ASC LIMIT 1000`) rather than loading
unbounded tables into memory, allowing safe execution across multi-gigabyte
production databases.

**Diagnostic Reporting:** Accumulates violations into structured diagnostic
models detailing the subject entity, property identifier, invalid value, and
specific rule failure without terminating audit runs prematurely.

```go
// Contract 1: Declaration-Storage Consistency
func AuditDeclarations(db *sql.DB, catalog *Catalog) (*AuditReport, error)

// Contract 2: Value-Type Consistency
func AuditValueTypes(db *sql.DB, catalog *Catalog) (*AuditReport, error)

// Contract 3: Value-Range Consistency
func AuditValueRanges(db *sql.DB, catalog *Catalog) (*AuditReport, error)
```

### Dynamic Schema Curation (`pkg/curate`)

Replaces hardcoded property arrays with dynamic generation driven by
`pkg/catalog`. Natural language gloss: a schema maintenance service that
provisions missing semantic property pages directly from declarative
specifications.

### Topological Drift Registry (`pkg/drift`)

Generalizes property drift tracking beyond single hardcoded identifiers.
Natural language gloss: a configurable registry mapping arbitrary property
identities to their expected fixed-property and data-item storage locations.

```go
type DriftTarget struct {
    PropertyName string
    FptTable     string
    DiTable      string
    Pid          int
}
```

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
    Given a compiled property catalog with property "Event type"
    And a database where "Event type" has no smw_object_ids entry
    Then an audit should report a declaration violation
```

```gherkin
Feature: Semantic Value-Type Consistency
  As a database engineer
  I want to verify that stored values reside in the correct data-item table
  So that SMW routing integrity is maintained

  Scenario: Historical type drift — Date values in smw_di_blob
    Given a property "Event start date" declared as Date in the catalog
    And a database with smw_di_blob rows for that property's p_id
    Then an audit should report a routing violation
    And the diagnostic should include the expected table (smw_di_time)
    And the diagnostic should include the actual table (smw_di_blob)
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

### Phase 1: Contract Serialization

Implement `bin/export-catalog.php` in the magazine repository to emit the sealed
`catalog.json` with an explicit semantic version header during build pipelines.

### Phase 2: In-Memory Go Loader

Implement `pkg/catalog` in dbtools, verifying deserialization with static JSON
test fixtures.

### Phase 3: Relational Auditors

Implement `AuditDeclarations`, `AuditValueTypes`, and `AuditValueRanges` in
`pkg/audit`, utilizing `testcontainers-go` to run integration tests against real
MySQL instances pre-loaded with valid and invalid semantic fixture sets.

### Phase 4: Curation and Drift Alignment

Update `pkg/curate` to consume the loaded catalog and configure `pkg/drift`
using the topological registry pattern.

### Phase 5: BDD Integration

Implement Godog scenarios in `test/bdd_test.go` executing end-to-end contract
validation against realistic MediaWiki schemas.

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
