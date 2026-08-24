package trend

import (
	"math"
	"testing"
)

func TestNelsonRule1(t *testing.T) {
	// Data with one point beyond 3 sigma.
	vals := []float64{10, 10, 10, 10, 10, 10, 25, 10, 10, 10}
	cl := 10.0
	sigma := 2.0
	violations := NelsonCheck(vals, cl, sigma)
	foundRule1 := false
	for _, v := range violations {
		if v.Rule == 1 && v.Start == 6 {
			foundRule1 = true
		}
	}
	if !foundRule1 {
		t.Fatal("expected Nelson rule 1 violation at index 6")
	}
}

func TestNelsonRule2(t *testing.T) {
	// Nine points above CL.
	vals := make([]float64, 15)
	for i := range vals {
		vals[i] = 11 // all above CL=10
	}
	violations := NelsonCheck(vals, 10, 2)
	foundRule2 := false
	for _, v := range violations {
		if v.Rule == 2 {
			foundRule2 = true
		}
	}
	if !foundRule2 {
		t.Fatal("expected Nelson rule 2 violation")
	}
}

func TestNelsonRule3(t *testing.T) {
	// Six consecutive increasing.
	vals := []float64{1, 2, 3, 4, 5, 6, 7, 3}
	violations := NelsonCheck(vals, 4, 3)
	foundRule3 := false
	for _, v := range violations {
		if v.Rule == 3 {
			foundRule3 = true
		}
	}
	if !foundRule3 {
		t.Fatal("expected Nelson rule 3 violation")
	}
}

func TestCUSUMDetectsShift(t *testing.T) {
	// 20 in-control points then 10 shifted points.
	vals := make([]float64, 30)
	for i := 0; i < 20; i++ {
		vals[i] = 10
	}
	for i := 20; i < 30; i++ {
		vals[i] = 12 // shift of 1 sigma
	}
	params := CUSUMParams{Target: 10, Sigma: 2, K: 0.5, H: 4}
	points := CUSUM(vals, params)
	signals := CUSUMSignals(points)
	if len(signals) == 0 {
		t.Fatal("CUSUM should detect the shift")
	}
	// First signal should be after the shift begins.
	if signals[0] < 20 {
		t.Fatalf("first signal at %d, expected >= 20", signals[0])
	}
}

func TestCUSUMNoFalseAlarm(t *testing.T) {
	// All in-control.
	vals := make([]float64, 50)
	for i := range vals {
		vals[i] = 10
	}
	params := CUSUMParams{Target: 10, Sigma: 2, K: 0.5, H: 4}
	points := CUSUM(vals, params)
	signals := CUSUMSignals(points)
	if len(signals) > 0 {
		t.Fatalf("no false alarms expected, got signals at %v", signals)
	}
}

func TestEWMADetectsShift(t *testing.T) {
	vals := make([]float64, 40)
	for i := 0; i < 20; i++ {
		vals[i] = 50
	}
	for i := 20; i < 40; i++ {
		vals[i] = 53 // shift of 1.5 sigma (sigma=2)
	}
	params := EWMAParams{Target: 50, Sigma: 2, Lambda: 0.2, L: 2.7}
	points := EWMA(vals, params)
	signals := EWMASignals(points)
	if len(signals) == 0 {
		t.Fatal("EWMA should detect the shift")
	}
}

func TestEWMASteadyStateLimits(t *testing.T) {
	params := EWMAParams{Target: 100, Sigma: 5, Lambda: 0.2, L: 3}
	ucl, lcl := EWMASteadyStateLimits(params)
	if ucl <= 100 || lcl >= 100 {
		t.Fatalf("ucl=%v lcl=%v invalid", ucl, lcl)
	}
	// Symmetric around target.
	if math.Abs((ucl-100)-(100-lcl)) > 0.001 {
		t.Fatalf("limits not symmetric: ucl=%v lcl=%v", ucl, lcl)
	}
}

func TestRunsTest(t *testing.T) {
	// Alternating pattern: many runs.
	vals := []float64{1, 3, 1, 3, 1, 3, 1, 3, 1, 3}
	result := RunsTest(vals, 2)
	if result.NumRuns != 10 {
		t.Fatalf("runs = %d, want 10", result.NumRuns)
	}
	// Too many runs relative to random: Z should be positive.
	if result.ZScore <= 0 {
		t.Fatalf("Z = %v, expected positive (too many runs)", result.ZScore)
	}
}

func TestLongestRun(t *testing.T) {
	vals := []float64{5, 5, 5, 1, 1, 5, 5, 5, 5, 1}
	longest := LongestRun(vals, 3)
	if longest != 4 {
		t.Fatalf("longest run = %d, want 4", longest)
	}
}

func TestConsecutiveIncreasing(t *testing.T) {
	vals := []float64{1, 2, 3, 4, 2, 3, 4, 5, 6, 1}
	inc := ConsecutiveIncreasing(vals)
	if inc != 5 {
		t.Fatalf("consecutive increasing = %d, want 5", inc)
	}
}
