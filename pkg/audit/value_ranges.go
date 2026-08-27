package audit

import (
	"database/sql"
	"fmt"

	"github.com/FlowFeel/observatory-dbtools/pkg/catalog"
)

// AuditValueRanges verifies Contract 3: all persisted literal values conform
// to range restrictions and syntactic formatting constraints.
//
// Three invariant topologies:
//  1. Discrete set membership — smw_di_blob literals checked against
//     declared AllowedValues (exact character match after canonical trim)
//  2. Scalar syntactic — ISO 8601 dates in smw_di_time, RFC 3986 URIs in
//     smw_di_uri
//  3. Relational reference — smw_di_wikipage.o_id must resolve to an
//     active entry in smw_object_ids
//
// Keyset streaming per property-targeted batch. No premature halting.
// Sampling bounds cap raw violation instances per property.
func AuditValueRanges(db *sql.DB, c *catalog.Catalog, opts AuditOptions) (*Report, error) {
	if opts.BatchSize == 0 {
		opts = DefaultOptions()
	}
	r := &Report{Summary: Summary{PerProperty: make(map[string]int)}}
	violationCount := make(map[string]int)

	// Build p_id lookup for all catalog properties.
	propByPID := make(map[int]propRef)
	for _, prop := range c.Properties {
		var pID int
		err := db.QueryRow(`
			SELECT smw_id FROM smw_object_ids
			WHERE smw_namespace = 102 AND smw_title = ?
		`, prop.Name).Scan(&pID)
		if err == sql.ErrNoRows {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("audit value ranges: resolve p_id for %q: %w", prop.Name, err)
		}
		propByPID[pID] = propRef{Name: prop.Name, Type: prop.Type, Allowed: prop.Allowed}
	}

	// 1. Discrete set membership — scan smw_di_blob for enumerated properties.
	auditBlobSetMembership(db, r, propByPID, opts, violationCount)

	// 2a. Scalar syntactic — ISO 8601 dates in smw_di_time.
	auditTimeSyntax(db, r, propByPID, opts, violationCount)

	// 2b. Scalar syntactic — RFC 3986 URIs in smw_di_uri.
	auditURISyntax(db, r, propByPID, opts, violationCount)

	// 3. Relational reference — smw_di_wikipage.o_id → smw_object_ids.
	auditWikipageReferences(db, r, propByPID, opts, violationCount)

	return r.Finalize(), nil
}

// withinSamplingBounds reports whether we should still collect raw violation
// instances for this property. The global counter always increments.
type propRef struct {
	Name    string
	Type    string
	Allowed []string
}

func withinSamplingBounds(propName string, count map[string]int, max int) bool {
	return count[propName] < max
}

func auditBlobSetMembership(db *sql.DB, r *Report, propByPID map[int]propRef, opts AuditOptions, violationCount map[string]int) {
	lastSeen := 0
	for {
		rows, err := db.Query(fmt.Sprintf(`
			SELECT s_id, p_id, IFNULL(o_hash, '') FROM smw_di_blob
			WHERE s_id > ? ORDER BY s_id ASC LIMIT %d
		`, opts.BatchSize), lastSeen)
		if err != nil {
			return
		}
		batchCount := 0
		for rows.Next() {
			var sID, pID int
			var hash string
			rows.Scan(&sID, &pID, &hash)
			batchCount++
			lastSeen = sID
			r.Summary.RowsScanned++

			ref, known := propByPID[pID]
			if !known || len(ref.Allowed) == 0 {
				continue // not enumerated or orphaned (Contract 2 handles)
			}
			if !InAllowedSet(hash, ref.Allowed) {
				violationCount[ref.Name]++
				if withinSamplingBounds(ref.Name, violationCount, opts.MaxViolationsPerProperty) {
					r.Add(Violation{
						Table:      "smw_di_blob",
						EntityID:   fmt.Sprintf("%d", sID),
						Property:   ref.Name,
						Value:      hash,
						Kind:       KindRange,
						Severity:   SevError,
						Rule:       "set_membership",
						Diagnostic: fmt.Sprintf("value %q not in declared allowed values %v", hash, ref.Allowed),
					})
				}
			}
		}
		rows.Close()
		if batchCount < opts.BatchSize {
			break
		}
	}
}

