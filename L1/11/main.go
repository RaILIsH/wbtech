package main

import "fmt"

// Функция для нахождения пересечения двух слайсов
func intersection(a, b []int) []int {
	// Создаем map для элементов первого слайса
	elements := make(map[int]bool)
	for _, num := range a {
		elements[num] = true
	}

	// Проверяем, какие элементы из второго слайса есть в map
	var result []int
	for _, num := range b {
		if elements[num] {
			result = append(result, num)
			// Чтобы избежать дубликатов, помечаем элемент как использованный
			elements[num] = false
		}
	}

	return result
}

func main() {
	A := []int{1, 2, 3}
	B := []int{2, 3, 4}

	fmt.Println(intersection(A, B)) // [2 3]
}
