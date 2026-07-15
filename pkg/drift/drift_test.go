package drift_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/FlowFeel/observatory-dbtools/pkg/connect"
	"github.com/FlowFeel/observatory-dbtools/pkg/drift"
	"github.com/google/go-cmp/cmp"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

var update = flag.Bool("update", false, "update golden files")

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	ctx := context.Background()

	container, err := mysql.Run(ctx, "mysql:8.4",
		mysql.WithDatabase("mediawiki"),
		mysql.WithUsername("root"),
		mysql.WithPassword("test"),
	)
	if err != nil {
		t.Fatalf("failed to start MySQL container: %v", err)
	}
	t.Cleanup(func() { container.Terminate(ctx) })

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get host: %v", err)
	}

	port, err := container.MappedPort(ctx, "3306")
	if err != nil {
		t.Fatalf("failed to get port: %v", err)
	}

	db, err := connect.Open(connect.Config{
		Host:     host,
		Port:     port.Port(),
		User:     "root",
		Password: "test",
		Database: "mediawiki",
	})
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	// Create SMW tables matching production schema
	_, err = db.Exec(`
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
		t.Fatalf("create tables failed: %v", err)
	}

	return db
}

func readGolden(t *testing.T, name string) drift.Report {
	t.Helper()
	path := filepath.Join("testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden %s: %v", path, err)
	}
	var r drift.Report
	if err := json.Unmarshal(data, &r); err != nil {
		t.Fatalf("parse golden %s: %v", path, err)
	}
	return r
}

func writeGolden(t *testing.T, name string, r *drift.Report) {
	t.Helper()
	path := filepath.Join("testdata", name)
	data, err := r.JSON()
	if err != nil {
		t.Fatalf("marshal golden: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		t.Fatalf("write golden %s: %v", path, err)
	}
	t.Logf("updated golden file: %s", path)
}

func TestDrift_ZeroOnCleanState(t *testing.T) {
	db := setupTestDB(t)

	// Insert matching data — 5 entries in both tables
	for i := 1; i <= 5; i++ {
		_, err := db.Exec(
			"INSERT INTO smw_fpt_mdat (s_id, o_serialized, o_sortkey) VALUES (?, ?, ?)",
			i, "1/2026/7/15/0/0/0/0", 2461237.0,
		)
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.Exec(
			"INSERT INTO smw_di_time (s_id, p_id, o_serialized, o_sortkey) VALUES (?, 29, ?, ?)",
			i, "1/2026/7/15/0/0/0/0", 2461237.0,
		)
		if err != nil {
			t.Fatal(err)
		}
	}

	report, err := drift.Check(db)
	if err != nil {
		t.Fatalf("Check failed: %v", err)
	}

	if *update {
		writeGolden(t, "drift_zero.json", report)
		return
	}

	golden := readGolden(t, "drift_zero.json")
	// Compare structural equality, not exact counts (counts depend on seed data)
	if report.HasDrift() {
		t.Errorf("expected zero drift, got: %s", report.String())
	}
	if golden.MissingInDI != 0 {
		t.Errorf("golden should have zero drift")
	}
}

func TestDrift_DetectsAndFixes(t *testing.T) {
	db := setupTestDB(t)

	// Insert 10 entries in FPT, only 7 in DI — 3 drifted
	for i := 1; i <= 10; i++ {
		_, err := db.Exec(
			"INSERT INTO smw_fpt_mdat (s_id, o_serialized, o_sortkey) VALUES (?, ?, ?)",
			i, "1/2026/7/15/0/0/0/0", 2461237.0,
		)
		if err != nil {
			t.Fatal(err)
		}
		if i <= 7 {
			_, err = db.Exec(
				"INSERT INTO smw_di_time (s_id, p_id, o_serialized, o_sortkey) VALUES (?, 29, ?, ?)",
				i, "1/2026/7/15/0/0/0/0", 2461237.0,
			)
			if err != nil {
				t.Fatal(err)
			}
		}
	}

	// Check: should detect drift
	report, err := drift.Check(db)
	if err != nil {
		t.Fatalf("Check failed: %v", err)
	}
	if !report.HasDrift() {
		t.Fatal("expected drift, got none")
	}
	if report.MissingInDI != 3 {
		t.Errorf("expected 3 missing, got %d", report.MissingInDI)
	}

	if *update {
		writeGolden(t, "drift_detected.json", report)
	}

	// Fix
	fixed, err := drift.Fix(db)
	if err != nil {
		t.Fatalf("Fix failed: %v", err)
	}
	if fixed != 3 {
		t.Errorf("expected 3 rows fixed, got %d", fixed)
	}

	// Verify: zero drift after fix
	after, err := drift.Check(db)
	if err != nil {
		t.Fatalf("Check after fix failed: %v", err)
	}
	if after.HasDrift() {
		t.Errorf("expected zero drift after fix, got: %s", after.String())
	}

	if *update {
		writeGolden(t, "drift_fixed.json", after)
	}

	if !*update {
		golden := readGolden(t, "drift_detected.json")
		if diff := cmp.Diff(golden.MissingInDI, report.MissingInDI); diff != "" {
			t.Errorf("drift count mismatch (-golden +actual):\n%s", diff)
		}
	}
}
