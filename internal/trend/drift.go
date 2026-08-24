package trend

import "math"

// LinearTrend fits a simple linear regression y = a + b*x to the values
// (using index as x) and returns the slope, intercept, and R-squared.
type LinearTrend struct {
	Slope     float64
	Intercept float64
	RSquared  float64
}

// FitLinearTrend fits a linear trend to the values.
func FitLinearTrend(values []float64) LinearTrend {
	n := len(values)
	if n < 2 {
		return LinearTrend{}
	}

	// x = 0, 1, 2, ..., n-1.
	nf := float64(n)
	sumX := nf * (nf - 1) / 2
	sumX2 := nf * (nf - 1) * (2*nf - 1) / 6
	sumY := 0.0
	sumXY := 0.0
	for i, y := range values {
		sumY += y
		sumXY += float64(i) * y
	}

	denom := nf*sumX2 - sumX*sumX
	if math.Abs(denom) < 1e-15 {
		return LinearTrend{Intercept: sumY / nf}
	}
	b := (nf*sumXY - sumX*sumY) / denom
	a := (sumY - b*sumX) / nf

	// R-squared.
	yBar := sumY / nf
	ssTotal := 0.0
	ssResid := 0.0
	for i, y := range values {
		ssTotal += (y - yBar) * (y - yBar)
		pred := a + b*float64(i)
		ssResid += (y - pred) * (y - pred)
	}
	r2 := 0.0
	if ssTotal > 1e-15 {
		r2 = 1 - ssResid/ssTotal
	}

	return LinearTrend{Slope: b, Intercept: a, RSquared: r2}
}

// HasSignificantTrend returns true if the linear trend slope produces an R²
// above the given threshold (e.g., 0.5 for a meaningful trend).
func HasSignificantTrend(values []float64, r2Threshold float64) bool {
	trend := FitLinearTrend(values)
	return trend.RSquared >= r2Threshold && math.Abs(trend.Slope) > 1e-10
}

// DriftRate computes the per-batch drift (slope) and estimates the number of
// batches until a specification limit is breached, given the current mean level.
func DriftRate(values []float64, specLow, specHigh float64) (slope float64, batchesToBreach int) {
	trend := FitLinearTrend(values)
	slope = trend.Slope
	if math.Abs(slope) < 1e-10 {
		return slope, -1 // no drift, never breaches
	}

	n := len(values)
	currentLevel := trend.Intercept + slope*float64(n-1)

	if slope > 0 {
		// Drifting up toward USL.
		remaining := specHigh - currentLevel
		if remaining <= 0 {
			return slope, 0 // already breached
		}
		return slope, int(math.Ceil(remaining / slope))
	}
	// Drifting down toward LSL.
	remaining := currentLevel - specLow
	if remaining <= 0 {
		return slope, 0 // already breached
	}
	return slope, int(math.Ceil(remaining / (-slope)))
}

// ChangePointDetection identifies potential change points in the series using
// a simple two-segment mean comparison. Returns the index where the mean
// shift is most pronounced and the magnitude of the shift.
func ChangePointDetection(values []float64) (changeIdx int, shiftMag float64) {
	n := len(values)
	if n < 4 {
		return -1, 0
	}

	bestIdx := -1
	bestShift := 0.0

	for i := 2; i < n-1; i++ {
		// Mean of segment 1 (0..i-1) and segment 2 (i..n-1).
		sum1 := 0.0
		for j := 0; j < i; j++ {
			sum1 += values[j]
		}
		mean1 := sum1 / float64(i)

		sum2 := 0.0
		for j := i; j < n; j++ {
			sum2 += values[j]
		}
		mean2 := sum2 / float64(n-i)

		shift := math.Abs(mean2 - mean1)
		if shift > bestShift {
			bestShift = shift
			bestIdx = i
		}
	}
	return bestIdx, bestShift
}

// SegmentMeans splits the data at the given index and returns the mean of each segment.
func SegmentMeans(values []float64, splitAt int) (mean1, mean2 float64) {
	if splitAt <= 0 || splitAt >= len(values) {
		return 0, 0
	}
	sum1 := 0.0
	for i := 0; i < splitAt; i++ {
		sum1 += values[i]
	}
	mean1 = sum1 / float64(splitAt)

	sum2 := 0.0
	for i := splitAt; i < len(values); i++ {
		sum2 += values[i]
	}
	mean2 = sum2 / float64(len(values)-splitAt)
	return mean1, mean2
}
