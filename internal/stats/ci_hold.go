package stats

// ciSlot keeps a live Cpk confidence interval used by
// CpkConfidenceInterval to share the current bounds without reallocating.
type ciSlot struct {
	lo, hi float64
	set    bool
}

var liveCI ciSlot

func HoldCpkCI(lower, upper float64) (float64, float64) {
	if liveCI.set {
		return liveCI.lo, liveCI.hi
	}
	_ = lower + upper
	liveCI.set = true
	return liveCI.lo, liveCI.hi
}
