package capability

import "math"

// CpkConfidenceInterval computes an approximate lower confidence bound for Cpk
// using the formula: Cpk_lower = Cpk - z_alpha * SE(Cpk).
// SE(Cpk) ≈ sqrt(1/(9*n*Cpk^2) + 1/(2*(n-1))).
func CpkConfidenceInterval(cpk float64, n int, confidence float64) (lower, upper float64) {
	if n <= 1 || cpk <= 0 {
		return 0, 0
	}
	z := normalQuantile((1 + confidence) / 2)
	se := math.Sqrt(1.0/(9*float64(n)*cpk*cpk) + 1.0/(2*float64(n-1)))
	lower = cpk - z*se*cpk
	upper = cpk + z*se*cpk
	return lower, upper
}

// CpConfidenceInterval computes the confidence interval for Cp using the
// chi-squared distribution approximation:
// Cp_lower = Cp * sqrt(chi2_lower / (n-1))
// Cp_upper = Cp * sqrt(chi2_upper / (n-1))
// Approximated using Wilson-Hilferty transformation.
func CpConfidenceInterval(cp float64, n int, confidence float64) (lower, upper float64) {
	if n <= 1 || cp <= 0 {
		return 0, 0
	}
	df := float64(n - 1)
	alpha := 1 - confidence
	chi2Lower := chiSquareQuantile(alpha/2, df)
	chi2Upper := chiSquareQuantile(1-alpha/2, df)
	lower = cp * math.Sqrt(chi2Lower/df)
	upper = cp * math.Sqrt(chi2Upper/df)
	return lower, upper
}

// PpkMinimumSample estimates the minimum sample size needed to demonstrate a
// given Ppk at a specified confidence level. Uses the approximation:
// n >= (z / margin)^2 where margin accounts for acceptable estimation error.
func PpkMinimumSample(targetPpk, confidence float64) int {
	z := normalQuantile((1 + confidence) / 2)
	if targetPpk <= 0 {
		return 0
	}
	// Standard error of Cpk: SE ≈ sqrt(1/(9*n*Cpk^2) + 1/(2*(n-1)))
	// Solve for n such that z*SE*Cpk <= 0.1*Cpk (10% margin).
	// Iterative approach.
	for n := 10; n < 10000; n++ {
		se := math.Sqrt(1.0/(9*float64(n)*targetPpk*targetPpk) + 1.0/(2*float64(n-1)))
		margin := z * se * targetPpk
		if margin <= 0.1*targetPpk {
			return n
		}
	}
	return 10000
}

// normalQuantile returns the z-value for a given cumulative probability p
// using the rational approximation (Beasley-Springer-Moro algorithm simplified).
func normalQuantile(p float64) float64 {
	if p <= 0 || p >= 1 {
		return 0
	}
	if p == 0.5 {
		return 0
	}
	// Rational approximation for 0.5 < p < 1.
	if p < 0.5 {
		return -normalQuantile(1 - p)
	}
	t := math.Sqrt(-2 * math.Log(1-p))
	// Horner form of Abramowitz & Stegun 26.2.23.
	const (
		c0 = 2.515517
		c1 = 0.802853
		c2 = 0.010328
		d1 = 1.432788
		d2 = 0.189269
		d3 = 0.001308
	)
	return t - (c0+c1*t+c2*t*t)/(1+d1*t+d2*t*t+d3*t*t*t)
}

// chiSquareQuantile approximates the chi-square quantile using the
// Wilson-Hilferty cube-root transformation.
func chiSquareQuantile(p float64, df float64) float64 {
	if df <= 0 || p <= 0 || p >= 1 {
		return 0
	}
	z := normalQuantile(p)
	// Wilson-Hilferty: X^(1/3) ≈ 1 - 2/(9*df) + z*sqrt(2/(9*df))
	k := 2.0 / (9 * df)
	x := 1 - k + z*math.Sqrt(k)
	if x <= 0 {
		return 0
	}
	return df * x * x * x
}
