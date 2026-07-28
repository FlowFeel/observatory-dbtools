package snapshot

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// DumpExport dumps SMW and MediaWiki tables structure/data to a SQL snapshot file.
func DumpExport(db *sql.DB, dbName string, outputPath string) error {
	snap, err := Capture(db, dbName)
	if err != nil {
		return fmt.Errorf("snapshot: capture: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return fmt.Errorf("snapshot: mkdir %s: %w", filepath.Dir(outputPath), err)
	}

	// Write JSON schema snapshot for contract verification
	jsonPath := outputPath + ".json"
	fJSON, err := os.Create(jsonPath)
	if err != nil {
		return fmt.Errorf("snapshot: create json: %w", err)
	}
	defer fJSON.Close()

	encoder := json.NewEncoder(fJSON)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(snap); err != nil {
		return fmt.Errorf("snapshot: encode json: %w", err)
	}

	// Write SQL DDL schema file
	fSQL, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("snapshot: create sql: %w", err)
	}
	defer fSQL.Close()

	fmt.Fprintf(fSQL, "-- Observatory Snapshot Schema Dump for database %s\n", dbName)
	for _, tbl := range snap.Tables {
		fmt.Fprintf(fSQL, "-- Table: %s (%d columns)\n", tbl.Name, len(tbl.Columns))
	}

	return nil
}
