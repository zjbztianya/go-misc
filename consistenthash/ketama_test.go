package consistenthash

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

var (
	nodes = []string{
		"test0.github.com",
		"test1.github.com",
		"test2.github.com",
		"test3.github.com",
		"test4.github.com",
	}
	replicas = 160
)

func BenchmarkHashRingGet(b *testing.B) {
	ring := NewHashRing(replicas)

	for _, node := range nodes {
		ring.AddNode(node, replicas)
	}

	for i := 0; i < b.N; i++ {
		key := "test_key_" + strconv.Itoa(i)
		ring.Get(key)
	}
}

func TestNewHashRing(t *testing.T) {
	ring := NewHashRing(minReplicas + 1)
	assert.Equal(t, minReplicas+1, ring.replicas)
	assert.NotNil(t, ring.hashFunc)

	ring = NewHashRing(3)
	assert.Equal(t, minReplicas, ring.replicas)
}

func TestHashRingAddNode(t *testing.T) {
	ring := NewHashRing(4)
	ring.AddNode(nodes[0], replicas)
	assert.Len(t, ring.nodes, replicas)

	for i := 0; i < replicas; i++ {
		key := strconv.Itoa(i) + nodes[0]
		host, ok := ring.Get(key)
		assert.True(t, ok)
		assert.Equal(t, nodes[0], host)
	}

	ring.AddNode(nodes[1], replicas)
	assert.Len(t, ring.nodes, 2*replicas)
}

func TestHashRingGet(t *testing.T) {
	ring := NewHashRing(4)
	for i := 0; i < len(nodes); i++ {
		ring.AddNode(nodes[i], 160)
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

func TestHashRingRemoveNode(t *testing.T) {
	ring := NewHashRing(4)
	ring.AddNode(nodes[0], replicas)
	assert.Len(t, ring.nodes, replicas)
	ring.AddNode(nodes[1], replicas)
	ring.RemoveNode(nodes[0])
	for i := 0; i < replicas; i++ {
		key := strconv.Itoa(i) + nodes[0]
		host, ok := ring.Get(key)
		assert.True(t, ok)
		assert.NotEqual(t, nodes[0], host)
	}
	assert.Len(t, ring.nodes, replicas)
}
