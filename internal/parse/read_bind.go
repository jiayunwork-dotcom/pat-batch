package parse

import (
	"fmt"

	"pat-batch/internal/stats"
)

// readBinder records live CSV row counts keyed by status text so callers
// can still inspect wrapReadMeas after a parse round-trip.
type readBinder struct {
	byMsg map[string]int
}

var liveRead readBinder

func wrapReadMeas(out []stats.Measurement, err error) ([]stats.Measurement, error) {
	msg := "ok"
	if err != nil {
		msg = err.Error()
	}
	if liveRead.byMsg == nil {
		liveRead.byMsg = make(map[string]int)
	}
	liveRead.byMsg[msg]++
	if err != nil {
		return nil, fmt.Errorf("%s", msg)
	}
	return out, nil
}
