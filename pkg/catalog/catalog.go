// Package catalog loads and parses catalog.json — the portable
// serialized contract artifact produced by the PHP semantic domain
// (bin/export-catalog.php in observatory-magazine-v2).
//
// The catalog is the shared contract surface between the PHP Producer
// and the Go Auditor. It is a stateless, read-only in-memory knowledge
// index: the Go audit engine reads property declarations and entity
// compositions from here and verifies the MySQL SMW persistence layer
// against them (Contracts 1-3 in pkg/audit).
//
// Design notes:
//   - Versioned artifact: Load validates the major version header before
//     exposing the domain schema. Rejects incompatible major versions.
//   - OWA compliance: unknown JSON fields are tolerated and reported as
//     warnings, never as hard failures (phos.arch/owa).
//   - Read-only: Catalog exposes only getters; no mutation methods.
package catalog

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// CurrentVersion is the supported major version of the catalog artifact.
const CurrentVersion = 1

// PropertyDeclaration is a single sovereign property declaration.
type PropertyDeclaration struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Allowed    []string `json:"allowed"`
	Equivalent string   `json:"equivalent"`
	Aliases    []string `json:"aliases"`
}

// EntityDefinition is a structural composition referencing canonical
// property names. Entities do not own declarations — predicate
// sovereignty is preserved.
type EntityDefinition struct {
	Name       string   `json:"name"`
	Gloss      string   `json:"gloss"`
	Properties []string `json:"properties"`
}

// Catalog is the in-memory knowledge index.
type Catalog struct {
	Version    int                   `json:"version"`
	Properties []PropertyDeclaration `json:"properties"`
	Entities   []EntityDefinition    `json:"entities"`
}

// Warnings collects OWA tolerance notes encountered during Load. A non-nil
// result always loads; warnings are informational.
type Warnings []string

// Load reads and parses a catalog.json artifact file.
//
// Returns an error if the file cannot be read, is malformed JSON, or
// carries an unsupported major version. Unknown JSON fields are preserved
// as warnings per OWA — they do not prevent a successful load.
func Load(path string) (*Catalog, Warnings, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("catalog: read %s: %w", path, err)
	}
	return Parse(data)
}

// Parse parses catalog.json bytes into a Catalog, validating the version
// header. See Load for error semantics.
func Parse(data []byte) (*Catalog, Warnings, error) {
	var raw struct {
		Version    *int                  `json:"version"`
		Properties []PropertyDeclaration `json:"properties"`
		Entities   []EntityDefinition    `json:"entities"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, nil, fmt.Errorf("catalog: parse json: %w", err)
	}

	// Capture unknown top-level fields for OWA warning reporting.
	var envelope map[string]json.RawMessage
	if err := json.Unmarshal(data, &envelope); err != nil {
		return nil, nil, fmt.Errorf("catalog: parse envelope: %w", err)
	}
	var warnings Warnings
	known := map[string]bool{"version": true, "properties": true, "entities": true}
	for key := range envelope {
		if !known[key] {
			warnings = append(warnings, fmt.Sprintf("unknown top-level field %q ignored (OWA)", key))
		}
	}

	if raw.Version == nil {
		return nil, nil, fmt.Errorf("catalog: missing version header (expected major version %d)", CurrentVersion)
	}
	if *raw.Version != CurrentVersion {
		return nil, nil, fmt.Errorf(
			"catalog: unsupported version %d (this loader supports major version %d)",
			*raw.Version, CurrentVersion,
		)
	}

	// Load-time cross-check: every declared type must be known to the Go
	// consumer. An unknown type would silently skip audit rows, so it is a
	// hard failure, not an OWA warning (see types.go).
	for _, p := range raw.Properties {
		if !KnownType(p.Type) {
			return nil, nil, fmt.Errorf(
				"catalog: property %q declares unknown type %q (known: %s)",
				p.Name, p.Type, strings.Join(KnownTypes(), ", "),
			)
		}
	}

	return &Catalog{
		Version:    *raw.Version,
		Properties: raw.Properties,
		Entities:   raw.Entities,
	}, warnings, nil
}

// PropertyByName returns the declaration for a canonical property name,
// or nil when absent.
func (c *Catalog) PropertyByName(name string) *PropertyDeclaration {
	for i := range c.Properties {
		if c.Properties[i].Name == name {
			return &c.Properties[i]
		}
	}
	return nil
}

// PropertyNames returns the canonical names of all declared properties.
func (c *Catalog) PropertyNames() []string {
	names := make([]string, 0, len(c.Properties))
	for i := range c.Properties {
		names = append(names, c.Properties[i].Name)
	}
	return names
}

// SMWTitle returns the canonical SMW storage title for a property name.
// SMW stores property page titles in smw_object_ids with underscores
// instead of spaces (e.g. "Event type" → "Event_type"). Both the audit
// queries and the BDD seeds must use this canonical form to avoid the
// divergent-parsing-semantics anti-pattern.
func SMWTitle(name string) string {
	return strings.ReplaceAll(name, " ", "_")
}

// EntityByName returns the entity definition for a name, or nil when absent.
func (c *Catalog) EntityByName(name string) *EntityDefinition {
	for i := range c.Entities {
		if c.Entities[i].Name == name {
			return &c.Entities[i]
		}
	}
	return nil
}

// EntityNames returns the names of all declared entity archetypes.
func (c *Catalog) EntityNames() []string {
	names := make([]string, 0, len(c.Entities))
	for i := range c.Entities {
		names = append(names, c.Entities[i].Name)
	}
	return names
}
