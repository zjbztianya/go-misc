package consistenthash

import (
	"math"
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newRing() *BoundedLoadHashRing {
	return NewBoundedLoadHashRing(1.25, 160)
}

func TestNewBoundedLoadHashRing(t *testing.T) {
	assert.NotNil(t, newRing())
}

func TestBoundedLoadHashRingGet(t *testing.T) {
	ring := newRing()
	for _, node := range nodes {
		ring.AddNode(node)
		delta := float64(rand.Intn(100))
		ring.loads[node].Add(delta)
		ring.totalLoad.Add(delta)
	}

	m := make(map[string]int)
	for i := 0; i < 1e6; i++ {
		key := "test_key_" + strconv.Itoa(i)
		host, _ := ring.Get(key)
		m[host]++
	}

	for i := 0; i < len(nodes); i++ {
		t.Log(nodes[i], m[nodes[i]])
	}
}

func TestBoundedLoadHashRingMaxLoad(t *testing.T) {
	ring := newRing()
	host := "test0.github.com"
	ring.AddNode(host)
	ring.Inc(host)
	assert.Equal(t, math.Ceil(ring.factor), ring.maxLoad())
}
