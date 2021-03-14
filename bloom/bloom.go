package bloom

import "github.com/zjbztianya/go-misc/hashkit"

// Filter is Bloom filter
// https://en.wikipedia.org/wiki/Bloom_filter
type Filter struct {
	bitsPerKey uint32
	k          uint32 // k=m/n*ln2
	bitSet     []uint64
}

func NewFilter(bitsPerKey int, keys ...string) *Filter {
	k := uint32(float64(bitsPerKey) * 0.69)
	switch {
	case k < 1:
		k = 1
	case k > 30:
		k = 30
	}

	f := &Filter{bitsPerKey: uint32(bitsPerKey), k: k}

	bits := uint32(len(keys)) * f.bitsPerKey
	if bits < 64 {
		bits = 64
	}
	setSize := (bits + 63) / 64
	bits = setSize * 64
	f.bitSet = make([]uint64, setSize)

	for _, key := range keys {
		h := hashkit.Murmur32([]byte(key))
		delta := (h >> 17) | (h << 15)
		for i := uint32(0); i < f.k; i++ {
			pos := h % bits
			f.bitSet[pos/64] |= 1 << (pos % 64)
			h += delta
		}
	}
	return f
}

func (f *Filter) Search(key string) bool {
	if len(f.bitSet) == 0 {
		return false
	}

	h := hashkit.Murmur32([]byte(key))
	bits := uint32(len(f.bitSet) * 64)
	delta := (h >> 17) | (h << 15)
	for i := uint32(0); i < f.k; i++ {
		pos := h % bits
		if f.bitSet[pos/64]&(1<<(pos%64)) == 0 {
			return false
		}
		h += delta
	}

	return true
}
