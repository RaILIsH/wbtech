package main

import (
	"fmt"
	"time"
)

func main() {
	// Устанавливаем время работы программы в секундах
	N := 5 // завершить через 5 секунд

	// Создаем канал для передачи данных
	dataChan := make(chan int)

	// Запускаем горутину для отправки данных в канал
	go func() {
		i := 0
		for {
			dataChan <- i
			i++
			time.Sleep(500 * time.Millisecond) // небольшая задержка между отправками
		}
	}()

	// Запускаем горутину для чтения данных из канала
	go func() {
		for value := range dataChan {
			fmt.Printf("Получено значение: %d\n", value)
		}
	}()

	// Ждем N секунд, затем завершаем программу
	fmt.Printf("Программа будет работать %d секунд...\n", N)
	<-time.After(time.Duration(N) * time.Second)
	fmt.Println("Время истекло, программа завершается.")
	close(dataChan) // закрываем канал для остановки горутин
}
