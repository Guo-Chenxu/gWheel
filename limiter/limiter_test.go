package limiter

import (
	"sync"
	"testing"
	"time"
)

func TestCounterLimiter(t *testing.T) {
	l := NewCounterLimiter(1000, 10)
	wg := sync.WaitGroup{}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			if l.Try(1) {
				t.Log(i, " success")
			}
			wg.Done()
		}(i)

		if i%100 == 0 {
			time.Sleep(time.Second * 2)
			t.Log("\n")
		}
	}

	wg.Wait()
}

func TestLeakBucketLimiter(t *testing.T) {
	l := NewLeakBucketLimiter(1, 5)
	wg := sync.WaitGroup{}

	for i := 0; i < 30; i++ {
		wg.Add(1)
		go func(i int) {
			if l.Try(1) {
				t.Log(i, " success")
			}
			wg.Done()
		}(i)
		time.Sleep(time.Millisecond * 300)
	}

	wg.Wait()
}

func TestTokenBucketLimiter(t *testing.T) {
	l := NewTokenBucketLimiter(1, 5)
	wg := sync.WaitGroup{}

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			if l.Try(1) {
				t.Log(i, " success")
			}
			wg.Done()
		}(i)
		time.Sleep(time.Millisecond * 300)
	}

	wg.Wait()
}
