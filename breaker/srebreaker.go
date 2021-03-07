package breaker

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/zjbztianya/go-misc/metrics"
)

type Config struct {
	k          float64
	threshold  int
	numBuckets int
	interval   time.Duration
}

// SreBreaker is a adaptive throttling technique mentioned in google SRE book
// https://sre.google/sre-book/handling-overload/#eq2101
type SreBreaker struct {
	mu        sync.Mutex // for rand concurrent thread safe
	k         float64
	threshold int // avoid low frequency requests to trigger throttling
	stats     metrics.SlidingCounter
	rand      *rand.Rand
}

func NewSreBreaker(c *Config) *SreBreaker {
	stats := metrics.NewSlidingCounter(c.numBuckets, c.interval)
	return &SreBreaker{
		k:         c.k,
		threshold: c.threshold,
		stats:     stats,
		rand:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *SreBreaker) summary() (float64, float64) {
	var total int64
	var success float64
	s.stats.Reduce(func(b *metrics.Bucket) {
		total += b.Count
		success += b.Sum
	})
	return float64(total), success
}

func (s *SreBreaker) getProb() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.rand.Float64()
}

func (s *SreBreaker) Allow() error {
	total, accepts := s.summary()
	accepts *= s.k
	if total < float64(s.threshold) || total <= accepts {
		return nil
	}
	dropProb := (total - accepts) / (total + 1)
	if s.getProb() <= dropProb {
		return errors.New("customer is out of quota")
	}
	return nil
}

func (s *SreBreaker) MarkSuccess() {
	s.stats.Add(1)
}

func (s *SreBreaker) MarkFailed() {
	s.stats.Add(0)
}
