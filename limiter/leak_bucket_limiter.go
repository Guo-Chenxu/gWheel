package limiter

import (
	"sync"
	"time"
)

type LeakBucketLimiter struct {
	lastTime int64 // 上次漏水时间
	rate     int32 // 每秒漏水速率
	capacity int32 // 桶的容量
	water    int32 // 当前桶中的水量
	mu       sync.Mutex
}

func NewLeakBucketLimiter(rate, capacity int32) ILimiter {
	return &LeakBucketLimiter{
		lastTime: time.Now().UnixMilli(),
		rate:     rate,
		capacity: capacity,
		water:    0,
	}
}

func (l *LeakBucketLimiter) Try(count int32) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.water == 0 {
		l.lastTime = time.Now().UnixMilli()
		l.water = count
		return true
	}

	now := time.Now().UnixMilli()
	// 计算漏了多少水和剩余水量
	waterLeaked := int32(now-l.lastTime) / 1000 * l.rate
	l.water -= waterLeaked
	if l.water < 0 {
		l.water = 0
	}

	// 超过一秒时间才会跳转
	if now-l.lastTime >= 1000 {
		l.lastTime = now
	}

	// 能继续加水的则直接处理了，加满则抛弃请求
	if l.water+count <= l.capacity {
		l.water += count
		return true
	}
	return false
}
