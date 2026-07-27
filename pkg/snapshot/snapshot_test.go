package snapshot_test

import (
	"testing"

	"github.com/FlowFeel/observatory-dbtools/pkg/snapshot"
)

func TestSnapshot_CanonicalSorting(t *testing.T) {
	snap := snapshot.Snapshot{
		DatabaseName: "mediawiki",
		Tables: []snapshot.Table{
			{
				Name: "user",
				Columns: []snapshot.Column{
					{Name: "user_id", Type: "int", IsNullable: "NO"},
					{Name: "user_name", Type: "varchar", IsNullable: "NO"},
				},
			},
		},
	}

	if len(snap.Tables) != 1 {
		t.Fatalf("expected 1 table, got %d", len(snap.Tables))
	}
	if snap.Tables[0].Name != "user" {
		t.Errorf("expected table user, got %s", snap.Tables[0].Name)
	}
}
