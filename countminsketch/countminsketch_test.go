package countminsketch

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountMinSketchBase(t *testing.T) {
	cms, err := NewCountMinSketch(uint64(10), uint64(4))
	if err != nil {
		t.Fatal(err)
	}

	cms.AddString("apple", 10)
	t.Log(cms.EstimateString("apple"))
	cms.AddString("banana", 20)
	t.Log(cms.EstimateString("banana"))
}

// based on https://github.com/jehiah/countmin/blob/master/sketch_test.go
func TestAccuracy(t *testing.T) {
	s, err := NewCountMinSketchWithEstimates(0.0001, 0.9999)
	if err != nil {
		t.Error(err)
	}

	iterations := 5500
	var diverged int
	for i := 1; i < iterations; i++ {
		v := uint64(i % 50)

		s.AddString(strconv.Itoa(i), v)
		vv := s.EstimateString(strconv.Itoa(i))
		if vv > v {
			diverged++
		}
	}

	var miss int
	for i := 1; i < iterations; i++ {
		vv := uint64(i % 50)

		v := s.EstimateString(strconv.Itoa(i))
		assert.Equal(t, v >= vv, true)
		if v != vv {
			t.Logf("real: %d, estimate: %d\n", vv, v)
			miss++
		}
	}
	t.Logf("missed %d of %d (%d diverged during adds)", miss, iterations, diverged)
}
