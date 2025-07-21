package main

import "fmt"

func removeElement(slice []int, i int) []int {
	// Проверка на корректность индекса
	if i < 0 || i >= len(slice) {
		return slice
	}

	// Сдвигаем хвост слайса на один элемент влево
	copy(slice[i:], slice[i+1:])

	// Обнуляем последний элемент (чтобы избежать утечки памяти)
	slice[len(slice)-1] = 0 // или nil для слайсов с ссылочными типами

	// Уменьшаем длину слайса
	return slice[:len(slice)-1]
}

func main() {
	s := []int{1, 2, 3, 4, 5}
	fmt.Println("До удаления:", s) // [1 2 3 4 5]

	s = removeElement(s, 2)           // Удаляем элемент с индексом 2 (значение 3)
	fmt.Println("После удаления:", s) // [1 2 4 5]

	// Проверка, что capacity не изменилась
	fmt.Println("Длина:", len(s), "Емкость:", cap(s))
}
