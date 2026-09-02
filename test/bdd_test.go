package test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/FlowFeel/observatory-dbtools/pkg/audit"
	"github.com/FlowFeel/observatory-dbtools/pkg/baseline"
	"github.com/FlowFeel/observatory-dbtools/pkg/catalog"
	"github.com/FlowFeel/observatory-dbtools/pkg/connect"
	"github.com/FlowFeel/observatory-dbtools/pkg/drift"
	"github.com/cucumber/godog"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

type bddSuite struct {
	container  *mysql.MySQLContainer
	db         *sql.DB
	driftRep   *drift.Report
	fixedCount int64

	// semantic audit state (T-538)
	catalog        *catalog.Catalog
	declarationRep *audit.Report
	valueTypeRep   *audit.Report
	valueRangeRep  *audit.Report
	pIDByProperty  map[string]int

	// projection discipline state (T-541)
	disciplineData   []byte
	disciplineReport *catalog.Report

	// derivation totality state (T-546)
	knownTypes []string
	totality   []string

	// load-time cross-check state (T-547)
	lastLoadErr error

	// currentProp is the property named by the most recent Given step
	// (e.g. "a compiled property catalog with property \"Event type\"").
	// Steps phrased "for that property" resolve the property from here.
	currentProp string
}

// propertyPageNS is the MediaWiki namespace for Property pages.
const propertyPageNS = 102

func (s *bddSuite) aCleanBaselineDatabaseIsLoaded() error {
	schemaPath := filepath.Join("..", "baselines", "schema.sql")
	if err := baseline.Import(s.db, schemaPath); err != nil {
		return fmt.Errorf("import baseline schema: %w", err)
	}
	return nil
}

func (s *bddSuite) theFollowingCoreTablesMustExist(tables *godog.Table) error {
	var req []string
	for _, row := range tables.Rows[1:] { // skip header
		req = append(req, row.Cells[0].Value)
	}
	return baseline.Verify(s.db, "mediawiki", 1, req)
}

func (s *bddSuite) theTotalTableCountMustBeAtLeast(min int) error {
	count, err := baseline.TableCount(s.db, "mediawiki")
	if err != nil {
		return err
	}
	if count < min {
		return fmt.Errorf("expected at least %d tables, got %d", min, count)
	}
	return nil
}

