package main

import "fmt"

func main() {
	// Исходная последовательность строк
	sequence := []string{"cat", "cat", "dog", "cat", "tree"}

	// Создаем множество с помощью map
	set := make(map[string]bool)

	// Добавляем элементы в множество
	for _, item := range sequence {
		set[item] = true
	}

	// Создаем слайс для хранения уникальных элементов
	uniqueItems := make([]string, 0, len(set))
	for item := range set {
		uniqueItems = append(uniqueItems, item)
	}

	// Выводим результат
	fmt.Println("Множество:", uniqueItems)
}
