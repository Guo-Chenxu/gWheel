package snowflakeid

import (
	"errors"
	"sync"
	"time"
)

const (
	WORKER_BITS = 10 // 10位机器ID
	NUMBER_BITS = 12 // 12位数据位

	WORKER_MAX int64 = -1 ^ (-1 << WORKER_BITS) // 最大机器ID
	NUMBER_MAX int64 = -1 ^ (-1 << NUMBER_BITS) // 最大数据ID

	TIME_SHIFT   = WORKER_BITS + NUMBER_BITS // 时间戳向左移动位数
	WORKER_SHIFT = NUMBER_BITS               // 机器ID向左移动位数

	EPOCH int64 = 1717064243667 // 起始时间戳
)

type IWorker interface {
	NextSnowflakeID() int64
}

type Worker struct {
	mu        sync.Mutex
	timestamp int64
	workerId  int64
	number    int64
}

// NewWorker 创建一个Worker
func NewWorker(workerId int64) (IWorker, error) {
	if workerId < 0 || workerId > WORKER_MAX {
		return nil, errors.New("Worker ID 不合法")
	}

	return &Worker{
		workerId:  workerId,
		timestamp: 0,
		number:    0,
	}, nil
}

// NextID 获取下一个雪花ID
func (w *Worker) NextSnowflakeID() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()

	now := w.now()
	if w.timestamp == now {
		w.number++
		if w.number > NUMBER_MAX {
			for now <= w.timestamp {
				now = w.now()
			}
			w.number = 0
			w.timestamp = now
		}
	} else if now > w.timestamp {
		w.number = 0
		w.timestamp = now
	} else {
		for now < w.timestamp {
			now = w.now()
		}
		w.timestamp = now
		w.number = 0
	}

	return int64((now-EPOCH)<<TIME_SHIFT | w.workerId<<WORKER_SHIFT | w.number)
}

// now 获取当前时间戳
func (w *Worker) now() int64 {
	return time.Now().Local().UnixMilli()
}
