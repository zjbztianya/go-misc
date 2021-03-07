package breaker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func getSreBreaker() *SreBreaker {
	config := Config{
		k:          2,
		threshold:  10,
		numBuckets: 10,
		interval:   100 * time.Millisecond,
	}

	return NewSreBreaker(&config)
}

func markSuccess(b *SreBreaker, count int) {
	for i := 0; i < count; i++ {
		b.MarkSuccess()
	}
}

func markFailed(b *SreBreaker, count int) {
	for i := 0; i < count; i++ {
		b.MarkFailed()
	}
}

func TestSreBreakerAllow(t *testing.T) {
	t.Run("total requests less than threshold", func(t *testing.T) {
		b := getSreBreaker()
		markSuccess(b, 9)
		assert.Nil(t, b.Allow())
	})

	t.Run("total requests large than threshold, less than K*accepts", func(t *testing.T) {
		b := getSreBreaker()
		markSuccess(b, 10)
		markFailed(b, 5)
		assert.Nil(t, b.Allow())
	})

	t.Run("total requests large than K*accepts", func(t *testing.T) {
		b := getSreBreaker()
		markFailed(b, 2000000)
		assert.NotNil(t, b.Allow())
	})

	t.Run("breaker close state", func(t *testing.T) {
		b := getSreBreaker()
		markSuccess(b, 15)
		assert.Nil(t, b.Allow())
	})
}

func BenchmarkSreBreakerAllow(b *testing.B) {
	breaker := getSreBreaker()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		breaker.Allow()
		if i%3 == 0 {
			breaker.MarkSuccess()
		} else {
			breaker.MarkFailed()
		}
	}
}
