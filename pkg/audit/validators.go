package audit

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/FlowFeel/observatory-dbtools/pkg/catalog"
)

// ---------------------------------------------------------------------------
// SMW type routing
//
// Maps catalog PropertyType names to their physical storage table and to the
// SMW internal type code used in smw_fpt_type.o_serialized. These tables are
// the authoritative routing for Contract 2.
// ---------------------------------------------------------------------------

// smwTypeToTable maps a catalog PropertyType to its physical data-item table.
//
// This map is the negotiation side of routing: it must be TOTAL over the
// catalog's KnownTypes (PHP truth) and must NOT route a type truth never
// declares. Both directions are enforced by RoutingTotality. "Code" was
// removed — PHP PropertyType truth does not declare it; the load-time
// cross-check now rejects any artifact that emits it.
var smwTypeToTable = map[string]string{
	"Text":    "smw_di_blob",
	"Date":    "smw_di_time",
	"Number":  "smw_di_number",
	"URL":     "smw_di_uri",
	"Email":   "smw_di_uri",
	"Boolean": "smw_di_bool",
	"Page":    "smw_di_wikipage",
}

// smwTypeToCode maps a catalog PropertyType to the SMW internal type id.
var smwTypeToCode = map[string]string{
	"Text":    "_txt",
	"Date":    "_dat",
	"Number":  "_num",
	"URL":     "_uri",
	"Email":   "_ema",
	"Boolean": "_boo",
	"Page":    "_wpg",
}

// smwCodeToHuman maps SMW internal codes to human type names (for diagnostics).
var smwCodeToHuman = map[string]string{
	"_txt": "Text", "_dat": "Date", "_num": "Number",
	"_uri": "URL", "_ema": "Email", "_boo": "Boolean", "_wpg": "Page",
}

// dataItemTables are all smw_di_* tables the auditor scans for routing.
var dataItemTables = []string{
	"smw_di_blob", "smw_di_time", "smw_di_number",
	"smw_di_uri", "smw_di_bool", "smw_di_wikipage",
}

// ExpectedTable returns the physical storage table for a catalog type,
// and whether the type is routable (i.e. not an unknown type).
func ExpectedTable(propType string) (string, bool) {
	t, ok := smwTypeToTable[propType]
	return t, ok
}

// RoutingTotality verifies the derivation tables are complete and total over
// the catalog's KnownTypes (PHP truth), and do not route anything truth never
// declares. Two directions:
//
//   - Total over truth: every known type has a storage table, an SMW code,
//     and round-trips through that code; the audit scan set equals the set of
//     routable tables.
//   - No extras: a routed type, a routed code, or a decoded human name that
//     PHP truth does not declare is a divergence, reported as a problem.
//
// Returns a list of problems; empty means the derivation is total.
func RoutingTotality() []string {
	var problems []string

	for _, t := range catalog.KnownTypes() {
		tbl, ok := smwTypeToTable[t]
		if !ok {
			problems = append(problems, fmt.Sprintf("known type %q has no storage table", t))
			continue
		}
		code, ok := smwTypeToCode[t]
		if !ok {
			problems = append(problems, fmt.Sprintf("known type %q has no SMW code", t))
			continue
		}
		if human, ok := smwCodeToHuman[code]; !ok || human != t {
			problems = append(problems, fmt.Sprintf("type %q does not round-trip through its SMW code %q", t, code))
		}
		if !inStrings(dataItemTables, tbl) {
			problems = append(problems, fmt.Sprintf("routed table %q for %q is not in the audit scan set", tbl, t))
		}
	}

	// No extras: routed types, codes, and decoded names must all be truth.
	for t := range smwTypeToTable {
		if !catalog.KnownType(t) {
			problems = append(problems, fmt.Sprintf("routing table knows type %q that PHP truth does not declare", t))
		}
	}
	for t := range smwTypeToCode {
		if !catalog.KnownType(t) {
			problems = append(problems, fmt.Sprintf("type-code table knows type %q that PHP truth does not declare", t))
		}
	}
	for code, human := range smwCodeToHuman {
		if !catalog.KnownType(human) {
			problems = append(problems, fmt.Sprintf("SMW code %q decodes to %q which PHP truth does not declare", code, human))
		}
	}

	// Scan set must not contain unreachable tables.
	for _, tbl := range dataItemTables {
		reachable := false
		for _, t := range catalog.KnownTypes() {
			if smwTypeToTable[t] == tbl {
				reachable = true
				break
			}
		}
		if !reachable {
			problems = append(problems, fmt.Sprintf("scanned table %q is not reachable from any known type", tbl))
		}
	}

	return problems
}

