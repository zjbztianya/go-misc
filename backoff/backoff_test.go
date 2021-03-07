package backoff

import (
	"testing"
	"time"
)

func getExpBackoff() *Exponential {
	cfg := &Config{
		BaseDelay:  time.Second,
		MaxDelay:   120 * time.Second,
		Multiplier: 2,
		Jitter:     0.2,
	}

	return &Exponential{cfg}
}

func TestExponentialBackOff(t *testing.T) {
	expBackoff := getExpBackoff()
	for i := 0; i < 10; i++ {
		t.Log(expBackoff.Backoff(i).Milliseconds())
	}
}

func BenchmarkExponentialBackOff(b *testing.B) {
	expBackoff := getExpBackoff()
	for i := 0; i < b.N; i++ {
		expBackoff.Backoff(i % 10)
	}
}
