package stats

import "testing"

func TestMean(t *testing.T) {
	if got := Mean(nil); got != 0 {
		t.Fatalf("mean of empty should be 0, got %v", got)
	}
	if got := Mean([]float64{1, 2, 3}); got != 2 {
		t.Fatalf("mean of [1,2,3] should be 2, got %v", got)
	}
}

func TestStd(t *testing.T) {
	if got := Std(nil); got != 0 {
		t.Fatalf("std of empty should be 0, got %v", got)
	}
	// population std of [1,2,3] = sqrt((1+0+1)/3) = sqrt(2/3)
	if got := Std([]float64{1, 2, 3}); got < 0.816 || got > 0.817 {
		t.Fatalf("std of [1,2,3] should be ~0.816, got %v", got)
	}
}

func TestCPK(t *testing.T) {
	// mean=120, std=0.5, band [118,122] -> (122-120)/1.5 = 1.333
	got := CPK(120, 0.5, 118, 122)
	if got < 1.33 || got > 1.34 {
		t.Fatalf("cpk should be ~1.333, got %v", got)
	}
	if got := CPK(120, 0, 118, 122); got != 999 {
		t.Fatalf("cpk with zero std should be 999, got %v", got)
	}
}

func TestEvaluateBatchOOS(t *testing.T) {
	specs := map[string]Spec{
		"Temp": {Parameter: "Temp", Target: 120, Low: 118, High: 122},
	}
	ms := []Measurement{
		{Batch: "B1", Parameter: "Temp", Value: 117.0}, // below low -> OOS
		{Batch: "B1", Parameter: "Temp", Value: 117.5},
	}
	r := EvaluateBatch(ms, specs)
	if len(r.OOS) != 1 || r.OOS[0] != "Temp" {
		t.Fatalf("expected OOS Temp, got %v", r.OOS)
	}
}

func TestEvaluateBatchOOT(t *testing.T) {
	specs := map[string]Spec{
		"Pressure": {Parameter: "Pressure", Target: 2.0, Low: 1.8, High: 2.2},
	}
	ms := []Measurement{
		{Batch: "B3", Parameter: "Pressure", Value: 2.15}, // within spec, mean 2.15 -> OOT
		{Batch: "B3", Parameter: "Pressure", Value: 2.16},
		{Batch: "B3", Parameter: "Pressure", Value: 2.14},
	}
	r := EvaluateBatch(ms, specs)
	if len(r.OOS) != 0 {
		t.Fatalf("expected no OOS, got %v", r.OOS)
	}
	if len(r.OOT) != 1 || r.OOT[0] != "Pressure" {
		t.Fatalf("expected OOT Pressure, got %v", r.OOT)
	}
}
