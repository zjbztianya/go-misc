package consistenthash

//Jump consistent hash also does a better job of splitting the keys evenly among the buckets,
//and of splitting the rebalancing workload among the shards.
//On the other hand, jump consistent hash does not support arbitrary server names,
//but only returns a shard number; it is thus primarily suitable for the data storage case.
//paper:https://arxiv.org/pdf/1406.2294.pdf

func JumpHash(key uint64, numBuckets int) int {
	var b, j int64 = -1, 0
	for j < int64(numBuckets) {
		b = j
		key = key*2862933555777941757 + 1
		j = int64(float64(b+1) * float64(1<<31) / float64((key>>33)+1))
	}
	return int(b)
}
