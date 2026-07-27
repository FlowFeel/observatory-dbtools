package migrate_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/FlowFeel/observatory-dbtools/pkg/connect"
	"github.com/FlowFeel/observatory-dbtools/pkg/migrate"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

func TestMigrate_UpAndDown(t *testing.T) {
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
	t.Cleanup(func() { db.Close() })

	// Pre-create SMW tables required by migration 000001
	_, err = db.Exec(`
		CREATE TABLE smw_fpt_mdat (
			s_id INT NOT NULL,
			o_serialized VARCHAR(255) NOT NULL,
			o_sortkey DOUBLE NOT NULL
		) ENGINE=InnoDB;

		CREATE TABLE smw_di_time (
			s_id INT NOT NULL,
			p_id INT NOT NULL,
			o_serialized VARCHAR(255) NOT NULL,
			o_sortkey DOUBLE NOT NULL
		) ENGINE=InnoDB;
	`)
	if err != nil {
		t.Fatalf("create tables: %v", err)
	}

	migrationsDir, err := filepath.Abs(filepath.Join("..", "..", "migrations"))
	if err != nil {
		t.Fatalf("abs path: %v", err)
	}

	runner := migrate.NewRunner(db, migrationsDir)

	if err := runner.Up(); err != nil {
		t.Fatalf("Runner.Up failed: %v", err)
	}

	// Verify schema migration tracking table exists using a fresh connection
	dbVerify, err := connect.Open(connect.Config{
		Host:     host,
		Port:     port.Port(),
		User:     "root",
		Password: "test",
		Database: "mediawiki",
	})
	if err != nil {
		t.Fatalf("connect verify: %v", err)
	}
	defer dbVerify.Close()

	var version int
	var dirty bool
	err = dbVerify.QueryRow("SELECT version, dirty FROM schema_migrations").Scan(&version, &dirty)
	if err != nil {
		t.Fatalf("query schema_migrations: %v", err)
	}
	if version != 1 || dirty {
		t.Errorf("expected version 1 (clean), got version %d (dirty: %t)", version, dirty)
	}

	runnerDown := migrate.NewRunner(dbVerify, migrationsDir)
	if err := runnerDown.Down(); err != nil {
		t.Fatalf("Runner.Down failed: %v", err)
	}
}
