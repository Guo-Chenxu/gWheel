package sort

func HeapSort(arr []int) {
	size := len(arr)
	for i := (size - 2) / 2; i >= 0; i-- {
		heapify(arr, i, size)
	}

	for size > 0 {
		arr[0], arr[size-1] = arr[size-1], arr[0]
		size--
		heapify(arr, 0, size)
	}
}

func heapify(arr []int, idx, size int) {
	parent, child := idx, idx*2+1
	for child < size {
		if child+1 < size && arr[child+1] > arr[child] {
			child += 1
		}

		if arr[child] > arr[parent] {
			arr[child], arr[parent] = arr[parent], arr[child]
			parent = child
			child = parent*2 + 1
		} else {
			break
		}
	}
}
