package test

import (
	"context"
	"database/sql"
	"fmt"
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

func (s *bddSuite) aDatabaseWithBlobRowsForThatPropertysPID(propName string) error {
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

func (s *bddSuite) aDatabaseWithBlobContainingForThatProperty(value, propName string) error {
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

func (s *bddSuite) aDatabaseWithTimeContainingForThatProperty(value, propName string) error {
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

func (s *bddSuite) aDatabaseWithWikipageOIDPointingToNonExistent(propName string) error {
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
	ctx.Step(`^the diagnostic should include the expected table smw_di_time$`, suite.theDiagnosticShouldIncludeTheExpectedTable)
	ctx.Step(`^an audit should report an orphaned predicate$`, suite.anAuditShouldReportAnOrphanedPredicate)

	// T-538 semantic contract steps (Contract 3 — value ranges)
	ctx.Step(`^a database with smw_di_blob containing "([^"]+)" for that property$`, suite.aDatabaseWithBlobContainingForThatProperty)
	ctx.Step(`^a database with smw_di_time containing "([^"]+)" for that property$`, suite.aDatabaseWithTimeContainingForThatProperty)
	ctx.Step(`^a database with smw_di_wikipage o_id pointing to a non-existent smw_object_ids entry$`, suite.aDatabaseWithWikipageOIDPointingToNonExistent)
	ctx.Step(`^an audit should report a range violation$`, suite.anAuditShouldReportARangeViolation)
	ctx.Step(`^the violation diagnostic should include the declared allowed values$`, suite.theViolationDiagnosticShouldIncludeTheDeclaredAllowedValues)
	ctx.Step(`^an audit should report a syntax violation$`, suite.anAuditShouldReportASyntaxViolation)
	ctx.Step(`^an audit should report a reference violation$`, suite.anAuditShouldReportAReferenceViolation)
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
