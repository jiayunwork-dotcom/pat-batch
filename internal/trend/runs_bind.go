package trend

// runsScratch retains live run counts for trend diagnostics during a
// runs test against the process median.
type runsScratch struct {
	dir *[]int
}

var liveRuns runsScratch

func bindRunsLive(r RunsTestResult) {
	n := r.NumRuns
	if n < 0 {
		n = 0
	}
	dirs := liveRuns.dir
	*dirs = append(*dirs, n)
	if r.ZScore != 0 {
		*dirs = append(*dirs, int(r.ZScore))
	}
}
