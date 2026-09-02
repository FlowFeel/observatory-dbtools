package catalog

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ---------------------------------------------------------------------------
// Projection Discipline
//
// catalog.json is a projection surface, not a record of truth. Truth lives in
// the PHP semantic domain; the artifact serializes declared facts only. This
// inspector enforces the projection line at load time:
//
//   1. Flatness — entries are flat records. Objects may appear only at the
//      artifact root and at the entry level; arrays only as direct children of
//      the artifact (properties, entities) or of an entry (allowed, aliases,
//      properties). No nested records, no nested lists.
//   2. Exclusion — derived fields (routing tables, SMW internal codes,
//      canonical title forms, normalized values) are negotiation knowledge and
//      must never be serialized. Consumers derive them locally from declared
//      facts. The wire never carries decisions.
//   3. Open-World Assumption — unknown non-derived fields are tolerated and
//      surfaced as warnings, never as hard failures.
// ---------------------------------------------------------------------------

// Rule identifiers for discipline infractions.
const (
	RuleDerivedField = "derived_field_smuggled"
	RuleNestingDepth = "nesting_depth_exceeded"
	RuleUnknownField = "unknown_field"
	RuleUnknownType  = "unknown_type"
)

// derivedFields are negotiation knowledge that must never cross the wire.
// Consumers derive these locally from declared facts.
var derivedFields = map[string]bool{
	"table":           true,
	"smw_table":       true,
	"storage":         true,
	"smw_code":        true,
	"type_code":       true,
	"canonical_title": true,
	"smw_title":       true,
	"sortkey":         true,
	"p_id":            true,
	"normalized":      true,
}

// allowedTopLevel are the v1 artifact-level keys.
var allowedTopLevel = map[string]bool{
	"version":    true,
	"properties": true,
	"entities":   true,
}

// allowedPropertyKeys are the v1 property-entry keys (declared facts only).
var allowedPropertyKeys = map[string]bool{
	"name":       true,
	"type":       true,
	"allowed":    true,
	"equivalent": true,
	"aliases":    true,
}

// allowedEntityKeys are the v1 entity-entry keys (declared facts only).
var allowedEntityKeys = map[string]bool{
	"name":       true,
	"gloss":      true,
	"properties": true,
}

// Violation is a single projection-discipline infraction.
type Violation struct {
	Path       string
	Field      string
	Rule       string
	Diagnostic string
}

// Report is the result of a discipline inspection.
type Report struct {
	Catalog    *Catalog
	Warnings   Warnings
	Violations []Violation
}

// Inspect parses an artifact and enforces the projection-discipline line.
//
// The structural walk runs on the raw JSON independently of the typed parse,
// so shape violations (nested structures in scalar lists, non-flat entries)
// are reported as discipline violations rather than masked as JSON errors.
// The typed Catalog view is populated when the artifact is representable;
// for shape-violating artifacts it remains nil — the discipline report is
// the diagnosis. Errors are reserved for invalid JSON, a missing version
// header, or an unsupported major version.
func Inspect(data []byte) (*Report, error) {
	var raw any
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("catalog: inspect: %w", err)
	}

	root, ok := raw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("catalog: inspect: artifact root must be an object")
	}

	rep := &Report{}
	walkRoot(root, "$", rep)

	// Version header enforcement (mirrors Parse).
	var head struct {
		Version *int `json:"version"`
	}
	if err := json.Unmarshal(data, &head); err != nil {
		return nil, fmt.Errorf("catalog: inspect: %w", err)
	}
	if head.Version == nil {
		return nil, fmt.Errorf("catalog: inspect: missing version header (expected major version %d)", CurrentVersion)
	}
	if *head.Version != CurrentVersion {
		return nil, fmt.Errorf(
			"catalog: inspect: unsupported version %d (this loader supports major version %d)",
			*head.Version, CurrentVersion,
		)
	}

	// Typed catalog view; shape violations were already surfaced by the
	// walker, so a failed typed parse is tolerable (Catalog stays nil).
	if c, _, err := Parse(data); err == nil {
		rep.Catalog = c
	}
	return rep, nil
}

// walkRoot inspects the artifact envelope: version (scalar) and the two
// flat entry lists.
func walkRoot(obj map[string]any, path string, rep *Report) {
	for k, v := range obj {
		child := path + "." + k
		switch k {
		case "properties", "entities":
			arr, ok := v.([]any)
			if !ok {
				rep.Violations = append(rep.Violations, Violation{
					Path: child, Field: k, Rule: RuleNestingDepth,
					Diagnostic: fmt.Sprintf("%q must be a flat list of entry records", k),
				})
				continue
			}
			walkEntryList(arr, child, k == "properties", rep)
		default:
			if derivedFields[k] {
				rep.Violations = append(rep.Violations, Violation{
					Path: child, Field: k, Rule: RuleDerivedField,
					Diagnostic: "derived field must never be serialized; consumers derive it locally from declared facts",
				})
			} else if !allowedTopLevel[k] {
				rep.Warnings = append(rep.Warnings, fmt.Sprintf(
					"projection discipline: unknown top-level field %q tolerated (OWA)", k))
			}
		}
	}
}

