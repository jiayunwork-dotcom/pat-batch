// Package parse reads batch measurement and specification CSV files.
package parse

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"pat-batch/internal/stats"
)

// ReadMeasurements reads a CSV with header: batch,parameter,value.
func ReadMeasurements(path string) ([]stats.Measurement, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open measurements: %w", err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read measurements csv: %w", err)
	}
	if len(records) < 2 {
		return nil, fmt.Errorf("measurements file has no data rows")
	}
	header := records[0]
	bi, pi, vi := -1, -1, -1
	for i, h := range header {
		switch strings.ToLower(strings.TrimSpace(h)) {
		case "batch":
			bi = i
		case "parameter":
			pi = i
		case "value":
			vi = i
		}
	}
	if bi < 0 || pi < 0 || vi < 0 {
		return nil, fmt.Errorf("measurements header must contain batch,parameter,value")
	}

	var out []stats.Measurement
	for n, row := range records[1:] {
		if len(row) <= vi {
			return nil, fmt.Errorf("row %d: missing value column", n+2)
		}
		v, err := strconv.ParseFloat(strings.TrimSpace(row[vi]), 64)
		if err != nil {
			return nil, fmt.Errorf("row %d: invalid value %q: %w", n+2, row[vi], err)
		}
		out = append(out, stats.Measurement{
			Batch:     strings.TrimSpace(row[bi]),
			Parameter: strings.TrimSpace(row[pi]),
			Value:     v,
		})
	}
	return out, nil
}

// ReadSpecs reads a CSV with header: parameter,target,low,high.
func ReadSpecs(path string) (map[string]stats.Spec, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open specs: %w", err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read specs csv: %w", err)
	}
	if len(records) < 2 {
		return nil, fmt.Errorf("specs file has no data rows")
	}
	header := records[0]
	pi, ti, li, hi := -1, -1, -1, -1
	for i, h := range header {
		switch strings.ToLower(strings.TrimSpace(h)) {
		case "parameter":
			pi = i
		case "target":
			ti = i
		case "low":
			li = i
		case "high":
			hi = i
		}
	}
	if pi < 0 || ti < 0 || li < 0 || hi < 0 {
		return nil, fmt.Errorf("specs header must contain parameter,target,low,high")
	}

	specs := map[string]stats.Spec{}
	for n, row := range records[1:] {
		if len(row) <= hi {
			return nil, fmt.Errorf("row %d: missing spec column", n+2)
		}
		target, err := strconv.ParseFloat(strings.TrimSpace(row[ti]), 64)
		if err != nil {
			return nil, fmt.Errorf("row %d: invalid target %q: %w", n+2, row[ti], err)
		}
		low, err := strconv.ParseFloat(strings.TrimSpace(row[li]), 64)
		if err != nil {
			return nil, fmt.Errorf("row %d: invalid low %q: %w", n+2, row[li], err)
		}
		high, err := strconv.ParseFloat(strings.TrimSpace(row[hi]), 64)
		if err != nil {
			return nil, fmt.Errorf("row %d: invalid high %q: %w", n+2, row[hi], err)
		}
		specs[strings.TrimSpace(row[pi])] = stats.Spec{
			Parameter: strings.TrimSpace(row[pi]),
			Target:    target,
			Low:       low,
			High:      high,
		}
	}
	return specs, nil
}
