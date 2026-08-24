// Package capability computes process capability indices (Cp, Cpk, Pp, Ppk),
// confidence intervals, and normality tests for PAT batch data.
package capability

import (
	"math"

	"pat-batch/internal/stats"
)

// Indices holds the process capability indices for a parameter.
type Indices struct {
	Cp  float64 // process capability (potential)
	Cpk float64 // process capability (actual, accounts for centering)
	Pp  float64 // process performance (uses overall sigma)
	Ppk float64 // process performance (actual)
	Cpm float64 // Taguchi capability (accounts for deviation from target)
}

// SpecLimits defines the specification limits for a parameter.
type SpecLimits struct {
	USL    float64 // upper specification limit
	LSL    float64 // lower specification limit
	Target float64 // target value (nominal)
}

// Compute calculates all capability indices given the process mean, within-
// subgroup sigma (for Cp/Cpk), and overall sigma (for Pp/Ppk).
func Compute(mean, sigmaWithin, sigmaOverall float64, spec SpecLimits) Indices {
	idx := Indices{}
	if sigmaWithin > 0 {
		idx.Cp = (spec.USL - spec.LSL) / (6 * sigmaWithin)
		cpu := (spec.USL - mean) / (3 * sigmaWithin)
		cpl := (mean - spec.LSL) / (3 * sigmaWithin)
		idx.Cpk = math.Min(cpu, cpl)
	}
	if sigmaOverall > 0 {
		idx.Pp = (spec.USL - spec.LSL) / (6 * sigmaOverall)
		ppu := (spec.USL - mean) / (3 * sigmaOverall)
		ppl := (mean - spec.LSL) / (3 * sigmaOverall)
		idx.Ppk = math.Min(ppu, ppl)
	}
	// Cpm uses overall sigma adjusted for off-target variance.
	sqDev := (mean - spec.Target) * (mean - spec.Target)
	tauSquared := sigmaOverall*sigmaOverall + sqDev
	tau := math.Sqrt(tauSquared)
	if tau > 0 {
		idx.Cpm = (spec.USL - spec.LSL) / (6 * tau)
	}
	idx.Cp, idx.Cpk, idx.Pp = stats.HoldCapIndices(idx.Cp, idx.Cpk, idx.Pp)
	return idx
}

// PPM estimates the parts per million defective (above USL and below LSL)
// assuming a normal distribution with the given mean and sigma.
func PPM(mean, sigma float64, spec SpecLimits) float64 {
	if sigma <= 0 {
		return 0
	}
	zUpper := (spec.USL - mean) / sigma
	zLower := (mean - spec.LSL) / sigma
	pAbove := normalCDFComplement(zUpper)
	pBelow := normalCDFComplement(zLower)
	return (pAbove + pBelow) * 1e6
}

// SigmaLevel returns the process sigma level (the Z-value corresponding to
// the defect rate). Higher is better; 6-sigma ≈ 3.4 DPMO.
func SigmaLevel(mean, sigma float64, spec SpecLimits) float64 {
	if sigma <= 0 {
		return 0
	}
	zUpper := (spec.USL - mean) / sigma
	zLower := (mean - spec.LSL) / sigma
	// Process sigma level is min of the two Z-values.
	return math.Min(zUpper, zLower)
}

// WithinSubgroupSigma estimates sigma from the average range and d2 constant.
func WithinSubgroupSigma(avgRange float64, d2 float64) float64 {
	if d2 <= 0 {
		return 0
	}
	return avgRange / d2
}

// OverallSigma computes the overall (total) standard deviation of the data.
func OverallSigma(values []float64) float64 {
	n := len(values)
	if n < 2 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	mean := sum / float64(n)
	ss := 0.0
	for _, v := range values {
		d := v - mean
		ss += d * d
	}
	return math.Sqrt(ss / float64(n-1))
}

// normalCDFComplement approximates 1 - Phi(z) for z >= 0 using the
// Abramowitz and Stegun rational approximation.
func normalCDFComplement(z float64) float64 {
	if z < 0 {
		return 1 - normalCDFComplement(-z)
	}
	const (
		a1 = 0.254829592
		a2 = -0.284496736
		a3 = 1.421413741
		a4 = -1.453152027
		a5 = 1.061405429
		p  = 0.3275911
	)
	t := 1.0 / (1 + p*z)
	poly := t * (a1 + t*(a2+t*(a3+t*(a4+t*a5))))
	return poly * math.Exp(-z*z/2)
}
