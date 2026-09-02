package catalog

import (
	"strings"
	"testing"
)

// --- KnownTypes (PHP PropertyType truth) ----------------------------------

func TestKnownType(t *testing.T) {
	for _, k := range []string{"Text", "Date", "URL", "Page", "Number", "Boolean", "Email"} {
		if !KnownType(k) {
			t.Errorf("KnownType(%q) = false, want true", k)
		}
	}
	for _, k := range []string{"Geo", "Enum", "Widget", "Code"} {
		if KnownType(k) {
			t.Errorf("KnownType(%q) = true, want false (not in PHP truth)", k)
		}
	}
}

func TestKnownTypesSorted(t *testing.T) {
	ks := KnownTypes()
	for i := 1; i < len(ks); i++ {
		if ks[i-1] > ks[i] {
			t.Fatalf("KnownTypes not sorted: %v", ks)
		}
	}
}

// --- Load-time cross-check (T-547) ----------------------------------------

func TestParseRejectsUnknownType(t *testing.T) {
	marker := `"type": "Text"`
	data := mutateFixture(t, marker, `"type": "Geo"`)

	_, _, err := Parse(data)
	if err == nil {
		t.Fatal("Parse must fail on an unknown type")
	}
	if !strings.Contains(err.Error(), "Geo") {
		t.Fatalf("error must identify the unknown type: %v", err)
	}
	if !strings.Contains(err.Error(), "Author") {
		t.Fatalf("error must identify the property name (first Text property is Author): %v", err)
	}
}

func TestInspectReportsUnknownTypeViolation(t *testing.T) {
	marker := `"type": "Text"`
	data := mutateFixture(t, marker, `"type": "Geo"`)

	rep := mustInspect(t, data)
	for _, v := range rep.Violations {
		if v.Rule == RuleUnknownType && strings.Contains(v.Diagnostic, "Geo") {
			return
		}
	}
	t.Fatalf("expected unknown-type violation, got: %+v", rep.Violations)
}

func TestUnknownTypeIsNotAnOWAWarning(t *testing.T) {
	marker := `"type": "Text"`
	data := mutateFixture(t, marker, `"type": "Geo"`)

	rep := mustInspect(t, data)
	for _, w := range rep.Warnings {
		if strings.Contains(w, "Geo") {
			t.Fatalf("unknown type must be a hard violation, not an OWA warning: %v", w)
		}
	}
	hasViolation := false
	for _, v := range rep.Violations {
		if v.Rule == RuleUnknownType {
			hasViolation = true
		}
	}
	if !hasViolation {
		t.Fatalf("expected unknown-type violation, got warnings only: %+v", rep.Warnings)
	}
}

// --- Title collision check (T-546) ----------------------------------------

func TestNoTitleCollisionsInFixture(t *testing.T) {
	c := mustInspect(t, loadFixture(t)).Catalog
	if c == nil {
		t.Fatal("expected catalog from clean fixture")
	}
	if cols := c.TitleCollisions(); len(cols) > 0 {
		t.Fatalf("fixture has SMWTitle collisions: %v", cols)
	}
}

func TestTitleCollisionsDetected(t *testing.T) {
	c := &Catalog{
		Version: 1,
		Properties: []PropertyDeclaration{
			{Name: "A B", Type: "Text"},
			{Name: "A_B", Type: "Text"},
		},
	}
	cols := c.TitleCollisions()
	if len(cols) != 1 {
		t.Fatalf("expected 1 collision, got %d: %v", len(cols), cols)
	}
}
