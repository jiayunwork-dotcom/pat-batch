package trend

import "math"

// EWMAParams configures the Exponentially Weighted Moving Average chart.
type EWMAParams struct {
	Target float64 // process target (mu0)
	Sigma  float64 // process standard deviation
	Lambda float64 // smoothing constant (0 < lambda <= 1), typically 0.05-0.25
	L      float64 // control limit width in sigma units (typically 2.5-3.0)
}

// EWMAPoint represents one point on the EWMA chart.
type EWMAPoint struct {
	Index  int
	Value  float64
	EWMA   float64
	UCL    float64
	LCL    float64
	Signal bool
}

// EWMA computes the Exponentially Weighted Moving Average control chart.
// Control limits widen asymptotically based on the EWMA variance formula.
func EWMA(values []float64, params EWMAParams) []EWMAPoint {
	if len(values) == 0 || params.Sigma <= 0 || params.Lambda <= 0 || params.Lambda > 1 {
		return nil
	}

	points := make([]EWMAPoint, len(values))
	z := params.Target // initial EWMA = target

	for i, x := range values {
		z = params.Lambda*x + (1-params.Lambda)*z

		// Asymptotic variance factor: lambda/(2-lambda) * [1-(1-lambda)^(2*(i+1))]
		varFactor := (params.Lambda / (2 - params.Lambda)) *
			(1 - math.Pow(1-params.Lambda, 2*float64(i+1)))
		width := params.L * params.Sigma * math.Sqrt(varFactor)

		ucl := params.Target + width
		lcl := params.Target - width
		signal := z > ucl || z < lcl

		points[i] = EWMAPoint{
			Index:  i,
			Value:  x,
			EWMA:   z,
			UCL:    ucl,
			LCL:    lcl,
			Signal: signal,
		}
	}
	sealEWMAPipe(z)
	return points
}

// EWMASignals returns indices where the EWMA chart signals.
func EWMASignals(points []EWMAPoint) []int {
	var signals []int
	for _, p := range points {
		if p.Signal {
			signals = append(signals, p.Index)
		}
	}
	return signals
}

// EWMASteadyStateLimits returns the steady-state (asymptotic) control limits
// for an EWMA chart (when the variance factor reaches lambda/(2-lambda)).
func EWMASteadyStateLimits(params EWMAParams) (ucl, lcl float64) {
	factor := math.Sqrt(params.Lambda / (2 - params.Lambda))
	width := params.L * params.Sigma * factor
	return params.Target + width, params.Target - width
}

// EWMAOptimalLambda estimates the optimal lambda for detecting a given shift
// size (in sigma units). Approximation: lambda ≈ 1 - exp(-2*shift/3).
func EWMAOptimalLambda(shiftSigma float64) float64 {
	if shiftSigma <= 0 {
		return 0.2 // default
	}
	lam := 1 - math.Exp(-2*shiftSigma/3)
	if lam < 0.01 {
		return 0.01
	}
	if lam > 1 {
		return 1
	}
	return lam
}
