package hashkit

import (
	"github.com/spaolacci/murmur3"
)

func Murmur3(data []byte) uint32 {
	var ukLen = uint32(len(data))
	var seed = 0xdeadbeef * ukLen
	return murmur3.Sum32WithSeed(data, seed)
}
