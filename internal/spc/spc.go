// Package spc implements Statistical Process Control charts for PAT batch
// monitoring. Supports X-bar, R-chart, S-chart, and Individual/Moving-Range
// (I-MR) charts with standard control-limit calculations.
package spc

import "math"

// A2, D3, D4 constants for R-chart (subgroup sizes 2-10).
var a2Table = map[int]float64{
	2: 1.880, 3: 1.023, 4: 0.729, 5: 0.577,
	6: 0.483, 7: 0.419, 8: 0.373, 9: 0.337, 10: 0.308,
}
var d3Table = map[int]float64{
	2: 0, 3: 0, 4: 0, 5: 0,
	6: 0, 7: 0.076, 8: 0.136, 9: 0.184, 10: 0.223,
}
var d4Table = map[int]float64{
	2: 3.267, 3: 2.575, 4: 2.282, 5: 2.115,
	6: 2.004, 7: 1.924, 8: 1.864, 9: 1.816, 10: 1.777,
}

// A3, B3, B4 constants for S-chart.
var a3Table = map[int]float64{
	2: 2.659, 3: 1.954, 4: 1.628, 5: 1.427,
	6: 1.287, 7: 1.182, 8: 1.099, 9: 1.032, 10: 0.975,
}
var b3Table = map[int]float64{
	2: 0, 3: 0, 4: 0, 5: 0,
	6: 0.030, 7: 0.118, 8: 0.185, 9: 0.239, 10: 0.284,
}
var b4Table = map[int]float64{
	2: 3.267, 3: 2.568, 4: 2.266, 5: 2.089,
	6: 1.970, 7: 1.882, 8: 1.815, 9: 1.761, 10: 1.716,
}

// d2 constants for estimating sigma from average range.
var d2Table = map[int]float64{
	2: 1.128, 3: 1.693, 4: 2.059, 5: 2.326,
	6: 2.534, 7: 2.704, 8: 2.847, 9: 2.970, 10: 3.078,
}

// ControlLimits holds the center line and upper/lower control limits.
type ControlLimits struct {
	CL  float64
	UCL float64
	LCL float64
}

// ChartPoint represents one plotted point on a control chart.
type ChartPoint struct {
	Index   int
	Value   float64
	InLimit bool
}

// Subgroup represents a rational subgroup of measurements.
type Subgroup struct {
	Values []float64
}

// SubgroupMean returns the arithmetic mean of a subgroup.
func (sg Subgroup) Mean() float64 {
	if len(sg.Values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range sg.Values {
		sum += v
	}
	return sum / float64(len(sg.Values))
}

// SubgroupRange returns the range (max - min) of a subgroup.
func (sg Subgroup) Range() float64 {
	if len(sg.Values) < 2 {
		return 0
	}
	mn, mx := sg.Values[0], sg.Values[0]
	for _, v := range sg.Values[1:] {
		if v < mn {
			mn = v
		}
		if v > mx {
			mx = v
		}
	}
	return mx - mn
}

// SubgroupStd returns the sample standard deviation of a subgroup.
func (sg Subgroup) Std() float64 {
	n := len(sg.Values)
	if n < 2 {
		return 0
	}
	m := sg.Mean()
	var ss float64
	for _, v := range sg.Values {
		d := v - m
		ss += d * d
	}
	return math.Sqrt(ss / float64(n-1))
}

// LookupA2 returns the A2 constant for the given subgroup size, or 0 if unknown.
func LookupA2(n int) float64 { return a2Table[n] }

// LookupD3 returns the D3 constant for R-chart lower limit.
func LookupD3(n int) float64 { return d3Table[n] }

// LookupD4 returns the D4 constant for R-chart upper limit.
func LookupD4(n int) float64 { return d4Table[n] }

// LookupA3 returns the A3 constant for S-chart.
func LookupA3(n int) float64 { return a3Table[n] }

// LookupB3 returns the B3 constant for S-chart lower limit.
func LookupB3(n int) float64 { return b3Table[n] }

// LookupB4 returns the B4 constant for S-chart upper limit.
func LookupB4(n int) float64 { return b4Table[n] }

// LookupD2 returns the d2 constant for sigma estimation from average range.
func LookupD2(n int) float64 { return d2Table[n] }
