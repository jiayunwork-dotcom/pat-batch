// Package report renders PAT statistics and batch conformance.
package report

import (
	"encoding/json"
	"fmt"
	"strings"

	"pat-batch/internal/stats"
)

// Format renders the analysis in the requested format (text or json).
func Format(paramStats []stats.ParamStat, batches []stats.BatchResult, format string) (string, error) {
	switch format {
	case "json", "":
		type payload struct {
			Params  []stats.ParamStat   `json:"params"`
			Batches []stats.BatchResult `json:"batches"`
		}
		b, err := json.MarshalIndent(payload{paramStats, batches}, "", "  ")
		if err != nil {
			return "", fmt.Errorf("marshal report: %w", err)
		}
		return string(b) + "\n", nil
	case "text":
		return textFormat(paramStats, batches), nil
	default:
		return "", fmt.Errorf("unsupported format %q (want text or json)", format)
	}
}

func textFormat(ps []stats.ParamStat, bs []stats.BatchResult) string {
	var b strings.Builder
	b.WriteString("parameter statistics:\n")
	for _, p := range ps {
		fmt.Fprintf(&b, "  %s mean=%.3f std=%.3f min=%.3f max=%.3f cpk=%.2f\n",
			p.Parameter, p.Mean, p.Std, p.Min, p.Max, p.CPK)
	}
	b.WriteString("batch conformance:\n")
	for _, r := range bs {
		if len(r.OOS) == 0 && len(r.OOT) == 0 {
			fmt.Fprintf(&b, "  %s: in control\n", r.Batch)
			continue
		}
		fmt.Fprintf(&b, "  %s:\n", r.Batch)
		for _, p := range r.OOS {
			fmt.Fprintf(&b, "    OOS %s\n", p)
		}
		for _, p := range r.OOT {
			fmt.Fprintf(&b, "    OOT %s\n", p)
		}
	}
	return b.String()
}
