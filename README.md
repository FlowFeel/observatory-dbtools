# observatory-dbtools

Contract-driven database tooling for [Observatory.wiki](https://observatory.wiki) — migrations, drift detection, and schema validation.

## Architecture

Standalone Go tool consumed by the [Observatory Magazine](https://github.com/janfrel/observatory-magazine-v2) CI pipeline. Tests database state using [testcontainers-go](https://golang.testcontainers.org/) against disposable MySQL instances. Golden files serve as contracts — committed, reviewed, never auto-updated in CI.

## Building Blocks

| Package | Responsibility |
|---------|---------------|
| `connect` | DB connection + health check |
| `snapshot` | Schema dump to structured output |
| `compare` | Diff two schema states |
| `migrate` | Wrapper around golang-migrate |
| `drift` | SMW-specific FPT↔DI mismatch detection |
| `report` | Structured JSON output for CI/humans |

## Quick Start

```bash
# Run contract tests (needs Docker)
task test

# Check drift on a target database
DB_HOST=localhost DB_PORT=3306 DB_PASS=secret task drift:check

# Apply pending migrations
DB_HOST=localhost DB_PORT=3306 DB_PASS=secret task migrate
```

## Golden Files

Golden files in `pkg/*/testdata/` are the contracts. Tests compare actual database state against golden files.

- **CI:** Golden file mismatch = build failure.
- **Local:** Run `task test:update` to refresh after intentional schema changes.
- **PR review:** Golden file changes are reviewed like any other code.

## Why Standalone?

Database tooling tests the DATABASE state — infrastructure, not application. The Observatory Magazine is one consumer. Any wiki on the same stack can use these tools. The tool versions independently from the application it serves.

## Initiative

Part of [I-005: Database Tooling](https://github.com/phosphene/woodchipper) — contract-driven migration testing for the Observatory platform.
