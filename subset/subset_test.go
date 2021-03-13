package subset

import (
	"strconv"
	"testing"
)

func TestSubset(t *testing.T) {
	n := 300
	subsetSize := 10
	backends := make([]string, n)
	for i := 0; i < n; i++ {
		backends[i] = strconv.Itoa(i)
	}

	m := make(map[string]int)
	for i := 0; i < n; i++ {
		tmp := make([]string, n)
		copy(tmp, backends)
		bs := Subset(tmp, i, subsetSize)
		for _, id := range bs {
			m[id]++
		}
	}
	t.Log(m)
}
