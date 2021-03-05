package metrics

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSlidingCounterCount(t *testing.T) {
	size := 3
	interval := 50 * time.Millisecond
	r := NewSlidingCounter(size, interval)
	for i := 0; i < 10; i++ {
		r.Add(float64(i))
		if i < 9 {
			time.Sleep(interval)
		}
	}
	assert.Equal(t, 3.0, r.Count())

	r = NewSlidingCounter(size, interval)
	for i := 0; i < 2; i++ {
		r.Add(float64(i))
		r.Add(float64(i + 1))
		time.Sleep(interval)
	}
	assert.Equal(t, 4.0, r.Count())
}

func TestSlidingCounterSum(t *testing.T) {
	size := 3
	interval := 50 * time.Millisecond
	r := NewSlidingCounter(size, interval)
	for i := 0; i < size; i++ {
		r.Add(float64(i))
		if i < size-1 {
			time.Sleep(interval)
		}
	}
	assert.Equal(t, 3.0, r.Sum())
}