func (s *bddSuite) anSMWDatabaseWithFPTEntriesAndDIEntries(fptCount, diCount int) error {
	_, err := s.db.Exec(`
		DROP TABLE IF EXISTS smw_fpt_mdat, smw_di_time;
		CREATE TABLE smw_fpt_mdat (
			s_id INT NOT NULL,
			o_serialized VARCHAR(255) NOT NULL,
			o_sortkey DOUBLE NOT NULL,
			KEY s_id (s_id)
		) ENGINE=InnoDB;

		CREATE TABLE smw_di_time (
			s_id INT NOT NULL,
			p_id INT NOT NULL,
			o_serialized VARCHAR(255) NOT NULL,
			o_sortkey DOUBLE NOT NULL,
			KEY s_id (s_id),
			KEY p_id (p_id)
		) ENGINE=InnoDB;
	`)
	if err != nil {
		return fmt.Errorf("create smw tables: %w", err)
	}

	for i := 1; i <= fptCount; i++ {
		_, err := s.db.Exec(
			"INSERT INTO smw_fpt_mdat (s_id, o_serialized, o_sortkey) VALUES (?, ?, ?)",
			i, "1/2026/7/15/0/0/0/0", 2461237.0,
		)
		if err != nil {
			return err
		}
		if i <= diCount {
			_, err = s.db.Exec(
				"INSERT INTO smw_di_time (s_id, p_id, o_serialized, o_sortkey) VALUES (?, 29, ?, ?)",
				i, "1/2026/7/15/0/0/0/0", 2461237.0,
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *bddSuite) aDriftCheckIsPerformed() error {
	rep, err := drift.Check(s.db, drift.DefaultRegistry())
	if err != nil {
		return err
	}
	s.driftRep = rep
	return nil
}

func (s *bddSuite) missingDIEntriesShouldBeDetected(expected int) error {
	if s.driftRep == nil {
		return fmt.Errorf("no drift check report found")
	}
	if len(s.driftRep.Targets) == 0 {
		return fmt.Errorf("no drift targets in report")
	}
	if s.driftRep.Targets[0].MissingInDI != expected {
		return fmt.Errorf("expected %d missing, got %d", expected, s.driftRep.Targets[0].MissingInDI)
	}
	return nil
}

func (s *bddSuite) aDriftFixIsExecuted() error {
	fixed, err := drift.Fix(s.db, drift.DefaultRegistry())
	if err != nil {
		return err
	}
	s.fixedCount = fixed
	return nil
}

func (s *bddSuite) rowsShouldBeInsertedIntoSmw_di_time(expected int64) error {
	if s.fixedCount != expected {
		return fmt.Errorf("expected %d fixed rows, got %d", expected, s.fixedCount)
	}
	return nil
}

func (s *bddSuite) aSubsequentDriftCheckShouldDetectZeroDrift() error {
	rep, err := drift.Check(s.db, drift.DefaultRegistry())
	if err != nil {
		return err
	}
	if rep.HasDrift {
		return fmt.Errorf("expected zero drift, got: %s", rep.String())
	}
	return nil
}

// ---------------------------------------------------------------------------
// T-538 semantic contract BDD steps
// ---------------------------------------------------------------------------

func (s *bddSuite) aCompiledPropertyCatalog() error {
	c, _, err := catalog.Load(filepath.Join("..", "pkg", "catalog", "testdata", "catalog_v1.json"))
	if err != nil {
		return fmt.Errorf("load catalog fixture: %w", err)
	}
	s.catalog = c

	// Build p_id lookup map for the loaded catalog (underscore-form titles).
	s.pIDByProperty = make(map[string]int)
	rows, err := s.db.Query(`SELECT smw_id, smw_title FROM smw_object_ids WHERE smw_namespace = ?`, propertyPageNS)
	if err != nil {
		return fmt.Errorf("query smw_object_ids: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var title string
		if err := rows.Scan(&id, &title); err != nil {
			return err
		}
		s.pIDByProperty[title] = id
	}
	return rows.Err()
}

func (s *bddSuite) aCompiledPropertyCatalogWithProperty(propName string) error {
	if err := s.aCompiledPropertyCatalog(); err != nil {
		return err
	}
	if s.catalog.PropertyByName(propName) == nil {
		return fmt.Errorf("property %q not in fixture catalog", propName)
	}
	// "that property" in subsequent steps resolves to this one.
	s.currentProp = propName
	return nil
}

func (s *bddSuite) aCompiledPropertyCatalogWithPropertyDeclaredAs(propName, declaredType string) error {
	if err := s.aCompiledPropertyCatalogWithProperty(propName); err != nil {
		return err
	}
	p := s.catalog.PropertyByName(propName)
	if !strings.EqualFold(p.Type, declaredType) {
		return fmt.Errorf("property %q declared type %q, want %q", propName, p.Type, declaredType)
	}
	return nil
}

// aCompiledPropertyCatalogWithAPageProperty builds a synthetic catalog with a
// single Page-typed property so the relational-reference topology (Contract 3)
// can be exercised even though the production catalog has no Page properties.
func (s *bddSuite) aCompiledPropertyCatalogWithAPageProperty(propName string) error {
	s.catalog = &catalog.Catalog{
		Version: 1,
		Properties: []catalog.PropertyDeclaration{
			{Name: propName, Type: "Page"},
		},
		Entities: nil,
	}
	s.pIDByProperty = make(map[string]int)
	// "that property" in subsequent steps resolves to this one.
	s.currentProp = propName
	return nil
}

func (s *bddSuite) aCompiledPropertyCatalogWithPropertyAndAllowedValues(propName, allowedCSV string) error {
	if err := s.aCompiledPropertyCatalogWithProperty(propName); err != nil {
		return err
	}
	p := s.catalog.PropertyByName(propName)
	if len(p.Allowed) == 0 {
		return fmt.Errorf("property %q has no allowed values in fixture", propName)
	}
	return nil
}

func (s *bddSuite) aDatabaseWithSMWTablesLoaded() error {
	schemaPath := filepath.Join("..", "baselines", "schema.sql")
	if err := baseline.Import(s.db, schemaPath); err != nil {
		return fmt.Errorf("import baseline schema: %w", err)
	}
	return nil
}

// theCatalogPropertyPagesAreSeededIntoSmwObjectIds populates smw_object_ids
// with every catalog property page (smw_namespace=102, underscore-form title)
// and a matching smw_fpt_type declaration, producing the realistic seeded
// schema that Contract 1's clean-state scenario expects.
func (s *bddSuite) theCatalogPropertyPagesAreSeededIntoSmwObjectIds() error {
	if s.catalog == nil {
		return fmt.Errorf("no catalog loaded")
	}
	for _, prop := range s.catalog.Properties {
		title := catalog.SMWTitle(prop.Name)
		if _, err := s.db.Exec(
			`INSERT INTO smw_object_ids (smw_namespace, smw_title, smw_iw, smw_subobject, smw_sortkey) VALUES (?, ?, '', '', ?)`,
			propertyPageNS, title, title,
		); err != nil {
			return fmt.Errorf("seed property %q: %w", prop.Name, err)
		}
		res, err := s.db.Exec(
			`INSERT INTO smw_fpt_type (s_id, o_serialized) SELECT smw_id, ? FROM smw_object_ids WHERE smw_namespace = ? AND smw_title = ?`,
			smwCodeForType(prop.Type), propertyPageNS, title,
		)
		if err != nil {
			return fmt.Errorf("seed fpt_type for %q: %w", prop.Name, err)
		}
		_ = res

		// Seed allowed values for enumerated properties.
		for _, allowed := range prop.Allowed {
			if _, err := s.db.Exec(
				`INSERT INTO smw_fpt_pval (s_id, o_hash) SELECT smw_id, ? FROM smw_object_ids WHERE smw_namespace = ? AND smw_title = ?`,
				allowed, propertyPageNS, title,
			); err != nil {
				return fmt.Errorf("seed fpt_pval for %q: %w", prop.Name, err)
			}
		}

		// Seed equivalence mapping for mapped properties.
		if prop.Equivalent != "" {
			if _, err := s.db.Exec(
				`INSERT INTO smw_fpt_impo (s_id, o_hash) SELECT smw_id, ? FROM smw_object_ids WHERE smw_namespace = ? AND smw_title = ?`,
				prop.Equivalent, propertyPageNS, title,
			); err != nil {
				return fmt.Errorf("seed fpt_impo for %q: %w", prop.Name, err)
			}
		}
	}
	return nil
}

func (s *bddSuite) aDatabaseWherePropertyHasNoObjectIDsEntry(propName string) error {
	if err := s.aDatabaseWithSMWTablesLoaded(); err != nil {
		return err
	}
	// Seed ALL catalog property pages except the target, so the audit reports
	// exactly one declaration violation (fidelity: not 83 blanket violations).
	if s.catalog == nil {
		return fmt.Errorf("no catalog loaded")
	}
	for _, prop := range s.catalog.Properties {
		if prop.Name == propName {
			continue
		}
		title := catalog.SMWTitle(prop.Name)
		if _, err := s.db.Exec(
			`INSERT INTO smw_object_ids (smw_namespace, smw_title, smw_iw, smw_subobject, smw_sortkey) VALUES (?, ?, '', '', ?)`,
			propertyPageNS, title, title,
		); err != nil {
			return fmt.Errorf("seed property %q: %w", prop.Name, err)
		}
	}
	// Ensure the target property is deliberately absent from smw_object_ids.
	_, err := s.db.Exec(`DELETE FROM smw_object_ids WHERE smw_namespace = ? AND smw_title = ?`, propertyPageNS, catalog.SMWTitle(propName))
	return err
}

// seeding helpers ----------------------------------------------------------

func (s *bddSuite) seedPropertyObjectID(propName, smwType string) (int, error) {
	title := catalog.SMWTitle(propName)
	res, err := s.db.Exec(
		`INSERT INTO smw_object_ids (smw_namespace, smw_title, smw_iw, smw_subobject, smw_sortkey) VALUES (?, ?, '', '', ?)`,
		propertyPageNS, title, title,
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	pID := int(id)
	s.pIDByProperty[title] = pID

	// Register the type in smw_fpt_type.
	if _, err := s.db.Exec(
		`INSERT INTO smw_fpt_type (s_id, o_serialized) VALUES (?, ?)`,
		pID, smwType,
	); err != nil {
		return 0, err
	}
	return pID, nil
}

func (s *bddSuite) aDatabaseWithADatePropertyStoredInSmwDiTime() error {
	if err := s.aDatabaseWithSMWTablesLoaded(); err != nil {
		return err
	}
	pID, err := s.seedPropertyObjectID("Event start date", "_dat")
	if err != nil {
		return err
	}
	_, err = s.db.Exec(
		`INSERT INTO smw_di_time (s_id, p_id, o_serialized, o_sortkey) VALUES (?, ?, ?, 1)`,
		pID, pID, "2026-08-27",
	)
	return err
}

// currentPropIs returns the property named by the preceding Given step, or an
// error when none has been set.
func (s *bddSuite) currentPropIs() (string, error) {
	if s.currentProp == "" {
		return "", fmt.Errorf("no current property set by a preceding Given step")
	}
	return s.currentProp, nil
}

func (s *bddSuite) aDatabaseWithBlobRowsForThatPropertysPID() error {
	propName, err := s.currentPropIs()
	if err != nil {
		return err
	}
	if err := s.aDatabaseWithSMWTablesLoaded(); err != nil {
		return err
	}
	// Seed the property as declared Date, then insert its value into the WRONG table.
	pID, err := s.seedPropertyObjectID(propName, "_dat")
	if err != nil {
		return err
	}
	_, err = s.db.Exec(
		`INSERT INTO smw_di_blob (s_id, p_id, o_hash) VALUES (?, ?, '2026-08-27')`,
		pID, pID,
	)
	return err
}

func (s *bddSuite) aDatabaseWithBlobRowsReferencingAnUnknownPID() error {
	if err := s.aDatabaseWithSMWTablesLoaded(); err != nil {
		return err
	}
	// Insert a row whose p_id matches no catalog property (e.g. 99999).
	_, err := s.db.Exec(
		`INSERT INTO smw_di_blob (s_id, p_id, o_hash) VALUES (?, ?, 'SomeValue')`,
		99999, 99999,
	)
	return err
}

func (s *bddSuite) aDatabaseWithBlobContainingForThatProperty(value string) error {
	propName, err := s.currentPropIs()
	if err != nil {
		return err
	}
	if err := s.aDatabaseWithSMWTablesLoaded(); err != nil {
		return err
	}
	p := s.catalog.PropertyByName(propName)
	if p == nil {
		return fmt.Errorf("property %q not in catalog", propName)
	}
	// Declared type per catalog; seed object id.
	code := smwCodeForType(p.Type)
	pID, err := s.seedPropertyObjectID(propName, code)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(
		`INSERT INTO smw_di_blob (s_id, p_id, o_hash) VALUES (?, ?, ?)`,
		pID, pID, value,
	)
	return err
}

func (s *bddSuite) aDatabaseWithTimeContainingForThatProperty(value string) error {
	propName, err := s.currentPropIs()
	if err != nil {
		return err
	}
	if err := s.aDatabaseWithSMWTablesLoaded(); err != nil {
		return err
	}
	p := s.catalog.PropertyByName(propName)
	if p == nil {
		return fmt.Errorf("property %q not in catalog", propName)
	}
	pID, err := s.seedPropertyObjectID(propName, smwCodeForType(p.Type))
	if err != nil {
		return err
	}
	_, err = s.db.Exec(
		`INSERT INTO smw_di_time (s_id, p_id, o_serialized, o_sortkey) VALUES (?, ?, ?, 1)`,
		pID, pID, value,
	)
	return err
}

func (s *bddSuite) aDatabaseWithWikipageOIDPointingToNonExistent() error {
	propName, err := s.currentPropIs()
	if err != nil {
		return err
	}
	if err := s.aDatabaseWithSMWTablesLoaded(); err != nil {
		return err
	}
	p := s.catalog.PropertyByName(propName)
	if p == nil {
		return fmt.Errorf("property %q not in catalog", propName)
	}
	pID, err := s.seedPropertyObjectID(propName, smwCodeForType(p.Type))
	if err != nil {
		return err
	}
	// o_id = 424242 does not exist in smw_object_ids.
	_, err = s.db.Exec(
		`INSERT INTO smw_di_wikipage (s_id, p_id, o_id) VALUES (?, ?, 424242)`,
		pID, pID,
	)
	return err
}

// audit execution ----------------------------------------------------------

func (s *bddSuite) anAuditShouldReport(rep *audit.Report, want bool) error {
	if rep == nil {
		return fmt.Errorf("audit not executed")
	}
	has := len(rep.Violations) > 0
	if want && !has {
		return fmt.Errorf("expected audit violations, got none: %s", rep.String())
	}
	if !want && has {
		return fmt.Errorf("expected no violations, got %d: %s", len(rep.Violations), rep.String())
	}
	return nil
}

func (s *bddSuite) everyCatalogPropertyMustHaveAnObjectIDsEntry() error {
	rep, err := audit.AuditDeclarations(s.db, s.catalog)
	if err != nil {
		return err
	}
	s.declarationRep = rep
	// In the clean fixture DB, all declared properties have object IDs.
	return s.anAuditShouldReport(rep, false)
}

func (s *bddSuite) theSMWFPTypeMustMatchTheDeclaredType() error {
	rep := s.declarationRep
	if rep == nil {
		return fmt.Errorf("no declaration audit ran")
	}
	for _, v := range rep.Violations {
		if v.Rule == "type_mismatch" {
			return fmt.Errorf("unexpected type mismatch: %s", v.Diagnostic)
		}
	}
	return nil
}

func (s *bddSuite) theSMWFPPvalMustMatchTheDeclaredAllowedValues() error {
	rep := s.declarationRep
	if rep == nil {
		return fmt.Errorf("no declaration audit ran")
	}
	for _, v := range rep.Violations {
		if v.Rule == "allowed_value_missing" {
			return fmt.Errorf("unexpected allowed-value mismatch: %s", v.Diagnostic)
		}
	}
	return nil
}

func (s *bddSuite) anAuditShouldReportADeclarationViolation() error {
	rep, err := audit.AuditDeclarations(s.db, s.catalog)
	if err != nil {
		return err
	}
	s.declarationRep = rep
	return s.anAuditShouldReport(rep, true)
}

func (s *bddSuite) theViolationDiagnosticShouldIdentifyThePropertyName() error {
	for _, v := range s.declarationRep.Violations {
		if strings.Contains(v.Diagnostic, "Event type") {
			return nil
		}
	}
	return fmt.Errorf("no violation diagnostic mentions 'Event type': %+v", s.declarationRep.Violations)
}

func (s *bddSuite) anAuditShouldReportZeroRoutingViolations() error {
	rep, err := audit.AuditValueTypes(s.db, s.catalog, audit.DefaultOptions())
	if err != nil {
		return err
	}
	s.valueTypeRep = rep
	return s.anAuditShouldReport(rep, false)
}

func (s *bddSuite) anAuditShouldReportARoutingViolation() error {
	rep, err := audit.AuditValueTypes(s.db, s.catalog, audit.DefaultOptions())
	if err != nil {
		return err
	}
	s.valueTypeRep = rep
	return s.anAuditShouldReport(rep, true)
}

func (s *bddSuite) theDiagnosticShouldIncludeTheExpectedTable(expectedTable string) error {
	for _, v := range s.valueTypeRep.Violations {
		if v.Kind == audit.KindRouting && strings.Contains(v.Diagnostic, expectedTable) {
			return nil
		}
	}
	return fmt.Errorf("no routing violation mentions %q: %+v", expectedTable, s.valueTypeRep.Violations)
}

func (s *bddSuite) anAuditShouldReportAnOrphanedPredicate() error {
	rep, err := audit.AuditValueTypes(s.db, s.catalog, audit.DefaultOptions())
	if err != nil {
		return err
	}
	for _, v := range rep.Violations {
		if v.Kind == audit.KindOrphanedPredicate {
			return nil
		}
	}
	return fmt.Errorf("expected orphaned predicate violation, got: %s", rep.String())
}

func (s *bddSuite) anAuditShouldReportARangeViolation() error {
	rep, err := audit.AuditValueRanges(s.db, s.catalog, audit.DefaultOptions())
	if err != nil {
		return err
	}
	s.valueRangeRep = rep
	return s.anAuditShouldReport(rep, true)
}

func (s *bddSuite) theViolationDiagnosticShouldIncludeTheDeclaredAllowedValues() error {
	for _, v := range s.valueRangeRep.Violations {
		if v.Kind == audit.KindRange && strings.Contains(v.Diagnostic, "In-Person") {
			return nil
		}
	}
	return fmt.Errorf("no range violation diagnostic mentions allowed values: %+v", s.valueRangeRep.Violations)
}

func (s *bddSuite) anAuditShouldReportASyntaxViolation() error {
	rep, err := audit.AuditValueRanges(s.db, s.catalog, audit.DefaultOptions())
	if err != nil {
		return err
	}
	for _, v := range rep.Violations {
		if v.Kind == audit.KindSyntax {
			return nil
		}
	}
	return fmt.Errorf("expected syntax violation, got: %s", rep.String())
}

func (s *bddSuite) anAuditShouldReportAReferenceViolation() error {
	rep, err := audit.AuditValueRanges(s.db, s.catalog, audit.DefaultOptions())
	if err != nil {
		return err
	}
	for _, v := range rep.Violations {
		if v.Kind == audit.KindReference {
			return nil
		}
	}
	return fmt.Errorf("expected reference violation, got: %s", rep.String())
}

// smwCodeForType maps a catalog type name to its SMW internal code.
func smwCodeForType(typ string) string {
	switch typ {
	case "Text", "Code":
		return "_txt"
	case "Date":
		return "_dat"
	case "Number":
		return "_num"
	case "URL", "Email":
		return "_uri"
	case "Boolean":
		return "_boo"
	case "Page":
		return "_wpg"
	default:
		return "_txt"
	}
}

// parseCSV splits a comma-separated list, trimming whitespace.
func parseCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}

func (s *bddSuite) stepErrIsNil(err error) error {
	return err
}

// ---------------------------------------------------------------------------
// T-541 projection discipline BDD steps — spec is the code
//
// These scenarios enforce the projection line: catalog.json carries declared
// facts only; derivation and negotiation live in consumers. The steps are
// pure-Go (no DB) — they run inside the shared testcontainers suite in CI.
// ---------------------------------------------------------------------------

func (s *bddSuite) aCatalogArtifactWithVersion(version int) error {
	data, err := os.ReadFile(filepath.Join("..", "pkg", "catalog", "testdata", "catalog_v1.json"))
	if err != nil {
		return fmt.Errorf("read fixture: %w", err)
	}
	var head struct {
		Version int `json:"version"`
	}
	if err := json.Unmarshal(data, &head); err != nil {
		return err
	}
	if head.Version != version {
		return fmt.Errorf("fixture version %d, want %d", head.Version, version)
	}
	s.disciplineData = data
	return nil
}

func (s *bddSuite) theDisciplineInspectionMustPass() error {
	if len(s.disciplineData) == 0 {
		return fmt.Errorf("no artifact loaded")
	}
	rep, err := catalog.Inspect(s.disciplineData)
	if err != nil {
		return fmt.Errorf("inspect: %w", err)
	}
	s.disciplineReport = rep
	if len(rep.Violations) > 0 {
		return fmt.Errorf("discipline inspection failed: %+v", rep.Violations)
	}
	return nil
}

func (s *bddSuite) theArtifactMustNotCarryAnyDerivedField() error {
	for _, v := range s.disciplineReport.Violations {
		if v.Rule == catalog.RuleDerivedField {
			return fmt.Errorf("derived field present: %+v", v)
		}
	}
	return nil
}

func (s *bddSuite) theSemanticNestingDepthMustNotExceedTheEntryLevel() error {
	for _, v := range s.disciplineReport.Violations {
		if v.Rule == catalog.RuleNestingDepth {
			return fmt.Errorf("nesting violation: %+v", v)
		}
	}
	return nil
}

// hostileArtifact builds a hostile catalog.json from the committed fixture by
// replacing the first occurrence of marker with replacement. It always starts
// from the pristine fixture so scenarios are independent — godog shares one
// suite instance across scenarios and mutated state would leak.
func (s *bddSuite) hostileArtifact(marker, replacement string) error {
	data, err := os.ReadFile(filepath.Join("..", "pkg", "catalog", "testdata", "catalog_v1.json"))
	if err != nil {
		return fmt.Errorf("read fixture: %w", err)
	}
	text := string(data)
	idx := strings.Index(text, marker)
	if idx < 0 {
		return fmt.Errorf("marker %q not found in fixture", marker)
	}
	s.disciplineData = []byte(text[:idx] + replacement + text[idx+len(marker):])
	return nil
}

func (s *bddSuite) aCatalogArtifactCarryingTheDerivedFieldOnAPropertyEntry(field string) error {
	return s.hostileArtifact(`"type": "Text"`, `"type": "Text", "`+field+`": "smw_di_time"`)
}

func (s *bddSuite) theDisciplineInspectionMustReportADerivedFieldViolation() error {
	rep, err := catalog.Inspect(s.disciplineData)
	if err != nil {
		return err
	}
	s.disciplineReport = rep
	for _, v := range rep.Violations {
		if v.Rule == catalog.RuleDerivedField {
			return nil
		}
	}
	return fmt.Errorf("expected derived-field violation, got: %+v", rep.Violations)
}

func (s *bddSuite) aCatalogArtifactCarryingTheUnknownTopLevelField(field string) error {
	return s.hostileArtifact(`"version": 1`, `"version": 1, "`+field+`": {}`)
}

func (s *bddSuite) theArtifactMustStillLoadSuccessfully() error {
	rep, err := catalog.Inspect(s.disciplineData)
	if err != nil {
		return fmt.Errorf("artifact failed to load: %w", err)
	}
	s.disciplineReport = rep
	return nil
}

func (s *bddSuite) theInspectionMustRecordAnOpenWorldWarningNotAViolation() error {
	if len(s.disciplineReport.Violations) > 0 {
		return fmt.Errorf("expected no violations under OWA, got: %+v", s.disciplineReport.Violations)
	}
	warned := false
	for _, w := range s.disciplineReport.Warnings {
		if strings.Contains(w, "extensions") {
			warned = true
		}
	}
	if !warned {
		return fmt.Errorf("expected OWA warning mentioning extensions, got: %+v", s.disciplineReport.Warnings)
	}
	return nil
}

func (s *bddSuite) aCatalogArtifactWhoseAllowedValuesAreNestedObjects() error {
	return s.hostileArtifact(`"allowed": null,`, `"allowed": [{"value": "A"}],`)
}

func (s *bddSuite) theDisciplineInspectionMustReportANestingDepthViolation() error {
	rep, err := catalog.Inspect(s.disciplineData)
	if err != nil {
		return err
	}
	s.disciplineReport = rep
	for _, v := range rep.Violations {
		if v.Rule == catalog.RuleNestingDepth {
			return nil
		}
	}
	return fmt.Errorf("expected nesting-depth violation, got: %+v", rep.Violations)
}

func (s *bddSuite) aCatalogArtifactDeclaringADateProperty() error {
	return s.hostileArtifact(`"type": "Text"`, `"type": "Date"`)
}

func (s *bddSuite) theExpectedStorageTableIsDerivedFromTheTypeAs(expected string) error {
	got, ok := audit.ExpectedTable("Date")
	if !ok {
		return fmt.Errorf("Date type not routable")
	}
	if got != expected {
		return fmt.Errorf("ExpectedTable(Date) = %q, want %q", got, expected)
	}
	return nil
}

// ---------------------------------------------------------------------------
// T-546 derivation totality steps
// ---------------------------------------------------------------------------

func (s *bddSuite) theCatalogKnownTypes() error {
	s.knownTypes = catalog.KnownTypes()
	if len(s.knownTypes) == 0 {
		return fmt.Errorf("no known types registered")
	}
	return nil
}

func (s *bddSuite) totalityProblems() []string {
	if s.totality == nil {
		s.totality = audit.RoutingTotality()
	}
	return s.totality
}

func (s *bddSuite) totalityMustBeClean() error {
	if ps := s.totalityProblems(); len(ps) > 0 {
		return fmt.Errorf("derivation totality violated: %+v", ps)
	}
	return nil
}

func (s *bddSuite) everyKnownTypeMapsToAStorageTable() error {
	return s.totalityMustBeClean()
}

func (s *bddSuite) everyKnownTypeMapsToAnSMWCode() error {
	return s.totalityMustBeClean()
}

func (s *bddSuite) everySMWCodeDecodesBackToItsKnownType() error {
	return s.totalityMustBeClean()
}

func (s *bddSuite) noRoutedTypeIsAbsentFromTheKnownSet() error {
	return s.totalityMustBeClean()
}

func (s *bddSuite) theScannedDataItemTablesEqualTheRoutableTableSet() error {
	return s.totalityMustBeClean()
}

func (s *bddSuite) aTypeOutsideTheKnownSetMustNotBeRoutable() error {
	for _, t := range []string{"Geo", "Enum", "Widget", "Code"} {
		if _, ok := audit.ExpectedTable(t); ok {
			return fmt.Errorf("type %q is routable but outside PHP truth", t)
		}
	}
	return nil
}

func (s *bddSuite) noTwoPropertyNamesShareAnSMWTitle() error {
	if s.catalog == nil {
		return fmt.Errorf("no catalog loaded")
	}
	if cols := s.catalog.TitleCollisions(); len(cols) > 0 {
		return fmt.Errorf("SMWTitle collisions: %v", cols)
	}
	return nil
}

// aCompiledPropertyCatalogIsLoadedFromTheFixture loads the catalog without
// touching the database — used by the DB-free discipline suite for title-
// collision checks.
func (s *bddSuite) aCompiledPropertyCatalogIsLoadedFromTheFixture() error {
	c, _, err := catalog.Load(filepath.Join("..", "pkg", "catalog", "testdata", "catalog_v1.json"))
	if err != nil {
		return fmt.Errorf("load catalog fixture: %w", err)
	}
	s.catalog = c
	return nil
}

// ---------------------------------------------------------------------------
// T-547 load-time cross-check steps
// ---------------------------------------------------------------------------

func (s *bddSuite) aCatalogArtifactWithAPropertyOfUnknownType(unknownType string) error {
	return s.hostileArtifact(`"type": "Text"`, `"type": "`+unknownType+`"`)
}

func (s *bddSuite) loadingTheArtifactMustFail() error {
	_, _, err := catalog.Parse(s.disciplineData)
	if err == nil {
		return fmt.Errorf("expected load to fail on unknown type")
	}
	s.lastLoadErr = err
	return nil
}

func (s *bddSuite) theErrorMustIdentifyThePropertyNameAndTheUnknownType() error {
	if s.lastLoadErr == nil {
		return fmt.Errorf("no load error captured")
	}
	if !strings.Contains(s.lastLoadErr.Error(), "Geo") {
		return fmt.Errorf("error does not name the unknown type: %v", s.lastLoadErr)
	}
	if !strings.Contains(s.lastLoadErr.Error(), "Author") {
		return fmt.Errorf("error does not name the property (expected Author): %v", s.lastLoadErr)
	}
	return nil
}

func (s *bddSuite) theDisciplineInspectionMustReportAnUnknownTypeViolation() error {
	rep, err := catalog.Inspect(s.disciplineData)
	if err != nil {
		return err
	}
	for _, v := range rep.Violations {
		if v.Rule == catalog.RuleUnknownType {
			return nil
		}
	}
	return fmt.Errorf("expected unknown-type violation, got: %+v", rep.Violations)
}

func (s *bddSuite) theInspectionMustNotClassifyTheUnknownTypeAsAnOpenWorldWarning() error {
	rep, err := catalog.Inspect(s.disciplineData)
	if err != nil {
		return err
	}
	for _, w := range rep.Warnings {
		if strings.Contains(w, "Geo") {
			return fmt.Errorf("unknown type surfaced as OWA warning: %v", w)
		}
	}
	for _, v := range rep.Violations {
		if v.Rule == catalog.RuleUnknownType {
			return nil
		}
	}
	return fmt.Errorf("expected unknown-type violation (not a warning): %+v", rep.Violations)
}

func (s *bddSuite) loadingTheArtifactMustSucceed() error {
	_, _, err := catalog.Parse(s.disciplineData)
	return err
}

// ---------------------------------------------------------------------------
// T-548 exclusion contract steps
// ---------------------------------------------------------------------------

func (s *bddSuite) aCatalogArtifactCarryingTheDerivedFieldOnAnEntityEntry(field string) error {
	// The entity block is pretty-printed: "name": "Area", then "gloss" on
	// the next line (12-space indent). The gloss continuation disambiguates
	// from the property entry of the same name.
	marker := "\"name\": \"Area\",\n            \"gloss\""
	replacement := "\"name\": \"Area\",\n            \"" + field + "\": \"Area\",\n            \"gloss\""
	return s.hostileArtifact(marker, replacement)
}

func (s *bddSuite) aCatalogArtifactCarryingTheDerivedFieldAtTheTopLevel(field string) error {
	return s.hostileArtifact(`"version": 1`, `"version": 1, "`+field+`": 999`)
}

func (s *bddSuite) theExclusionListCoversTheDerivedFields(fields string) error {
	for _, f := range strings.Split(fields, ",") {
		f = strings.TrimSpace(f)
		if f == "" {
			continue
		}
		if !catalog.IsDerivedField(f) {
			return fmt.Errorf("derived field %q not covered by the exclusion list", f)
		}
	}
	return nil
}

func (s *bddSuite) expectedTableStillDerivesRoutingFromTheDeclaredType() error {
	got, ok := audit.ExpectedTable("Text")
	if !ok {
		return fmt.Errorf("Text type not routable")
	}
	if got != "smw_di_blob" {
		return fmt.Errorf("ExpectedTable(Text) = %q, want smw_di_blob (wire field must not steer routing)", got)
	}
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext, suite *bddSuite) {
	ctx.Step(`^a clean baseline database is loaded$`, suite.aCleanBaselineDatabaseIsLoaded)
	ctx.Step(`^the following core tables must exist:$`, suite.theFollowingCoreTablesMustExist)
	ctx.Step(`^the total table count must be at least (\d+)$`, suite.theTotalTableCountMustBeAtLeast)
	ctx.Step(`^an SMW database with (\d+) FPT entries and (\d+) DI entries$`, suite.anSMWDatabaseWithFPTEntriesAndDIEntries)
	ctx.Step(`^a drift check is performed$`, suite.aDriftCheckIsPerformed)
	ctx.Step(`^(\d+) missing DI entries should be detected$`, suite.missingDIEntriesShouldBeDetected)
	ctx.Step(`^a drift fix is executed$`, suite.aDriftFixIsExecuted)
	ctx.Step(`^(\d+) rows should be inserted into smw_di_time$`, suite.rowsShouldBeInsertedIntoSmw_di_time)
	ctx.Step(`^a subsequent drift check should detect zero drift$`, suite.aSubsequentDriftCheckShouldDetectZeroDrift)

	// T-538 semantic contract steps (Contract 1 — declarations)
	ctx.Step(`^a compiled property catalog$`, suite.aCompiledPropertyCatalog)
	ctx.Step(`^a compiled property catalog with property "([^"]+)"$`, suite.aCompiledPropertyCatalogWithProperty)
	ctx.Step(`^a compiled property catalog with property "([^"]+)" declared as (\w+)$`, suite.aCompiledPropertyCatalogWithPropertyDeclaredAs)
	ctx.Step(`^a compiled property catalog with property "([^"]+)" and allowed values "([^"]+)"$`, suite.aCompiledPropertyCatalogWithPropertyAndAllowedValues)
	ctx.Step(`^a compiled property catalog with a Page property "([^"]+)"$`, suite.aCompiledPropertyCatalogWithAPageProperty)
	ctx.Step(`^the catalog property pages are seeded into smw_object_ids$`, suite.theCatalogPropertyPagesAreSeededIntoSmwObjectIds)
	ctx.Step(`^a database with SMW tables loaded$`, suite.aDatabaseWithSMWTablesLoaded)
	ctx.Step(`^a database where "([^"]+)" has no smw_object_ids entry$`, suite.aDatabaseWherePropertyHasNoObjectIDsEntry)
	ctx.Step(`^every catalog property must have an smw_object_ids entry$`, suite.everyCatalogPropertyMustHaveAnObjectIDsEntry)
	ctx.Step(`^the smw_fpt_type must match the declared type$`, suite.theSMWFPTypeMustMatchTheDeclaredType)
	ctx.Step(`^the smw_fpt_pval must match the declared allowed values$`, suite.theSMWFPPvalMustMatchTheDeclaredAllowedValues)
	ctx.Step(`^an audit should report a declaration violation$`, suite.anAuditShouldReportADeclarationViolation)
	ctx.Step(`^the violation diagnostic should identify the property name$`, suite.theViolationDiagnosticShouldIdentifyThePropertyName)

	// T-538 semantic contract steps (Contract 2 — value types)
	ctx.Step(`^a database with a Date property stored in smw_di_time$`, suite.aDatabaseWithADatePropertyStoredInSmwDiTime)
	ctx.Step(`^a database with smw_di_blob rows for that property's p_id$`, suite.aDatabaseWithBlobRowsForThatPropertysPID)
	ctx.Step(`^a database with smw_di_blob rows referencing an unknown p_id$`, suite.aDatabaseWithBlobRowsReferencingAnUnknownPID)
	ctx.Step(`^an audit should report zero routing violations$`, suite.anAuditShouldReportZeroRoutingViolations)
	ctx.Step(`^an audit should report a routing violation$`, suite.anAuditShouldReportARoutingViolation)
	ctx.Step(`^the diagnostic should include the expected table (\w+)$`, suite.theDiagnosticShouldIncludeTheExpectedTable)
	ctx.Step(`^an audit should report an orphaned predicate$`, suite.anAuditShouldReportAnOrphanedPredicate)

	// T-538 semantic contract steps (Contract 3 — value ranges)
	ctx.Step(`^a database with smw_di_blob containing "([^"]+)" for that property$`, suite.aDatabaseWithBlobContainingForThatProperty)
	ctx.Step(`^a database with smw_di_time containing "([^"]+)" for that property$`, suite.aDatabaseWithTimeContainingForThatProperty)
	ctx.Step(`^a database with smw_di_wikipage o_id pointing to a non-existent smw_object_ids entry$`, suite.aDatabaseWithWikipageOIDPointingToNonExistent)
	ctx.Step(`^an audit should report a range violation$`, suite.anAuditShouldReportARangeViolation)
	ctx.Step(`^the violation diagnostic should include the declared allowed values$`, suite.theViolationDiagnosticShouldIncludeTheDeclaredAllowedValues)
	ctx.Step(`^an audit should report a syntax violation$`, suite.anAuditShouldReportASyntaxViolation)
	ctx.Step(`^an audit should report a reference violation$`, suite.anAuditShouldReportAReferenceViolation)

	// T-541+ projection discipline steps (pure-Go, no DB)
	registerDisciplineSteps(ctx, suite)
}

// registerDisciplineSteps binds all projection-discipline scenarios — the
// projection line, derivation totality, load-time cross-check, and exclusion
// contract. These steps are pure Go: they never touch the database, so the
// discipline suite runs without Docker (TestBDDDiscipline) while also running
// inside the full container suite (TestBDDFeatures).
func registerDisciplineSteps(ctx *godog.ScenarioContext, suite *bddSuite) {
	ctx.Step(`^a catalog artifact with version (\d+)$`, suite.aCatalogArtifactWithVersion)
	ctx.Step(`^the discipline inspection must pass$`, suite.theDisciplineInspectionMustPass)
	ctx.Step(`^the artifact must not carry any derived field$`, suite.theArtifactMustNotCarryAnyDerivedField)
	ctx.Step(`^the semantic nesting depth must not exceed the entry level$`, suite.theSemanticNestingDepthMustNotExceedTheEntryLevel)
	ctx.Step(`^a catalog artifact carrying the derived field "([^"]+)" on a property entry$`, suite.aCatalogArtifactCarryingTheDerivedFieldOnAPropertyEntry)
	ctx.Step(`^a catalog artifact carrying the derived field "([^"]+)" on an entity entry$`, suite.aCatalogArtifactCarryingTheDerivedFieldOnAnEntityEntry)
	ctx.Step(`^a catalog artifact carrying the derived field "([^"]+)" at the top level$`, suite.aCatalogArtifactCarryingTheDerivedFieldAtTheTopLevel)
	ctx.Step(`^the discipline inspection must report a derived-field violation$`, suite.theDisciplineInspectionMustReportADerivedFieldViolation)
	ctx.Step(`^a catalog artifact carrying the unknown top-level field "([^"]+)"$`, suite.aCatalogArtifactCarryingTheUnknownTopLevelField)
	ctx.Step(`^the artifact must still load successfully$`, suite.theArtifactMustStillLoadSuccessfully)
	ctx.Step(`^the inspection must record an open-world warning, not a violation$`, suite.theInspectionMustRecordAnOpenWorldWarningNotAViolation)
	ctx.Step(`^a catalog artifact whose allowed values are nested objects$`, suite.aCatalogArtifactWhoseAllowedValuesAreNestedObjects)
	ctx.Step(`^the discipline inspection must report a nesting-depth violation$`, suite.theDisciplineInspectionMustReportANestingDepthViolation)
	ctx.Step(`^a catalog artifact declaring a Date property$`, suite.aCatalogArtifactDeclaringADateProperty)
	ctx.Step(`^the expected storage table is derived from the type as (\w+)$`, suite.theExpectedStorageTableIsDerivedFromTheTypeAs)

	// T-546 derivation totality
	ctx.Step(`^the catalog known types$`, suite.theCatalogKnownTypes)
	ctx.Step(`^every known type maps to a storage table$`, suite.everyKnownTypeMapsToAStorageTable)
	ctx.Step(`^every known type maps to an SMW code$`, suite.everyKnownTypeMapsToAnSMWCode)
	ctx.Step(`^every SMW code decodes back to its known type$`, suite.everySMWCodeDecodesBackToItsKnownType)
	ctx.Step(`^no routed type is absent from the known set$`, suite.noRoutedTypeIsAbsentFromTheKnownSet)
	ctx.Step(`^the scanned data-item tables equal the routable table set$`, suite.theScannedDataItemTablesEqualTheRoutableTableSet)
	ctx.Step(`^a type outside the known set must not be routable$`, suite.aTypeOutsideTheKnownSetMustNotBeRoutable)

	// T-546 title collisions (needs the fixture catalog, no DB)
	ctx.Step(`^a compiled property catalog is loaded from the fixture$`, suite.aCompiledPropertyCatalogIsLoadedFromTheFixture)
	ctx.Step(`^no two property names share an SMWTitle$`, suite.noTwoPropertyNamesShareAnSMWTitle)

	// T-547 load-time cross-check
	ctx.Step(`^a catalog artifact with a property of unknown type "([^"]+)"$`, suite.aCatalogArtifactWithAPropertyOfUnknownType)
	ctx.Step(`^loading the artifact must fail$`, suite.loadingTheArtifactMustFail)
	ctx.Step(`^the error must identify the property name and the unknown type$`, suite.theErrorMustIdentifyThePropertyNameAndTheUnknownType)
	ctx.Step(`^the discipline inspection must report an unknown-type violation$`, suite.theDisciplineInspectionMustReportAnUnknownTypeViolation)
	ctx.Step(`^the inspection must not classify the unknown type as an open-world warning$`, suite.theInspectionMustNotClassifyTheUnknownTypeAsAnOpenWorldWarning)
	ctx.Step(`^loading the artifact must succeed$`, suite.loadingTheArtifactMustSucceed)

	// T-548 exclusion contract
	ctx.Step(`^the exclusion list covers the derived fields (.*)$`, suite.theExclusionListCoversTheDerivedFields)
	ctx.Step(`^ExpectedTable still derives routing from the declared type$`, suite.expectedTableStillDerivesRoutingFromTheDeclaredType)
}

func TestBDDFeatures(t *testing.T) {
	ctx := context.Background()

	container, err := mysql.Run(ctx, "mysql:8.4",
		mysql.WithDatabase("mediawiki"),
		mysql.WithUsername("root"),
		mysql.WithPassword("test"),
	)
	if err != nil {
		t.Fatalf("start container: %v", err)
	}
	defer container.Terminate(ctx)

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "3306")

	db, err := connect.Open(connect.Config{
		Host:     host,
		Port:     port.Port(),
		User:     "root",
		Password: "test",
		Database: "mediawiki",
	})
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer db.Close()

	suite := &bddSuite{
		container: container,
		db:        db,
	}

	opts := godog.Options{
		Format:   "pretty",
		Paths:    []string{"../features"},
		TestingT: t,
	}

	status := godog.TestSuite{
		Name:                "observatory-dbtools BDD Compliance Suite",
		ScenarioInitializer: func(c *godog.ScenarioContext) { InitializeScenario(c, suite) },
		Options:             &opts,
	}.Run()

	if status != 0 {
		t.Fatalf("godog suite failed with status code %d", status)
	}
}

// TestBDDDiscipline runs the projection-discipline scenarios in isolation —
// no Docker, no database. Totality, load-time cross-check, and exclusion are
// pure-Go contracts, so they are fully developable independently of the DB
// features and of the PHP producer.
func TestBDDDiscipline(t *testing.T) {
	suite := &bddSuite{}
	opts := godog.Options{
		Format: "pretty",
		Paths: []string{
			"../features/projection_discipline.feature",
			"../features/derivation_totality.feature",
			"../features/load_time_cross_check.feature",
			"../features/exclusion_contract.feature",
		},
		TestingT: t,
	}

	status := godog.TestSuite{
		Name:                "observatory-dbtools Projection Discipline BDD (no Docker)",
		ScenarioInitializer: func(c *godog.ScenarioContext) { registerDisciplineSteps(c, suite) },
		Options:             &opts,
	}.Run()

	if status != 0 {
		t.Fatalf("discipline suite failed with status code %d", status)
	}
}
