package connect_test

import (
	"context"
	"testing"

	"github.com/FlowFeel/observatory-dbtools/pkg/connect"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

func TestConnect(t *testing.T) {
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

	cfg := connect.Config{
		Host:     host,
		Port:     port.Port(),
		User:     "root",
		Password: "test",
		Database: "mediawiki",
	}

	db, err := connect.Open(cfg)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer db.Close()

	version, err := connect.Health(db)
	if err != nil {
		t.Fatalf("Health failed: %v", err)
	}

	if version == "" {
		t.Error("expected non-empty MySQL version")
	}

	t.Logf("MySQL version: %s", version)
}
