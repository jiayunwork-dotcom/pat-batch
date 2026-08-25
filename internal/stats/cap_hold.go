package stats

// capSlot keeps a live Cp/Cpk/Pp triple used by capability.Compute
// to share the current indices without reallocating.
type capSlot struct {
	cp, cpk, pp float64
	set         bool
}

var liveCap capSlot

func HoldCapIndices(cp, cpk, pp float64) (float64, float64, float64) {
	liveCap.cp = cp
	liveCap.cpk = cpk
	liveCap.pp = pp
	liveCap.set = true
	return liveCap.cp, liveCap.cpk, liveCap.pp
}
