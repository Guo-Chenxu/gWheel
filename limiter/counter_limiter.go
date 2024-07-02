package limiter

import (
	"sync"
	"time"
)

type CounterLimiter struct {
	startTime int64
	interval  int64 // 间隔时间
	maxCount  int32 // 间隔时间内最大请求数
	counter   int32 // 计数器
	mu        sync.Mutex
}

// 间隔时间 (ms)  最大请求数(个)
func NewCounterLimiter(interval int64, maxCount int32) ILimiter {
	return &CounterLimiter{
		startTime: time.Now().UnixMilli(),
		interval:  int64(interval),
		maxCount:  maxCount,
		counter:   0,
	}
}

func (l *CounterLimiter) Try(count int32) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now().UnixMilli()
	if now-l.startTime < l.interval {
		return l.tryGet(count)
	} else {
		l.startTime = now
		l.counter = 0
		return l.tryGet(count)
	}
}

func (l *CounterLimiter) tryGet(count int32) bool {
	nowCount := l.counter + count
	if nowCount <= l.maxCount {
		l.counter = nowCount
		return true
	} else {
		return false
	}
}
