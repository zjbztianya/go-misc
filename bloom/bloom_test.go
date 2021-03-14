package bloom

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyFilter(t *testing.T) {
	filter := NewFilter(10)
	assert.False(t, filter.Search("bloom"))
	assert.False(t, filter.Search("filter"))
}

func TestFilterSearchSmall(t *testing.T) {
	keys := []string{"bloom", "filter"}
	filter := NewFilter(10, keys...)
	assert.NotNil(t, filter)
	assert.True(t, filter.Search("bloom"))
	assert.True(t, filter.Search("filter"))
	assert.False(t, filter.Search("hello"))
	assert.False(t, filter.Search("world"))
}

func nextLen(l int) int {
	if l < 10 {
		l++
	} else if l < 100 {
		l += 10
	} else if l < 1000 {
		l += 100
	} else {
		l += 1000
	}
	return l
}

func falsePositiveRate(filter *Filter) float64 {
	var res int
	for i := 0; i < 10000; i++ {
		if filter.Search(strconv.Itoa(i + 1000000000)) {
			res++
		}
	}
	return float64(res) / 10000.0
}

func TestFilterVaryingLengths(t *testing.T) {
	var mediocreFilters, goodFilters int

	for l := 1; l <= 10000; l = nextLen(l) {
		keys := make([]string, l)
		for i := 0; i < l; i++ {
			keys[i] = strconv.Itoa(i)
		}
		filter := NewFilter(10, keys...)
		assert.LessOrEqual(t, len(filter.bitSet), (l*10+63)/64)

		for i := 0; i < l; i++ {
			assert.True(t, filter.Search(strconv.Itoa(i)))
		}

		rate := falsePositiveRate(filter)
		assert.LessOrEqual(t, rate, 0.025)
		if rate > 0.0125 {
			mediocreFilters++
		} else {
			goodFilters++
		}
	}
	assert.LessOrEqual(t, mediocreFilters, goodFilters/5)
}
