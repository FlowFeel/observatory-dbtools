package catalog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// fixturePath returns the committed v1 catalog artifact.
func fixturePath(t *testing.T) string {
	t.Helper()
	return filepath.Join("testdata", "catalog_v1.json")
}

func loadFixture(t *testing.T) []byte {
	t.Helper()
	data, err := os.ReadFile(fixturePath(t))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	return data
}

// mutateFixture is a tiny JSON document mutator used to build hostile
// artifacts for discipline tests. It replaces the first occurrence of the
// literal marker string with the replacement. Markers are chosen to be
// unambiguous in the fixture (e.g. a full property object).
func mutateFixture(t *testing.T, marker, replacement string) []byte {
	t.Helper()
	data := string(loadFixture(t))
	idx := strings.Index(data, marker)
	if idx < 0 {
		t.Fatalf("marker %q not found in fixture", marker)
	}
	return []byte(data[:idx] + replacement + data[idx+len(marker):])
}

func mustInspect(t *testing.T, data []byte) *Report {
	t.Helper()
	rep, err := Inspect(data)
	if err != nil {
		t.Fatalf("Inspect: %v", err)
	}
	return rep
}

// --- clean fixture ---------------------------------------------------------

func TestInspectCleanFixturePasses(t *testing.T) {
	rep := mustInspect(t, loadFixture(t))
	if len(rep.Violations) != 0 {
		t.Fatalf("clean fixture must have zero discipline violations, got: %+v", rep.Violations)
	}
	if len(rep.Warnings) != 0 {
		t.Fatalf("clean fixture must have zero OWA warnings, got: %+v", rep.Warnings)
	}
	if rep.Catalog == nil || rep.Catalog.Version != 1 {
		t.Fatalf("inspected catalog missing or wrong version: %+v", rep.Catalog)
	}
	if len(rep.Catalog.Properties) != 83 {
		t.Fatalf("expected 83 properties, got %d", len(rep.Catalog.Properties))
	}
	if len(rep.Catalog.Entities) != 5 {
		t.Fatalf("expected 5 entities, got %d", len(rep.Catalog.Entities))
	}
}

// --- derived-field exclusion -----------------------------------------------

func TestInspectSmuggledDerivedFieldFlagged(t *testing.T) {
	// Smuggle a routing table into a property entry. The discipline line
	// refuses it: consumers derive routing from declared type, never the wire.
	marker := `"type": "Text"`
	replacement := `"type": "Text", "table": "smw_di_time"`
	data := mutateFixture(t, marker, replacement)

	rep := mustInspect(t, data)
	found := false
	for _, v := range rep.Violations {
		if v.Rule == RuleDerivedField && v.Field == "table" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected derived-field violation for %q, got: %+v", "table", rep.Violations)
	}
}

func TestInspectSmuggledSMWCodeFlagged(t *testing.T) {
	marker := `"type": "Date"`
	replacement := `"type": "Date", "smw_code": "_dat"`
	data := mutateFixture(t, marker, replacement)

	rep := mustInspect(t, data)
	for _, v := range rep.Violations {
		if v.Rule == RuleDerivedField && v.Field == "smw_code" {
			return
		}
	}
	t.Fatalf("expected derived-field violation for smw_code, got: %+v", rep.Violations)
}

// --- OWA tolerance ---------------------------------------------------------

func TestInspectUnknownFieldIsWarningNotViolation(t *testing.T) {
	// Inject an unknown top-level field. OWA: tolerated, warned, not failed.
	data := mutateFixture(t, `"version": 1`, `"version": 1, "extensions": {}`)

	rep := mustInspect(t, data)
	if len(rep.Violations) != 0 {
		t.Fatalf("unknown field must not be a violation, got: %+v", rep.Violations)
	}
	warned := false
	for _, w := range rep.Warnings {
		if strings.Contains(w, "extensions") {
			warned = true
		}
	}
	if !warned {
		t.Fatalf("expected OWA warning mentioning extensions, got: %+v", rep.Warnings)
	}
}

// --- nesting depth ---------------------------------------------------------

func TestInspectNestedAllowedValuesRejected(t *testing.T) {
	// Replace an existing flat allowed-value list (null here) with nested
	// objects — the JSON-hell shape. Avoids duplicate-key ambiguity.
	marker := `"allowed": null,`
	replacement := `"allowed": [{"value": "A"}],`
	data := mutateFixture(t, marker, replacement)

	rep := mustInspect(t, data)
	found := false
	for _, v := range rep.Violations {
		if v.Rule == RuleNestingDepth {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected nesting-depth violation, got: %+v", rep.Violations)
	}
}

func TestInspectNestedListRejected(t *testing.T) {
	// A list inside a list is equally rejected.
	marker := `"aliases": []`
	replacement := `"aliases": [["a"]]`
	data := mutateFixture(t, marker, replacement)

	rep := mustInspect(t, data)
	for _, v := range rep.Violations {
		if v.Rule == RuleNestingDepth {
			return
		}
	}
	t.Fatalf("expected nesting-depth violation for nested list, got: %+v", rep.Violations)
}

func TestInspectNestedEntityPropertiesRejected(t *testing.T) {
	// Entity property lists are pretty-printed with newlines; match the raw
	// text and smuggle an object as the first element.
	marker := "\"properties\": [\n                \"Area\""
	replacement := "\"properties\": [\n                {\"name\": \"Area\"}"
	data := mutateFixture(t, marker, replacement)

	rep := mustInspect(t, data)
	for _, v := range rep.Violations {
		if v.Rule == RuleNestingDepth {
			return
		}
	}
	t.Fatalf("expected nesting-depth violation for entity properties, got: %+v", rep.Violations)
}

// --- fixture remains loadable as a Catalog ---------------------------------

func TestInspectPreservesCatalogSemantics(t *testing.T) {
	data := loadFixture(t)
	rep := mustInspect(t, data)
	c := rep.Catalog
	if c.PropertyByName("Event type") == nil {
		t.Fatal("expected Event type property in inspected catalog")
	}
	if c.EntityByName("Area") == nil {
		t.Fatal("expected Area entity in inspected catalog")
	}
}
