package sort

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHeapSort(t *testing.T) {
	arr := []int{5, 3, 8, 4, 2}
	HeapSort(arr)
	assert.Equal(t, []int{2, 3, 4, 5, 8}, arr)
}
func TestQuickSort(t *testing.T) {
	arr := []int{5, 3, 8, 4, 4, 1, 9, 2}
	QuickSort(arr)
	t.Log(arr)
	assert.Equal(t, []int{1, 2, 3, 4, 4, 5, 8, 9}, arr)
}

func TestMergeSort(t *testing.T) {
	arr := []int{5, 3, 8, 4, 4, 1, 9, 2}
	arr = MergeSort(arr)
	t.Log(arr)
	assert.Equal(t, []int{1, 2, 3, 4, 4, 5, 8, 9}, arr)
}
