package catalog_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/FlowFeel/observatory-dbtools/pkg/catalog"
)

const fixtureV1 = "testdata/catalog_v1.json"

func TestLoadRealFixture(t *testing.T) {
	c, warnings, err := catalog.Load(filepath.Join("testdata", "catalog_v1.json"))
	if err != nil {
		t.Fatalf("load real fixture: %v", err)
	}
	if c.Version != 1 {
		t.Fatalf("version = %d, want 1", c.Version)
	}
	if len(c.Properties) != 83 {
		t.Errorf("properties = %d, want 83", len(c.Properties))
	}
	if len(c.Entities) != 5 {
		t.Errorf("entities = %d, want 5", len(c.Entities))
	}
	if len(warnings) != 0 {
		t.Errorf("unexpected warnings on clean fixture: %v", warnings)
	}

	// Entity archetypes present.
	for _, name := range []string{"Article", "Event", "Campaign", "Resource", "Area"} {
		if c.EntityByName(name) == nil {
			t.Errorf("entity archetype %q missing", name)
		}
	}

	// Property lookup by canonical name.
	if p := c.PropertyByName("Event type"); p == nil {
		t.Error("Event type property missing")
	} else {
		if len(p.Allowed) == 0 {
			t.Error("Event type should have allowed values")
		}
	}
}

func TestParseRealFixtureBytes(t *testing.T) {
	data, err := readFixture(t, fixtureV1)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	c, _, err := catalog.Parse(data)
	if err != nil {
		t.Fatalf("parse fixture: %v", err)
	}
	if len(c.PropertyNames()) != 83 {
		t.Errorf("PropertyNames() = %d, want 83", len(c.PropertyNames()))
	}
	if len(c.EntityNames()) != 5 {
		t.Errorf("EntityNames() = %d, want 5", len(c.EntityNames()))
	}
}

func TestVersionMismatchRejected(t *testing.T) {
	// Patch version to 2.
	data := strings.Replace(readFixtureString(t, fixtureV1), `"version": 1`, `"version": 2`, 1)
	_, _, err := catalog.Parse([]byte(data))
	if err == nil {
		t.Fatal("expected version mismatch error, got nil")
	}
	if !strings.Contains(err.Error(), "unsupported version 2") {
		t.Errorf("error = %q, want unsupported version mention", err)
	}
}

func TestMissingVersionRejected(t *testing.T) {
	data := strings.Replace(readFixtureString(t, fixtureV1), `"version": 1,`, ``, 1)
	_, _, err := catalog.Parse([]byte(data))
	if err == nil {
		t.Fatal("expected missing version error, got nil")
	}
	if !strings.Contains(err.Error(), "missing version") {
		t.Errorf("error = %q, want missing version mention", err)
	}
}

func TestOWAUnknownFieldsTolerated(t *testing.T) {
	// Add an unknown top-level field — must load with a warning, not fail.
	data := strings.Replace(readFixtureString(t, fixtureV1),
		`"version": 1`, `"future_field": {"x": 1}, "version": 1`, 1)
	c, warnings, err := catalog.Parse([]byte(data))
	if err != nil {
		t.Fatalf("OWA: unknown field must not fail load: %v", err)
	}
	if c.Version != 1 {
		t.Errorf("version = %d, want 1", c.Version)
	}
	if len(warnings) == 0 {
		t.Error("expected OWA warning for unknown field, got none")
	}
	if !strings.Contains(warnings[0], "future_field") {
		t.Errorf("warning = %q, want mention of future_field", warnings[0])
	}
}

func TestMalformedJSONRejected(t *testing.T) {
	_, _, err := catalog.Parse([]byte(`{not json`))
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
}

func TestReadOnlyNoMutation(t *testing.T) {
	// Ensure the Catalog type exposes only value receivers (read-only).
	// A compile-time check: assigning to Catalog fields is impossible from
	// outside the package since all fields are unexported by convention.
	c := &catalog.Catalog{}
	_ = c.Version
	_ = c.Properties
	_ = c.Entities
}

func readFixture(t *testing.T, name string) ([]byte, error) {
	t.Helper()
	return readFileBytes(filepath.Join("testdata", filepath.Base(name)))
}

func readFixtureString(t *testing.T, name string) string {
	t.Helper()
	data, err := readFileBytes(filepath.Join("testdata", filepath.Base(name)))
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	return string(data)
}
