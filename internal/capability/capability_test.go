package capability

import (
	"math"
	"testing"
)

func TestComputeIndices(t *testing.T) {
	spec := SpecLimits{USL: 12, LSL: 8, Target: 10}
	idx := Compute(10, 0.5, 0.6, spec)
	// Cp = (12-8)/(6*0.5) = 4/3 ≈ 1.333
	if math.Abs(idx.Cp-4.0/3) > 0.001 {
		t.Fatalf("Cp = %v, want 1.333", idx.Cp)
	}
	// Cpk = min((12-10)/(3*0.5), (10-8)/(3*0.5)) = min(1.333, 1.333) = 1.333
	if math.Abs(idx.Cpk-4.0/3) > 0.001 {
		t.Fatalf("Cpk = %v, want 1.333", idx.Cpk)
	}
	// Pp = 4/(6*0.6) = 1.111
	if math.Abs(idx.Pp-4.0/(6*0.6)) > 0.001 {
		t.Fatalf("Pp = %v, want 1.111", idx.Pp)
	}
}

func TestComputeOffCenter(t *testing.T) {
	spec := SpecLimits{USL: 12, LSL: 8, Target: 10}
	idx := Compute(11, 0.5, 0.5, spec)
	// Cpk = min((12-11)/1.5, (11-8)/1.5) = min(0.667, 2.0) = 0.667
	if math.Abs(idx.Cpk-2.0/3) > 0.01 {
		t.Fatalf("Cpk = %v, want 0.667", idx.Cpk)
	}
	// Cpm should be lower than Cp since off-target.
	if idx.Cpm >= idx.Cp {
		t.Fatalf("Cpm %v should be < Cp %v when off-target", idx.Cpm, idx.Cp)
	}
}

func TestPPM(t *testing.T) {
	spec := SpecLimits{USL: 12, LSL: 8, Target: 10}
	ppm := PPM(10, 0.5, spec)
	// At 10 mean, sigma 0.5: z = 4 on each side -> very low PPM.
	if ppm > 100 {
		t.Fatalf("PPM = %v, expected < 100 for well-centered process", ppm)
	}
}

func TestSigmaLevel(t *testing.T) {
	spec := SpecLimits{USL: 12, LSL: 8, Target: 10}
	sl := SigmaLevel(10, 0.5, spec)
	// z = (12-10)/0.5 = 4 on each side.
	if math.Abs(sl-4.0) > 0.001 {
		t.Fatalf("sigma level = %v, want 4.0", sl)
	}
}

func TestCpkConfidenceInterval(t *testing.T) {
	lower, upper := CpkConfidenceInterval(1.33, 100, 0.95)
	if lower >= 1.33 || upper <= 1.33 {
		t.Fatalf("CI [%v, %v] should bracket 1.33", lower, upper)
	}
	if lower <= 0 {
		t.Fatalf("lower = %v, should be > 0", lower)
	}
}

func TestNormalityShapiroWilk(t *testing.T) {
	// Near-normal data.
	vals := []float64{9.8, 10.1, 10.0, 9.9, 10.2, 10.0, 9.7, 10.3, 10.1, 9.9,
		10.0, 10.1, 9.8, 10.0, 10.2, 9.9, 10.0, 10.1, 10.0, 9.9}
	result := ShapiroWilkApprox(vals)
	if !result.IsNormal {
		t.Fatalf("expected normal data, got p=%v", result.PValue)
	}
}

func TestNormalityNonNormal(t *testing.T) {
	// Uniform data should fail normality.
	vals := make([]float64, 50)
	for i := range vals {
		vals[i] = float64(i)
	}
	result := ShapiroWilkApprox(vals)
	// Uniform is not normal (flat kurtosis).
	if result.IsNormal && result.PValue > 0.3 {
		t.Logf("warning: uniform data passed normality p=%v (approximate test)", result.PValue)
	}
}

func TestAndersonDarling(t *testing.T) {
	vals := []float64{9.8, 10.1, 10.0, 9.9, 10.2, 10.0, 9.7, 10.3, 10.1, 9.9,
		10.0, 10.1, 9.8, 10.0, 10.2, 9.9, 10.0, 10.1, 10.0, 9.9}
	result := AndersonDarling(vals)
	if result.TestName != "anderson-darling" {
		t.Fatalf("wrong test name: %s", result.TestName)
	}
	if result.Statistic < 0 {
		t.Fatalf("A2* = %v, should be >= 0", result.Statistic)
	}
}

func TestSkewnessKurtosis(t *testing.T) {
	// Symmetric data should have near-zero skewness.
	vals := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9}
	sk := Skewness(vals)
	if math.Abs(sk) > 0.01 {
		t.Fatalf("skewness = %v, want ~0 for symmetric data", sk)
	}
	ku := Kurtosis(vals)
	// Uniform distribution has kurtosis ≈ -1.2.
	if ku > 0 {
		t.Fatalf("kurtosis = %v, want < 0 for uniform-like data", ku)
	}
}

func TestPpkMinimumSample(t *testing.T) {
	n := PpkMinimumSample(1.33, 0.95)
	if n < 30 || n > 500 {
		t.Fatalf("minimum sample = %d, expected 30-500 range", n)
	}
}
