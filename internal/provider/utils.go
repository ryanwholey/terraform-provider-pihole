package provider

import "hash/fnv"

// hash hashes string contents
func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))

	return h.Sum32()
}
