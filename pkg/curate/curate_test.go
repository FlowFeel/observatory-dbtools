package curate

import (
	"testing"
)

func TestNewPlan(t *testing.T) {
	plan, err := NewPlan("test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(plan.RequiredPages) != 5 {
		t.Errorf("expected 5 required pages, got %d", len(plan.RequiredPages))
	}
	if !plan.AnonymizeUsers {
		t.Errorf("expected AnonymizeUsers to be true for test tier")
	}

	_, err = NewPlan("invalid-tier")
	if err == nil {
		t.Errorf("expected error for invalid tier")
	}
}

func TestValidateSeed(t *testing.T) {
	plan, _ := NewPlan("test")

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
