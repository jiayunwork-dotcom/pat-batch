package spc

import (
	"math"
	"testing"
)

func TestXBarRChartBasic(t *testing.T) {
	subgroups := []Subgroup{
		{Values: []float64{25.0, 26.0, 24.5, 25.5, 25.0}},
		{Values: []float64{25.2, 24.8, 25.3, 25.0, 25.1}},
		{Values: []float64{24.5, 25.5, 25.0, 24.8, 25.2}},
		{Values: []float64{25.1, 25.0, 24.9, 25.3, 25.0}},
		{Values: []float64{25.4, 24.6, 25.0, 25.2, 24.8}},
	}
	chart := ComputeXBarR(subgroups)
	if chart.XBar.CL == 0 {
		t.Fatal("X-bar CL should not be zero")
	}
	if chart.XBar.UCL <= chart.XBar.CL {
		t.Fatalf("UCL %v should be > CL %v", chart.XBar.UCL, chart.XBar.CL)
	}
	if chart.XBar.LCL >= chart.XBar.CL {
		t.Fatalf("LCL %v should be < CL %v", chart.XBar.LCL, chart.XBar.CL)
	}
	if chart.R.CL <= 0 {
		t.Fatal("R-bar should be > 0")
	}
}

func TestXBarROutOfControl(t *testing.T) {
	subgroups := []Subgroup{
		{Values: []float64{10, 10, 10, 10, 10}},
		{Values: []float64{10, 10, 10, 10, 10}},
		{Values: []float64{10, 10, 10, 10, 10}},
		{Values: []float64{50, 50, 50, 50, 50}}, // big shift
	}
	chart := ComputeXBarR(subgroups)
	ooc := chart.OutOfControl()
	if len(ooc) == 0 {
		t.Fatal("expected at least one out-of-control point")
	}
}

func TestXBarSChart(t *testing.T) {
	subgroups := []Subgroup{
		{Values: []float64{20, 21, 19, 20, 20}},
		{Values: []float64{20, 20, 20, 20, 20}},
		{Values: []float64{19, 21, 20, 20, 20}},
	}
	chart := ComputeXBarS(subgroups)
	if chart.XBar.CL == 0 {
		t.Fatal("X-bar CL = 0")
	}
	if chart.S.CL <= 0 {
		t.Fatal("S-bar should be > 0")
	}
}

func TestIMRChart(t *testing.T) {
	values := []float64{10.1, 10.0, 10.2, 9.9, 10.0, 10.1, 9.8, 10.3, 10.0, 10.1}
	chart := ComputeIMR(values)
	if chart.I.CL == 0 {
		t.Fatal("I CL = 0")
	}
	if chart.MR.CL <= 0 {
		t.Fatal("MR-bar should be > 0")
	}
	// All values should be in control for this stable data.
	ooc := chart.IOutOfControl()
	if len(ooc) > 0 {
		t.Fatalf("expected all in control, got OOC at %v", ooc)
	}
}

func TestIMROutOfControl(t *testing.T) {
	values := []float64{10, 10, 10, 10, 10, 10, 10, 10, 10, 50}
	chart := ComputeIMR(values)
	ooc := chart.IOutOfControl()
	if len(ooc) == 0 {
		t.Fatal("expected OOC for the large value")
	}
}

func TestEstimateSigma(t *testing.T) {
	subgroups := []Subgroup{
		{Values: []float64{10, 12}},
		{Values: []float64{11, 13}},
		{Values: []float64{10, 14}},
		{Values: []float64{12, 12}},
	}
	chart := ComputeXBarR(subgroups)
	sigma := chart.EstimateSigma()
	if sigma <= 0 || math.IsNaN(sigma) {
		t.Fatalf("sigma = %v, want > 0", sigma)
	}
}

func TestSubgroupOperations(t *testing.T) {
	sg := Subgroup{Values: []float64{10, 20, 30}}
	if sg.Mean() != 20 {
		t.Fatalf("mean = %v, want 20", sg.Mean())
	}
	if sg.Range() != 20 {
		t.Fatalf("range = %v, want 20", sg.Range())
	}
	if sg.Std() <= 0 {
		t.Fatalf("std = %v, want > 0", sg.Std())
	}
}
