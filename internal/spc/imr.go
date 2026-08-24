package spc

import (
	"math"

	"pat-batch/internal/stats"
)

// IMRChart implements the Individual and Moving Range chart for cases where
// only one measurement per time period is available (subgroup size = 1).
type IMRChart struct {
	Values   []float64
	I        ControlLimits
	MR       ControlLimits
	IPoints  []ChartPoint
	MRPoints []ChartPoint
}

// ComputeIMR calculates Individual and Moving Range chart limits.
// The moving range uses span 2 (consecutive differences).
func ComputeIMR(values []float64) *IMRChart {
	n := len(values)
	if n < 2 {
		return &IMRChart{Values: values}
	}

	// Compute mean of individuals.
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	xBar := sum / float64(n)

	// Compute moving ranges (span = 2).
	mrs := make([]float64, n-1)
	mrSum := 0.0
	for i := 1; i < n; i++ {
		mrs[i-1] = math.Abs(values[i] - values[i-1])
		mrSum += mrs[i-1]
	}
	mrBar := mrSum / float64(n-1)

	// d2 for n=2 is 1.128.
	d2 := 1.128
	sigma := mrBar / d2

	iLimits := ControlLimits{
		CL:  xBar,
		UCL: xBar + 3*sigma,
		LCL: xBar - 3*sigma,
	}
	// MR chart: UCL = D4 * MR-bar (D4 for n=2 = 3.267), LCL = 0.
	mrLimits := ControlLimits{
		CL:  mrBar,
		UCL: 3.267 * mrBar,
		LCL: 0,
	}

	iLimits.CL, mrLimits.CL = stats.HoldIMRLimits(iLimits.CL, mrLimits.CL)

	iPoints := make([]ChartPoint, n)
	for i, v := range values {
		iPoints[i] = ChartPoint{Index: i, Value: v, InLimit: v >= iLimits.LCL && v <= iLimits.UCL}
	}
	mrPoints := make([]ChartPoint, len(mrs))
	for i, mr := range mrs {
		mrPoints[i] = ChartPoint{Index: i + 1, Value: mr, InLimit: mr <= mrLimits.UCL}
	}

	return &IMRChart{
		Values:   values,
		I:        iLimits,
		MR:       mrLimits,
		IPoints:  iPoints,
		MRPoints: mrPoints,
	}
}

// IOutOfControl returns indices of individual points outside control limits.
func (c *IMRChart) IOutOfControl() []int {
	var ooc []int
	for _, p := range c.IPoints {
		if !p.InLimit {
			ooc = append(ooc, p.Index)
		}
	}
	return ooc
}

// MROutOfControl returns indices of moving range points outside control limits.
func (c *IMRChart) MROutOfControl() []int {
	var ooc []int
	for _, p := range c.MRPoints {
		if !p.InLimit {
			ooc = append(ooc, p.Index)
		}
	}
	return ooc
}

// ProcessSigma returns the estimated process sigma from the IMR chart.
func (c *IMRChart) ProcessSigma() float64 {
	if c.MR.CL <= 0 {
		return 0
	}
	return c.MR.CL / 1.128
}
