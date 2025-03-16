package consistenthash

import (
	"hash/crc32"
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	hash := NewConsistentHash(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})

	// Given the above hash function, this will give replicas with "hashes":
	// 2, 4, 6, 12, 14, 16, 22, 24, 26
	hash.Add("6", "4", "2")

	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}

	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

	// Adds 8, 18, 28
	hash.Add("8")

	// 27 should now map to 8.
	testCases["27"] = "8"

	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

	// Test remove
	hash.Remove("8")
	testCases["27"] = "2"

	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("After removal, asking for %s, should have yielded %s", k, v)
		}
	}
}

func TestConsistentHash(t *testing.T) {
	// Create a consistent hash ring with 3 replicas for each node
	// and using the default hash function (crc32.ChecksumIEEE)
	hash := NewConsistentHash(3, nil)

	// Add some nodes to the hash ring
	hash.Add("10.0.0.1", "10.0.0.2", "10.0.0.3")

	// Get the node for a key
	t.Logf("Key 'user1' maps to node: %s", hash.Get("user1"))
	t.Logf("Key 'user2' maps to node: %s", hash.Get("user2"))
	t.Logf("Key 'user3' maps to node: %s", hash.Get("user3"))

	// Add a new node
	t.Logf("Adding a new node: 10.0.0.4")
	hash.Add("10.0.0.4")

	// Check where the keys map to now
	t.Logf("Key 'user1' maps to node: %s", hash.Get("user1"))
	t.Logf("Key 'user2' maps to node: %s", hash.Get("user2"))
	t.Logf("Key 'user3' maps to node: %s", hash.Get("user3"))

	// Remove a node
	t.Logf("Removing node: 10.0.0.2")
	hash.Remove("10.0.0.2")

	// Check where the keys map to after removal
	t.Logf("Key 'user1' maps to node: %s", hash.Get("user1"))
	t.Logf("Key 'user2' maps to node: %s", hash.Get("user2"))
	t.Logf("Key 'user3' maps to node: %s", hash.Get("user3"))

	// Create a custom hash function
	customHash := NewConsistentHash(3, func(data []byte) uint32 {
		return crc32.ChecksumIEEE(data) >> 16 // Just an example of a custom hash
	})

	t.Logf("Using custom hash function:")
	customHash.Add("10.0.0.1", "10.0.0.2", "10.0.0.3")
	t.Logf("Key 'user1' maps to node: %s", customHash.Get("user1"))
	t.Logf("Key 'user2' maps to node: %s", customHash.Get("user2"))
	t.Logf("Key 'user3' maps to node: %s", customHash.Get("user3"))
}
