package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func worker(ctx context.Context, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done(): // Получаем сигнал о завершении
			fmt.Printf("Worker %d shutting down...\n", id)
			return
		default:
			// Имитация работы
			fmt.Printf("Worker %d is working...\n", id)
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// Запускаем воркеров
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(ctx, i, &wg)
	}

	// Ожидаем SIGINT (Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Блокируемся до получения сигнала
	<-sigChan
	fmt.Println("\nReceived interrupt signal. Shutting down...")

	// Отменяем контекст и ждём завершения воркеров
	cancel()
	wg.Wait()

	fmt.Println("All workers shut down. Exiting.")
}
