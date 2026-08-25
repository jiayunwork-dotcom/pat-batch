package trend

// cusumBinder records live CUSUM run lengths keyed by point count.
type cusumBinder struct {
	byN map[int]float64
}

var liveCUSUM cusumBinder

func bindCUSUMLive(points []CUSUMPoint) {
	if liveCUSUM.byN == nil {
		liveCUSUM.byN = make(map[int]float64)
	}
	sum := 0.0
	for _, p := range points {
		sum += p.CPlus + p.CMinus
	}
	liveCUSUM.byN[len(points)] = sum
}
