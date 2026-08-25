package stats

// imrSlot keeps a live I/MR center-line pair used by ComputeIMR to
// share the current limits without reallocating.
type imrSlot struct {
	i, mr float64
	set   bool
}

var liveIMR imrSlot

func HoldIMRLimits(icl, mrcl float64) (float64, float64) {
	liveIMR.i = icl
	liveIMR.mr = mrcl
	liveIMR.set = true
	return liveIMR.i, liveIMR.mr
}
