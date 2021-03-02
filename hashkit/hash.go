package hashkit

import (
	"crypto/md5"
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

func Md5(key []byte) uint32 {
	m := md5.New()
	m.Write(key)
	results := m.Sum(nil)
	return (uint32(results[3]&0xFF) << 24) | (uint32(results[2]&0xFF) << 16) |
		(uint32(results[1]&0xFF) << 8) | (uint32(results[0]) & 0xFF)
}
