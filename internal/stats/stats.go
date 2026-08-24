// Package stats computes process capability and batch conformance.
package stats

import "math"

// Spec defines the acceptable band for a process parameter.
type Spec struct {
	Parameter string
	Target    float64
	Low       float64
	High      float64
}

// Measurement is a single reading of one parameter in one batch.
type Measurement struct {
	Batch     string
	Parameter string
	Value     float64
}

// ParamStat summarizes one parameter across all batches.
type ParamStat struct {
	Parameter string
	Mean      float64
	Std       float64
	Min       float64
	Max       float64
	CPK       float64
}

// BatchResult summarizes conformance of a single batch.
type BatchResult struct {
	Batch string
	OOS   []string // parameters with out-of-specification values
	OOT   []string // parameters trending out of band (within spec but off-target)
}

// Mean returns the arithmetic mean; 0 for an empty slice.
func Mean(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	var sum float64
	for _, v := range vals {
		sum += v
	}
	return sum / float64(len(vals))
}

// Std returns the population standard deviation; 0 for empty or single value.
func Std(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	m := Mean(vals)
	var s float64
	for _, v := range vals {
		d := v - m
		s += d * d
	}
	return math.Sqrt(s / float64(len(vals)))
}

// CPK returns the process capability index. When std is zero, capability is
// considered effectively unbounded and a large finite value is returned.
func CPK(mean, std, low, high float64) float64 {
	if std <= 0 {
		return 999
	}
	upper := (high - mean) / (3 * std)
	lower := (mean - low) / (3 * std)
	if upper < lower {
		return upper
	}
	return lower
}

// ComputeParamStats groups measurements by parameter and summarizes them.
func ComputeParamStats(ms []Measurement, specs map[string]Spec) []ParamStat {
	byParam := map[string][]float64{}
	var order []string
	for _, m := range ms {
		if _, ok := byParam[m.Parameter]; !ok {
			order = append(order, m.Parameter)
		}
		byParam[m.Parameter] = append(byParam[m.Parameter], m.Value)
	}
	out := make([]ParamStat, 0, len(order))
	for _, p := range order {
		vals := byParam[p]
		st := ParamStat{
			Parameter: p,
			Mean:      Mean(vals),
			Std:       Std(vals),
			Min:       minSlice(vals),
			Max:       maxSlice(vals),
		}
		if sp, ok := specs[p]; ok {
			st.CPK = CPK(st.Mean, st.Std, sp.Low, sp.High)
		}
		out = append(out, st)
	}
	return out
}

// EvaluateBatch checks one batch for out-of-spec (OOS) and out-of-trend (OOT).
func EvaluateBatch(ms []Measurement, specs map[string]Spec) BatchResult {
	res := BatchResult{}
	if len(ms) > 0 {
		res.Batch = ms[0].Batch
	}
	byParam := map[string][]float64{}
	for _, m := range ms {
		sp, ok := specs[m.Parameter]
		if !ok {
			continue
		}
		if m.Value < sp.Low || m.Value > sp.High {
			res.OOS = appendUnique(res.OOS, m.Parameter)
		}
		byParam[m.Parameter] = append(byParam[m.Parameter], m.Value)
	}
	for p, vals := range byParam {
		sp := specs[p]
		band := (sp.High - sp.Low) / 3.0
		mean := Mean(vals)
		if mean < sp.Target-band || mean > sp.Target+band {
			res.OOT = appendUnique(res.OOT, p)
		}
	}
	return res
}

func appendUnique(s []string, v string) []string {
	for _, x := range s {
		if x == v {
			return s
		}
	}
	return append(s, v)
}

func minSlice(v []float64) float64 {
	if len(v) == 0 {
		return 0
	}
	m := v[0]
	for _, x := range v[1:] {
		if x < m {
			m = x
		}
	}
	return m
}

func maxSlice(v []float64) float64 {
	if len(v) == 0 {
		return 0
	}
	m := v[0]
	for _, x := range v[1:] {
		if x > m {
			m = x
		}
	}
	return m
}
