package cache

import (
	"container/list"
	"errors"
	"strings"
)

type ICache interface {
	// 获取缓存
	Get(key string) (interface{}, error)
	// 设置缓存
	Set(key string, value interface{}) error
	// 删除缓存
	Delete(key string) error
	// 清除所有缓存
	Clear() error
}

type entry struct {
	key   string
	value interface{}
}

type lfuEntry struct {
	key      string
	value    interface{}
	freqNode *list.Element
}

const (
	KEY_NOT_FOUND    = "key not found"
	CAPACITY_INVALID = "capacity is invalid"
	KEY_INVALID      = "key is invalid"

	LRU  = "lru"
	LFU  = "lfu"
	FIFO = "fifo"
)

func CacheFactory(cacheType string, capacity int) (ICache, error) {
	switch strings.ToLower(cacheType) {
	case LRU:
		return NewLRUCache(capacity)
	case LFU:
		return NewLFUCache(capacity)
	case FIFO:
		return NewFIFOCache(capacity)
	default:
		return nil, errors.New("unsupported cache type")
	}
}
