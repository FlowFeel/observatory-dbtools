// Package baseline manages database schema baselines — import, export, and verification.
//
// A baseline is a SQL dump representing the known-good schema state of the database.
// Baselines are committed to version control, reviewed in PRs, and never auto-generated in CI.
//
// Import loads a baseline into a target database. Export captures the current schema.
// Verify checks that a database matches the expected baseline.
package baseline

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
)

// Import loads a SQL file into the target database.
// The SQL file should contain CREATE TABLE and INSERT statements.
// Existing tables are NOT dropped — use on clean databases or with --force.
func Import(db *sql.DB, sqlPath string) error {
	data, err := os.ReadFile(sqlPath)
	if err != nil {
		return fmt.Errorf("baseline: read %s: %w", sqlPath, err)
	}

	content := string(data)
	if strings.TrimSpace(content) == "" {
		return fmt.Errorf("baseline: %s is empty", sqlPath)
	}

	_, err = db.Exec(content)
	if err != nil {
		return fmt.Errorf("baseline: import %s: %w", sqlPath, err)
	}

	return nil
}

// TableCount returns the number of tables in the database.
func TableCount(db *sql.DB, dbName string) (int, error) {
	var count int
	err := db.QueryRow(
		"SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ?",
		dbName,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("baseline: table count: %w", err)
	}
	return count, nil
}

// Tables returns the sorted list of table names in the database.
func Tables(db *sql.DB, dbName string) ([]string, error) {
	rows, err := db.Query(
		"SELECT table_name FROM information_schema.tables WHERE table_schema = ? ORDER BY table_name",
		dbName,
	)
	if err != nil {
		return nil, fmt.Errorf("baseline: list tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("baseline: scan table: %w", err)
		}
		tables = append(tables, name)
	}
	return tables, rows.Err()
}

// Verify checks that the database has at least the expected number of tables
// and that all expected tables exist.
func Verify(db *sql.DB, dbName string, expectedMinTables int, requiredTables []string) error {
	count, err := TableCount(db, dbName)
	if err != nil {
		return err
	}
	if count < expectedMinTables {
		return fmt.Errorf("baseline: expected at least %d tables, got %d", expectedMinTables, count)
	}

	tables, err := Tables(db, dbName)
	if err != nil {
		return err
	}

	tableSet := make(map[string]bool, len(tables))
	for _, t := range tables {
		tableSet[t] = true
	}

	var missing []string
	for _, req := range requiredTables {
		if !tableSet[req] {
			missing = append(missing, req)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("baseline: missing required tables: %v", missing)
	}

	return nil
}
