package metrics

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewSlidingWindow(t *testing.T) {
	assert.NotNil(t, NewSlidingWindow(5, 10))
	assert.Panics(t, func() {
		NewSlidingWindow(0, 10)
	})
}

func TestSlidingWindowAdd(t *testing.T) {
	size := 3
	interval := 50 * time.Millisecond
	r := NewSlidingWindow(size, interval)
	listBuckets := func() []float64 {
		buckets := make([]float64, 0)
		r.Reduce(func(b *Bucket) {
			buckets = append(buckets, b.Sum)
		})
		return buckets
	}
	assert.Equal(t, []float64{0, 0, 0}, listBuckets())
	r.Add(1)
	assert.Equal(t, []float64{0, 0, 1}, listBuckets())
	time.Sleep(interval)
	r.Add(2)
	r.Add(3)
	assert.Equal(t, []float64{0, 1, 5}, listBuckets())
	time.Sleep(interval)
	r.Add(4)
	r.Add(5)
	assert.Equal(t, []float64{1, 5, 9}, listBuckets())
	time.Sleep(interval)
	r.Add(6)
	r.Add(7)
	assert.Equal(t, []float64{5, 9, 13}, listBuckets())

}

func TestSlidingWindowBucketTimeBoundary(t *testing.T) {
	const size = 3
	interval := time.Millisecond * 30
	r := NewSlidingWindow(size, interval)
	listBuckets := func() []float64 {
		var buckets []float64
		r.Reduce(func(b *Bucket) {
			buckets = append(buckets, b.Sum)
		})
		return buckets
	}
	assert.Equal(t, []float64{0, 0, 0}, listBuckets())
	r.Add(1)
	assert.Equal(t, []float64{0, 0, 1}, listBuckets())
	time.Sleep(time.Millisecond * 45)
	r.Add(2)
	r.Add(3)
	assert.Equal(t, []float64{0, 1, 5}, listBuckets())
	time.Sleep(time.Millisecond * 20)
	r.Add(4)
	r.Add(5)
	r.Add(6)
	assert.Equal(t, []float64{1, 5, 15}, listBuckets())
}

func TestSlidingWindowReduce(t *testing.T) {
	size := 4
	interval := 100 * time.Millisecond
	r := NewSlidingWindow(size, interval)
	for x := 0; x < size; x++ {
		for i := 0; i <= x; i++ {
			r.Add(1)
		}
		if x < size-1 {
			time.Sleep(interval)
		}
	}
	var result float64
	r.Reduce(func(b *Bucket) {
		result += b.Sum
	})
	assert.Equal(t, 10.0, result)
}

func BenchmarkSlidingWindowInc(b *testing.B) {
	size := 3
	interval := 100 * time.Millisecond
	r := NewSlidingWindow(size, interval)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Inc()
	}
}
