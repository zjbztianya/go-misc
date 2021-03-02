package consistenthash

import (
	"github.com/zjbztianya/go-misc/hashkit"

	"errors"
	"math/big"
	"sort"
)

//Maglev consistent hashing algorithm paper:https://static.googleusercontent.com/media/research.google.com/zh-CN//pubs/archive/44824.pdf
type Maglev struct {
	permutation map[string][]uint32
	entry       []int
	nodes       []string
	numBuckets  uint64
	h1, h2      hashkit.HashFunc64
}

func NewMaglev(nodes []string, numBuckets uint64, h1, h2 hashkit.HashFunc64) (*Maglev, error) {
	if !big.NewInt(int64(numBuckets)).ProbablyPrime(0) {
		return nil, errors.New("lookup table size must be prime")
	}

	if len(nodes) == 0 {
		return nil, errors.New("node nums must be greater than zero")
	}

	m := &Maglev{
		permutation: make(map[string][]uint32),
		numBuckets:  numBuckets,
		nodes:       make([]string, len(nodes)),
		h1:          h1,
		h2:          h2,
	}
	copy(m.nodes, nodes)
	sort.Strings(m.nodes)

	for _, node := range m.nodes {
		m.permutation[node] = m.generatePermutation(node)
	}
	m.populate()

	return m, nil
}

//generatePermutation guarantee permutation array to be a full permutation,proof as follows:
//Suppose that permutation[] is not a full permutation of 0,1... ,M-1,
//then there exists permutation[i] which is equal to permutation[j].
//Then the following equations hold.
//1.(offset + i * skip) % m == (offset + j * skip) % m
//2.(i * skip) % m == (j * skip) % m
//3.(i - j) * skip == x * m , assuming i > j and x >= 1 (congruence modulo)
//4.(i - j ) * skip / x == m
//Since 1 <= skip < m, 1 <= (i - j) < m, and m is a prime number, Equation 4 cannot hold.
func (m *Maglev) generatePermutation(node string) []uint32 {
	offset := m.h1([]byte(node)) % m.numBuckets
	skip := m.h2([]byte(node))%(m.numBuckets-1) + 1
	permutation := make([]uint32, m.numBuckets)

	for j := uint64(0); j < m.numBuckets; j++ {
		permutation[j] = uint32((offset + j*skip) % m.numBuckets)
	}

	return permutation
}

func (m *Maglev) populate() {
	next := make([]uint32, len(m.nodes))
	m.entry = make([]int, m.numBuckets)
	for i := uint64(0); i < m.numBuckets; i++ {
		m.entry[i] = -1
	}

	var n uint64
	for {
		for i, node := range m.nodes {
			permutation := m.permutation[node]
			c := permutation[next[i]]
			for m.entry[c] >= 0 {
				next[i]++
				c = permutation[next[i]]
			}
			m.entry[c] = i
			next[i]++
			n++
			if n == m.numBuckets {
				return
			}
		}
	}
}

func (m *Maglev) AddNode(node string) error {
	idx := sort.SearchStrings(m.nodes, node)
	if idx < len(m.nodes) && m.nodes[idx] == node {
		return errors.New("node already exist")
	}

	m.nodes = append(m.nodes[:idx], append([]string{node}, m.nodes[idx:]...)...)
	m.permutation[node] = m.generatePermutation(node)
	m.populate()
	return nil
}

func (m *Maglev) RemoveNode(node string) error {
	idx := sort.SearchStrings(m.nodes, node)
	if idx >= len(m.nodes) || m.nodes[idx] != node {
		return errors.New("node not find")
	}
	m.nodes = append(m.nodes[:idx], m.nodes[idx+1:]...)
	delete(m.permutation, node)
	m.populate()
	return nil
}

func (m *Maglev) Lookup(key uint64) string {
	if len(m.nodes) == 0 {
		return ""
	}
	return m.nodes[m.entry[key%m.numBuckets]]
}
