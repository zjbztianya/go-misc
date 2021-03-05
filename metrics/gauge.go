package metrics

import (
	"math"
	"sync/atomic"
)

var _ Gauge = (*gauge)(nil)

type gauge struct {
	valBits uint64
}

func NewGauge() Gauge {
	return &gauge{}
}

func (g *gauge) Set(val float64) {
	atomic.StoreUint64(&g.valBits, math.Float64bits(val))
}

func (g *gauge) Inc() {
	g.Add(1)
}

func (g *gauge) Dec() {
	g.Add(-1)
}

func (g *gauge) Add(delta float64) {
	for {
		oldBits := atomic.LoadUint64(&g.valBits)
		newBits := math.Float64bits(delta + math.Float64frombits(oldBits))
		if atomic.CompareAndSwapUint64(&g.valBits, oldBits, newBits) {
			return
		}
	}
}

func (g *gauge) Sub(delta float64) {
	g.Add(delta * -1)
}

func (g *gauge) Value() float64 {
	return math.Float64frombits(atomic.LoadUint64(&g.valBits))
}
