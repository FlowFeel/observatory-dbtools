package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("observatory-dbtools — contract-driven database tooling")
		fmt.Println()
		fmt.Println("Commands:")
		fmt.Println("  migrate   Apply/rollback SQL migrations")
		fmt.Println("  drift     Check or fix SMW FPT↔DI drift")
		fmt.Println("  snapshot  Dump current schema to file")
		fmt.Println("  compare   Diff two schema states")
		fmt.Println()
		fmt.Println("Use 'task --list' for the go-task interface.")
		os.Exit(0)
	}

	switch os.Args[1] {
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
