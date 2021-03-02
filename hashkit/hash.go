package hashkit

import (
	"github.com/spaolacci/murmur3"
	"hash/fnv"
)

type HashFunc func([]byte) uint32

func Murmur3(data []byte) uint32 {
	var ukLen = uint32(len(data))
	var seed = 0xdeadbeef * ukLen
	return murmur3.Sum32WithSeed(data, seed)
}

func Fnv32(data []byte) uint32 {
	f := fnv.New32()
	f.Write(data)
	return f.Sum32()
}
