package slidingwindow

import (
	"math"
	"sync"
	"time"
)

type WindowOption func(*SlidingWindow)

type SlidingWindow struct {
	mu       sync.RWMutex
	win      *window
	interval time.Duration
	lastTime time.Time
}

func NewSlidingWindow(size int, interval time.Duration, options ...WindowOption) *SlidingWindow {
	if size <= 0 {
		panic("rolling window size must greater than 0")
	}
	w := &SlidingWindow{
		mu:       sync.RWMutex{},
		win:      newWindow(size),
		interval: interval,
		lastTime: time.Now(),
	}

	for _, opt := range options {
		opt(w)
	}
	return w
}

func (r *SlidingWindow) timeSpan() int {
	return int(time.Since(r.lastTime) / r.interval)
}

func (r *SlidingWindow) Increment() {
	r.Add(1)
}

func (r *SlidingWindow) Add(v float64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	span := r.timeSpan()
	if span > 0 {
		r.lastTime = r.lastTime.Add(time.Duration(int(r.interval) * span))
		r.win.adjustOffset(span)
	}
	r.win.add(v)
}

func (r *SlidingWindow) Reduce(fn func(b *Bucket)) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	span := r.timeSpan()
	if count := r.win.size - span; count > 0 {
		for i := 0; i < count; i++ {
			offset := (r.win.offset + span + i + 1) % r.win.size
			fn(r.win.buckets[offset])
		}
	}
}

func (r *SlidingWindow) Sum() float64 {
	var v float64
	r.Reduce(func(b *Bucket) {
		v += b.Sum
	})
	return v
}

func (r *SlidingWindow) Max() float64 {
	v := math.SmallestNonzeroFloat64
	r.Reduce(func(b *Bucket) {
		if v < b.Sum {
			v = b.Sum
		}
	})
	return v
}

func (r *SlidingWindow) Min() float64 {
	v := math.MaxFloat64
	r.Reduce(func(b *Bucket) {
		if v > b.Sum {
			v = b.Sum
		}
	})
	return v
}

func (r *SlidingWindow) Avg() float64 {
	var v float64
	r.Reduce(func(b *Bucket) {
		v += b.Sum
	})
	return v / float64(r.win.size)
}

func (r *SlidingWindow) Count() float64 {
	var v int64
	r.Reduce(func(b *Bucket) {
		v += b.Count
	})
	return float64(v)
}

// window is acting as a circular array(ring buffer)
type window struct {
	buckets []*Bucket
	size    int //bucket numbers
	offset  int
}

func newWindow(size int) *window {
	w := &window{size: size, buckets: make([]*Bucket, size)}
	for i := 0; i < size; i++ {
		w.buckets[i] = new(Bucket)
	}
	return w
}

func (w *window) adjustOffset(span int) {
	if span > w.size {
		span = w.size
	}

	for i := 0; i < span; i++ {
		w.buckets[(w.offset+1+i)%w.size].reset()
	}
	w.offset = (w.offset + span) % w.size
}

func (w *window) add(v float64) {
	w.buckets[w.offset].add(v)
}

type Bucket struct {
	Sum   float64
	Count int64
}

func (b *Bucket) add(v float64) {
	b.Sum += v
	b.Count++
}

func (b *Bucket) reset() {
	b.Count = 0
	b.Sum = 0
}
