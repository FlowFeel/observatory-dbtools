package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	"github.com/FlowFeel/observatory-dbtools/pkg/baseline"
	"github.com/FlowFeel/observatory-dbtools/pkg/connect"
)

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func openDB() (*sql.DB, string, error) {
	cfg := connect.Config{
		Host:     envOrDefault("DB_HOST", "localhost"),
		Port:     envOrDefault("DB_PORT", "3306"),
		User:     envOrDefault("DB_USER", "root"),
		Password: envOrDefault("DB_PASS", ""),
		Database: envOrDefault("DB_NAME", "mediawiki"),
	}
	db, err := connect.Open(cfg)
	return db, cfg.Database, err
}

func cmdBaseline(args []string) {
	fs := flag.NewFlagSet("baseline", flag.ExitOnError)
	schema := fs.String("schema", "", "Path to schema SQL file")
	seed := fs.String("seed", "", "Path to seed data SQL file")
	fs.Parse(args)

	if *schema == "" {
		fmt.Fprintln(os.Stderr, "baseline: --schema is required")
		os.Exit(1)
	}

	db, dbName, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "baseline: connect: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	fmt.Printf("Importing schema from %s...\n", *schema)
	if err := baseline.Import(db, *schema); err != nil {
		fmt.Fprintf(os.Stderr, "baseline: schema import: %v\n", err)
		os.Exit(1)
	}

	count, _ := baseline.TableCount(db, dbName)
	fmt.Printf("✓ Schema imported (%d tables)\n", count)

	if *seed != "" {
		fmt.Printf("Importing seed data from %s...\n", *seed)
		if err := baseline.Import(db, *seed); err != nil {
			fmt.Fprintf(os.Stderr, "baseline: seed import: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✓ Seed data imported")
	}

	// Verify
	err = baseline.Verify(db, dbName, 50, []string{
		"user", "page", "revision", "text", "site_stats",
		"smw_object_ids", "smw_fpt_mdat", "smw_di_time",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "baseline: verify: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Baseline verified")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("observatory-dbtools — contract-driven database tooling")
		fmt.Println()
		fmt.Println("Commands:")
		fmt.Println("  baseline  Import schema + seed into target DB")
		fmt.Println("  migrate   Apply/rollback SQL migrations")
		fmt.Println("  drift     Check or fix SMW FPT↔DI drift")
		fmt.Println("  snapshot  Dump current schema to file")
		fmt.Println("  compare   Diff two schema states")
		fmt.Println()
		fmt.Println("Use 'task --list' for the go-task interface.")
		os.Exit(0)
	}

	switch os.Args[1] {
	case "baseline":
		cmdBaseline(os.Args[2:])
	case "migrate":
		fmt.Println("migrate: not yet implemented")
	case "drift":
		fmt.Println("drift: not yet implemented")
	case "snapshot":
		fmt.Println("snapshot: not yet implemented")
	case "compare":
		fmt.Println("compare: not yet implemented")
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
