package reentrantlock

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReentrantLock(t *testing.T) {
	lock := NewReentrantLock()
	cnt, m, n := 1000, 0, 0
	wg := sync.WaitGroup{}
	for i := 0; i < cnt; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			lock.Lock()
			defer lock.Unlock()
			n++
			lock.Lock()
			defer lock.Unlock()
			m++
			n--
			t.Log(m, n)
		}(i)
	}

	wg.Wait()
	assert.Equal(t, cnt, m)
	assert.Equal(t, 0, n)
}

func TestInvalidUnlock(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			t.Log(err)
		}
	}()

	lock := NewReentrantLock()
	lock.Unlock()
}

func BenchmarkReentrantLock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lock := NewReentrantLock()
		cnt, m, n := 1000, 0, 0
		wg := sync.WaitGroup{}
		for i := 0; i < cnt; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				lock.Lock()
				defer lock.Unlock()
				n++
				lock.Lock()
				defer lock.Unlock()
				m++
				n--
			}(i)
		}

		wg.Wait()
	}
}
