package audit

import (
	"testing"

	"github.com/FlowFeel/observatory-dbtools/pkg/catalog"
)

// RoutingTotality is the totality contract (T-546): the derivation tables are
// total over PHP truth and free of truth-divergent entries. A non-empty result
// is a defect, not a warning.
func TestRoutingTotalityClean(t *testing.T) {
	if problems := RoutingTotality(); len(problems) > 0 {
		t.Fatalf("routing derivation is not total: %+v", problems)
	}
}

// No routed type may exist outside the KnownTypes set (PHP truth). This is the
// regression guard for the removed "Code" divergence — Go had dead routing for
// a type the PHP PropertyType enum never declared.
func TestRoutingTotalityForbidsTruthDivergentTypes(t *testing.T) {
	if catalog.KnownType("Code") {
		t.Skip("Code is declared truth; not a divergence")
	}
	// Simulate re-introducing the divergence.
	orig := smwTypeToTable
	defer func() { smwTypeToTable = orig }()
	smwTypeToTable = map[string]string{"Code": "smw_di_blob"}

	problems := RoutingTotality()
	for _, p := range problems {
		if containsStr(p, "Code") && containsStr(p, "does not declare") {
			return
		}
	}
	t.Fatalf("expected Code routed-but-not-declared to be flagged, got: %+v", problems)
}

// Every KnownType must round-trip through its SMW code (type -> code -> type).
func TestRoutingTotalityRoundTrips(t *testing.T) {
	for _, k := range catalog.KnownTypes() {
		code, ok := smwTypeToCode[k]
		if !ok {
			t.Fatalf("known type %q has no SMW code", k)
		}
		human, ok := smwCodeToHuman[code]
		if !ok || human != k {
			t.Fatalf("type %q does not round-trip through code %q (got %q)", k, code, human)
		}
	}
}

func containsStr(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
