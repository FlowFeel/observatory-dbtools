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
		HasDrift: true,
		Targets: []drift.TargetReport{
			{
				Target: drift.DriftTarget{
					PropertyName: "_MDAT",
					FptTable:     "smw_fpt_mdat",
					DiTable:      "smw_di_time",
					Pid:          29,
				},
				FPTCount:    10,
				DICount:     7,
				MissingInDI: 3,
			},
		},
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
	if !strings.Contains(out, `"_MDAT"`) {
		t.Errorf("expected _MDAT target in json output, got:\n%s", out)
	}
	if !strings.Contains(out, `"has_drift": true`) {
		t.Errorf("expected has_drift true in json output, got:\n%s", out)
	}
}
