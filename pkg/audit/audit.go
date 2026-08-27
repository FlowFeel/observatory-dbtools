// Package audit implements the three semantic contracts that verify the
// MySQL SMW persistence layer against the compiled property catalog:
//
//	Contract 1 — Declaration-Storage Consistency (catalog → DB)
//	Contract 2 — Value-Type Consistency and Table Routing (DB → catalog)
//	Contract 3 — Value-Range Consistency and Historical Invariants (DB → catalog)
//
// Design:
//   - Pure validation logic lives in validators.go (no DB) and is unit-tested
//     locally. DB-querying auditors live in declarations.go, value_types.go,
//     value_ranges.go and use testcontainers integration tests (CI-runnable).
//   - Keyset streaming (WHERE s_id > last_seen ORDER BY s_id ASC LIMIT N)
//     keeps memory bounded across multi-gigabyte production databases.
//   - No premature halting: audits accumulate all violations.
//   - Sampling bounds prevent log flooding on systematic corruption.
package audit

import (
	"fmt"
	"strings"
)

// ViolationKind classifies a detected non-conformity for graduated reports.
type ViolationKind string

const (
	// KindDeclaration is a Contract 1 violation: a declared property is
	// missing or mismatched in the SMW store.
	KindDeclaration ViolationKind = "declaration"
	// KindRouting is a Contract 2 violation: a value resides in the wrong
	// data-item table for its property's declared type.
	KindRouting ViolationKind = "routing"
	// KindOrphanedPredicate is a Contract 2/3 diagnostic: a data-item row
	// references a property identifier with no catalog declaration.
	KindOrphanedPredicate ViolationKind = "orphaned_predicate"
	// KindRange is a Contract 3 violation: a value falls outside the
	// declared allowed-value enumeration.
	KindRange ViolationKind = "range"
	// KindSyntax is a Contract 3 violation: a value fails its scalar
	// syntactic invariant (ISO 8601 date, RFC 3986 URI).
	KindSyntax ViolationKind = "syntax"
	// KindReference is a Contract 3 violation: a smw_di_wikipage edge
	// points to a non-existent subject.
	KindReference ViolationKind = "reference"
)

// Severity distinguishes true violations from transient states.
type Severity string

const (
	// SevError is a true violation requiring remediation.
	SevError Severity = "error"
	// SevWarning is a transient or advisory finding (e.g. job-queue
	// propagation lag where a page exists but metadata not yet parsed).
	SevWarning Severity = "warning"
)

// Violation is the immutable diagnostic envelope for a single non-conformity.
type Violation struct {
	Table      string        `json:"table"`
	EntityID   string        `json:"entity_id,omitempty"`
	EntityNS   int           `json:"entity_ns,omitempty"`
	Property   string        `json:"property,omitempty"`
	Value      string        `json:"value,omitempty"`
	Kind       ViolationKind `json:"kind"`
	Severity   Severity      `json:"severity"`
	Rule       string        `json:"rule"`
	Diagnostic string        `json:"diagnostic"`
}

// Summary is the high-level result of an audit run.
type Summary struct {
	RowsScanned  int            `json:"rows_scanned"`
	ErrorCount   int            `json:"error_count"`
	WarningCount int            `json:"warning_count"`
	PerProperty  map[string]int `json:"per_property_error_rate"`
	Pass         bool           `json:"pass"`
}

// CategoryBreakdown groups violations by constraint kind.
type CategoryBreakdown struct {
	DeclarationCount int `json:"declaration_count"`
	RoutingCount     int `json:"routing_count"`
	OrphanedCount    int `json:"orphaned_count"`
	RangeCount       int `json:"range_count"`
	SyntaxCount      int `json:"syntax_count"`
	ReferenceCount   int `json:"reference_count"`
}

// Report is the graduated aggregation of an audit run: summary → categorical
// breakdown → itemized violation log.
type Report struct {
	Summary       Summary           `json:"summary"`
	Categories    CategoryBreakdown `json:"categories"`
	Violations    []Violation       `json:"violations"`
	TotalViolated int               `json:"total_violated"`
}

// Add records a violation and updates the rolling counters. Sampling bounds
// are applied by the caller via opts before Add is called for raw instances.
func (r *Report) Add(v Violation) {
	r.Violations = append(r.Violations, v)
	switch v.Severity {
	case SevError:
		r.Summary.ErrorCount++
		if v.Property != "" {
			if r.Summary.PerProperty == nil {
				r.Summary.PerProperty = make(map[string]int)
			}
			r.Summary.PerProperty[v.Property]++
		}
	case SevWarning:
		r.Summary.WarningCount++
	}
	r.TotalViolated++
	switch v.Kind {
	case KindDeclaration:
		r.Categories.DeclarationCount++
	case KindRouting:
		r.Categories.RoutingCount++
	case KindOrphanedPredicate:
		r.Categories.OrphanedCount++
	case KindRange:
		r.Categories.RangeCount++
	case KindSyntax:
		r.Categories.SyntaxCount++
	case KindReference:
		r.Categories.ReferenceCount++
	}
}

// Finalize marks the report pass/fail and returns it.
func (r *Report) Finalize() *Report {
	r.Summary.Pass = r.Summary.ErrorCount == 0
	return r
}

// AuditOptions configures audit behavior.
type AuditOptions struct {
	// BatchSize is the keyset page size for streaming queries.
	BatchSize int
	// MaxViolationsPerProperty caps collected raw instances per property;
	// the global counter continues beyond the cap.
	MaxViolationsPerProperty int
}

// DefaultOptions returns the production-sane audit options.
func DefaultOptions() AuditOptions {
	return AuditOptions{
		BatchSize:                1000,
		MaxViolationsPerProperty: 100,
	}
}

// String returns a compact human-readable summary.
func (r *Report) String() string {
	if r.Summary.Pass {
		return fmt.Sprintf("AUDIT PASS: %d rows scanned, 0 errors (%d warnings)",
			r.Summary.RowsScanned, r.Summary.WarningCount)
	}
	return fmt.Sprintf("AUDIT FAIL: %d rows scanned, %d errors, %d warnings — %s",
		r.Summary.RowsScanned, r.Summary.ErrorCount, r.Summary.WarningCount,
		summarizeKinds(r.Categories))
}

func summarizeKinds(c CategoryBreakdown) string {
	var parts []string
	if c.DeclarationCount > 0 {
		parts = append(parts, fmt.Sprintf("%d declaration", c.DeclarationCount))
	}
	if c.RoutingCount > 0 {
		parts = append(parts, fmt.Sprintf("%d routing", c.RoutingCount))
	}
	if c.OrphanedCount > 0 {
		parts = append(parts, fmt.Sprintf("%d orphaned", c.OrphanedCount))
	}
	if c.RangeCount > 0 {
		parts = append(parts, fmt.Sprintf("%d range", c.RangeCount))
	}
	if c.SyntaxCount > 0 {
		parts = append(parts, fmt.Sprintf("%d syntax", c.SyntaxCount))
	}
	if c.ReferenceCount > 0 {
		parts = append(parts, fmt.Sprintf("%d reference", c.ReferenceCount))
	}
	if len(parts) == 0 {
		return "no violations"
	}
	return strings.Join(parts, ", ")
}
