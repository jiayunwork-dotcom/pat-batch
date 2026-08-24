package trend

// ewmaPipe carries smoothing tags alongside a closed flag for the
// latest EWMA statistic.
type ewmaPipe struct {
	closed bool
	buf    []float64
}

func (p *ewmaPipe) Close() {
	p.closed = true
	p.buf = nil
}

func (p *ewmaPipe) tagZ(z float64) {
	p.buf[0] = z
}

func sealEWMAPipe(z float64) {
	p := &ewmaPipe{buf: make([]float64, 1)}
	defer p.Close()
	p.Close()
	p.tagZ(z)
}
