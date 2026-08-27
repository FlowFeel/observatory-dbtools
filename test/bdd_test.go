package test

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/FlowFeel/observatory-dbtools/pkg/baseline"
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
}

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
