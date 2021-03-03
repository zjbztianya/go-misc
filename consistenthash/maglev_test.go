package consistenthash

import (
	"sort"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zjbztianya/go-misc/hashkit"
)

var (
	smallTableSize uint64 = 65537
	bigTableSize   uint64 = 655373
)

func genNodes(n int) []string {
	nodes := make([]string, n)
	for i := 0; i < n; i++ {
		nodes[i] = "test" + strconv.Itoa(i) + ".github.com"
	}
	return nodes
}

func TestNewMaglev(t *testing.T) {
	N := 5
	nodes := make([]string, N)
	for i := 0; i < N; i++ {
		nodes[i] = "test" + strconv.Itoa(i) + ".github.com"
	}
	mgv, _ := NewMaglev(nodes, 10000, hashkit.Murmur64, hashkit.Fnv64)
	assert.Nil(t, mgv)
	mgv, _ = NewMaglev(nodes, smallTableSize, hashkit.Murmur64, hashkit.Fnv64)
	assert.NotNil(t, mgv)
	mgv, _ = NewMaglev([]string{}, smallTableSize, hashkit.Murmur64, hashkit.Fnv64)
	assert.Nil(t, mgv)
}

func TestMaglevAddNode(t *testing.T) {
	nodes := genNodes(5)
	prefix := "test"
	suffix := ".github.com"
	mgv, _ := NewMaglev(nodes, smallTableSize, hashkit.Murmur64, hashkit.Fnv64)
	assert.NotNil(t, mgv)

	node := prefix + "0" + suffix
	err := mgv.AddNode(node)
	assert.NotNil(t, err)

	node = prefix + "6" + suffix
	err = mgv.AddNode(node)
	assert.Nil(t, err)
	idx := sort.SearchStrings(mgv.nodes, node)
	assert.Less(t, idx, len(mgv.nodes))
	var cnt int
	for _, c := range mgv.entry {
		if c == idx {
			cnt++
		}
	}
	assert.Greater(t, cnt, 0)
}

func TestMaglevLookupSmallTableSize(t *testing.T) {
	nodes := genNodes(100)
	mgv, _ := NewMaglev(nodes, smallTableSize, hashkit.Murmur64, hashkit.Fnv64)

	m := make(map[string]int)
	for i := 0; i < 1e6; i++ {
		key := "test_key_" + strconv.Itoa(i)
		host := mgv.Lookup(uint64(hashkit.Md5([]byte(key))))
		m[host]++
	}

	for i := 0; i < len(nodes); i++ {
		t.Log(nodes[i], m[nodes[i]])
	}
}

func TestMaglevLookupBigTableSize(t *testing.T) {
	nodes := genNodes(100)
	mgv, _ := NewMaglev(nodes, bigTableSize, hashkit.Murmur64, hashkit.Fnv64)

	m := make(map[string]int)
	for i := 0; i < 1e6; i++ {
		key := "test_key_" + strconv.Itoa(i)
		host := mgv.Lookup(uint64(hashkit.Md5([]byte(key))))
		m[host]++
	}

	for i := 0; i < len(nodes); i++ {
		t.Log(nodes[i], m[nodes[i]])
	}
}

func TestMaglevRemoveNode(t *testing.T) {
	nodes := genNodes(10)
	mgv, _ := NewMaglev(nodes, 13, hashkit.Murmur64, hashkit.Fnv64)

	err := mgv.RemoveNode("test11.github.com")
	assert.NotNil(t, err)

	err = mgv.RemoveNode(nodes[5])
	assert.Nil(t, err)

	idx := sort.SearchStrings(mgv.nodes, nodes[5])
	assert.NotEqual(t, nodes[5], mgv.nodes[idx])

	for i := 0; i < 13; i++ {
		assert.Less(t, mgv.entry[i], 9)
	}
}

func TestMaglevGeneratePermutation(t *testing.T) {
	nodes := genNodes(10)
	mgv, _ := NewMaglev(nodes, smallTableSize, hashkit.Murmur64, hashkit.Fnv64)
	for _, p := range mgv.permutation {
		m := make(map[uint32]struct{})
		for _, v := range p {
			_, ok := m[v]
			assert.False(t, ok)
			m[v] = struct{}{}
		}
	}
}

func BenchmarkMaglevLookupSmallTableSize(b *testing.B) {
	nodes := genNodes(100)
	mgv, _ := NewMaglev(nodes, smallTableSize, hashkit.Murmur64, hashkit.Fnv64)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mgv.Lookup(uint64(i))
	}
}

func BenchmarkMaglevLookupBigTableSize(b *testing.B) {
	nodes := genNodes(100)
	mgv, _ := NewMaglev(nodes, bigTableSize, hashkit.Murmur64, hashkit.Fnv64)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mgv.Lookup(uint64(i))
	}
}
