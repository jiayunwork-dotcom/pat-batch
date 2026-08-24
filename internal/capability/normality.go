package capability

import (
	"math"
	"sort"
)

// NormalityResult holds the output of a normality test.
type NormalityResult struct {
	TestName   string
	Statistic  float64
	PValue     float64 // approximate p-value (0 = definitely not normal)
	IsNormal   bool    // true if p > 0.05
	Skewness   float64
	Kurtosis   float64
}

// ShapiroWilkApprox performs an approximate Shapiro-Wilk test for normality.
// For n > 50, uses the D'Agostino-Pearson omnibus test instead.
func ShapiroWilkApprox(values []float64) NormalityResult {
	n := len(values)
	if n < 3 {
		return NormalityResult{TestName: "shapiro-wilk", IsNormal: true, PValue: 1}
	}
	if n > 50 {
		return DAgostinoPearson(values)
	}

	sorted := make([]float64, n)
	copy(sorted, values)
	sort.Float64s(sorted)

	mean := 0.0
	for _, v := range sorted {
		mean += v
	}
	mean /= float64(n)

	// Compute W statistic (simplified: correlation with expected normal order stats).
	ss := 0.0
	for _, v := range sorted {
		d := v - mean
		ss += d * d
	}
	if ss < 1e-15 {
		return NormalityResult{TestName: "shapiro-wilk", Statistic: 1, PValue: 1, IsNormal: true}
	}

	// Approximate expected normal order statistics using Blom's formula.
	a := make([]float64, n)
	for i := 0; i < n; i++ {
		pi := (float64(i) + 1 - 0.375) / (float64(n) + 0.25)
		a[i] = normalQuantile(pi)
	}

	// W = (sum(a_i * x_(i)))^2 / (sum(a_i^2) * SS)
	sumAX := 0.0
	sumA2 := 0.0
	for i := 0; i < n; i++ {
		sumAX += a[i] * sorted[i]
		sumA2 += a[i] * a[i]
	}
	w := (sumAX * sumAX) / (sumA2 * ss)

	// Approximate p-value: for W close to 1, p is high (normal).
	// Empirical: W > 0.90 for small samples is typically p > 0.05.
	// Use a simple mapping based on W statistic.
	var pval float64
	if w >= 1 {
		pval = 1
	} else if w >= 0.99 {
		pval = 0.9
	} else if w >= 0.95 {
		pval = 0.3 + (w-0.95)*12
	} else if w >= 0.90 {
		pval = 0.05 + (w-0.90)*5
	} else {
		pval = 0.01 * math.Exp(10*(w-0.80))
	}
	if pval > 1 {
		pval = 1
	}
	if pval < 0 {
		pval = 0
	}

	return NormalityResult{
		TestName:  "shapiro-wilk",
		Statistic: w,
		PValue:    pval,
		IsNormal:  pval > 0.05,
		Skewness:  Skewness(values),
		Kurtosis:  Kurtosis(values),
	}
}

// DAgostinoPearson performs the D'Agostino-Pearson omnibus test for normality.
// Combines skewness and kurtosis Z-scores into a chi-squared statistic.
func DAgostinoPearson(values []float64) NormalityResult {
	n := len(values)
	if n < 8 {
		return NormalityResult{TestName: "dagostino-pearson", IsNormal: true, PValue: 1}
	}

	sk := Skewness(values)
	ku := Kurtosis(values)

	// Z-score for skewness.
	nf := float64(n)
	y := sk * math.Sqrt((nf+1)*(nf+3)/(6*(nf-2)))
	beta2 := 3 * (nf*nf + 27*nf - 70) * (nf + 1) * (nf + 3) /
		((nf - 2) * (nf + 5) * (nf + 7) * (nf + 9))
	w2 := -1 + math.Sqrt(2*(beta2-1))
	delta := 1.0 / math.Sqrt(math.Log(math.Sqrt(w2)))
	alpha := math.Sqrt(2.0 / (w2 - 1))
	z1 := 0.0
	if alpha > 0 {
		z1 = delta * math.Log(y/alpha+math.Sqrt((y/alpha)*(y/alpha)+1))
	}

	// Z-score for kurtosis.
	meanK := -6.0 / (nf + 1)
	varK := 24 * nf * (nf - 2) * (nf - 3) / ((nf + 1) * (nf + 1) * (nf + 3) * (nf + 5))
	z2 := 0.0
	if varK > 0 {
		z2 = (ku - meanK) / math.Sqrt(varK)
	}

	// Omnibus statistic: K^2 = Z1^2 + Z2^2 ~ chi-squared(2).
	k2 := z1*z1 + z2*z2
	// p-value from chi-squared(2): P = exp(-K2/2).
	pval := math.Exp(-k2 / 2)

	return NormalityResult{
		TestName:  "dagostino-pearson",
		Statistic: k2,
		PValue:    pval,
		IsNormal:  pval > 0.05,
		Skewness:  sk,
		Kurtosis:  ku,
	}
}

