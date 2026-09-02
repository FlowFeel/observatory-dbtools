package audit

import (
	"database/sql"
	"fmt"

	"github.com/FlowFeel/observatory-dbtools/pkg/catalog"
)

// AuditValueTypes verifies Contract 2: every stored value must reside in
// the data-item table designated for its property's declared type.
//
// For each catalog property:
//   - Resolve p_id from smw_object_ids
//   - Check all smw_di_* tables EXCEPT the expected one for orphaned rows
//     with that p_id (historical type drift — values trapped in deprecated
//     physical tables after a type change)
//
// Orphaned predicate detection: rows with p_id that doesn't resolve to any
// catalog property are captured and reported, never silently discarded.
//
// Uses keyset streaming: WHERE s_id > last_seen ORDER BY s_id ASC LIMIT N
func AuditValueTypes(db *sql.DB, c *catalog.Catalog, opts AuditOptions) (*Report, error) {
	if opts.BatchSize == 0 {
		opts = DefaultOptions()
	}
	r := &Report{Summary: Summary{PerProperty: make(map[string]int)}}

	// Build p_id lookup map from smw_object_ids for all catalog properties.
	type propRef struct {
		Name     string
		Type     string
		Expected string // expected table
	}
	propByPID := make(map[int]propRef)
	for _, prop := range c.Properties {
		var pID int
		err := db.QueryRow(`
			SELECT smw_id FROM smw_object_ids
			WHERE smw_namespace = 102 AND smw_title = ?
		`, catalog.SMWTitle(prop.Name)).Scan(&pID)
		if err == sql.ErrNoRows {
			continue // property not in DB — Contract 1 will report this
		}
		if err != nil {
			return nil, fmt.Errorf("audit value types: resolve p_id for %q: %w", prop.Name, err)
		}
		expected, ok := ExpectedTable(prop.Type)
		if !ok {
			continue // unknown type — skip routing check
		}
		propByPID[pID] = propRef{Name: prop.Name, Type: prop.Type, Expected: expected}
	}

	// Scan each data-item table for routing violations and orphaned predicates.
	for _, table := range dataItemTables {
		lastSeen := 0
		for {
			rows, err := db.Query(fmt.Sprintf(`
				SELECT s_id, p_id FROM %s
				WHERE s_id > ? ORDER BY s_id ASC LIMIT %d
			`, table, opts.BatchSize), lastSeen)
			if err != nil {
				return nil, fmt.Errorf("audit value types: query %s: %w", table, err)
			}

			batchCount := 0
			for rows.Next() {
				var sID, pID int
				if err := rows.Scan(&sID, &pID); err != nil {
					rows.Close()
					return nil, fmt.Errorf("audit value types: scan %s: %w", table, err)
				}
				batchCount++
				lastSeen = sID
				r.Summary.RowsScanned++

				ref, known := propByPID[pID]
				if !known {
					// Orphaned predicate: p_id doesn't resolve to any catalog property.
					r.Add(Violation{
						Table:      table,
						EntityID:   fmt.Sprintf("%d", sID),
						Property:   fmt.Sprintf("p_id=%d", pID),
						Kind:       KindOrphanedPredicate,
						Severity:   SevError,
						Rule:       "orphaned_predicate",
						Diagnostic: fmt.Sprintf("row in %s (s_id=%d) references p_id=%d which has no catalog declaration", table, sID, pID),
					})
					continue
				}

				if table != ref.Expected {
					// Routing violation: value stored in wrong table.
					r.Add(Violation{
						Table:      table,
						EntityID:   fmt.Sprintf("%d", sID),
						Property:   ref.Name,
						Kind:       KindRouting,
						Severity:   SevError,
						Rule:       "table_routing",
						Diagnostic: fmt.Sprintf("property %q (declared %q) has value in %s but expected %s — historical type drift", ref.Name, ref.Type, table, ref.Expected),
					})
				}
			}
			rows.Close()

			if batchCount < opts.BatchSize {
				break // table exhausted
			}
		}
	}

	return r.Finalize(), nil
}
