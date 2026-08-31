# observatory-dbtools

Contract-driven, isolated Go database tooling for [Observatory.wiki](https://observatory.wiki) — migrations, drift detection, snapshot comparison, and schema validation.

---

## Architectural Philosophy & Standards

`observatory-dbtools` strictly adheres to **Phosphene Engineering Standards** (specifically the [`phosphene/go_bdd_reference`](https://github.com/phosphene/go_bdd_reference) design pattern).

### 1. Domain Decoupling & Serialization Isolation
- **JSON is strictly an external serialization and reporting format.**
- Core packages operate entirely on strongly-typed Go domain models (`Snapshot`, `Table`, `Column`, `Diff`, `Report`).
- The [`pkg/report`](file:///Users/edphillips/projects/new/observatory-dbtools/pkg/report/report.go) package encapsulates boundary formatting (`WriteJSON`), preventing JSON tags, formatting concerns, and serialization logic from leaking into core domain algorithms.

### 2. Standalone Database Tooling
- Database tooling validates and manages the **database state as infrastructure**, independently of application code or PHP maintenance scripts (`update.php`).
- Consumed by the [Observatory Magazine](https://github.com/janfrel/observatory-magazine-v2) CI/CD pipeline, but versioned independently to support any wiki on the same stack.

### 3. Golden Files as Invariant Contracts
- Golden files located in `pkg/*/testdata/` serve as strict contracts representing expected database structure and drift checks.
- **In CI:** Any mismatch against golden files results in an immediate build failure.
- **In PR Reviews:** Golden files are version-controlled and reviewed explicitly alongside schema changes.

---

## The Full Testing Pyramid

Our testing methodology enforces strict quality guarantees at every layer of the software stack:

```
          / \
         /   \       Layer 3: BDD Behavioral Specs (Godog)
        / BDD \      features/*.feature <-> test/bdd_test.go
       /-------\
      / Integr. \    Layer 2: Testcontainers Integration
     /           \   pkg/*/*_test.go (Live MySQL 8.4 Containers)
    /-------------\
   /     Unit      \ Layer 1: Pure In-Memory Unit Tests
  /-----------------\ pkg/compare, pkg/snapshot, pkg/report
```

### Layer 1: Pure In-Memory Unit Tests
- Fast, zero-dependency unit tests validating domain algorithms:
  - `pkg/compare`: Diff calculations between schema snapshots (`TestDiffSnapshots_Identical`, `TestDiffSnapshots_DetectsTableAddedAndRemoved`).
  - `pkg/snapshot`: Canonical table/column sorting and schema structure guarantees (`TestSnapshot_CanonicalSorting`).
  - `pkg/report`: Boundary JSON serialization and formatting (`TestWriteJSON`).

### Layer 2: Testcontainers Integration Tests
- Executed against ephemeral, real MySQL 8.4 container instances via [`testcontainers-go`](https://golang.testcontainers.org/):
  - `pkg/connect`: Validates live MySQL connection handshakes and ping responses (`TestConnect`).
  - `pkg/baseline`: Validates baseline SQL schema imports and seed data insertion (`TestImportSchema`, `TestImportSchemaAndSeed`).
  - `pkg/drift`: Validates Semantic MediaWiki (SMW) `_MDAT` property table (`smw_fpt_mdat`) vs Data Item (`smw_di_time`) mismatch detection and automatic repair (`TestDrift_ZeroOnCleanState`, `TestDrift_DetectsAndFixes`).
  - `pkg/migrate`: Validates SQL migration execution, rollback capability (`Up` / `Down`), and `schema_migrations` table tracking (`TestMigrate_UpAndDown`).

### Layer 3: BDD Behavioral Compliance Suite
- Human-readable Gherkin features (`features/`) paired with Go step definitions (`test/bdd_test.go`) via Godog:
  - [`features/db_invariants.feature`](file:///Users/edphillips/projects/new/observatory-dbtools/features/db_invariants.feature): Guarantees MediaWiki core tables (`user`, `page`, `revision`, `text`, `site_stats`) and SMW tables (`smw_object_ids`, `smw_fpt_mdat`, `smw_di_time`) exist and meet minimum volume thresholds.
  - [`features/smw_drift.feature`](file:///Users/edphillips/projects/new/observatory-dbtools/features/smw_drift.feature): Behavioral spec covering SMW storage drift detection, repair execution, and post-repair verification.

---

## Package Responsibilities

| Package | Path | Responsibility |
|---|---|---|
| `connect` | [`pkg/connect`](file:///Users/edphillips/projects/new/observatory-dbtools/pkg/connect) | Database connection creation, DSN parsing, and ping health verification. |
| `baseline` | [`pkg/baseline`](file:///Users/edphillips/projects/new/observatory-dbtools/pkg/baseline) | Importing canonical SQL baselines and seed datasets into target databases. |
| `snapshot` | [`pkg/snapshot`](file:///Users/edphillips/projects/new/observatory-dbtools/pkg/snapshot) | Introspecting database metadata (`information_schema`) into canonical `Snapshot` domain structs. |
| `compare` | [`pkg/compare`](file:///Users/edphillips/projects/new/observatory-dbtools/pkg/compare) | Computing structural diffs (`Diff`) between schema snapshots (detecting added/removed tables, modified columns, data types). |
| `migrate` | [`pkg/migrate`](file:///Users/edphillips/projects/new/observatory-dbtools/pkg/migrate) | Wrapper around `golang-migrate` for applying versioned SQL migration scripts. |
| `drift` | [`pkg/drift`](file:///Users/edphillips/projects/new/observatory-dbtools/pkg/drift) | SMW-specific storage drift detection and automated data reconciliation. |
| `report` | [`pkg/report`](file:///Users/edphillips/projects/new/observatory-dbtools/pkg/report) | JSON and text boundary formatting for CLI and CI/CD reporting outputs. |

---

## Quick Start & Usage

```bash
# Run unit, integration, and BDD tests (requires Docker / OrbStack)
go test -v -count=1 -timeout=300s ./pkg/... ./test/...

# Run only fast in-memory unit tests
go test -v ./pkg/compare ./pkg/snapshot ./pkg/report

# Update golden files after intentional schema changes
go test -v -count=1 -timeout=300s -update ./pkg/...

# Check drift on a target database using CLI
DB_HOST=localhost DB_PORT=3306 DB_PASS=secret go run ./cmd/dbtools drift --check

# Apply pending migrations
DB_HOST=localhost DB_PORT=3306 DB_PASS=secret go run ./cmd/dbtools migrate --dir=./migrations
```

---

## Initiative Context

Part of **[I-005: Database Tooling](https://github.com/phosphene/woodchipper)** — contract-driven migration testing and database management for the Observatory platform.

## Documentation Index

- **[Semantic DB Contracts Design](docs/semantic-db-contracts-design.md)** — the two-pipe architecture spec: PHP producer → `catalog.json` → Go auditor (authoritative).
- **[PHP Producer Side — Slice 1 State](docs/php-producer-slice1-state.md)** — T-424: the PHP raw SMW schema writer (`SmwBulkSink`/`MysqlSmwDb`) — dry-run-first, CI proof, landmines. The producer half of the two pipes.
- **[Architecture & Testing Strategy](docs/architecture.md)** — design principles, test pyramid, SMW drift reconciliation.
- **[Projection Discipline](docs/projection-discipline.md)** — page-generation projection rules as an executable spec.
