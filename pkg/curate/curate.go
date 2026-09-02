// Package curate generates contract-tested, sanitized database seeds for Test, Stage, and Prod.
//
// Curation is catalog-driven: required property pages are derived from the
// compiled property catalog, not hardcoded lists. This eliminates the
// "magic string" anti-pattern where domain knowledge was embedded in code
// rather than derived from the source of truth.
package curate

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/FlowFeel/observatory-dbtools/pkg/catalog"
)

// RequiredPortalPages lists mandatory area portals that MUST exist in canonical seeds.
// These are structural pages, not semantic properties — they come from the
// MW page structure, not the property catalog.
var RequiredPortalPages = []string{
	"Main_Page",
	"Animals",
	"Classics",
	"Dig_Labs",
	"Human_Bridges",
}

// Plan represents dataset curation rules for a tier.
type Plan struct {
	Tier               string   `json:"tier"`
	RequiredPages      []string `json:"required_pages"`
	RequiredProperties []string `json:"required_properties"`
	AnonymizeUsers     bool     `json:"anonymize_users"`
	IncludeJobQueue    bool     `json:"include_job_queue"`
}

// NewPlan returns the curation policy for a given tier, with property pages
// derived from the loaded catalog.
func NewPlan(tier string, c *catalog.Catalog) (*Plan, error) {
	switch strings.ToLower(tier) {
	case "test":
		return &Plan{
			Tier:               "test",
			RequiredPages:      RequiredPortalPages,
			RequiredProperties: catalogPropertyPages(c),
			AnonymizeUsers:     true,
			IncludeJobQueue:    false,
		}, nil
	case "stage":
		return &Plan{
			Tier:               "stage",
			RequiredPages:      RequiredPortalPages,
			RequiredProperties: catalogPropertyPages(c),
			AnonymizeUsers:     true,
			IncludeJobQueue:    false,
		}, nil
	case "prod":
		return &Plan{
			Tier:               "prod",
			RequiredPages:      RequiredPortalPages,
			RequiredProperties: catalogPropertyPages(c),
			AnonymizeUsers:     false,
			IncludeJobQueue:    false,
		}, nil
	default:
		return nil, fmt.Errorf("curate: unknown tier %q (valid: test, stage, prod)", tier)
	}
}

// catalogPropertyPages generates the list of Property: pages from the catalog.
// Returns an empty list when no catalog is loaded (backwards-compatible fallback).
func catalogPropertyPages(c *catalog.Catalog) []string {
	if c == nil || len(c.Properties) == 0 {
		return []string{}
	}
	pages := make([]string, 0, len(c.Properties))
	for _, prop := range c.Properties {
		pages = append(pages, "Property:"+strings.ReplaceAll(prop.Name, " ", "_"))
	}
	return pages
}

// ValidateSeed validates that a SQL dump contains all required pages.
func ValidateSeed(content string, plan *Plan) error {
	var missing []string
	for _, page := range plan.RequiredPages {
		cleanName := strings.ReplaceAll(page, "_", " ")
		if !strings.Contains(content, cleanName) && !strings.Contains(content, page) {
			missing = append(missing, page)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("curate: seed missing required portal pages: %v", missing)
	}
	return nil
}

// ExportCleanSeed writes a curated minimal SQL baseline file.
func ExportCleanSeed(outputPath string, content string, plan *Plan) error {
	if err := ValidateSeed(content, plan); err != nil {
		return err
	}
	return os.WriteFile(outputPath, []byte(content), 0o644)
}

// VerifyDatabase verifies live DB meets the tier plan requirements.
func VerifyDatabase(db *sql.DB, plan *Plan) error {
	for _, page := range plan.RequiredPages {
		cleanName := strings.ReplaceAll(page, "_", " ")
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM page WHERE page_title = ? OR page_title = ?", page, cleanName).Scan(&count)
		if err != nil {
			return fmt.Errorf("curate verify: query page %s: %w", page, err)
		}
		if count == 0 {
			return fmt.Errorf("curate verify: required page %s missing from database", page)
		}
	}
	return nil
}
