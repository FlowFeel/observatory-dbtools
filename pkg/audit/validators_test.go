package audit

import (
	"testing"
)

func TestExpectedTable(t *testing.T) {
	cases := map[string]string{
		"Text": "smw_di_blob",
		"Date": "smw_di_time", "Number": "smw_di_number",
		"URL": "smw_di_uri", "Email": "smw_di_uri",
		"Boolean": "smw_di_bool", "Page": "smw_di_wikipage",
	}
	for typ, want := range cases {
		got, ok := ExpectedTable(typ)
		if !ok {
			t.Errorf("ExpectedTable(%q) not routable", typ)
		}
		if got != want {
			t.Errorf("ExpectedTable(%q) = %q, want %q", typ, got, want)
		}
	}
	if _, ok := ExpectedTable("Bogus"); ok {
		t.Error("ExpectedTable(Bogus) should not be routable")
	}
}

func TestTypeMatchesStoredValue(t *testing.T) {
	// Exact human name.
	if !TypeMatchesStoredValue("Date", "Date") {
		t.Error("Date == Date should match")
	}
	// SMW internal code.
	if !TypeMatchesStoredValue("Date", "_dat") {
		t.Error("Date == _dat should match")
	}
	if !TypeMatchesStoredValue("Text", "_txt") {
		t.Error("Text == _txt should match")
	}
	if !TypeMatchesStoredValue("Page", "_wpg") {
		t.Error("Page == _wpg should match")
	}
	// Mismatch.
	if TypeMatchesStoredValue("Date", "Text") {
		t.Error("Date != Text should not match")
	}
	if TypeMatchesStoredValue("Date", "_txt") {
		t.Error("Date != _txt should not match")
	}
	// Empty.
	if TypeMatchesStoredValue("Date", "") {
		t.Error("empty stored value should not match")
	}
}

func TestNormalizeLiteral(t *testing.T) {
	if got := NormalizeLiteral("  Webinar  "); got != "Webinar" {
		t.Errorf("NormalizeLiteral = %q, want Webinar", got)
	}
	if got := NormalizeLiteral("In-Person\r\n"); got != "In-Person" {
		t.Errorf("NormalizeLiteral = %q, want In-Person", got)
	}
	if got := NormalizeLiteral(" Hybrid "); got != "Hybrid" {
		t.Errorf("NormalizeLiteral = %q, want Hybrid", got)
	}
}

func TestValidISODate(t *testing.T) {
	valid := []string{
		"2026-08-27",
		"2026-08-27T14:00:00",
		"2026-08-27T14:00:00Z",
		"2026-08-27T14:00:00+00:00",
		"2026-08-27 14:00:00",
		"1999-12-31T23:59:59.999",
	}
	for _, s := range valid {
		if !validISODate(s) {
			t.Errorf("validISODate(%q) = false, want true", s)
		}
	}
	invalid := []string{
		"", "August 27, 2026", "2026-13-01", "2026-08-32",
		"0000-01-01", "not a date", "27/08/2026", "2026-08-27T25:00:00",
	}
	for _, s := range invalid {
		if validISODate(s) {
			t.Errorf("validISODate(%q) = true, want false", s)
		}
	}
}

func TestValidURI(t *testing.T) {
	valid := []string{
		"https://schema.org/startDate",
		"https://observatory.wiki/w/index.php?title=Main_Page",
		"mailto:flow@ind.media",
		"http://example.com",
	}
	for _, s := range valid {
		if !validURI(s) {
			t.Errorf("validURI(%q) = false, want true", s)
		}
	}
	invalid := []string{
		"", "startDate", "//schema.org/startDate", "schema.org/startDate",
		"https://", "foo bar",
	}
	for _, s := range invalid {
		if validURI(s) {
			t.Errorf("validURI(%q) = true, want false", s)
		}
	}
}

func TestInAllowedSet(t *testing.T) {
	allowed := []string{"In-Person", "Virtual", "Hybrid"}
	if !InAllowedSet("In-Person", allowed) {
		t.Error("In-Person should be in set")
	}
	if !InAllowedSet("  Virtual  ", allowed) {
		t.Error("trimmed Virtual should be in set")
	}
	if InAllowedSet("Webinar", allowed) {
		t.Error("Webinar should NOT be in set")
	}
	if InAllowedSet("in-person", allowed) {
		t.Error("case-sensitive match should fail")
	}
}

func TestReportAccumulation(t *testing.T) {
	r := &Report{}
	r.Add(Violation{Kind: KindRange, Severity: SevError, Property: "Event type"})
	r.Add(Violation{Kind: KindRange, Severity: SevError, Property: "Event type"})
	r.Add(Violation{Kind: KindSyntax, Severity: SevError, Property: "Event start date"})
	r.Add(Violation{Kind: KindDeclaration, Severity: SevWarning})
	r.Finalize()

	if r.Summary.ErrorCount != 3 {
		t.Errorf("ErrorCount = %d, want 3", r.Summary.ErrorCount)
	}
	if r.Summary.WarningCount != 1 {
		t.Errorf("WarningCount = %d, want 1", r.Summary.WarningCount)
	}
	if r.Summary.Pass {
		t.Error("report with errors should not pass")
	}
	if r.Categories.RangeCount != 2 || r.Categories.SyntaxCount != 1 || r.Categories.DeclarationCount != 1 {
		t.Errorf("category counts wrong: %+v", r.Categories)
	}
	if r.Summary.PerProperty["Event type"] != 2 {
		t.Errorf("per-property count = %d, want 2", r.Summary.PerProperty["Event type"])
	}
}
