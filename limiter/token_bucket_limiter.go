package limiter

import (
	"sync"
	"time"
)

type TokenBucketLimiter struct {
	lastTime int64 // 上次发放令牌的时间
	rate     int32 // 每秒产生的令牌数
	capacity int32 // 桶的大小
	tokens   int32 // 当前令牌数
	mu       sync.Mutex
}

func NewTokenBucketLimiter(rate, capacity int32) ILimiter {
	return &TokenBucketLimiter{
		lastTime: time.Now().UnixMilli(),
		rate:     rate,
		capacity: capacity,
		tokens:   0,
	}
}

func (l *TokenBucketLimiter) Try(count int32) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now().UnixMilli()
	l.tokens += int32(now-l.lastTime) / 1000 * l.rate
	if l.tokens > l.capacity {
		l.tokens = l.capacity
	}

	// 超过一秒时间才会跳转
	if now-l.lastTime >= 1000 {
		l.lastTime = now
	}

	if l.tokens >= count {
		l.tokens -= count
		return true
	} else {
		return false
	}
}
