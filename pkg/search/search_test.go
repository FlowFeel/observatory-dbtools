package search_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FlowFeel/observatory-dbtools/pkg/connect"
	"github.com/FlowFeel/observatory-dbtools/pkg/search"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	ctx := context.Background()

	container, err := mysql.Run(ctx, "mysql:8.4",
		mysql.WithDatabase("mediawiki"),
		mysql.WithUsername("root"),
		mysql.WithPassword("test"),
	)
	if err != nil {
		t.Fatalf("failed to start MySQL container: %v", err)
	}
	t.Cleanup(func() { container.Terminate(ctx) })

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get host: %v", err)
	}

	port, err := container.MappedPort(ctx, "3306")
	if err != nil {
		t.Fatalf("failed to get port: %v", err)
	}

	db, err := connect.Open(connect.Config{
		Host:     host,
		Port:     port.Port(),
		User:     "root",
		Password: "test",
		Database: "mediawiki",
	})
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	_, err = db.Exec(`
		CREATE TABLE page (
			page_id INT AUTO_INCREMENT PRIMARY KEY,
			page_namespace INT NOT NULL,
			page_title VARCHAR(255) NOT NULL
		);
		CREATE TABLE job (
			job_id INT AUTO_INCREMENT PRIMARY KEY,
			job_cmd VARCHAR(255) NOT NULL
		);
	`)
	if err != nil {
		t.Fatalf("create tables failed: %v", err)
	}

	return db
}

func TestSearch_Audit(t *testing.T) {
	db := setupTestDB(t)

	// Seed 10 pages into MySQL
	for i := 1; i <= 10; i++ {
		_, _ = db.Exec("INSERT INTO page (page_namespace, page_title) VALUES (0, ?)", "Page_"+string(rune(i)))
	}

	// Mock Elasticsearch HTTP Server
	mockES := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]map[string]string{
			{"index": "mediawiki_content_1770474050", "docs.count": "10"},
			{"index": "mediawiki_general_1770474050", "docs.count": "15"},
		})
	}))
	defer mockES.Close()

	rpt, err := search.Audit(db, mockES.URL)
	if err != nil {
		t.Fatalf("Audit failed: %v", err)
	}

	if !rpt.ESConnected {
		t.Errorf("expected ES to be connected")
	}
	if rpt.ContentDocs != 10 {
		t.Errorf("expected 10 content docs, got %d", rpt.ContentDocs)
	}
	if rpt.HasDrift {
		t.Errorf("expected no drift when docs match pages, got drift: %s", rpt.String())
	}
}
