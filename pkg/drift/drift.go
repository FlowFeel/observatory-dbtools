// Package drift detects and fixes SMW _MDAT routing mismatches
// between smw_fpt_mdat (fixed property table) and smw_di_time (data item table).
package drift

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

// Report holds the result of a drift check.
type Report struct {
	FPTCount    int    `json:"fpt_count"`
	DICount     int    `json:"di_count"`
	MissingInDI int    `json:"missing_in_di"`
	Description string `json:"description"`
}

// HasDrift returns true if there are FPT entries missing from DI.
func (r Report) HasDrift() bool {
	return r.MissingInDI > 0
}

// String returns a human-readable summary.
func (r Report) String() string {
	if r.HasDrift() {
		return fmt.Sprintf("DRIFT DETECTED: %d entries in smw_fpt_mdat missing from smw_di_time (p_id=29). FPT: %d, DI: %d",
			r.MissingInDI, r.FPTCount, r.DICount)
	}
	return fmt.Sprintf("No drift. FPT: %d, DI: %d", r.FPTCount, r.DICount)
}

// JSON returns the report as indented JSON.
func (r Report) JSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

// Check compares smw_fpt_mdat against smw_di_time for _MDAT (p_id=29) drift.
func Check(db *sql.DB) (*Report, error) {
	r := &Report{}

	err := db.QueryRow("SELECT COUNT(*) FROM smw_fpt_mdat").Scan(&r.FPTCount)
	if err != nil {
		return nil, fmt.Errorf("drift: count fpt: %w", err)
	}

	err = db.QueryRow("SELECT COUNT(*) FROM smw_di_time WHERE p_id = 29").Scan(&r.DICount)
	if err != nil {
		return nil, fmt.Errorf("drift: count di: %w", err)
	}

	err = db.QueryRow(`
		SELECT COUNT(*) FROM smw_fpt_mdat
		WHERE s_id NOT IN (SELECT s_id FROM smw_di_time WHERE p_id = 29)
	`).Scan(&r.MissingInDI)
	if err != nil {
		return nil, fmt.Errorf("drift: count missing: %w", err)
	}

	if r.HasDrift() {
		r.Description = fmt.Sprintf("%d _MDAT entries in FPT not mirrored to DI", r.MissingInDI)
	} else {
		r.Description = "All _MDAT entries in FPT are also in DI"
	}

	return r, nil
}

// Fix migrates missing _MDAT entries from smw_fpt_mdat to smw_di_time.
// Returns the number of rows inserted.
func Fix(db *sql.DB) (int64, error) {
	result, err := db.Exec(`
		INSERT INTO smw_di_time (s_id, p_id, o_serialized, o_sortkey)
		SELECT s_id, 29, o_serialized, o_sortkey FROM smw_fpt_mdat
		WHERE s_id NOT IN (SELECT s_id FROM smw_di_time WHERE p_id = 29)
	`)
	if err != nil {
		return 0, fmt.Errorf("drift: fix: %w", err)
	}

	return result.RowsAffected()
}
