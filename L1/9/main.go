package main

import (
	"fmt"
	"sync"
)

func main() {
	numbers := make(chan int)
	results := make(chan int)
	var wg sync.WaitGroup

	// Горутина для генерации чисел
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(numbers)

		nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

		// Отправляем числа в канал
		for _, x := range nums {
			numbers <- x
		}
	}()

	// Горутина для обработки чисел (умножение на 2)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(results)

		// Читаем числа из первого канала, умножаем на 2 и отправляем во второй
		for x := range numbers {
			results <- x * 2
		}
	}()

	// Горутина для вывода результатов
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Читаем результаты из второго канала и выводим их
		for res := range results {
			fmt.Println(res)
		}
	}()
	wg.Wait()
}
