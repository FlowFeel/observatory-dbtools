package compare_test

import (
	"testing"

	"github.com/FlowFeel/observatory-dbtools/pkg/compare"
	"github.com/FlowFeel/observatory-dbtools/pkg/snapshot"
)

func TestDiffSnapshots_Identical(t *testing.T) {
	s1 := &snapshot.Snapshot{
		DatabaseName: "mediawiki",
		Tables: []snapshot.Table{
			{
				Name: "page",
				Columns: []snapshot.Column{
					{Name: "page_id", Type: "int", IsNullable: "NO"},
				},
			},
		},
	}

	diff := compare.DiffSnapshots(s1, s1)
	if diff.HasChanges() {
		t.Errorf("expected no changes, got: %v", diff)
	}
}

func TestDiffSnapshots_DetectsTableAddedAndRemoved(t *testing.T) {
	s1 := &snapshot.Snapshot{
		DatabaseName: "mediawiki",
		Tables: []snapshot.Table{
			{Name: "table_a"},
		},
	}

	s2 := &snapshot.Snapshot{
		DatabaseName: "mediawiki",
		Tables: []snapshot.Table{
			{Name: "table_b"},
		},
	}

	diff := compare.DiffSnapshots(s1, s2)
	if !diff.HasChanges() {
		t.Fatal("expected changes, got none")
	}

	if len(diff.Tables) != 2 {
		t.Fatalf("expected 2 table diffs, got %d", len(diff.Tables))
	}
}
