---
uri: docs/php-producer-slice1-state
title: PHP Producer Side — SmwBulkSink Slice 1 State (T-424)
owner: edphos
status: committed
updated: 2026-08-31
relates_to:
  - docs/semantic-db-contracts-design.md
  - janfrel/observatory-magazine-v2 docs/architecture/smw-bulk-sink.md
---

# PHP Producer Side — SmwBulkSink Slice 1 State

This document records the state of the **left-hand side of the two-pipe
architecture** (the PHP producer) as of 2026-08-31. It is the counterpart
document to `docs/semantic-db-contracts-design.md` (which specifies the Go
auditor side and the `catalog.json` contract). The dbtools auditor verifies
semantics of the SMW store; this doc records that the **write path** into that
store is now real, direct, and dry-run-first.

## 1. What landed

On `janfrel/observatory-magazine-v2` branch `feat/php-intelligence-layer`
(PR #268), ticket **T-424**:

| Commit | Content |
|--------|---------|
| `b884a9e6` | `SmwBulkSink` dry-run/diff + `PropertyTypeMap` (230 props) + `SmwValueEncoder` + `SmwDbGateway` port + `FakeSmwDb` + 10 tests |
| `0395c24c` | `MysqlSmwDb` real gateway (MySQL 8.4) + CI proof job (`smw-sink-mysql`) + port gap + `smw_fpt_inst` diff fix |

The PHP producer can now write derived facts **directly to the SMW schema
tables** (`smw_object_ids`, `smw_di_*`, `smw_fpt_inst`) — bypassing the Lua
`{{#set:...}}` stamping pipeline entirely. This is the concrete realization of
the `PHP Semantic Domain (The Producer)` box in the two-pipe diagram.

## 2. Why it matters to dbtools

The Go auditor's contracts (1/2/3) assume the SMW store is populated and
well-formed. This slice establishes:

- **A defined write contract into `smw_object_ids`** — no UNIQUE index on
  `(smw_namespace, smw_title)` ⇒ upsert must be SELECT-then-INSERT, never
  `ON DUPLICATE KEY UPDATE` (silent no-op). dbtools' own writes/migrations
  should respect the same constraint.
- **Typed row shapes per `smw_di_*` table** are now exercised against real
  MySQL in CI (`smw-sink-mysql` job), so the auditor's expected shapes have a
  producer-side proof to agree with.
- **The `smw_fpt_inst` shape** (categories: `s_id`, `o_id` — no `p_id`) is
  captured as a load-bearing fact; the auditor's drift/catalog logic must not
  assume a `p_id` on that table.
- **Known landmines the auditor should also track:** `smw_proptable_hash`
  recompute after bulk writes, `smw_rev`/`smw_touched` for change propagation,
  `smw_prop_stats` + `smw_object_aux` rebuild, and `smw_di_wikipage`
  denormalization (subject + object title/ns/sortkey, not just `o_id` FK).

## 3. CI proof (`smw-sink-mysql` on the magazine)

Three passes against real MySQL 8.4 seeded from the baseline snapshot:

1. **Dry-run** → `derived: 2`, **zero writes** (everything classifies insert).
2. **Live apply** (`--live`) → real `smw_*` INSERTs commit (column lists
   proven against actual schema).
3. **Re-dry-run** → `changed rows: 0`, **zero writes** (idempotency).

## 4. Producer-side test state

- 1,650 pure PHPUnit tests / 5,235 assertions green (magazine `phpunit-pure`),
  incl. 11 `SmwBulkSinkTest` (round-trip through fake DB, category unchanged
  classification, live-mode apply, report title resolution).
- Dry-run over the canonical corpus: 2 subjects, 31 typed rows, zero writes.

## 5. Open landmines before prod live apply

See the magazine's [`docs/architecture/smw-bulk-sink.md`](https://github.com/janfrel/observatory-magazine-v2/blob/feat/php-intelligence-layer/docs/architecture/smw-bulk-sink.md)
§8. The auditor and the producer share these constraints; a future dbtools
contract (or drift rule) could make `smw_proptable_hash` freshness enforceable
from the Go side.

## 6. Next slices (agreed order)

1. ✅ **Slice 1 — raw schema writer + dry-run/diff + real-gateway CI proof**
2. **Slice 2 — PHP parser-level query layer** replacing the 17 Lua ask modules
   (76 call-sites); ResultPrinter Architecture (AC1–AC4).
3. **Hash/stat landmines** wiring before any prod live apply.
4. Lua → pure view → delete dead set logic.
