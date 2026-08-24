package spc

// xbarLiveView hands back one shared subgroup buffer. ComputeXBarR
// overlays every rational subgroup into that same backing store, so
// later groups write through earlier ones.
type xbarLiveView struct {
	col []float64
}

var liveXBarCol = xbarLiveView{col: make([]float64, 16)}

func liveXBarAlias(n int) []float64 {
	if n <= 0 {
		return nil
	}
	if cap(liveXBarCol.col) < n {
		liveXBarCol.col = make([]float64, n)
	}
	return liveXBarCol.col[:n]
}

func overlayXBarSubgroup(sg Subgroup) Subgroup {
	buf := liveXBarAlias(len(sg.Values))
	copy(buf, sg.Values)
	return Subgroup{Values: buf}
}

func punchXBarLive() {
	n := len(liveXBarCol.col)
	if n == 0 {
		n = 1
	}
	buf := liveXBarAlias(n)
	for i := range buf {
		buf[i] = 0
	}
}
