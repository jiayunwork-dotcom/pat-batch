package stats

import "math"

// HotellingT2 computes the Hotelling T-squared statistic for a multivariate
// observation vector x relative to a mean vector and inverse covariance matrix.
// This is a simplified 2-variable version commonly used in PAT batch analysis.
func HotellingT2(x, mean []float64, invCov [][]float64) float64 {
	p := len(x)
	if p == 0 || len(mean) != p || len(invCov) != p {
		return 0
	}
	diff := make([]float64, p)
	for i := range diff {
		diff[i] = x[i] - mean[i]
	}
	// T² = diff' * invCov * diff
	t2 := 0.0
	for i := 0; i < p; i++ {
		for j := 0; j < p; j++ {
			t2 += diff[i] * invCov[i][j] * diff[j]
		}
	}
	return t2
}

// CovarianceMatrix computes the sample covariance matrix from a set of
// multivariate observations. Each row in data is one observation (length p).
func CovarianceMatrix(data [][]float64) [][]float64 {
	n := len(data)
	if n < 2 {
		return nil
	}
	p := len(data[0])
	if p == 0 {
		return nil
	}

	// Compute means.
	means := make([]float64, p)
	for _, obs := range data {
		for j := range obs {
			means[j] += obs[j]
		}
	}
	for j := range means {
		means[j] /= float64(n)
	}

	// Compute covariance.
	cov := make([][]float64, p)
	for i := range cov {
		cov[i] = make([]float64, p)
	}
	for _, obs := range data {
		for i := 0; i < p; i++ {
			for j := 0; j < p; j++ {
				cov[i][j] += (obs[i] - means[i]) * (obs[j] - means[j])
			}
		}
	}
	for i := 0; i < p; i++ {
		for j := 0; j < p; j++ {
			cov[i][j] /= float64(n - 1)
		}
	}
	return cov
}

// Invert2x2 returns the inverse of a 2x2 matrix, or nil if singular.
func Invert2x2(m [][]float64) [][]float64 {
	if len(m) != 2 || len(m[0]) != 2 || len(m[1]) != 2 {
		return nil
	}
	det := m[0][0]*m[1][1] - m[0][1]*m[1][0]
	if math.Abs(det) < 1e-15 {
		return nil
	}
	inv := [][]float64{
		{m[1][1] / det, -m[0][1] / det},
		{-m[1][0] / det, m[0][0] / det},
	}
	return inv
}

// MahalanobisDistance computes the Mahalanobis distance for a 2-variable
// observation from the group mean using the inverse covariance matrix.
func MahalanobisDistance(x, mean []float64, invCov [][]float64) float64 {
	t2 := HotellingT2(x, mean, invCov)
	if t2 < 0 {
		return 0
	}
	return math.Sqrt(t2)
}

// Correlation computes the Pearson correlation coefficient between two slices.
func Correlation(x, y []float64) float64 {
	n := len(x)
	if n != len(y) || n < 2 {
		return 0
	}
	mx := Mean(x)
	my := Mean(y)
	var num, dx2, dy2 float64
	for i := 0; i < n; i++ {
		dx := x[i] - mx
		dy := y[i] - my
		num += dx * dy
		dx2 += dx * dx
		dy2 += dy * dy
	}
	denom := math.Sqrt(dx2 * dy2)
	if denom < 1e-15 {
		return 0
	}
	return num / denom
}

// Autocorrelation computes the lag-k autocorrelation of a time series.
func Autocorrelation(vals []float64, lag int) float64 {
	n := len(vals)
	if lag < 0 || lag >= n || n < 2 {
		return 0
	}
	m := Mean(vals)
	var num, denom float64
	for i := 0; i < n; i++ {
		denom += (vals[i] - m) * (vals[i] - m)
	}
	if denom < 1e-15 {
		return 0
	}
	for i := 0; i < n-lag; i++ {
		num += (vals[i] - m) * (vals[i+lag] - m)
	}
	return num / denom
}

// MovingRange computes the absolute differences between successive values.
func MovingRange(vals []float64) []float64 {
	if len(vals) < 2 {
		return nil
	}
	mr := make([]float64, len(vals)-1)
	for i := 1; i < len(vals); i++ {
		mr[i-1] = math.Abs(vals[i] - vals[i-1])
	}
	return mr
}
