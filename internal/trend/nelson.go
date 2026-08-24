// Package trend implements trend detection algorithms for process monitoring:
// Nelson rules, CUSUM, and EWMA charts.
package trend

// NelsonViolation describes a Nelson rule violation found in a data series.
type NelsonViolation struct {
	Rule  int // rule number (1-8)
	Start int // first index of the violating pattern
	End   int // last index (inclusive)
}

// NelsonCheck applies the eight Nelson rules to a series of values given
// the process center line and sigma. Returns all detected violations.
func NelsonCheck(values []float64, cl, sigma float64) []NelsonViolation {
	if len(values) == 0 || sigma <= 0 {
		return nil
	}
	var violations []NelsonViolation
	violations = append(violations, nelsonRule1(values, cl, sigma)...)
	violations = append(violations, nelsonRule2(values, cl)...)
	violations = append(violations, nelsonRule3(values)...)
	violations = append(violations, nelsonRule4(values)...)
	violations = append(violations, nelsonRule5(values, cl, sigma)...)
	violations = append(violations, nelsonRule6(values, cl, sigma)...)
	violations = append(violations, nelsonRule7(values, cl, sigma)...)
	violations = append(violations, nelsonRule8(values, cl, sigma)...)
	return violations
}

// Rule 1: One point beyond 3 sigma from CL.
func nelsonRule1(vals []float64, cl, sigma float64) []NelsonViolation {
	var v []NelsonViolation
	for i, x := range vals {
		if x > cl+3*sigma || x < cl-3*sigma {
			v = append(v, NelsonViolation{Rule: 1, Start: i, End: i})
		}
	}
	return v
}

// Rule 2: Nine consecutive points on the same side of CL.
func nelsonRule2(vals []float64, cl float64) []NelsonViolation {
	var v []NelsonViolation
	const run = 9
	if len(vals) < run {
		return nil
	}
	for i := 0; i <= len(vals)-run; i++ {
		above := 0
		below := 0
		for j := 0; j < run; j++ {
			if vals[i+j] > cl {
				above++
			} else if vals[i+j] < cl {
				below++
			}
		}
		if above == run || below == run {
			v = append(v, NelsonViolation{Rule: 2, Start: i, End: i + run - 1})
		}
	}
	return v
}

// Rule 3: Six consecutive points steadily increasing or decreasing.
func nelsonRule3(vals []float64) []NelsonViolation {
	var v []NelsonViolation
	const run = 6
	if len(vals) < run {
		return nil
	}
	for i := 0; i <= len(vals)-run; i++ {
		inc := true
		dec := true
		for j := 1; j < run; j++ {
			if vals[i+j] <= vals[i+j-1] {
				inc = false
			}
			if vals[i+j] >= vals[i+j-1] {
				dec = false
			}
		}
		if inc || dec {
			v = append(v, NelsonViolation{Rule: 3, Start: i, End: i + run - 1})
		}
	}
	return v
}

// Rule 4: Fourteen consecutive points alternating up and down.
func nelsonRule4(vals []float64) []NelsonViolation {
	var v []NelsonViolation
	const run = 14
	if len(vals) < run {
		return nil
	}
	for i := 0; i <= len(vals)-run; i++ {
		alt := true
		for j := 1; j < run-1; j++ {
			mid := vals[i+j]
			prev := vals[i+j-1]
			next := vals[i+j+1]
			if !((mid > prev && mid > next) || (mid < prev && mid < next)) {
				alt = false
				break
			}
		}
		if alt {
			v = append(v, NelsonViolation{Rule: 4, Start: i, End: i + run - 1})
		}
	}
	return v
}

// Rule 5: Two of three consecutive points beyond 2 sigma on same side.
func nelsonRule5(vals []float64, cl, sigma float64) []NelsonViolation {
	var v []NelsonViolation
	if len(vals) < 3 {
		return nil
	}
	for i := 0; i <= len(vals)-3; i++ {
		above := 0
		below := 0
		for j := 0; j < 3; j++ {
			if vals[i+j] > cl+2*sigma {
				above++
			}
			if vals[i+j] < cl-2*sigma {
				below++
			}
		}
		if above >= 2 || below >= 2 {
			v = append(v, NelsonViolation{Rule: 5, Start: i, End: i + 2})
		}
	}
	return v
}

// Rule 6: Four of five consecutive points beyond 1 sigma on same side.
func nelsonRule6(vals []float64, cl, sigma float64) []NelsonViolation {
	var v []NelsonViolation
	const run = 5
	if len(vals) < run {
		return nil
	}
	for i := 0; i <= len(vals)-run; i++ {
		above := 0
		below := 0
		for j := 0; j < run; j++ {
			if vals[i+j] > cl+sigma {
				above++
			}
			if vals[i+j] < cl-sigma {
				below++
			}
		}
		if above >= 4 || below >= 4 {
			v = append(v, NelsonViolation{Rule: 6, Start: i, End: i + run - 1})
		}
	}
	return v
}

// Rule 7: Fifteen consecutive points within 1 sigma of CL (too little variability).
func nelsonRule7(vals []float64, cl, sigma float64) []NelsonViolation {
	var v []NelsonViolation
	const run = 15
	if len(vals) < run {
		return nil
	}
	for i := 0; i <= len(vals)-run; i++ {
		all := true
		for j := 0; j < run; j++ {
			if vals[i+j] > cl+sigma || vals[i+j] < cl-sigma {
				all = false
				break
			}
		}
		if all {
			v = append(v, NelsonViolation{Rule: 7, Start: i, End: i + run - 1})
		}
	}
	return v
}

// Rule 8: Eight consecutive points beyond 1 sigma on either side (mixture).
func nelsonRule8(vals []float64, cl, sigma float64) []NelsonViolation {
	var v []NelsonViolation
	const run = 8
	if len(vals) < run {
		return nil
	}
	for i := 0; i <= len(vals)-run; i++ {
		all := true
		for j := 0; j < run; j++ {
			if vals[i+j] >= cl-sigma && vals[i+j] <= cl+sigma {
				all = false
				break
			}
		}
		if all {
			v = append(v, NelsonViolation{Rule: 8, Start: i, End: i + run - 1})
		}
	}
	return v
}
