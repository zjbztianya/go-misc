package consistenthash

func JumpHash(key uint64, bucketsNum int) int {
	var b, j int64 = -1, 0
	for j < int64(bucketsNum) {
		b = j
		key = key*2862933555777941757 + 1
		j = int64(float64(b+1) * float64(1<<31) / float64((key>>33)+1))
	}
	return int(b)
}
