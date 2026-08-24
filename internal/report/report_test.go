package report

import (
	"strings"
	"testing"

	"pat-batch/internal/stats"
)

func TestFormatText(t *testing.T) {
	ps := []stats.ParamStat{{Parameter: "Temp", Mean: 120, Std: 0.5, CPK: 1.33}}
	bs := []stats.BatchResult{{Batch: "B1", OOS: []string{"Temp"}}}
	out, err := Format(ps, bs, "text")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Temp mean=120.000") {
		t.Fatalf("missing param stat line: %q", out)
	}
	if !strings.Contains(out, "OOS Temp") {
		t.Fatalf("missing OOS line: %q", out)
	}
}

func TestFormatJSON(t *testing.T) {
	ps := []stats.ParamStat{{Parameter: "Temp", Mean: 120}}
	bs := []stats.BatchResult{{Batch: "B1"}}
	out, err := Format(ps, bs, "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "\"params\"") || !strings.Contains(out, "\"batches\"") {
		t.Fatalf("missing json fields: %q", out)
	}
}

func TestFormatUnsupported(t *testing.T) {
	if _, err := Format(nil, nil, "yaml"); err == nil {
		t.Fatal("expected error for unsupported format")
	}
}
