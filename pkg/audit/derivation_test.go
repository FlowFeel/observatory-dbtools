package audit

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/FlowFeel/observatory-dbtools/pkg/catalog"
)

func loadFixtureBytes(t *testing.T) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "catalog", "testdata", "catalog_v1.json"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	return data
}

// Derivation is negotiation: the auditor computes routing from declared facts
// (the property's type), never from fields on the wire. The projection
// discipline line forbids derived fields in catalog.json (see
// pkg/catalog/discipline.go), so derivation has exactly one input: the type.
func TestExpectedTableIsTypeDerived(t *testing.T) {
	checks := map[string]string{
		"Text":    "smw_di_blob",
		"Code":    "smw_di_blob",
		"Date":    "smw_di_time",
		"Number":  "smw_di_number",
		"URL":     "smw_di_uri",
		"Email":   "smw_di_uri",
		"Boolean": "smw_di_bool",
		"Page":    "smw_di_wikipage",
	}
	for typ, want := range checks {
		got, ok := ExpectedTable(typ)
		if !ok {
			t.Fatalf("ExpectedTable(%q) not routable", typ)
		}
		if got != want {
			t.Fatalf("ExpectedTable(%q) = %q, want %q", typ, got, want)
		}
	}
}

// Every type declared in the committed fixture must be routable — a type the
// auditor cannot derive a storage table for would silently skip rows. This is
// the totality contract: no declared type is left undecided.
func TestEveryFixtureTypeIsRoutable(t *testing.T) {
	rep, err := catalog.Inspect(loadFixtureBytes(t))
	if err != nil {
		t.Fatalf("inspect fixture: %v", err)
	}
	for _, p := range rep.Catalog.Properties {
		if _, ok := ExpectedTable(p.Type); !ok {
			t.Fatalf("fixture property %q has type %q with no Go derivation — silent divergence", p.Name, p.Type)
		}
	}
}

// The exclusion contract holds structurally: ExpectedTable accepts only the
// declared type string, so a hostile artifact that smuggles "table" onto a
// property (flagged by the discipline inspector as derived_field_smuggled)
// can never steer routing. The wire carries facts; Go carries the decision.
func TestWireFieldCannotInfluenceDerivation(t *testing.T) {
	if got, _ := ExpectedTable("Text"); got != "smw_di_blob" {
		t.Fatalf("ExpectedTable(Text) = %q, want smw_di_blob", got)
	}
	// Assert the discipline line rejects the smuggling vector itself.
	marker := `"type": "Text"`
	data := string(loadFixtureBytes(t))
	idx := indexOf(t, data, marker)
	hostile := []byte(data[:idx] + `"type": "Text", "table": "smw_di_time"` + data[idx+len(marker):])
	rep, err := catalog.Inspect(hostile)
	if err != nil {
		t.Fatalf("inspect hostile artifact: %v", err)
	}
	flagged := false
	for _, v := range rep.Violations {
		if v.Rule == catalog.RuleDerivedField && v.Field == "table" {
			flagged = true
		}
	}
	if !flagged {
		t.Fatalf("hostile artifact not flagged: %+v", rep.Violations)
	}
}

func indexOf(t *testing.T, s, sub string) int {
	t.Helper()
	idx := indexOfStr(s, sub)
	if idx < 0 {
		t.Fatalf("marker %q not found", sub)
	}
	return idx
}

func indexOfStr(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
