package capability

import "math"

// YieldResult summarizes process yield calculations.
type YieldResult struct {
	FirstPassYield float64 // proportion within spec on first pass
	DPMO           float64 // defects per million opportunities
	SigmaLevel     float64 // process sigma level
	RolledYield    float64 // multiplicative yield across multiple parameters
}

// FirstPassYield estimates the yield (proportion within spec) from the Z-values.
func FirstPassYield(mean, sigma float64, spec SpecLimits) float64 {
	if sigma <= 0 {
		return 1
	}
	zUpper := (spec.USL - mean) / sigma
	zLower := (mean - spec.LSL) / sigma
	pAbove := normalCDFComplement(zUpper)
	pBelow := normalCDFComplement(zLower)
	yield := 1 - pAbove - pBelow
	if yield < 0 {
		yield = 0
	}
	return yield
}

// RolledThroughputYield computes the multiplicative yield across multiple
// independent parameters (each with their own yield).
func RolledThroughputYield(yields []float64) float64 {
	if len(yields) == 0 {
		return 0
	}
	rty := 1.0
	for _, y := range yields {
		rty *= y
	}
	return rty
}

// DPMOToSigma converts DPMO to process sigma level.
// Uses the approximate inverse: sigma ≈ 0.8406 + sqrt(29.37 - 2.221*ln(DPMO)).
func DPMOToSigma(dpmo float64) float64 {
	if dpmo <= 0 {
		return 6.0 // perfect process
	}
	if dpmo >= 1e6 {
		return 0 // completely defective
	}
	// Numerical inverse of the normal CDF relationship.
	pDefect := dpmo / 1e6
	zVal := -normalQuantile(pDefect)
	return zVal
}

// SigmaToDPMO converts process sigma level to DPMO.
func SigmaToDPMO(sigma float64) float64 {
	pDefect := normalCDFComplement(sigma)
	return pDefect * 1e6
}

// BatchYield computes the overall yield for a batch given measurements and specs.
// Each parameter contributes independently.
func BatchYield(paramMeans, paramSigmas []float64, specs []SpecLimits) YieldResult {
	if len(paramMeans) == 0 || len(paramMeans) != len(paramSigmas) || len(paramMeans) != len(specs) {
		return YieldResult{}
	}

	yields := make([]float64, len(paramMeans))
	totalDPMO := 0.0
	minSigma := math.MaxFloat64

	for i, mean := range paramMeans {
		sigma := paramSigmas[i]
		spec := specs[i]
		yields[i] = FirstPassYield(mean, sigma, spec)
		dpmo := PPM(mean, sigma, spec)
		totalDPMO += dpmo
		sl := SigmaLevel(mean, sigma, spec)
		if sl < minSigma {
			minSigma = sl
		}
	}

	avgDPMO := totalDPMO / float64(len(paramMeans))
	rty := RolledThroughputYield(yields)

	return YieldResult{
		FirstPassYield: yields[0], // first parameter's yield
		DPMO:           avgDPMO,
		SigmaLevel:     minSigma,
		RolledYield:    rty,
	}
}
