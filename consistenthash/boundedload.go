package consistenthash

import (
	"math"

	"github.com/zjbztianya/go-misc/metrics"
)

// Bounded loads consistent hash algorithm
// paper:https://arxiv.org/pdf/1608.01350.pdf
// reference:https://medium.com/vimeo-engineering-blog/improving-load-balancing-with-a-new-consistent-hashing-algorithm-9f1bd75709ed
type BoundedLoadHashRing struct {
	*HashRing
	factor    float64 //values between 1.25 and 2 are good for practical use
	loads     map[string]metrics.Gauge
	totalLoad metrics.Gauge
}

func NewBoundedLoadHashRing(factor float64, replicas int, opts ...HashRingOption) *BoundedLoadHashRing {
	return &BoundedLoadHashRing{
		HashRing:  NewHashRing(replicas, opts...),
		factor:    factor,
		loads:     make(map[string]metrics.Gauge),
		totalLoad: metrics.NewGauge(),
	}
}

func (b *BoundedLoadHashRing) AddNode(key string) {
	if _, ok := b.loads[key]; ok {
		return
	}
	b.loads[key] = metrics.NewGauge()
	b.HashRing.AddNode(key, b.replicas) // TODO:add weight
}

func (b *BoundedLoadHashRing) RemoveNode(key string) {
	delete(b.loads, key)
	b.HashRing.RemoveNode(key)
}

func (b *BoundedLoadHashRing) maxLoad() float64 {
	return math.Ceil(b.totalLoad.Value() * b.factor / float64(len(b.loads)))
}

func (b *BoundedLoadHashRing) Get(key string) (string, error) {
	idx, err := b.search(key)
	if err != nil {
		return "", nil
	}

	pos := idx
	maxLoad := b.maxLoad()
	for {
		if ld, ok := b.loads[b.nodes[pos].key]; ok && ld.Value() < maxLoad {
			break
		}
		pos++
		if pos == len(b.nodes) {
			pos = 0
		}
		if pos == idx {
			break
		}
	}
	return b.nodes[pos].key, nil
}

func (b *BoundedLoadHashRing) Inc(key string) {
	if ld, ok := b.loads[key]; ok {
		ld.Inc()
		b.totalLoad.Inc()
	}
}

func (b *BoundedLoadHashRing) Dec(key string) {
	if ld, ok := b.loads[key]; ok {
		ld.Dec()
		b.totalLoad.Dec()
	}
}
