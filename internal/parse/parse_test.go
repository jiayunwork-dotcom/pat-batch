package parse

import (
	"os"
	"path/filepath"
	"testing"
)

const measCSV = `batch,parameter,value
B1,Temperature,120.0
B1,Temperature,120.5
B2,Pressure,2.1
`

const specsCSV = `parameter,target,low,high
Temperature,120,118,122
Pressure,2.0,1.8,2.2
`

func TestReadMeasurementsValid(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "m.csv")
	if err := os.WriteFile(p, []byte(measCSV), 0o644); err != nil {
		t.Fatal(err)
	}
	ms, err := ReadMeasurements(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ms) != 3 {
		t.Fatalf("expected 3 measurements, got %d", len(ms))
	}
	if ms[0].Batch != "B1" || ms[0].Parameter != "Temperature" || ms[0].Value != 120.0 {
		t.Fatalf("unexpected first measurement: %+v", ms[0])
	}
}

func TestReadMeasurementsMissingFile(t *testing.T) {
	if _, err := ReadMeasurements(filepath.Join(t.TempDir(), "nope.csv")); err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestReadMeasurementsBadValue(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "bad.csv")
	if err := os.WriteFile(p, []byte("batch,parameter,value\nB1,Temp,abc\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadMeasurements(p); err == nil {
		t.Fatal("expected error for non-numeric value")
	}
}

func TestReadSpecsValid(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "s.csv")
	if err := os.WriteFile(p, []byte(specsCSV), 0o644); err != nil {
		t.Fatal(err)
	}
	specs, err := ReadSpecs(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(specs) != 2 {
		t.Fatalf("expected 2 specs, got %d", len(specs))
	}
	sp := specs["Temperature"]
	if sp.Target != 120 || sp.Low != 118 || sp.High != 122 {
		t.Fatalf("unexpected spec: %+v", sp)
	}
}
