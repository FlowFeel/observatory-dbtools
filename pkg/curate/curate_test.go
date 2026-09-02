package curate

import (
	"testing"

	"github.com/FlowFeel/observatory-dbtools/pkg/catalog"
)

func TestNewPlan(t *testing.T) {
	plan, err := NewPlan("test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(plan.RequiredPages) != 5 {
		t.Errorf("expected 5 required pages, got %d", len(plan.RequiredPages))
	}
	if !plan.AnonymizeUsers {
		t.Errorf("expected AnonymizeUsers to be true for test tier")
	}

	_, err = NewPlan("invalid-tier", nil)
	if err == nil {
		t.Errorf("expected error for invalid tier")
	}
}

func TestNewPlanCatalogDriven(t *testing.T) {
	// A loaded catalog drives property page generation (no hardcoded list).
	c := &catalog.Catalog{
		Version: 1,
		Properties: []catalog.PropertyDeclaration{
			{Name: "Event type", Type: "Text", Allowed: []string{"In-Person", "Virtual"}},
			{Name: "Author", Type: "Text"},
		},
	}
	plan, err := NewPlan("prod", c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"Property:Event_type", "Property:Author"}
	for i, w := range want {
		if plan.RequiredProperties[i] != w {
			t.Errorf("RequiredProperties[%d] = %q, want %q", i, plan.RequiredProperties[i], w)
		}
	}
}

func TestValidateSeed(t *testing.T) {
	plan, _ := NewPlan("test", nil)

	validSeed := `
INSERT INTO page (page_id, page_title) VALUES
(1, 'Main_Page'),
(2, 'Animals'),
(3, 'Classics'),
(4, 'Dig_Labs'),
(5, 'Human_Bridges');
`
	if err := ValidateSeed(validSeed, plan); err != nil {
		t.Errorf("expected valid seed to pass, got error: %v", err)
	}

	invalidSeed := `
INSERT INTO page (page_id, page_title) VALUES
(1, 'Main_Page');
`
	if err := ValidateSeed(invalidSeed, plan); err == nil {
		t.Errorf("expected invalid seed to fail validation")
	}
}
