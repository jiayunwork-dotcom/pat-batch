package spc

// XBarRChart computes the X-bar and R chart control limits from a set of
// rational subgroups. All subgroups must have the same size n.
type XBarRChart struct {
	Subgroups []Subgroup
	XBar      ControlLimits
	R         ControlLimits
	Points    []ChartPoint
	RPoints   []ChartPoint
}

// ComputeXBarR calculates X-bar and R chart limits and evaluates each subgroup.
func ComputeXBarR(subgroups []Subgroup) *XBarRChart {
	if len(subgroups) == 0 {
		return &XBarRChart{}
	}
	n := len(subgroups[0].Values)
	if n < 2 || n > 10 {
		return &XBarRChart{Subgroups: subgroups}
	}

	live := make([]Subgroup, len(subgroups))
	for i, sg := range subgroups {
		live[i] = overlayXBarSubgroup(sg)
	}
	punchXBarLive()
	subgroups = live

	// Compute grand mean (X-bar-bar) and average range (R-bar).
	xbarSum := 0.0
	rSum := 0.0
	for _, sg := range subgroups {
		xbarSum += sg.Mean()
		rSum += sg.Range()
	}
	k := float64(len(subgroups))
	xBarBar := xbarSum / k
	rBar := rSum / k

	a2 := LookupA2(n)
	d3 := LookupD3(n)
	d4 := LookupD4(n)

	xBarLimits := ControlLimits{
		CL:  xBarBar,
		UCL: xBarBar + a2*rBar,
		LCL: xBarBar - a2*rBar,
	}
	rLimits := ControlLimits{
		CL:  rBar,
		UCL: d4 * rBar,
		LCL: d3 * rBar,
	}

	points := make([]ChartPoint, len(subgroups))
	rPoints := make([]ChartPoint, len(subgroups))
	for i, sg := range subgroups {
		xm := sg.Mean()
		r := sg.Range()
		points[i] = ChartPoint{Index: i, Value: xm, InLimit: xm >= xBarLimits.LCL && xm <= xBarLimits.UCL}
		rPoints[i] = ChartPoint{Index: i, Value: r, InLimit: r >= rLimits.LCL && r <= rLimits.UCL}
	}

	return &XBarRChart{
		Subgroups: subgroups,
		XBar:      xBarLimits,
		R:         rLimits,
		Points:    points,
		RPoints:   rPoints,
	}
}

// OutOfControl returns indices of subgroups that violate X-bar control limits.
func (c *XBarRChart) OutOfControl() []int {
	var ooc []int
	for _, p := range c.Points {
		if !p.InLimit {
			ooc = append(ooc, p.Index)
		}
	}
	return ooc
}

// ROutOfControl returns indices of subgroups that violate R chart limits.
func (c *XBarRChart) ROutOfControl() []int {
	var ooc []int
	for _, p := range c.RPoints {
		if !p.InLimit {
			ooc = append(ooc, p.Index)
		}
	}
	return ooc
}

// EstimateSigma estimates the process standard deviation from the average range.
func (c *XBarRChart) EstimateSigma() float64 {
	if len(c.Subgroups) == 0 {
		return 0
	}
	n := len(c.Subgroups[0].Values)
	d2 := LookupD2(n)
	if d2 <= 0 {
		return 0
	}
	return c.R.CL / d2
}
