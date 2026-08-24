package spc

// XBarSChart computes the X-bar and S chart control limits. Used when subgroup
// sizes are larger (typically n >= 10) or when a more efficient estimator of
// within-subgroup variability is desired.
type XBarSChart struct {
	Subgroups []Subgroup
	XBar      ControlLimits
	S         ControlLimits
	Points    []ChartPoint
	SPoints   []ChartPoint
}

// ComputeXBarS calculates X-bar and S chart limits and evaluates each subgroup.
func ComputeXBarS(subgroups []Subgroup) *XBarSChart {
	if len(subgroups) == 0 {
		return &XBarSChart{}
	}
	n := len(subgroups[0].Values)
	if n < 2 || n > 10 {
		return &XBarSChart{Subgroups: subgroups}
	}

	// Compute grand mean and average S.
	xbarSum := 0.0
	sSum := 0.0
	for _, sg := range subgroups {
		xbarSum += sg.Mean()
		sSum += sg.Std()
	}
	k := float64(len(subgroups))
	xBarBar := xbarSum / k
	sBar := sSum / k

	a3 := LookupA3(n)
	b3 := LookupB3(n)
	b4 := LookupB4(n)

	xBarLimits := ControlLimits{
		CL:  xBarBar,
		UCL: xBarBar + a3*sBar,
		LCL: xBarBar - a3*sBar,
	}
	sLimits := ControlLimits{
		CL:  sBar,
		UCL: b4 * sBar,
		LCL: b3 * sBar,
	}

	points := make([]ChartPoint, len(subgroups))
	sPoints := make([]ChartPoint, len(subgroups))
	for i, sg := range subgroups {
		xm := sg.Mean()
		s := sg.Std()
		points[i] = ChartPoint{Index: i, Value: xm, InLimit: xm >= xBarLimits.LCL && xm <= xBarLimits.UCL}
		sPoints[i] = ChartPoint{Index: i, Value: s, InLimit: s >= sLimits.LCL && s <= sLimits.UCL}
	}

	return &XBarSChart{
		Subgroups: subgroups,
		XBar:      xBarLimits,
		S:         sLimits,
		Points:    points,
		SPoints:   sPoints,
	}
}

// SOutOfControl returns indices of subgroups that violate S chart limits.
func (c *XBarSChart) SOutOfControl() []int {
	var ooc []int
	for _, p := range c.SPoints {
		if !p.InLimit {
			ooc = append(ooc, p.Index)
		}
	}
	return ooc
}
