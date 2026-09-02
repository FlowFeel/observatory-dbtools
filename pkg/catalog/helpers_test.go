package catalog_test

import "os"

// readFileBytes reads a file relative to the repo root helper.
// Kept here to avoid colliding with the catalog package internals.
func readFileBytes(path string) ([]byte, error) {
	return os.ReadFile(path)
}
