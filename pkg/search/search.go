// Package search audits Elasticsearch index document counts, aliases, and queue state vs MySQL.
package search

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Report holds search health and index document count comparison details.
type Report struct {
	ESHost         string `json:"es_host"`
	ESConnected    bool   `json:"es_connected"`
	ContentDocs    int    `json:"content_docs"`
	GeneralDocs    int    `json:"general_docs"`
	MySQLPages     int    `json:"mysql_pages"`
	HasDrift       bool   `json:"has_drift"`
	QueuedJobCount int    `json:"queued_job_count"`
	Description    string `json:"description"`
}

func (r Report) String() string {
	if !r.ESConnected {
		return fmt.Sprintf("SEARCH ERROR: Cannot connect to Elasticsearch at %s", r.ESHost)
	}
	if r.HasDrift {
		return fmt.Sprintf("SEARCH DRIFT DETECTED: ES Content Docs=%d, ES General Docs=%d, MySQL Pages=%d (Queued jobs: %d)",
			r.ContentDocs, r.GeneralDocs, r.MySQLPages, r.QueuedJobCount)
	}
	return fmt.Sprintf("Search Healthy. ES Content Docs: %d, General Docs: %d, MySQL Pages: %d (Queued jobs: %d)",
		r.ContentDocs, r.GeneralDocs, r.MySQLPages, r.QueuedJobCount)
}

type catIndexResponse struct {
	Index     string `json:"index"`
	DocsCount string `json:"docs.count"`
}

// Audit compares Elasticsearch indexed documents with MySQL page counts.
func Audit(db *sql.DB, esHost string) (*Report, error) {
	r := &Report{
		ESHost: esHost,
	}

	// 1. Check MySQL page table count
	err := db.QueryRow("SELECT COUNT(*) FROM page WHERE page_namespace = 0").Scan(&r.MySQLPages)
	if err != nil {
		return nil, fmt.Errorf("search: count mysql pages: %w", err)
	}

	// 2. Check pending CirrusSearch job queue depth
	_ = db.QueryRow("SELECT COUNT(*) FROM job WHERE job_cmd LIKE 'cirrusSearch%'").Scan(&r.QueuedJobCount)

	// 3. Query Elasticsearch cat indices API
	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("%s/_cat/indices?format=json", esHost)
	resp, err := client.Get(url)
	if err != nil {
		r.ESConnected = false
		r.Description = fmt.Sprintf("Elasticsearch unreachable: %v", err)
		return r, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		r.ESConnected = false
		r.Description = fmt.Sprintf("Elasticsearch returned status %d", resp.StatusCode)
		return r, nil
	}

	r.ESConnected = true

	var indices []catIndexResponse
	if err := json.NewDecoder(resp.Body).Decode(&indices); err != nil {
		return nil, fmt.Errorf("search: decode ES cat indices: %w", err)
	}

	for _, idx := range indices {
		var cnt int
		fmt.Sscanf(idx.DocsCount, "%d", &cnt)
		if idx.Index == "mediawiki_content_1770474050" || idx.Index == "mediawiki_content" {
			r.ContentDocs += cnt
		}
		if idx.Index == "mediawiki_general_1770474050" || idx.Index == "mediawiki_general" {
			r.GeneralDocs += cnt
		}
	}

	// Drift condition: 0 content docs when MySQL pages > 0, or ES content docs < 50% of pages when queue is 0
	if (r.ContentDocs == 0 && r.MySQLPages > 0) || (r.QueuedJobCount == 0 && r.ContentDocs < r.MySQLPages/2) {
		r.HasDrift = true
		r.Description = fmt.Sprintf("Index document deficit: ES content docs (%d) vs MySQL pages (%d)", r.ContentDocs, r.MySQLPages)
	} else {
		r.HasDrift = false
		r.Description = "Search indices match target thresholds"
	}

	return r, nil
}
