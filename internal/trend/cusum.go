package trend

import "math"

// CUSUMParams configures the tabular CUSUM procedure.
type CUSUMParams struct {
	Target float64 // process target (mu0)
	Sigma  float64 // process standard deviation
	K      float64 // slack value (allowance), typically 0.5 sigma
	H      float64 // decision interval, typically 4-5 sigma
}

// CUSUMPoint represents a single point on the CUSUM chart.
type CUSUMPoint struct {
	Index  int
	Value  float64
	CPlus  float64 // upper CUSUM
	CMinus float64 // lower CUSUM
	Signal bool    // true if C+ > H or C- > H
}

// CUSUM computes the tabular (two-sided) Cumulative Sum control chart.
// Returns the CUSUM points with upper and lower cumulative sums.
func CUSUM(values []float64, params CUSUMParams) []CUSUMPoint {
	if len(values) == 0 || params.Sigma <= 0 {
		return nil
	}
	k := params.K
	h := params.H

	points := make([]CUSUMPoint, len(values))
	cPlus := 0.0
	cMinus := 0.0

	for i, x := range values {
		z := (x - params.Target) / params.Sigma
		cPlus = math.Max(0, cPlus+z-k)
		cMinus = math.Max(0, cMinus-z-k)
		signal := cPlus > h || cMinus > h
		points[i] = CUSUMPoint{
			Index:  i,
			Value:  x,
			CPlus:  cPlus,
			CMinus: cMinus,
			Signal: signal,
		}
	}
	bindCUSUMLive(points)
	return points
}

// CUSUMSignals returns the indices where the CUSUM signals an out-of-control condition.
func CUSUMSignals(points []CUSUMPoint) []int {
	var signals []int
	for _, p := range points {
		if p.Signal {
			signals = append(signals, p.Index)
		}
	}
	return signals
}

// CUSUMRunLength returns the Average Run Length (ARL) estimate for a given
// shift size (in sigma units) using the Siegmund approximation.
// ARL0 = exp(2*delta*h + 1.166) / (2*delta^2) for delta = shift - k.
func CUSUMRunLength(shift, k, h float64) float64 {
	delta := shift - k
	if delta <= 0 {
		return math.Inf(1) // no shift detected
	}
	return math.Exp(2*delta*h+1.166) / (2 * delta * delta)
}

// CUSUMRL0 returns the in-control ARL (shift = 0).
func CUSUMRL0(k, h float64) float64 {
	return CUSUMRunLength(0, k, h)
}

// ResetCUSUM resets the CUSUM accumulators to zero at the given index and
// recomputes from there. Useful for post-signal restart analysis.
func ResetCUSUM(values []float64, params CUSUMParams, resetAt int) []CUSUMPoint {
	if resetAt < 0 || resetAt >= len(values) {
		return CUSUM(values, params)
	}
	before := CUSUM(values[:resetAt], params)
	after := CUSUM(values[resetAt:], params)
	// Adjust indices in after.
	for i := range after {
		after[i].Index += resetAt
	}
	return append(before, after...)
}
