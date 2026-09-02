package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	"github.com/FlowFeel/observatory-dbtools/pkg/asset"
	"github.com/FlowFeel/observatory-dbtools/pkg/baseline"
	"github.com/FlowFeel/observatory-dbtools/pkg/connect"
	"github.com/FlowFeel/observatory-dbtools/pkg/curate"
	"github.com/FlowFeel/observatory-dbtools/pkg/drift"
	"github.com/FlowFeel/observatory-dbtools/pkg/search"
	"github.com/FlowFeel/observatory-dbtools/pkg/snapshot"
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

func cmdAsset(args []string) {
	fs := flag.NewFlagSet("asset", flag.ExitOnError)
	auditFlag := fs.Bool("audit", false, "Audit database for hardcoded host URLs and missing assets")
	remediateFlag := fs.Bool("remediate-urls", false, "Rewrite hardcoded host URLs in wikitext to canonical filenames")
	dryRunFlag := fs.Bool("dry-run", false, "Simulate URL remediation without writing changes")
	fs.Parse(args)

	if !*auditFlag && !*remediateFlag {
		fmt.Fprintln(os.Stderr, "asset: --audit or --remediate-urls is required")
		os.Exit(1)
	}

	db, _, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "asset: connect error: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if *auditFlag {
		rpt, err := asset.AuditDatabase(db)
		if err != nil {
			fmt.Fprintf(os.Stderr, "asset: audit error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(rpt.String())
		if rpt.HasViolations() {
			os.Exit(2)
		}
	}

	if *remediateFlag {
		count, err := asset.RemediateDatabase(db, *dryRunFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "asset: remediation error: %v\n", err)
			os.Exit(1)
		}
		if *dryRunFlag {
			fmt.Printf("✓ [Dry Run] Would remediate %d text records with hardcoded URLs\n", count)
		} else {
			fmt.Printf("✓ Remediated %d text records with canonical filenames\n", count)
		}
	}
}

func cmdCurate(args []string) {
	fs := flag.NewFlagSet("curate", flag.ExitOnError)
	tier := fs.String("tier", "test", "Dataset tier to curate: test, stage, or prod")
	verifyFlag := fs.Bool("verify", false, "Verify target database satisfies tier plan")
	fs.Parse(args)

	plan, err := curate.NewPlan(*tier, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "curate: plan error: %v\n", err)
		os.Exit(1)
	}

	db, _, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "curate: connect error: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if *verifyFlag {
		if err := curate.VerifyDatabase(db, plan); err != nil {
			fmt.Fprintf(os.Stderr, "curate: tier verification failed: %v\n", err)
			os.Exit(2)
		}
		fmt.Printf("✓ Database satisfies tier %q curation contracts (%d portal pages verified)\n", *tier, len(plan.RequiredPages))
	}
}

func cmdDrift(args []string) {
	fs := flag.NewFlagSet("drift", flag.ExitOnError)
	checkFlag := fs.Bool("check", false, "Check for SMW property storage drift")
	fixFlag := fs.Bool("fix", false, "Fix detected SMW property storage drift")
	fs.Parse(args)

	if !*checkFlag && !*fixFlag {
		fmt.Fprintln(os.Stderr, "drift: --check or --fix is required")
		os.Exit(1)
	}

	db, _, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "drift: connect error: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if *checkFlag {
		rpt, err := drift.Check(db, drift.DefaultRegistry())
		if err != nil {
			fmt.Fprintf(os.Stderr, "drift: check error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(rpt.String())
		if rpt.HasDrift {
			os.Exit(2)
		}
	}

	if *fixFlag {
		rows, err := drift.Fix(db, drift.DefaultRegistry())
		if err != nil {
			fmt.Fprintf(os.Stderr, "drift: fix error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Fixed SMW drift: inserted %d missing rows into smw_di_time\n", rows)
	}
}

func cmdSearch(args []string) {
	fs := flag.NewFlagSet("search", flag.ExitOnError)
	checkFlag := fs.Bool("check", false, "Audit Elasticsearch search index health vs MySQL")
	esHost := fs.String("es-host", envOrDefault("ES_HOST", "http://localhost:9200"), "Elasticsearch endpoint host")
	fs.Parse(args)

	if !*checkFlag {
		fmt.Fprintln(os.Stderr, "search: --check is required")
		os.Exit(1)
	}

	db, _, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "search: connect error: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	rpt, err := search.Audit(db, *esHost)
	if err != nil {
		fmt.Fprintf(os.Stderr, "search: audit error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(rpt.String())
	if rpt.HasDrift || !rpt.ESConnected {
		os.Exit(2)
	}
}

func cmdSnapshot(args []string) {
	fs := flag.NewFlagSet("snapshot", flag.ExitOnError)
	output := fs.String("output", "", "Output path for schema dump file")
	fs.Parse(args)

	if *output == "" {
		fmt.Fprintln(os.Stderr, "snapshot: --output is required")
		os.Exit(1)
	}

	db, dbName, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "snapshot: connect error: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := snapshot.DumpExport(db, dbName, *output); err != nil {
		fmt.Fprintf(os.Stderr, "snapshot: export error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Database schema snapshot captured to %s\n", *output)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("observatory-dbtools — contract-driven database tooling")
		fmt.Println()
		fmt.Println("Commands:")
		fmt.Println("  baseline  Import schema + seed into target DB")
		fmt.Println("  asset     Audit and remediate hardcoded asset URLs in database")
		fmt.Println("  curate    Validate dataset curation requirements for tier (test, stage, prod)")
		fmt.Println("  drift     Check or fix SMW FPT↔DI drift")
		fmt.Println("  search    Audit search index document counts vs MySQL")
		fmt.Println("  snapshot  Dump current schema to file")
		fmt.Println("  migrate   Apply/rollback SQL migrations")
		fmt.Println("  compare   Diff two schema states")
		fmt.Println()
		fmt.Println("Use 'task --list' for the go-task interface.")
		os.Exit(0)
	}

	switch os.Args[1] {
	case "baseline":
		cmdBaseline(os.Args[2:])
	case "asset":
		cmdAsset(os.Args[2:])
	case "curate":
		cmdCurate(os.Args[2:])
	case "drift":
		cmdDrift(os.Args[2:])
	case "search":
		cmdSearch(os.Args[2:])
	case "snapshot":
		cmdSnapshot(os.Args[2:])
	case "migrate":
		fmt.Println("migrate: not yet implemented")
	case "compare":
		fmt.Println("compare: not yet implemented")
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
