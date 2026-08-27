package catalog

import "sort"

// KnownTypes mirrors the PHP PropertyType enum — the source of truth for what
// a property may declare. A property whose type is outside this set cannot be
// audited (no routing, no validation), so the loader fails loudly at load time
// rather than silently skipping it.
//
// OWA tolerance applies to unknown *fields*, never to unknown *types*: an
// undecodable type is a semantic gap, not an open world. If PHP truth grows a
// type, this set grows with it (versioned alongside the artifact).
var knownTypes = map[string]bool{
	"Text":    true,
	"Date":    true,
	"URL":     true,
	"Page":    true,
	"Number":  true,
	"Boolean": true,
	"Email":   true,
}

// KnownType reports whether name is a declared property type per PHP truth.
func KnownType(name string) bool { return knownTypes[name] }

// KnownTypes returns the declared property types, sorted for determinism.
func KnownTypes() []string {
	out := make([]string, 0, len(knownTypes))
	for t := range knownTypes {
		out = append(out, t)
	}
	sort.Strings(out)
	return out
}
