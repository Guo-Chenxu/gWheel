package sort

import "math/rand"

func QuickSort(arr []int) {
	quickSort(arr, 0, len(arr)-1)
}

func quickSort(arr []int, low, high int) {
	if low >= high {
		return
	}

	mid := partition(arr, low, high)
	quickSort(arr, low, mid-1)
	quickSort(arr, mid+1, high)
}

func partition(arr []int, low, high int) int {
	pivot := rand.Int()%(high-low+1) + low
	arr[pivot], arr[low] = arr[low], arr[pivot]
	pivot = low

	for low < high {
		for low < high && arr[high] >= arr[pivot] {
			high--
		}
		for low < high && arr[low] <= arr[pivot] {
			low++
		}
		if low < high {
			arr[low], arr[high] = arr[high], arr[low]
		}
	}
	arr[low], arr[pivot] = arr[pivot], arr[low]

	return low
}
