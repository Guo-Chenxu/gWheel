package cache

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var kinds = []string{"lru", "lfu", "fifo"}

func TestCacheCreateFail(t *testing.T) {
	_, err := CacheFactory("unknown", 3)
	assert.Error(t, err)

	for _, kind := range kinds {
		_, err = CacheFactory(kind, -3)
		assert.Error(t, err)
	}
}

func TestCacheSetAndGet(t *testing.T) {
	for _, kind := range kinds {
		cache, _ := CacheFactory(kind, 3)
		t.Log("\n", kind)

		cache.Set("key1", "value1")
		cache.Set("key2", "value2")
		cache.Set("key3", "value3")

		t.Log(cache.Get("key1"))
		t.Log(cache.Get("key1"))
		t.Log(cache.Get("key1"))

		t.Log(cache.Get("key2"))
		t.Log(cache.Get("key2"))
		t.Log(cache.Get("key3"))

		cache.Set("key4", "value4")

		t.Log(cache.Get("key4"))
		t.Log(cache.Get("key3"))
		t.Log(cache.Get("key1"))

		cache.Set("key4", "value44")
		t.Log(cache.Get("key4"))
		t.Log(cache.Get("key3"))
		t.Log(cache.Get("key2"))

		cache.Set("key5", "value5")
		t.Log(cache.Get("key5"))
		t.Log(cache.Get("key2"))
		t.Log(cache.Get("key4"))
	}
}

func TestCacheSetValidKey(t *testing.T) {
	for _, kind := range kinds {
		cache, _ := CacheFactory(kind, 3)
		cache.Set("", "")
	}
}

func TestCacheDelete(t *testing.T) {
	for _, kind := range kinds {
		cache, _ := CacheFactory(kind, 3)
		cache.Set("key1", "value1")
		cache.Set("key2", "value2")
		cache.Set("key3", "value3")

		cache.Delete("key1")
		cache.Delete("key2")
		cache.Delete("key3")
		cache.Delete("key4")
		
		t.Log(cache.Get("key1"))

		cache.Set("key4", "value4")
		cache.Set("key5", "value5")

		cache.Clear()
		t.Log(cache.Get("key4"))
	}
}
