package backoff

import (
	"math/rand"
	"time"
)

type Config struct {
	BaseDelay  time.Duration
	MaxDelay   time.Duration
	Multiplier float64
	Jitter     float64
}

// Exponential implements exponential backoff algorithm as defined in
// https://aws.amazon.com/cn/blogs/architecture/exponential-backoff-and-jitter/
type Exponential struct {
	cfg *Config
}

func (e *Exponential) BackOff(retries int) time.Duration {
	if retries == 0 {
		return e.cfg.BaseDelay
	}

	backoff, maxDelay := float64(e.cfg.BaseDelay), float64(e.cfg.MaxDelay)
	//for loop performance better than math.Pow
	for retries > 0 && backoff < maxDelay {
		backoff *= e.cfg.Multiplier
		retries--
	}
	if backoff > maxDelay {
		backoff = maxDelay
	}

	// backoff range:[backoff*(1-jitter),backoff*(1+jitter))
	backoff = backoff * (1 + e.cfg.Jitter*(rand.Float64()*2-1))
	if backoff < 0 {
		backoff = 0
	}
	return time.Duration(backoff)
}