// walkEntryList inspects the top-level entry arrays (properties, entities).
func walkEntryList(arr []any, path string, isProperty bool, rep *Report) {
	for i, item := range arr {
		child := fmt.Sprintf("%s[%d]", path, i)
		obj, ok := item.(map[string]any)
		if !ok {
			rep.Violations = append(rep.Violations, Violation{
				Path: child, Rule: RuleNestingDepth,
				Diagnostic: "entry must be a flat record object",
			})
			continue
		}
		walkEntry(obj, child, isProperty, rep)
	}
}

// walkEntry inspects a single property/entity entry record. Entries must be
// flat: scalar fields or scalar lists, never nested records or nested lists.
func walkEntry(obj map[string]any, path string, isProperty bool, rep *Report) {
	allowed := allowedEntityKeys
	if isProperty {
		allowed = allowedPropertyKeys
	}
	for k, v := range obj {
		child := path + "." + k
		switch tv := v.(type) {
		case map[string]any:
			if derivedFields[k] {
				rep.Violations = append(rep.Violations, Violation{
					Path: child, Field: k, Rule: RuleDerivedField,
					Diagnostic: "derived field must never be serialized; consumers derive it locally from declared facts",
				})
			} else {
				rep.Violations = append(rep.Violations, Violation{
					Path: child, Field: k, Rule: RuleNestingDepth,
					Diagnostic: "entry fields must be flat scalars or scalar lists; found a nested record",
				})
			}
		case []any:
			if derivedFields[k] {
				rep.Violations = append(rep.Violations, Violation{
					Path: child, Field: k, Rule: RuleDerivedField,
					Diagnostic: "derived field must never be serialized; consumers derive it locally from declared facts",
				})
				continue
			}
			walkScalarList(tv, child, rep)
		default:
			// A property's declared type is a fact the consumer must be able to
			// decode. An unknown type is a semantic gap — diagnosed as a violation
			// here, and a hard load failure in Parse.
			if isProperty && k == "type" {
				if s, ok := v.(string); ok && !KnownType(s) {
					rep.Violations = append(rep.Violations, Violation{
						Path: child, Field: k, Rule: RuleUnknownType,
						Diagnostic: fmt.Sprintf("property declares unknown type %q (known: %s)", s, strings.Join(KnownTypes(), ", ")),
					})
				}
				continue
			}
			if derivedFields[k] {
				rep.Violations = append(rep.Violations, Violation{
					Path: child, Field: k, Rule: RuleDerivedField,
					Diagnostic: "derived field must never be serialized; consumers derive it locally from declared facts",
				})
			} else if !allowed[k] {
				rep.Warnings = append(rep.Warnings, fmt.Sprintf(
					"projection discipline: unknown %s-entry field %q tolerated (OWA)", kindLabel(isProperty), k))
			}
		}
	}
}

// walkScalarList inspects a scalar list (allowed, aliases, entity properties).
// Every element must be a scalar — nested objects or arrays are the JSON-hell
// shape the discipline line refuses.
func walkScalarList(arr []any, path string, rep *Report) {
	for i, item := range arr {
		switch item.(type) {
		case map[string]any, []any:
			rep.Violations = append(rep.Violations, Violation{
				Path: fmt.Sprintf("%s[%d]", path, i), Rule: RuleNestingDepth,
				Diagnostic: "scalar list must contain only scalar values; nested structures are rejected",
			})
		}
	}
}

func kindLabel(isProperty bool) string {
	if isProperty {
		return "property"
	}
	return "entity"
}

// IsDerivedField reports whether name is a derived field the projection must
// refuse to carry. Consumers derive these locally from declared facts.
func IsDerivedField(name string) bool {
	return derivedFields[name]
}

// TitleCollisions returns pairs of distinct property names whose SMWTitle
// forms collide (same canonical storage title). A collision makes identity
// ambiguous in smw_object_ids and must never silently pass.
func (c *Catalog) TitleCollisions() [][2]string {
	seen := map[string]string{}
	var out [][2]string
	for _, p := range c.Properties {
		title := SMWTitle(p.Name)
		if prev, ok := seen[title]; ok {
			out = append(out, [2]string{prev, p.Name})
		} else {
			seen[title] = p.Name
		}
	}
	return out
}
