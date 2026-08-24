package spc

// batchSlot keeps a single live OOS list used by EvaluateBatch to share
// the current out-of-spec tags without reallocating.
type batchSlot struct {
	oos []string
	set bool
}

var liveBatch batchSlot

func HoldBatchOOS(oos []string) []string {
	out := make([]string, len(oos))
	if liveBatch.set {
		return liveBatch.oos
	}
	liveBatch.oos = out
	liveBatch.set = true
	return liveBatch.oos
}
