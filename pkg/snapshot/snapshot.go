// Package snapshot handles schema introspection and snapshot representations.
package snapshot

import (
	"database/sql"
	"fmt"
	"sort"
)

// Column represents a database table column definition.
type Column struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	IsNullable string `json:"is_nullable"`
}

// Table represents a database table schema structure.
type Table struct {
	Name    string   `json:"name"`
	Columns []Column `json:"columns"`
}

// Snapshot represents a full database schema snapshot.
type Snapshot struct {
	DatabaseName string  `json:"database_name"`
	Tables       []Table `json:"tables"`
}

// Capture inspects the target database schema and produces a canonical Snapshot.
func Capture(db *sql.DB, dbName string) (*Snapshot, error) {
	rows, err := db.Query(`
		SELECT table_name, column_name, data_type, is_nullable
		FROM information_schema.columns
		WHERE table_schema = ?
		ORDER BY table_name, ordinal_position
	`, dbName)
	if err != nil {
		return nil, fmt.Errorf("snapshot: query columns: %w", err)
	}
	defer rows.Close()

	tableMap := make(map[string]*Table)
	var tableNames []string

	for rows.Next() {
		var tableName, colName, dataType, isNullable string
		if err := rows.Scan(&tableName, &colName, &dataType, &isNullable); err != nil {
			return nil, fmt.Errorf("snapshot: scan column: %w", err)
		}

		tbl, exists := tableMap[tableName]
		if !exists {
			tbl = &Table{Name: tableName}
			tableMap[tableName] = tbl
			tableNames = append(tableNames, tableName)
		}

		tbl.Columns = append(tbl.Columns, Column{
			Name:       colName,
			Type:       dataType,
			IsNullable: isNullable,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("snapshot: iterate columns: %w", err)
	}

	sort.Strings(tableNames)
	snap := &Snapshot{
		DatabaseName: dbName,
		Tables:       make([]Table, 0, len(tableNames)),
	}

	for _, name := range tableNames {
		snap.Tables = append(snap.Tables, *tableMap[name])
	}

	return snap, nil
}
