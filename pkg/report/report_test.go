package report_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/FlowFeel/observatory-dbtools/pkg/drift"
	"github.com/FlowFeel/observatory-dbtools/pkg/report"
)

func TestWriteJSON(t *testing.T) {
	rep := drift.Report{
		FPTCount:    10,
		DICount:     7,
		MissingInDI: 3,
		Description: "3 entries missing",
	}

	var buf bytes.Buffer
	fmt := report.NewFormatter(&buf)

	if err := fmt.WriteJSON(rep); err != nil {
		t.Fatalf("WriteJSON failed: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, `"missing_in_di": 3`) {
		t.Errorf("expected missing_in_di in json output, got:\n%s", out)
	}
}
