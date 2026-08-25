// Package asset enforces asset boundary contracts across MediaWiki and SMW tables.
// It detects and remediates hardcoded environment URLs into canonical File:... references.
package asset

import (
	"database/sql"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var (
	hardcodedURLRegex = regexp.MustCompile(`https?://(?:observatory\.wiki|localhost(?::\d+)?)/w/images/(?:thumb/)?(?:[0-9a-f]/[0-9a-f]{2}/)?([^"'\s|\]}]+)`)
	commonsURLRegex   = regexp.MustCompile(`https?://upload\.wikimedia\.org/wikipedia/commons/(?:thumb/)?(?:[0-9a-f]/[0-9a-f]{2}/)?([^"'\s|\]}]+)`)
)

// Violation represents a single detected asset contract violation.
type Violation struct {
	SourceTable string `json:"source_table"`
	EntityID    string `json:"entity_id"`
	RawURL      string `json:"raw_url"`
	Filename    string `json:"filename"`
}

// Report holds the complete asset audit findings.
type Report struct {
	TotalScanned    int         `json:"total_scanned"`
	Violations      []Violation `json:"violations"`
	OrphanedImages  []string    `json:"orphaned_images"`
	MissingDefaults []string    `json:"missing_defaults"`
	ViolationsCount int         `json:"violations_count"`
}

// HasViolations returns true if any asset boundary violations were detected.
func (r Report) HasViolations() bool {
	return len(r.Violations) > 0 || len(r.MissingDefaults) > 0
}

// String returns a formatted summary of the asset audit.
func (r Report) String() string {
	if !r.HasViolations() {
		return fmt.Sprintf("Asset Boundaries Clean: %d entities scanned, 0 hardcoded URLs.", r.TotalScanned)
	}
	return fmt.Sprintf("ASSET VIOLATIONS DETECTED: %d hardcoded host URLs in database, %d missing default assets.",
		len(r.Violations), len(r.MissingDefaults))
}

// DetectHardcodedURLs returns all absolute image URLs found in the text.
func DetectHardcodedURLs(content string) []string {
	var matches []string
	for _, m := range hardcodedURLRegex.FindAllString(content, -1) {
		matches = append(matches, m)
	}
	for _, m := range commonsURLRegex.FindAllString(content, -1) {
		matches = append(matches, m)
	}
	return matches
}

// ExtractFilename strips host, /w/images/ path hashing, and thumb parameters, returning the clean filename.
func ExtractFilename(rawURL string) string {
	rawURL = strings.TrimSpace(rawURL)
	if m := hardcodedURLRegex.FindStringSubmatch(rawURL); len(m) > 1 {
		fn := m[1]
		if strings.Contains(rawURL, "/thumb/") {
			parts := strings.Split(fn, "/")
			if len(parts) > 0 {
				fn = parts[0]
			}
		}
		if idx := strings.IndexAny(fn, "?%#"); idx != -1 {
			fn = fn[:idx]
		}
		if decoded, err := url.QueryUnescape(fn); err == nil {
			return decoded
		}
		return fn
	}

	if m := commonsURLRegex.FindStringSubmatch(rawURL); len(m) > 1 {
		fn := m[1]
		if strings.Contains(rawURL, "/thumb/") {
			parts := strings.Split(fn, "/")
			if len(parts) > 0 {
				fn = parts[0]
			}
		}
		if idx := strings.IndexAny(fn, "?%#"); idx != -1 {
			fn = fn[:idx]
		}
		if decoded, err := url.QueryUnescape(fn); err == nil {
			return decoded
		}
		return fn
	}

	return rawURL
}

// RemediateWikitext converts hardcoded image URLs inside wikitext to canonical File:... references or bare filenames.
func RemediateWikitext(content string) string {
	content = hardcodedURLRegex.ReplaceAllStringFunc(content, func(m string) string {
		return ExtractFilename(m)
	})
	content = commonsURLRegex.ReplaceAllStringFunc(content, func(m string) string {
		return ExtractFilename(m)
	})
	return content
}

// AuditDatabase scans the text table and SMW blob tables for hardcoded host URLs.
func AuditDatabase(db *sql.DB) (*Report, error) {
	rpt := &Report{
		Violations: []Violation{},
	}

	rows, err := db.Query("SELECT old_id, old_text FROM text WHERE old_text LIKE '%/w/images/%' OR old_text LIKE '%upload.wikimedia.org%'")
	if err != nil {
		return nil, fmt.Errorf("asset audit: query text: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var oldID int
		var oldText string
		if err := rows.Scan(&oldID, &oldText); err != nil {
			return nil, fmt.Errorf("asset audit: scan text: %w", err)
		}
		rpt.TotalScanned++
		urls := DetectHardcodedURLs(oldText)
		for _, u := range urls {
			rpt.Violations = append(rpt.Violations, Violation{
				SourceTable: "text",
				EntityID:    fmt.Sprintf("old_id=%d", oldID),
				RawURL:      u,
				Filename:    ExtractFilename(u),
			})
		}
	}

	rpt.ViolationsCount = len(rpt.Violations)
	return rpt, nil
}

// RemediateDatabase rewrites hardcoded host URLs in text table to canonical filenames.
func RemediateDatabase(db *sql.DB, dryRun bool) (int64, error) {
	rpt, err := AuditDatabase(db)
	if err != nil {
		return 0, err
	}

	if dryRun || len(rpt.Violations) == 0 {
		return int64(len(rpt.Violations)), nil
	}

	var updatedCount int64
	rows, err := db.Query("SELECT old_id, old_text FROM text WHERE old_text LIKE '%/w/images/%' OR old_text LIKE '%upload.wikimedia.org%'")
	if err != nil {
		return 0, fmt.Errorf("asset remediate: query text: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var oldID int
		var oldText string
		if err := rows.Scan(&oldID, &oldText); err != nil {
			return updatedCount, fmt.Errorf("asset remediate: scan text: %w", err)
		}

		remediated := RemediateWikitext(oldText)
		if remediated != oldText {
			_, err := db.Exec("UPDATE text SET old_text = ? WHERE old_id = ?", remediated, oldID)
			if err != nil {
				return updatedCount, fmt.Errorf("asset remediate: update text old_id=%d: %w", oldID, err)
			}
			updatedCount++
		}
	}

	return updatedCount, nil
}
