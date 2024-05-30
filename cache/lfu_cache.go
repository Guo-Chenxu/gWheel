package cache

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

type LFUCache struct {
	mu       sync.RWMutex
	capacity int
	cache    map[string]*lfuEntry
	freqMap  map[string]int
	freq2Key map[int]*list.List
	minFreq  int
}

// Clear implements ICache.
func (l *LFUCache) Clear() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	for k, v := range l.cache {
		l.removeEntry(k, v)
	}

	return nil
}

// Delete implements ICache.
func (l *LFUCache) Delete(key string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if entry, ok := l.cache[key]; key != "" && ok {
		l.removeEntry(key, entry)
		return nil
	}

	return errors.New(KEY_NOT_FOUND)
}

// Get implements ICache.
func (l *LFUCache) Get(key string) (interface{}, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if entry, ok := l.cache[key]; key != "" && ok {
		l.increseFreq(key)
		return entry.value, nil
	}

	return nil, errors.New(KEY_NOT_FOUND)
}

// Set implements ICache.
func (l *LFUCache) Set(key string, value interface{}) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if entry, ok := l.cache[key]; key != "" && ok {
		l.increseFreq(key)
		entry.value = value
		return nil
	}

	for len(l.cache) >= l.capacity {
		time.Sleep(time.Millisecond * 100)
		l.removeMinFreq()
	}

	element := l.freq2Key[0].PushBack(&lfuEntry{key: key, value: value})
	l.cache[key] = &lfuEntry{key: key, value: value, freqNode: element}
	l.increseFreq(key)
	l.minFreq = 1

	return nil
}

func (l *LFUCache) removeEntry(key string, entry *lfuEntry) {
	if key == "" {
		key = entry.key
	}

	cnt := l.freqMap[key]
	l.freq2Key[cnt].Remove(entry.freqNode)
	delete(l.freqMap, key)

	delete(l.cache, key)
	// l.list.Remove(element)
}

func (l *LFUCache) increseFreq(key string) {
	// l.mu.Lock()
	// defer l.mu.Unlock()

	cnt := l.freqMap[key]
	l.freq2Key[cnt].Remove(l.cache[key].freqNode)
	cnt++
	l.freqMap[key] = cnt
	if l.freq2Key[cnt] == nil {
		l.freq2Key[cnt] = list.New()
	}
	newElement := l.freq2Key[cnt].PushBack(l.cache[key])
	l.cache[key].freqNode = newElement
}

func (l *LFUCache) removeMinFreq() {
	for freq := l.minFreq; ; {
		if l.freq2Key[freq].Len() == 0 {
			freq++
			continue
		}

		e := l.freq2Key[freq].Front().Value.(*lfuEntry)
		l.removeEntry("", e)
		return
	}
}

func NewLFUCache(capacity int) (ICache, error) {
	if capacity <= 0 {
		return nil, errors.New(CAPACITY_INVALID)
	}

	l := &LFUCache{
		capacity: capacity,
		cache:    make(map[string]*lfuEntry, capacity),
		freqMap:  make(map[string]int, capacity),
		freq2Key: make(map[int]*list.List, capacity),
		minFreq:  0,
	}
	l.freq2Key[0] = list.New()
	return l, nil
}
