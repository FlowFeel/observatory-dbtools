package audit

import (
	"database/sql"
	"fmt"

	"github.com/FlowFeel/observatory-dbtools/pkg/catalog"
)

// AuditDeclarations verifies Contract 1: every property declared in the
// catalog must be materialized in the SMW store with matching metadata.
//
// Checks per property:
//  1. smw_object_ids entry exists (smw_namespace=102, smw_title=name)
//  2. smw_fpt_type.o_serialized matches declared PropertyType
//  3. smw_fpt_pval rows match declared AllowedValues
//  4. smw_fpt_impo matches declared equivalentProperty
//
// Transactional state awareness: if a property page exists in the MW page
// table but smw_fpt_type is empty, the audit reports a warning (job-queue
// propagation lag), not an error (true declaration violation).
func AuditDeclarations(db *sql.DB, c *catalog.Catalog) (*Report, error) {
	r := &Report{Summary: Summary{PerProperty: make(map[string]int)}}

	for _, prop := range c.Properties {
		// 1. Identity allocation in smw_object_ids.
		var pID int
		err := db.QueryRow(`
			SELECT smw_id FROM smw_object_ids
			WHERE smw_namespace = 102 AND smw_title = ?
		`, catalog.SMWTitle(prop.Name)).Scan(&pID)

		if err == sql.ErrNoRows {
			// Check if a page exists but hasn't been parsed into SMW yet.
			var pageExists bool
			db.QueryRow(`
				SELECT EXISTS(SELECT 1 FROM page
				WHERE page_namespace = 102 AND page_title = ?)
			`, prop.Name).Scan(&pageExists)

			if pageExists {
				r.Add(Violation{
					Property:   prop.Name,
					Kind:       KindDeclaration,
					Severity:   SevWarning,
					Rule:       "identity_allocation",
					Diagnostic: fmt.Sprintf("property %q has a wiki page but no smw_object_ids entry — job-queue propagation lag", prop.Name),
				})
			} else {
				r.Add(Violation{
					Property:   prop.Name,
					Kind:       KindDeclaration,
					Severity:   SevError,
					Rule:       "identity_allocation",
					Diagnostic: fmt.Sprintf("property %q not found in smw_object_ids (smw_namespace=102)", prop.Name),
				})
			}
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("audit declarations: query smw_object_ids for %q: %w", prop.Name, err)
		}

		// 2. Type declaration in smw_fpt_type.
		var storedType string
		err = db.QueryRow(`SELECT o_serialized FROM smw_fpt_type WHERE s_id = ?`, pID).Scan(&storedType)
		if err == sql.ErrNoRows {
			// Transactional state awareness: page exists, smw_object_ids exists,
			// but smw_fpt_type not populated — job queue hasn't flushed.
			r.Add(Violation{
				Property:   prop.Name,
				Kind:       KindDeclaration,
				Severity:   SevWarning,
				Rule:       "type_declaration",
				Diagnostic: fmt.Sprintf("property %q has smw_object_ids entry but no smw_fpt_type row — SMW parser has not yet processed this property", prop.Name),
			})
		} else if err != nil {
			return nil, fmt.Errorf("audit declarations: query smw_fpt_type for %q: %w", prop.Name, err)
		} else if !TypeMatchesStoredValue(prop.Type, storedType) {
			r.Add(Violation{
				Property:   prop.Name,
				Kind:       KindDeclaration,
				Severity:   SevError,
				Rule:       "type_mismatch",
				Diagnostic: fmt.Sprintf("property %q declared type %q but smw_fpt_type stores %q", prop.Name, prop.Type, storedType),
			})
		}

		// 3. Allowed values in smw_fpt_pval (only when declared).
		if len(prop.Allowed) > 0 {
			rows, err := db.Query(`SELECT o_hash FROM smw_fpt_pval WHERE s_id = ?`, pID)
			if err != nil {
				return nil, fmt.Errorf("audit declarations: query smw_fpt_pval for %q: %w", prop.Name, err)
			}
			declaredSet := make(map[string]bool)
			for _, a := range prop.Allowed {
				declaredSet[NormalizeLiteral(a)] = true
			}
			foundSet := make(map[string]bool)
			for rows.Next() {
				var hash string
				if err := rows.Scan(&hash); err != nil {
					rows.Close()
					return nil, fmt.Errorf("audit declarations: scan smw_fpt_pval: %w", err)
				}
				foundSet[NormalizeLiteral(hash)] = true
			}
			rows.Close()

			for declared := range declaredSet {
				if !foundSet[declared] {
					r.Add(Violation{
						Property:   prop.Name,
						Kind:       KindDeclaration,
						Severity:   SevError,
						Rule:       "allowed_value_missing",
						Diagnostic: fmt.Sprintf("property %q declared allowed value %q but no smw_fpt_pval row found", prop.Name, declared),
					})
				}
			}
		}

		// 4. Equivalent property in smw_fpt_impo (only when declared).
		if prop.Equivalent != "" {
			var storedImmo string
			err = db.QueryRow(`SELECT o_hash FROM smw_fpt_impo WHERE s_id = ?`, pID).Scan(&storedImmo)
			if err == sql.ErrNoRows {
				r.Add(Violation{
					Property:   prop.Name,
					Kind:       KindDeclaration,
					Severity:   SevError,
					Rule:       "equivalent_missing",
					Diagnostic: fmt.Sprintf("property %q declares equivalentProperty %q but no smw_fpt_impo row found", prop.Name, prop.Equivalent),
				})
			} else if err != nil {
				return nil, fmt.Errorf("audit declarations: query smw_fpt_impo for %q: %w", prop.Name, err)
			} else if NormalizeLiteral(storedImmo) != NormalizeLiteral(prop.Equivalent) {
				r.Add(Violation{
					Property:   prop.Name,
					Kind:       KindDeclaration,
					Severity:   SevError,
					Rule:       "equivalent_mismatch",
					Diagnostic: fmt.Sprintf("property %q declares equivalentProperty %q but smw_fpt_impo stores %q", prop.Name, prop.Equivalent, storedImmo),
				})
			}
		}

		r.Summary.RowsScanned++
	}

	return r.Finalize(), nil
}
