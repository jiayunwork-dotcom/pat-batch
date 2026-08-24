package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestAnalyzeEndpoint(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	payload := analyzeRequest{
		Measurements: []measurementInput{
			{Batch: "B1", Parameter: "temp", Value: 100},
			{Batch: "B1", Parameter: "temp", Value: 102},
			{Batch: "B1", Parameter: "temp", Value: 99},
			{Batch: "B2", Parameter: "temp", Value: 101},
		},
		Specs: []specInput{
			{Parameter: "temp", Target: 100, Low: 95, High: 105},
		},
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/analyze", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var resp struct {
		Params []paramStatOutput `json:"params"`
	}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if len(resp.Params) == 0 {
		t.Error("expected params")
	}
}

func TestAnalyzeEndpoint_Empty(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	body := []byte(`{"measurements":[],"specs":[]}`)
	req := httptest.NewRequest(http.MethodPost, "/api/analyze", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestBatchEndpoint(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	payload := batchRequest{
		Measurements: []measurementInput{
			{Batch: "B1", Parameter: "pH", Value: 7.0},
			{Batch: "B1", Parameter: "pH", Value: 7.2},
			{Batch: "B2", Parameter: "pH", Value: 3.0},
		},
		Specs: []specInput{
			{Parameter: "pH", Target: 7.0, Low: 6.5, High: 7.5},
		},
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/batch", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var resp struct {
		Batches []batchResultOutput `json:"batches"`
	}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if len(resp.Batches) != 2 {
		t.Fatalf("expected 2 batches, got %d", len(resp.Batches))
	}
}

func TestMethodNotAllowed(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	endpoints := []string{"/api/analyze", "/api/batch"}
	for _, ep := range endpoints {
		req := httptest.NewRequest(http.MethodGet, ep, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("%s: expected 405, got %d", ep, rec.Code)
		}
	}
}
