package baseline_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/FlowFeel/observatory-dbtools/pkg/baseline"
	"github.com/FlowFeel/observatory-dbtools/pkg/connect"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

func TestImportSchema(t *testing.T) {
	ctx := context.Background()

	container, err := mysql.Run(ctx, "mysql:8.4",
		mysql.WithDatabase("mediawiki"),
		mysql.WithUsername("root"),
		mysql.WithPassword("test"),
	)
	if err != nil {
		t.Fatalf("start container: %v", err)
	}
	t.Cleanup(func() { container.Terminate(ctx) })

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "3306")

	db, err := connect.Open(connect.Config{
		Host:     host,
		Port:     port.Port(),
		User:     "root",
		Password: "test",
		Database: "mediawiki",
	})
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer db.Close()

	// Import the production schema baseline
	schemaPath := filepath.Join("..", "..", "baselines", "schema.sql")
	if _, err := os.Stat(schemaPath); os.IsNotExist(err) {
		t.Skip("baselines/schema.sql not found — run 'task snapshot' first")
	}

	if err := baseline.Import(db, schemaPath); err != nil {
		t.Fatalf("Import schema: %v", err)
	}

	// Contract: schema creates the expected tables
	err = baseline.Verify(db, "mediawiki", 100, []string{
		// MW core tables
		"user",
		"page",
		"revision",
		"text",
		"site_stats",
		"interwiki",
		// SMW tables
		"smw_object_ids",
		"smw_fpt_mdat",
		"smw_di_time",
		"smw_di_blob",
	})
	if err != nil {
		t.Errorf("Verify: %v", err)
	}

	count, err := baseline.TableCount(db, "mediawiki")
	if err != nil {
		t.Fatalf("TableCount: %v", err)
	}
	t.Logf("Imported %d tables from baseline", count)
}

func TestImportSchemaAndSeed(t *testing.T) {
	ctx := context.Background()

	container, err := mysql.Run(ctx, "mysql:8.4",
		mysql.WithDatabase("mediawiki"),
		mysql.WithUsername("root"),
		mysql.WithPassword("test"),
	)
	if err != nil {
		t.Fatalf("start container: %v", err)
	}
	t.Cleanup(func() { container.Terminate(ctx) })

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "3306")

	db, err := connect.Open(connect.Config{
		Host:     host,
		Port:     port.Port(),
		User:     "root",
		Password: "test",
		Database: "mediawiki",
	})
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer db.Close()

	schemaPath := filepath.Join("..", "..", "baselines", "schema.sql")
	seedPath := filepath.Join("..", "..", "baselines", "seed.sql")

	if _, err := os.Stat(schemaPath); os.IsNotExist(err) {
		t.Skip("baselines/schema.sql not found")
	}
	if _, err := os.Stat(seedPath); os.IsNotExist(err) {
		t.Skip("baselines/seed.sql not found")
	}

	// Import schema
	if err := baseline.Import(db, schemaPath); err != nil {
		t.Fatalf("Import schema: %v", err)
	}

	// Import seed data
	if err := baseline.Import(db, seedPath); err != nil {
		t.Fatalf("Import seed: %v", err)
	}

	// Contract: seed data creates admin user
	var userCount int
	err = db.QueryRow("SELECT COUNT(*) FROM user").Scan(&userCount)
	if err != nil {
		t.Fatalf("count users: %v", err)
	}
	if userCount == 0 {
		t.Error("expected at least 1 user after seed import")
	}
	t.Logf("Seed imported %d users", userCount)
}
