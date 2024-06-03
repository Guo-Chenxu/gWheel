package countminsketch

import (
	"errors"
	"hash"
	"hash/fnv"
	"math"
	"sync"
)

// 本文参考：
// https://florian.github.io/count-min-sketch/

type CountMinSketch struct {
	w      uint64
	d      uint64
	table  [][]uint64
	hasher hash.Hash64
	mu     sync.RWMutex
}

type ICountMinSketch interface {
	Add(key []byte, count uint64)
	AddString(key string, count uint64)
	Estimate(key []byte) uint64
	EstimateString(key string) uint64
}

// 每行 w 列，一共 d 行
func NewCountMinSketch(w, d uint64) (ICountMinSketch, error) {
	if w <= 0 || d <= 0 {
		return nil, errors.New("w and d must be greater than 0")
	}

	c := &CountMinSketch{
		w:      w,
		d:      d,
		table:  make([][]uint64, d),
		hasher: fnv.New64(),
	}
	for i := range c.table {
		c.table[i] = make([]uint64, w)
	}

	return c, nil
}

// 概率为 1-delta，误差至多为 epsilon*每行总和
func NewCountMinSketchWithEstimates(epsilon, delta float64) (ICountMinSketch, error) {
	if epsilon <= 0 || delta <= 0 || epsilon >= 1 || delta >= 1 {
		return nil, errors.New("epsilon and delta must be greater than 0 and less than 1")
	}

	// w=⌈e/ϵ⌉  d=⌈ln(1/(1-δ))⌉
	w := uint64(math.Ceil(math.E / epsilon))
	d := uint64(math.Ceil(-math.Log(1 - delta)))
	// w := uint64(math.Ceil(2 / epsilon))
	// d := uint64(math.Ceil(math.Log(1-delta) / math.Log(0.5)))
	return NewCountMinSketch(w, d)
}

func (c *CountMinSketch) AddString(key string, count uint64) {
	c.Add([]byte(key), count)
}

func (c *CountMinSketch) Add(key []byte, count uint64) {
	c.increase(key, count)
}

func (c *CountMinSketch) EstimateString(key string) uint64 {
	return c.Estimate([]byte(key))
}

func (c *CountMinSketch) Estimate(key []byte) uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	min := uint64(math.MaxUint64)
	for i, j := range c.indexOf(key) {
		if c.table[i][j] < min {
			min = c.table[i][j]
		}
	}
	return min
}

func (c *CountMinSketch) increase(key []byte, count uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i, j := range c.indexOf(key) {
		c.table[i][j] += count
	}
}

// 计算每个值在每一行的位置
func (c *CountMinSketch) indexOf(key []byte) []uint64 {
	upper, lower := c.hash(key)
	locs := make([]uint64, c.d)
	for i := uint64(0); i < c.d; i++ {
		locs[i] = (upper + i*lower) % c.w
	}
	return locs
}

// 计算hash, 获取结果的高32位和低32位
func (c *CountMinSketch) hash(key []byte) (uint64, uint64) {
	c.hasher.Reset()
	c.hasher.Write(key)
	hash := c.hasher.Sum64()
	return hash >> 32, hash << 32 >> 32
}
