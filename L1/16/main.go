package main

import "fmt"

func quickSort(arr []int) []int {
	// Базовый случай: если длина среза меньше 2, он уже отсортирован
	if len(arr) < 2 {
		return arr
	}

	left, right := 0, len(arr)-1
	pivot := arr[len(arr)/2]

	// Разделяем срез на элементы меньше опорного и больше опорного
	for left <= right {
		// Ищем элемент слева, который больше или равен опорному
		for arr[left] < pivot {
			left++
		}
		// Ищем элемент справа, который меньше или равен опорному
		for arr[right] > pivot {
			right--
		}
		// Если нашли элементы не на своих местах, меняем их местами
		if left <= right {
			arr[left], arr[right] = arr[right], arr[left]
			left++
			right--
		}
	}

	// Рекурсивно сортируем подмассивы
	if right > 0 {
		quickSort(arr[:right+1])
	}
	if left < len(arr)-1 {
		quickSort(arr[left:])
	}

	return arr
}

func main() {
	arr := []int{10, 5, 2, 3, 7, 8, 1, 4, 6, 9}
	fmt.Println("До сортировки:", arr)
	sortedArr := quickSort(arr)
	fmt.Println("После сортировки:", sortedArr)
}
