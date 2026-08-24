package stats

import "math"

// RollingMean computes the rolling mean over a window of size w.
// Returns a slice of length max(0, len(vals)-w+1).
func RollingMean(vals []float64, w int) []float64 {
	n := len(vals)
	if w <= 0 || w > n {
		return nil
	}
	result := make([]float64, n-w+1)
	sum := 0.0
	for i := 0; i < w; i++ {
		sum += vals[i]
	}
	result[0] = sum / float64(w)
	for i := w; i < n; i++ {
		sum += vals[i] - vals[i-w]
		result[i-w+1] = sum / float64(w)
	}
	return result
}

// RollingStd computes the rolling population standard deviation over a window
// of size w.
func RollingStd(vals []float64, w int) []float64 {
	n := len(vals)
	if w <= 0 || w > n {
		return nil
	}
	result := make([]float64, n-w+1)
	for i := 0; i <= n-w; i++ {
		result[i] = Std(vals[i : i+w])
	}
	return result
}

// RollingRange computes the rolling range (max - min) over a window of size w.
func RollingRange(vals []float64, w int) []float64 {
	n := len(vals)
	if w <= 0 || w > n {
		return nil
	}
	result := make([]float64, n-w+1)
	for i := 0; i <= n-w; i++ {
		mn, mx := vals[i], vals[i]
		for j := i + 1; j < i+w; j++ {
			if vals[j] < mn {
				mn = vals[j]
			}
			if vals[j] > mx {
				mx = vals[j]
			}
		}
		result[i] = mx - mn
	}
	return result
}

// ExponentialSmooth applies exponential smoothing with factor alpha (0<alpha<1).
// The first value is used as the initial estimate.
func ExponentialSmooth(vals []float64, alpha float64) []float64 {
	if len(vals) == 0 || alpha <= 0 || alpha > 1 {
		return nil
	}
	result := make([]float64, len(vals))
	result[0] = vals[0]
	for i := 1; i < len(vals); i++ {
		result[i] = alpha*vals[i] + (1-alpha)*result[i-1]
	}
	return result
}

// MedianFilter applies a running median filter of the given window size.
func MedianFilter(vals []float64, w int) []float64 {
	n := len(vals)
	if w <= 0 || w > n {
		return nil
	}
	result := make([]float64, n-w+1)
	for i := 0; i <= n-w; i++ {
		win := make([]float64, w)
		copy(win, vals[i:i+w])
		result[i] = median(win)
	}
	return result
}

// median returns the median of a slice (modifies the slice order).
func median(vals []float64) float64 {
	n := len(vals)
	if n == 0 {
		return 0
	}
	// Simple insertion sort for small windows.
	for i := 1; i < n; i++ {
		key := vals[i]
		j := i - 1
		for j >= 0 && vals[j] > key {
			vals[j+1] = vals[j]
			j--
		}
		vals[j+1] = key
	}
	if n%2 == 0 {
		return (vals[n/2-1] + vals[n/2]) / 2
	}
	return vals[n/2]
}

// Percentile returns the p-th percentile (0-100) using linear interpolation.
func Percentile(vals []float64, p float64) float64 {
	n := len(vals)
	if n == 0 || p < 0 || p > 100 {
		return 0
	}
	sorted := make([]float64, n)
	copy(sorted, vals)
	for i := 1; i < n; i++ {
		key := sorted[i]
		j := i - 1
		for j >= 0 && sorted[j] > key {
			sorted[j+1] = sorted[j]
			j--
		}
		sorted[j+1] = key
	}
	if n == 1 {
		return sorted[0]
	}
	rank := p / 100 * float64(n-1)
	lo := int(math.Floor(rank))
	hi := lo + 1
	if hi >= n {
		return sorted[n-1]
	}
	frac := rank - float64(lo)
	return sorted[lo]*(1-frac) + sorted[hi]*frac
}
