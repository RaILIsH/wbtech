package main

import (
	"fmt"
	"sync"
)

func main() {
	numbers := []int{2, 4, 6, 8, 10}

	var wg sync.WaitGroup
	// Возведение в квадрат с использованием WaitGroup для синхронизации горутин
	for _, num := range numbers {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			square := n * n
			fmt.Println(square)
		}(num)
	}

	wg.Wait()
}
