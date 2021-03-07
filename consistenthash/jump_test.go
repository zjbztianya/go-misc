package consistenthash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkJumpHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		JumpHash(uint64(i), 100)
	}
}

func TestJumpHash(t *testing.T) {
	assert.Equal(t, 0, JumpHash(3, 1))
}

func TestJumpHashMove(t *testing.T) {
	bucketsNum := 10
	keySize := uint64(1000000)
	buckets := make(map[int]int, bucketsNum)
	for i := uint64(0); i < keySize; i++ {
		b := JumpHash(i, bucketsNum)
		buckets[b]++
	}
	t.Log("buckets:", buckets)
	bucketsNum = 12
	for i := uint64(0); i < keySize; i++ {
		oldBucket := JumpHash(i, bucketsNum-2)
		newBucket := JumpHash(i, bucketsNum)
		if oldBucket != newBucket {
			buckets[oldBucket]--
			buckets[newBucket]++
		}
	}
	t.Log("buckets:", buckets)
}
