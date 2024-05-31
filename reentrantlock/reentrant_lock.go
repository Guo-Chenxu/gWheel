package reentrantlock

import (
	"bytes"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
)

type ReentrantLock struct {
	mu    sync.Mutex
	owner int64
	count int64
}

func getGoroutineID() int64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	gid, _ := strconv.ParseInt(string(b), 10, 64)
	return gid
}

func NewReentrantLock() *ReentrantLock {
	return &ReentrantLock{}
}

func (r *ReentrantLock) Lock() {
	gid := getGoroutineID()
	if atomic.LoadInt64(&r.owner) == gid {
		// r.count++
		atomic.AddInt64(&r.count, 1)
		return
	}
	r.mu.Lock()

	atomic.StoreInt64(&r.owner, gid)
	// r.count = 1
	atomic.StoreInt64(&r.count, 1)
}

func (r *ReentrantLock) Unlock() {
	gid := getGoroutineID()
	if atomic.LoadInt64(&r.owner) != gid {
		panic("ReentrantLock: not locked by current goroutine, expected goroutine " +
			strconv.FormatInt(r.owner, 10) + ", got " + strconv.FormatInt(gid, 10))
	}

	// r.count--
	atomic.AddInt64(&r.count, -1)
	if atomic.LoadInt64(&r.count) <= 0 {
		atomic.StoreInt64(&r.owner, 0)
		r.mu.Unlock()
	}
}
