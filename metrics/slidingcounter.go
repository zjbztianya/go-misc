package metrics

import (
	"math"
	"time"
)

var _ SlidingCounter = (*slidingCounter)(nil)

type SlidingCounter interface {
	Counter
	Reduce(func(b *Bucket))
	Min() float64
	Max() float64
	Count() float64
	Avg() float64
	Sum() float64
}

type slidingCounter struct {
	win *SlidingWindow
}

func NewSlidingCounter(size int, interval time.Duration) SlidingCounter {
	return &slidingCounter{win: NewSlidingWindow(size, interval)}
}

func (s *slidingCounter) Inc() {
	s.win.Inc()
}

func (s *slidingCounter) Add(delta float64) {
	s.win.Add(delta)
}

func (s *slidingCounter) Reduce(fn func(b *Bucket)) {
	s.win.Reduce(fn)
}

func (s *slidingCounter) Min() float64 {
	v := math.MaxFloat64
	s.Reduce(func(b *Bucket) {
		if v > b.Sum {
			v = b.Sum
		}
	})
	return v
}

func (s *slidingCounter) Max() float64 {
	v := math.SmallestNonzeroFloat64
	s.Reduce(func(b *Bucket) {
		if v < b.Sum {
			v = b.Sum
		}
	})
	return v
}

func (s *slidingCounter) Count() float64 {
	var v int64
	s.Reduce(func(b *Bucket) {
		v += b.Count
	})
	return float64(v)
}

func (s *slidingCounter) Avg() float64 {
	var v float64
	s.Reduce(func(b *Bucket) {
		v += b.Sum
	})
	return v / float64(s.win.Size())
}

func (s *slidingCounter) Sum() float64 {
	var v float64
	s.Reduce(func(b *Bucket) {
		v += b.Sum
	})
	return v
}
