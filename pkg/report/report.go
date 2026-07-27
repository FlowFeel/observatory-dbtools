// Package report provides serialization and formatting of DB tools metrics, diffs, and drift checks.
package report

import (
	"encoding/json"
	"fmt"
	"io"
)

// Formatter handles boundary serialization of domain reports.
type Formatter struct {
	w io.Writer
}

// NewFormatter creates a report Formatter writing to w.
func NewFormatter(w io.Writer) *Formatter {
	return &Formatter{w: w}
}

// WriteJSON serializes any domain report to indented JSON output.
func (f *Formatter) WriteJSON(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("report: marshal json: %w", err)
	}
	_, err = f.w.Write(append(data, '\n'))
	return err
}
