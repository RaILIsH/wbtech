package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

func worker(id int, ch <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for data := range ch {
		fmt.Printf("Worker %d: %s\n", id, data)
	}
}

func main() {
	numWorkers, _ := strconv.Atoi(os.Args[1])
	ch := make(chan string)
	var wg sync.WaitGroup
	// Запускаем воркеров
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, ch, &wg)
	}
	// Главная горутина записывает данные в канал
	dataCounter := 1
	for {
		data := fmt.Sprintf("Data %d", dataCounter)
		ch <- data
		dataCounter++
	}
	// Закрываем канал и ждем завершения воркеров
	// close(ch)
	// wg.Wait()
}
