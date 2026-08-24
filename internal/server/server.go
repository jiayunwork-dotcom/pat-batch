// Package server exposes PAT batch analysis via HTTP/JSON endpoints.
package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"pat-batch/internal/stats"
)

// Config holds server configuration.
type Config struct {
	Addr string
}

// New creates a configured http.ServeMux with all routes registered.
func New(cfg Config) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/api/analyze", handleAnalyze)
	mux.HandleFunc("/api/batch", handleBatch)
	return mux
}

// ListenAndServe starts the HTTP server.
func ListenAndServe(cfg Config) error {
	mux := New(cfg)
	return http.ListenAndServe(cfg.Addr, mux)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

type specInput struct {
	Parameter string  `json:"parameter"`
	Target    float64 `json:"target"`
	Low       float64 `json:"low"`
	High      float64 `json:"high"`
}

type measurementInput struct {
	Batch     string  `json:"batch"`
	Parameter string  `json:"parameter"`
	Value     float64 `json:"value"`
}

type analyzeRequest struct {
	Measurements []measurementInput `json:"measurements"`
	Specs        []specInput        `json:"specs"`
}

type paramStatOutput struct {
	Parameter string  `json:"parameter"`
	Mean      float64 `json:"mean"`
	Std       float64 `json:"std"`
	Min       float64 `json:"min"`
	Max       float64 `json:"max"`
	CPK       float64 `json:"cpk"`
}

func handleAnalyze(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var req analyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if len(req.Measurements) == 0 {
		httpError(w, http.StatusBadRequest, "measurements array is empty")
		return
	}
	ms := toMeasurements(req.Measurements)
	specs := toSpecs(req.Specs)
	paramStats := stats.ComputeParamStats(ms, specs)

	out := make([]paramStatOutput, len(paramStats))
	for i, ps := range paramStats {
		out[i] = paramStatOutput{
			Parameter: ps.Parameter,
			Mean:      ps.Mean,
			Std:       ps.Std,
			Min:       ps.Min,
			Max:       ps.Max,
			CPK:       ps.CPK,
		}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"params": out})
}

type batchRequest struct {
	Measurements []measurementInput `json:"measurements"`
	Specs        []specInput        `json:"specs"`
}

type batchResultOutput struct {
	Batch  string   `json:"batch"`
	OOS    []string `json:"oos,omitempty"`
	OOT    []string `json:"oot,omitempty"`
}

func handleBatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var req batchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if len(req.Measurements) == 0 {
		httpError(w, http.StatusBadRequest, "measurements array is empty")
		return
	}
	ms := toMeasurements(req.Measurements)
	specs := toSpecs(req.Specs)

	// Group by batch
	byBatch := map[string][]stats.Measurement{}
	var order []string
	for _, m := range ms {
		if _, ok := byBatch[m.Batch]; !ok {
			order = append(order, m.Batch)
		}
		byBatch[m.Batch] = append(byBatch[m.Batch], m)
	}

	results := make([]batchResultOutput, 0, len(order))
	for _, b := range order {
		br := stats.EvaluateBatch(byBatch[b], specs)
		results = append(results, batchResultOutput{
			Batch: b,
			OOS:   br.OOS,
			OOT:   br.OOT,
		})
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"batches": results})
}

func toMeasurements(inputs []measurementInput) []stats.Measurement {
	out := make([]stats.Measurement, len(inputs))
	for i, mi := range inputs {
		out[i] = stats.Measurement{Batch: mi.Batch, Parameter: mi.Parameter, Value: mi.Value}
	}
	return out
}

func toSpecs(inputs []specInput) map[string]stats.Spec {
	m := make(map[string]stats.Spec, len(inputs))
	for _, si := range inputs {
		m[si.Parameter] = stats.Spec{Parameter: si.Parameter, Target: si.Target, Low: si.Low, High: si.High}
	}
	return m
}

func httpError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func writeJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}

// ParsePort extracts port number from addr string.
func ParsePort(addr string) int {
	parts := strings.Split(addr, ":")
	if len(parts) < 2 {
		return 0
	}
	p, _ := strconv.Atoi(parts[len(parts)-1])
	return p
}

// FormatAddr produces a display-friendly address.
func FormatAddr(addr string) string {
	port := ParsePort(addr)
	if port == 0 {
		return addr
	}
	return fmt.Sprintf("http://0.0.0.0:%d", port)
}
