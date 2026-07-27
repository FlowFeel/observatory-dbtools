// Package compare provides schema snapshot diffing capabilities.
package compare

import (
	"github.com/FlowFeel/observatory-dbtools/pkg/snapshot"
)

// ColumnDiff describes a changed or added/removed column.
type ColumnDiff struct {
	ColumnName string `json:"column_name"`
	ChangeType string `json:"change_type"` // "added", "removed", "type_changed"
	OldType    string `json:"old_type,omitempty"`
	NewType    string `json:"new_type,omitempty"`
}

// TableDiff describes differences between two table definitions.
type TableDiff struct {
	TableName string       `json:"table_name"`
	Status    string       `json:"status"` // "added", "removed", "modified"
	Columns   []ColumnDiff `json:"columns,omitempty"`
}

// Diff represents the comparison result between two schema snapshots.
type Diff struct {
	Tables []TableDiff `json:"tables"`
}

// HasChanges returns true if any differences were detected.
func (d Diff) HasChanges() bool {
	return len(d.Tables) > 0
}

// DiffSnapshots compares source and target snapshots to compute structural diffs.
func DiffSnapshots(src, tgt *snapshot.Snapshot) Diff {
	diff := Diff{}

	srcTables := make(map[string]snapshot.Table)
	for _, t := range src.Tables {
		srcTables[t.Name] = t
	}

	tgtTables := make(map[string]snapshot.Table)
	for _, t := range tgt.Tables {
		tgtTables[t.Name] = t
	}

	// Detect removed or modified tables
	for name, srcTbl := range srcTables {
		tgtTbl, exists := tgtTables[name]
		if !exists {
			diff.Tables = append(diff.Tables, TableDiff{
				TableName: name,
				Status:    "removed",
			})
			continue
		}

		tblDiff := diffTables(srcTbl, tgtTbl)
		if len(tblDiff.Columns) > 0 {
			diff.Tables = append(diff.Tables, tblDiff)
		}
	}

	// Detect added tables
	for name := range tgtTables {
		if _, exists := srcTables[name]; !exists {
			diff.Tables = append(diff.Tables, TableDiff{
				TableName: name,
				Status:    "added",
			})
		}
	}

	return diff
}

func diffTables(src, tgt snapshot.Table) TableDiff {
	td := TableDiff{
		TableName: src.Name,
		Status:    "modified",
	}

	srcCols := make(map[string]snapshot.Column)
	for _, c := range src.Columns {
		srcCols[c.Name] = c
	}

	tgtCols := make(map[string]snapshot.Column)
	for _, c := range tgt.Columns {
		tgtCols[c.Name] = c
	}

	for name, srcCol := range srcCols {
		tgtCol, exists := tgtCols[name]
		if !exists {
			td.Columns = append(td.Columns, ColumnDiff{
				ColumnName: name,
				ChangeType: "removed",
			})
			continue
		}
		if srcCol.Type != tgtCol.Type {
			td.Columns = append(td.Columns, ColumnDiff{
				ColumnName: name,
				ChangeType: "type_changed",
				OldType:    srcCol.Type,
				NewType:    tgtCol.Type,
			})
		}
	}

	for name := range tgtCols {
		if _, exists := srcCols[name]; !exists {
			td.Columns = append(td.Columns, ColumnDiff{
				ColumnName: name,
				ChangeType: "added",
			})
		}
	}

	return td
}
