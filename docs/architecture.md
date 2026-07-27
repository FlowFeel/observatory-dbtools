# Observatory DB Tools — Architecture & Testing Strategy

## Executive Summary

`observatory-dbtools` provides a contract-driven framework for managing MediaWiki databases. It isolates database operations from PHP runtime environments and maintenance scripts, using Go binaries and containerized test instances (`testcontainers-go`).

---

## 1. Design Principles & Phosphene Standards

1. **Domain Model Primacy:**
   Domain entities (`Snapshot`, `Diff`, `Report`) represent system truth. Business logic functions consume and return domain types without coupling to external protocols or serialization tags.

2. **JSON as Boundary Serialization Only:**
   JSON is relegated strictly to communication boundaries (CLI JSON outputs, CI artifacts, golden files). Formatting and marshaling are handled exclusively by `pkg/report`.

3. **Hermetic Testcontainers Isolation:**
   All database integration tests execute against real MySQL 8.4 containers created on demand. No shared, persistent, or pre-existing local database state is assumed.

4. **Behavioral Invariant Enforcement:**
   Gherkin feature files in `features/` document core data layer guarantees. The BDD test runner (`test/bdd_test.go`) validates these rules against real database instances.

---

## 2. Test Pyramid Matrix

| Layer | Target Packages | Execution Context | Execution Time | Purpose |
|---|---|---|---|---|
| **Unit** | `pkg/compare`, `pkg/snapshot`, `pkg/report` | Pure Go memory | < 1s | Fast feedback on diffing algorithms & serialization |
| **Integration** | `pkg/connect`, `pkg/baseline`, `pkg/drift`, `pkg/migrate` | `testcontainers-go` (MySQL 8.4) | ~1-2 min | Live MySQL schema validation & drift fix proofs |
| **BDD Compliance** | `test/bdd_test.go` (`features/*`) | Godog + `testcontainers-go` | ~30s | End-to-end behavioral contract validation |

---

## 3. Storage Layer Drift Reconciliation (SMW)

Semantic MediaWiki stores property statistics in both Fixed Property Tables (`smw_fpt_*`) and Data Item tables (`smw_di_*`). Under high-write or migration scenarios, routing mismatches (`_MDAT`, `p_id=29`) can occur.

The `pkg/drift` package implements a two-phase reconciliation protocol:
1. **Check Phase (`drift.Check`):** Queries row counts and executes set subtraction (`s_id NOT IN (...)`) to detect unmirrored entries.
2. **Fix Phase (`drift.Fix`):** Atomic batch insertion migrating missing records from `smw_fpt_mdat` into `smw_di_time`.

---

## 4. Maintenance & Operations

- **Golden Files:** Stored in `pkg/*/testdata/*.json`. Used by contract tests to prevent unintended regression. Update using `go test -update ./pkg/...`.
- **Migration Scripts:** Managed via `golang-migrate` in `migrations/`. Naming convention: `<version>_<name>.[up|down].sql`.
