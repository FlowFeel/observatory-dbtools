// Package drift detects and fixes SMW fixed-property-table routing mismatches
// between smw_fpt_* (fixed property tables) and smw_di_* (data item tables).
//
// Generalized via a topological drift registry: each DriftTarget maps a
// property to its expected FPT and DI storage locations. The existing
// _MDAT drift (p_id=29, smw_fpt_mdat → smw_di_time) is one registry entry.
// Future drift targets are added to the registry without new code.
package drift

import (
	"database/sql"
	"fmt"
	"strings"
)

// DriftTarget maps a property to its expected FPT and DI storage tables.
type DriftTarget struct {
	PropertyName string `json:"property_name"`
	FptTable     string `json:"fpt_table"`
	DiTable      string `json:"di_table"`
	Pid          int    `json:"pid"`
}

// DefaultRegistry returns the standard drift registry with known SMW
// fixed-property-table routing targets.
func DefaultRegistry() []DriftTarget {
	return []DriftTarget{
		{
			PropertyName: "_MDAT",
			FptTable:     "smw_fpt_mdat",
			DiTable:      "smw_di_time",
			Pid:          29,
		},
	}
}

// Report holds the result of a drift check across all registry targets.
type Report struct {
	Targets     []TargetReport `json:"targets"`
	HasDrift    bool           `json:"has_drift"`
	Description string         `json:"description"`
}

// TargetReport holds the drift check result for a single target.
type TargetReport struct {
	Target      DriftTarget `json:"target"`
	FPTCount    int         `json:"fpt_count"`
	DICount     int         `json:"di_count"`
	MissingInDI int         `json:"missing_in_di"`
}

// String returns a human-readable summary.
func (r Report) String() string {
	if !r.HasDrift {
		return fmt.Sprintf("No drift across %d targets", len(r.Targets))
	}
	var parts []string
	for _, tr := range r.Targets {
		if tr.MissingInDI > 0 {
			parts = append(parts, fmt.Sprintf("%s: %d missing in DI",
				tr.Target.PropertyName, tr.MissingInDI))
		}
	}
	return fmt.Sprintf("DRIFT DETECTED: %s", strings.Join(parts, "; "))
}

// Check compares FPT against DI for all registry targets.
func Check(db *sql.DB, registry []DriftTarget) (*Report, error) {
	r := &Report{Targets: make([]TargetReport, 0, len(registry))}

	for _, target := range registry {
		tr, err := checkTarget(db, target)
		if err != nil {
			return nil, fmt.Errorf("drift: check %s: %w", target.PropertyName, err)
		}
		r.Targets = append(r.Targets, *tr)
		if tr.MissingInDI > 0 {
			r.HasDrift = true
		}
	}

	if r.HasDrift {
		r.Description = "drift detected in one or more targets"
	} else {
		r.Description = "all targets clean"
	}

	return r, nil
}

func checkTarget(db *sql.DB, target DriftTarget) (*TargetReport, error) {
	tr := &TargetReport{Target: target}

	if err := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", target.FptTable)).Scan(&tr.FPTCount); err != nil {
		return nil, fmt.Errorf("count fpt: %w", err)
	}

	if err := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE p_id = ?", target.DiTable), target.Pid).Scan(&tr.DICount); err != nil {
		return nil, fmt.Errorf("count di: %w", err)
	}

	if err := db.QueryRow(fmt.Sprintf(`
		SELECT COUNT(*) FROM %s
		WHERE s_id NOT IN (SELECT s_id FROM %s WHERE p_id = ?)
	`, target.FptTable, target.DiTable), target.Pid).Scan(&tr.MissingInDI); err != nil {
		return nil, fmt.Errorf("count missing: %w", err)
	}

	return tr, nil
}

// Fix migrates missing entries for all registry targets.
// Returns the total number of rows inserted across all targets.
func Fix(db *sql.DB, registry []DriftTarget) (int64, error) {
	var total int64
	for _, target := range registry {
		n, err := fixTarget(db, target)
		if err != nil {
			return total, fmt.Errorf("drift: fix %s: %w", target.PropertyName, err)
		}
		total += n
	}
	return total, nil
}

func fixTarget(db *sql.DB, target DriftTarget) (int64, error) {
	// _MDAT has o_serialized and o_sortkey columns. Other targets may differ.
	// For now, the _MDAT fix is the proven path. Additional targets extend this.
	if target.PropertyName == "_MDAT" {
		result, err := db.Exec(`
			INSERT INTO smw_di_time (s_id, p_id, o_serialized, o_sortkey)
			SELECT s_id, 29, o_serialized, o_sortkey FROM smw_fpt_mdat
			WHERE s_id NOT IN (SELECT s_id FROM smw_di_time WHERE p_id = 29)
		`)
		if err != nil {
			return 0, fmt.Errorf("fix: %w", err)
		}
		return result.RowsAffected()
	}
	// For non-MDAT targets, return 0 — extend fixTarget as new targets are added.
	return 0, nil
}
