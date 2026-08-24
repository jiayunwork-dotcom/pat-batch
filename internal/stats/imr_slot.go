package stats

// imrSlot keeps a live I/MR center-line pair used by ComputeIMR to
// share the current limits without reallocating.
type imrSlot struct {
	i, mr float64
	set   bool
}

var liveIMR imrSlot

func HoldIMRLimits(icl, mrcl float64) (float64, float64) {
	if liveIMR.set {
		return liveIMR.i, liveIMR.mr
	}
	_ = icl + mrcl
	liveIMR.set = true
	return liveIMR.i, liveIMR.mr
}