func inStrings(xs []string, want string) bool {
	for _, x := range xs {
		if x == want {
			return true
		}
	}
	return false
}

// TypeMatchesStoredValue reports whether a stored smw_fpt_type serialized
// value agrees with a catalog-declared type. Accepts both the SMW internal
// code (_txt) and the human name (Text) — divergence from either is a
// declaration mismatch.
func TypeMatchesStoredValue(declaredType, stored string) bool {
	stored = strings.TrimSpace(stored)
	if stored == "" {
		return false
	}
	if stored == declaredType {
		return true
	}
	// Human name in DB.
	if human, ok := smwCodeToHuman[stored]; ok && human == declaredType {
		return true
	}
	// Internal code in DB.
	if code, ok := smwTypeToCode[declaredType]; ok && code == stored {
		return true
	}
	return false
}

// ---------------------------------------------------------------------------
// Normalization (anti-pattern: Normalization Mismatch)
//
// Stored wikitext values may carry trailing CR/LF or leading whitespace from
// template formatting. The auditor applies the same trimming rules the PHP
// domain value objects apply before asserting set membership — canonical
// whitespace trimming, never partial or environment-specific normalization.
// ---------------------------------------------------------------------------

// NormalizeLiteral canonicalizes a stored literal for comparison against the
// declared allowed-value set. Exact character matching AFTER canonical trim.
func NormalizeLiteral(s string) string {
	return strings.TrimSpace(s)
}

// ---------------------------------------------------------------------------
// Scalar syntactic invariants (anti-pattern: Divergent Parsing Semantics)
//
// Both engines anchor parsing in strict formal standards:
//   - ISO 8601 for temporal instances
//   - RFC 3986 for resource identifiers
// ---------------------------------------------------------------------------

// iso8601Pattern covers the ISO 8601 calendar forms SMW emits
// (YYYY-MM-DD, YYYY-MM-DDThh:mm:ss, with optional timezone).
var iso8601Pattern = regexp.MustCompile(
	`^\d{4}-\d{2}-\d{2}([T ]\d{2}:\d{2}(:\d{2})?(\.\d+)?(Z|[+-]\d{2}:?\d{2})?)?$`,
)

// validISODate reports whether s is a strict ISO 8601 temporal string that
// parses to a real timestamp without overflow (year 0000, month 13, etc.).
func validISODate(s string) bool {
	s = strings.TrimSpace(s)
	if !iso8601Pattern.MatchString(s) {
		return false
	}
	// Normalize a space separator and accept a bare date.
	layout := s
	if strings.ContainsRune(layout, ' ') {
		layout = strings.Replace(layout, " ", "T", 1)
	}
	if !strings.ContainsAny(layout, "Tt") {
		layout += "T00:00:00Z"
	}
	// Reject year 0000 (not a valid historical timestamp).
	if len(s) >= 4 && s[:4] == "0000" {
		return false
	}
	for _, fmt := range []string{time.RFC3339, time.RFC3339Nano, "2006-01-02T15:04:05", "2006-01-02T15:04"} {
		if _, err := time.Parse(fmt, layout); err == nil {
			return true
		}
	}
	return false
}

// validURI reports whether s is an absolute URI with a valid scheme prefix
// per RFC 3986.
func validURI(s string) bool {
	s = strings.TrimSpace(s)
	u, err := url.Parse(s)
	if err != nil {
		return false
	}
	return u.IsAbs() && u.Scheme != "" && (u.Host != "" || u.Opaque != "")
}

// ---------------------------------------------------------------------------
// Set membership (Discrete Set Membership Invariant)
// ---------------------------------------------------------------------------

// InAllowedSet reports whether a normalized literal belongs to the declared
// allowed-value enumeration (exact character matching after canonical trim).
func InAllowedSet(literal string, allowed []string) bool {
	lit := NormalizeLiteral(literal)
	for _, a := range allowed {
		if NormalizeLiteral(a) == lit {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------------------------
// Catalog helpers
// ---------------------------------------------------------------------------

// EnumeratedProperties returns the subset of catalog properties carrying an
// active AllowedValues declaration, keyed by canonical name.
func EnumeratedProperties(c *catalog.Catalog) []catalog.PropertyDeclaration {
	var out []catalog.PropertyDeclaration
	for _, p := range c.Properties {
		if len(p.Allowed) > 0 {
			out = append(out, p)
		}
	}
	return out
}