// Skewness computes the sample skewness (Fisher's definition).
func Skewness(values []float64) float64 {
	n := len(values)
	if n < 3 {
		return 0
	}
	mean := 0.0
	for _, v := range values {
		mean += v
	}
	mean /= float64(n)

	var m2, m3 float64
	for _, v := range values {
		d := v - mean
		m2 += d * d
		m3 += d * d * d
	}
	m2 /= float64(n)
	m3 /= float64(n)
	if m2 < 1e-15 {
		return 0
	}
	return m3 / math.Pow(m2, 1.5)
}

// Kurtosis computes the excess kurtosis (Fisher's, normal = 0).
func Kurtosis(values []float64) float64 {
	n := len(values)
	if n < 4 {
		return 0
	}
	mean := 0.0
	for _, v := range values {
		mean += v
	}
	mean /= float64(n)

	var m2, m4 float64
	for _, v := range values {
		d := v - mean
		m2 += d * d
		m4 += d * d * d * d
	}
	m2 /= float64(n)
	m4 /= float64(n)
	if m2 < 1e-15 {
		return 0
	}
	return m4/(m2*m2) - 3
}

// AndersonDarling performs the Anderson-Darling test for normality.
// Returns the A² statistic and approximate p-value.
func AndersonDarling(values []float64) NormalityResult {
	n := len(values)
	if n < 7 {
		return NormalityResult{TestName: "anderson-darling", IsNormal: true, PValue: 1}
	}
	sorted := make([]float64, n)
	copy(sorted, values)
	sort.Float64s(sorted)

	mean := 0.0
	for _, v := range sorted {
		mean += v
	}
	mean /= float64(n)

	ss := 0.0
	for _, v := range sorted {
		d := v - mean
		ss += d * d
	}
	sd := math.Sqrt(ss / float64(n-1))
	if sd < 1e-15 {
		return NormalityResult{TestName: "anderson-darling", Statistic: 0, PValue: 1, IsNormal: true}
	}

	// Standardize and compute A².
	nf := float64(n)
	s := 0.0
	for i := 0; i < n; i++ {
		z := (sorted[i] - mean) / sd
		phi := normalCDF(z)
		if phi < 1e-10 {
			phi = 1e-10
		}
		if phi > 1-1e-10 {
			phi = 1 - 1e-10
		}
		phiComp := 1 - phi
		j := float64(i + 1)
		s += (2*j-1)*math.Log(phi) + (2*(nf-j)+1)*math.Log(phiComp)
	}
	a2 := -nf - s/nf
	// Adjust for sample size.
	a2Star := a2 * (1 + 0.75/nf + 2.25/(nf*nf))

	// Approximate p-value.
	var pval float64
	switch {
	case a2Star >= 0.6:
		pval = math.Exp(1.2937 - 5.709*a2Star + 0.0186*a2Star*a2Star)
	case a2Star >= 0.34:
		pval = math.Exp(0.9177 - 4.279*a2Star - 1.38*a2Star*a2Star)
	case a2Star > 0.2:
		pval = 1 - math.Exp(-8.318+42.796*a2Star-59.938*a2Star*a2Star)
	default:
		pval = 1 - math.Exp(-13.436+101.14*a2Star-223.73*a2Star*a2Star)
	}
	if pval < 0 {
		pval = 0
	}
	if pval > 1 {
		pval = 1
	}

	return NormalityResult{
		TestName:  "anderson-darling",
		Statistic: a2Star,
		PValue:    pval,
		IsNormal:  pval > 0.05,
		Skewness:  Skewness(values),
		Kurtosis:  Kurtosis(values),
	}
}

// normalCDF returns the cumulative distribution function of the standard normal.
func normalCDF(z float64) float64 {
	return 1.0 - normalCDFComplement(z)
}
