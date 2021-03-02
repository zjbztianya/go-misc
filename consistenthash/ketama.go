package consistenthash

import (
	"github.com/zjbztianya/go-misc/hashkit"
	"sort"
	"strconv"
	"sync"
)

const (
	minReplicas = 160
)

var defaultHash = hashkit.Murmur32

type node struct {
	hash uint32
	key  string
}

type HashRing struct {
	mu       sync.RWMutex
	nodes    []node
	replicas int
	hashFunc hashkit.HashFunc32
}

type HashRingOption func(*HashRing)

func WithHashFunc(hash hashkit.HashFunc32) HashRingOption {
	return func(ring *HashRing) {
		ring.hashFunc = hash
	}
}

func NewHashRing(replicas int, opts ...HashRingOption) *HashRing {
	if replicas < minReplicas {
		replicas = minReplicas
	}
	h := &HashRing{replicas: replicas}
	for _, opt := range opts {
		opt(h)
	}
	if h.hashFunc == nil {
		h.hashFunc = defaultHash
	}
	return h
}

func (h *HashRing) AddNode(key string, replicas int) {
	if replicas > h.replicas {
		replicas = h.replicas
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for i := 0; i < replicas; i++ {
		hash := h.hashFunc([]byte(strconv.Itoa(i) + key))
		h.nodes = append(h.nodes, node{
			hash: hash,
			key:  key,
		})

	}
	sort.Slice(h.nodes, func(i, j int) bool {
		return h.nodes[i].hash < h.nodes[j].hash
	})
}

func (h *HashRing) RemoveNode(key string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	nodes := make([]node, len(h.nodes))
	var n int
	for _, node := range h.nodes {
		if node.key != key {
			nodes[n] = node
			n++
		}
	}
	h.nodes = nodes[:n]
}

func (h *HashRing) Get(key string) (string, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.nodes) == 0 {
		return "", false
	}

	hash := h.hashFunc([]byte(key))
	idx := sort.Search(len(h.nodes), func(i int) bool {
		return h.nodes[i].hash >= hash
	}) % len(h.nodes)

	return h.nodes[idx].key, true
}
