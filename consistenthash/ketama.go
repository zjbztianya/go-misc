package consistenthash

import (
	"errors"
	"sort"
	"strconv"

	"github.com/zjbztianya/go-misc/hashkit"
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

func (h *HashRing) search(key string) (int, error) {
	if len(h.nodes) == 0 {
		return -1, errors.New("empty ring")
	}

	hash := h.hashFunc([]byte(key))
	idx := sort.Search(len(h.nodes), func(i int) bool {
		return h.nodes[i].hash >= hash
	}) % len(h.nodes)

	return idx, nil
}

func (h *HashRing) Get(key string) (string, error) {
	idx, err := h.search(key)
	if err != nil {
		return "", err
	}
	return h.nodes[idx].key, nil
}