func auditTimeSyntax(db *sql.DB, r *Report, propByPID map[int]propRef, opts AuditOptions, violationCount map[string]int) {
	lastSeen := 0
	for {
		rows, err := db.Query(fmt.Sprintf(`
			SELECT s_id, p_id, IFNULL(o_serialized, '') FROM smw_di_time
			WHERE s_id > ? ORDER BY s_id ASC LIMIT %d
		`, opts.BatchSize), lastSeen)
		if err != nil {
			return
		}
		batchCount := 0
		for rows.Next() {
			var sID, pID int
			var serialized string
			rows.Scan(&sID, &pID, &serialized)
			batchCount++
			lastSeen = sID
			r.Summary.RowsScanned++

			ref, known := propByPID[pID]
			if !known {
				continue
			}
			if !validISODate(serialized) {
				violationCount[ref.Name]++
				if withinSamplingBounds(ref.Name, violationCount, opts.MaxViolationsPerProperty) {
					r.Add(Violation{
						Table:      "smw_di_time",
						EntityID:   fmt.Sprintf("%d", sID),
						Property:   ref.Name,
						Value:      serialized,
						Kind:       KindSyntax,
						Severity:   SevError,
						Rule:       "iso8601_temporal",
						Diagnostic: fmt.Sprintf("value %q is not a valid ISO 8601 date", serialized),
					})
				}
			}
		}
		rows.Close()
		if batchCount < opts.BatchSize {
			break
		}
	}
}

func auditURISyntax(db *sql.DB, r *Report, propByPID map[int]propRef, opts AuditOptions, violationCount map[string]int) {
	lastSeen := 0
	for {
		rows, err := db.Query(fmt.Sprintf(`
			SELECT s_id, p_id, IFNULL(o_serialized, '') FROM smw_di_uri
			WHERE s_id > ? ORDER BY s_id ASC LIMIT %d
		`, opts.BatchSize), lastSeen)
		if err != nil {
			return
		}
		batchCount := 0
		for rows.Next() {
			var sID, pID int
			var serialized string
			rows.Scan(&sID, &pID, &serialized)
			batchCount++
			lastSeen = sID
			r.Summary.RowsScanned++

			ref, known := propByPID[pID]
			if !known {
				continue
			}
			if !validURI(serialized) {
				violationCount[ref.Name]++
				if withinSamplingBounds(ref.Name, violationCount, opts.MaxViolationsPerProperty) {
					r.Add(Violation{
						Table:      "smw_di_uri",
						EntityID:   fmt.Sprintf("%d", sID),
						Property:   ref.Name,
						Value:      serialized,
						Kind:       KindSyntax,
						Severity:   SevError,
						Rule:       "rfc3986_uri",
						Diagnostic: fmt.Sprintf("value %q is not a valid RFC 3986 absolute URI", serialized),
					})
				}
			}
		}
		rows.Close()
		if batchCount < opts.BatchSize {
			break
		}
	}
}

func auditWikipageReferences(db *sql.DB, r *Report, propByPID map[int]propRef, opts AuditOptions, violationCount map[string]int) {
	lastSeen := 0
	for {
		rows, err := db.Query(fmt.Sprintf(`
			SELECT s_id, p_id, o_id FROM smw_di_wikipage
			WHERE s_id > ? ORDER BY s_id ASC LIMIT %d
		`, opts.BatchSize), lastSeen)
		if err != nil {
			return
		}
		batchCount := 0
		for rows.Next() {
			var sID, pID, oID int
			rows.Scan(&sID, &pID, &oID)
			batchCount++
			lastSeen = sID
			r.Summary.RowsScanned++

			ref, known := propByPID[pID]
			if !known {
				continue
			}
			// Check o_id resolves to an active entry in smw_object_ids.
			var exists bool
			db.QueryRow(`SELECT EXISTS(SELECT 1 FROM smw_object_ids WHERE smw_id = ?)`, oID).Scan(&exists)
			if !exists {
				violationCount[ref.Name]++
				if withinSamplingBounds(ref.Name, violationCount, opts.MaxViolationsPerProperty) {
					r.Add(Violation{
						Table:      "smw_di_wikipage",
						EntityID:   fmt.Sprintf("%d", sID),
						Property:   ref.Name,
						Value:      fmt.Sprintf("o_id=%d", oID),
						Kind:       KindReference,
						Severity:   SevError,
						Rule:       "relational_reference",
						Diagnostic: fmt.Sprintf("smw_di_wikipage o_id=%d does not resolve to any smw_object_ids entry", oID),
					})
				}
			}
		}
		rows.Close()
		if batchCount < opts.BatchSize {
			break
		}
	}
}
