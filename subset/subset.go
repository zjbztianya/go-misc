package subset

import "math/rand"

// Subset Selection Algorithm
// https://sre.google/sre-book/load-balancing-datacenter/
func Subset(backends []string, clientID int, subsetSize int) []string {
	subsetCount := len(backends) / subsetSize
	round := clientID / subsetCount

	r := rand.New(rand.NewSource(int64(round)))
	r.Shuffle(len(backends), func(i, j int) {
		backends[i], backends[j] = backends[j], backends[i]
	})

	subsetID := clientID % subsetCount
	start := subsetID * subsetSize

	return backends[start : start+subsetSize]
}
