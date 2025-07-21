package main

import "fmt"

func BinarySearch(arr []int, target int) int {
	left, right := 0, len(arr)-1

	for left <= right {
		mid := left + (right-left)/2 // Чтобы избежать переполнения
		if arr[mid] == target {
			return mid
		} else if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return -1
}

func main() {
	var array []int = []int{5, 3, 4, 2, 1, 10, 8, 7, 9, 6}
	cnt := BinarySearch(array, 9)
	fmt.Println(cnt)
}
