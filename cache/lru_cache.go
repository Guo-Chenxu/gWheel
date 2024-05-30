package cache

import (
	"container/list"
	"errors"
	"sync"
)

type LRUCache struct {
	mu       sync.RWMutex
	capacity int
	cache    map[string]*list.Element
	list     *list.List
}

// Clear implements ICache.
func (l *LRUCache) Clear() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	for k, v := range l.cache {
		l.removeElement(k, v)
	}

	return nil
}

// Delete implements ICache.
func (l *LRUCache) Delete(key string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if element, ok := l.cache[key]; key != "" && ok {
		l.removeElement(key, element)
		return nil
	}

	return errors.New(KEY_NOT_FOUND)
}

// Get implements ICache.
func (l *LRUCache) Get(key string) (interface{}, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if element, ok := l.cache[key]; key != "" && ok {
		l.list.MoveToFront(element)
		return element.Value.(*entry).value, nil
	}

	return nil, errors.New(KEY_NOT_FOUND)
}

// Set implements ICache.
func (l *LRUCache) Set(key string, value interface{}) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if key == "" {
		return errors.New(KEY_INVALID)
	}

	if element, ok := l.cache[key]; ok {
		element.Value.(*entry).value = value
		l.list.MoveToFront(element)
		return nil
	}

	element := l.list.PushFront(&entry{key, value})
	l.cache[key] = element

	for l.list.Len() > l.capacity {
		e := l.list.Back()
		if e != nil {
			l.removeElement("", e)
		}
	}

	return nil
}

func (l *LRUCache) removeElement(k string, e *list.Element) {
	if k == "" {
		k = e.Value.(*entry).key
	}

	l.list.Remove(e)
	delete(l.cache, k)
}

// NewLRUCache 新建LRU缓存
func NewLRUCache(capacity int) (ICache, error) {
	if capacity <= 0 {
		return nil, errors.New(CAPACITY_INVALID)
	}

	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*list.Element, capacity),
		list:     list.New(),
	}, nil
}
