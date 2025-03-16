package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash maps bytes to uint32
type Hash func(data []byte) uint32

// Map contains all hashed keys
type ConsistentHash struct {
	hash     Hash
	replicas int
	keys     []int // Sorted
	hashMap  map[int]string
}

// NewConsistentHash creates a Map instance
func NewConsistentHash(replicas int, fn Hash) *ConsistentHash {
	m := &ConsistentHash{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add adds some keys to the hash
func (ch *ConsistentHash) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < ch.replicas; i++ {
			hash := int(ch.hash([]byte(strconv.Itoa(i) + key)))
			ch.keys = append(ch.keys, hash)
			ch.hashMap[hash] = key
		}
	}
	sort.Ints(ch.keys)
}

// Get gets the closest item in the hash to the provided key
func (ch *ConsistentHash) Get(key string) string {
	if len(ch.keys) == 0 {
		return ""
	}

	hash := int(ch.hash([]byte(key)))

	// Binary search for appropriate replica
	idx := sort.Search(len(ch.keys), func(i int) bool {
		return ch.keys[i] >= hash
	})

	// If we have reached the end of the ring, return the first node
	if idx == len(ch.keys) {
		idx = 0
	}

	return ch.hashMap[ch.keys[idx]]
}

// Remove removes a key from the hash
func (ch *ConsistentHash) Remove(key string) {
	for i := 0; i < ch.replicas; i++ {
		hash := int(ch.hash([]byte(strconv.Itoa(i) + key)))
		idx := sort.SearchInts(ch.keys, hash)
		if idx < len(ch.keys) && ch.keys[idx] == hash {
			ch.keys = append(ch.keys[:idx], ch.keys[idx+1:]...)
		}
		delete(ch.hashMap, hash)
	}
}
