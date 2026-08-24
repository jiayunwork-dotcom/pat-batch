package trend

import "math"

// RunsTestResult holds the result of a runs test for randomness.
type RunsTestResult struct {
	NumRuns        int
	NumPositive    int
	NumNegative    int
	ExpectedRuns   float64
	StdRuns        float64
	ZScore         float64
	SignificantAt5 bool // |Z| > 1.96
}

// RunsTest performs a runs test for randomness relative to the median (or a
// specified center). A run is a consecutive sequence of values above or below
// the center. Too few runs indicate trending; too many indicate oscillation.
func RunsTest(values []float64, center float64) RunsTestResult {
	n := len(values)
	if n < 2 {
		return RunsTestResult{}
	}

	nPos := 0
	nNeg := 0
	runs := 1
	prevAbove := values[0] >= center

	for i, v := range values {
		above := v >= center
		if above {
			nPos++
		} else {
			nNeg++
		}
		if i > 0 && above != prevAbove {
			runs++
		}
		prevAbove = above
	}

	// Expected runs and standard deviation under H0 (random).
	n1 := float64(nPos)
	n2 := float64(nNeg)
	N := n1 + n2
	expectedRuns := 0.0
	stdRuns := 0.0
	if N > 1 {
		expectedRuns = (2*n1*n2)/N + 1
		num := 2 * n1 * n2 * (2*n1*n2 - N)
		denom := N * N * (N - 1)
		if denom > 0 && num > 0 {
			stdRuns = math.Sqrt(num / denom)
		}
	}

	z := 0.0
	if stdRuns > 0 {
		z = (float64(runs) - expectedRuns) / stdRuns
	}

	return RunsTestResult{
		NumRuns:        runs,
		NumPositive:    nPos,
		NumNegative:    nNeg,
		ExpectedRuns:   expectedRuns,
		StdRuns:        stdRuns,
		ZScore:         z,
		SignificantAt5: math.Abs(z) > 1.96,
	}
}

// RunsAboveBelow classifies each value as above or below center and returns
// the run lengths. Useful for visualizing clustering patterns.
func RunsAboveBelow(values []float64, center float64) []int {
	if len(values) == 0 {
		return nil
	}
	var lengths []int
	currentLen := 1
	prevAbove := values[0] >= center

	for i := 1; i < len(values); i++ {
		above := values[i] >= center
		if above == prevAbove {
			currentLen++
		} else {
			lengths = append(lengths, currentLen)
			currentLen = 1
			prevAbove = above
		}
	}
	lengths = append(lengths, currentLen)
	return lengths
}

// LongestRun returns the length of the longest run in the series relative to center.
func LongestRun(values []float64, center float64) int {
	runs := RunsAboveBelow(values, center)
	max := 0
	for _, r := range runs {
		if r > max {
			max = r
		}
	}
	return max
}

// ConsecutiveIncreasing returns the length of the longest strictly increasing
// subsequence of consecutive values.
func ConsecutiveIncreasing(values []float64) int {
	if len(values) < 2 {
		return len(values)
	}
	max := 1
	cur := 1
	for i := 1; i < len(values); i++ {
		if values[i] > values[i-1] {
			cur++
			if cur > max {
				max = cur
			}
		} else {
			cur = 1
		}
	}
	return max
}

// ConsecutiveDecreasing returns the length of the longest strictly decreasing
// subsequence of consecutive values.
func ConsecutiveDecreasing(values []float64) int {
	if len(values) < 2 {
		return len(values)
	}
	max := 1
	cur := 1
	for i := 1; i < len(values); i++ {
		if values[i] < values[i-1] {
			cur++
			if cur > max {
				max = cur
			}
		} else {
			cur = 1
		}
	}
	return max
}
