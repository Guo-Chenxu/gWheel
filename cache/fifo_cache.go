package cache

import (
	"container/list"
	"errors"
	"sync"
)

type FIFOCache struct {
	mu       sync.RWMutex
	capacity int
	cache    map[string]*list.Element
	list     *list.List
}


// Clear implements ICache.
func (f *FIFOCache) Clear() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	for k, v := range f.cache {
		f.list.Remove(v)
		delete(f.cache, k)
	}

	return nil
}

// Delete implements ICache.
func (f *FIFOCache) Delete(key string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if element, ok := f.cache[key]; key != "" && ok {
		f.removeElement(key, element)
		return nil
	}

	return errors.New(KEY_NOT_FOUND)
}

// Get implements ICache.
func (f *FIFOCache) Get(key string) (interface{}, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if element, ok := f.cache[key]; key != "" && ok {
		return element.Value.(*entry).value, nil
	}

	return nil, errors.New(KEY_NOT_FOUND)
}

// Set implements ICache.
func (f *FIFOCache) Set(key string, value interface{}) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if key == "" {
		return errors.New(KEY_INVALID)
	}

	if element, ok := f.cache[key]; ok {
		element.Value.(*entry).value = value
		return nil
	}

	element := f.list.PushBack(&entry{key, value})
	f.cache[key] = element

	for f.list.Len() > f.capacity {
		e := f.list.Front()
		if e != nil {
			f.removeElement("", e)
		}
	}

	return nil
}

func (f *FIFOCache) removeElement(k string, e *list.Element) {
	if k == "" {
		k = e.Value.(*entry).key
	}
	f.list.Remove(e)
	delete(f.cache, k)
}

func NewFIFOCache(capacity int) (ICache, error) {
	if capacity <= 0 {
		return nil, errors.New(CAPACITY_INVALID)
	}

	return &FIFOCache{
		capacity: capacity,
		cache:    make(map[string]*list.Element, capacity),
		list:     list.New(),
	}, nil
}
