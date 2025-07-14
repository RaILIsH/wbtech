package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

// 1. Остановка по флагу
func workerFlag(stopFlag *bool, wg *sync.WaitGroup) {
	defer wg.Done()
	for !*stopFlag {
		fmt.Println("WorkerFlag: работаю...")
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Println("WorkerFlag: остановлен по флагу")
}

// 2. Остановка через канал
func workerChan(stopChan chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-stopChan:
			fmt.Println("WorkerChan: остановлен по каналу")
			return
		default:
			fmt.Println("WorkerChan: работаю...")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

// 3. Остановка через контекст
func workerCtx(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("WorkerCtx: остановлен через контекст:", ctx.Err())
			return
		default:
			fmt.Println("WorkerCtx: работаю...")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

// 4. Остановка через runtime.Goexit()
func workerGoexit(wg *sync.WaitGroup) {
	defer wg.Done()
	defer fmt.Println("WorkerGoexit: завершена через Goexit")

	for i := 0; i < 5; i++ {
		fmt.Println("WorkerGoexit: работаю...")
		time.Sleep(500 * time.Millisecond)

		if i == 3 {
			runtime.Goexit()
		}
	}
}

// 5. Остановка по таймауту
func workerTimeout(timeoutChan <-chan time.Time, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-timeoutChan:
			fmt.Println("WorkerTimeout: остановлен по таймауту")
			return
		default:
			fmt.Println("WorkerTimeout: работаю...")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func main() {
	var wg sync.WaitGroup

	// 1. Worker с флагом
	wg.Add(1)
	flagStop := false
	go workerFlag(&flagStop, &wg)

	// 2. Worker с каналом
	wg.Add(1)
	chanStop := make(chan struct{})
	go workerChan(chanStop, &wg)

	// 3. Worker с контекстом
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())
	go workerCtx(ctx, &wg)

	// 4. Worker с runtime.Goexit()
	wg.Add(1)
	go workerGoexit(&wg)

	// 5. Worker с таймаутом
	wg.Add(1)
	timeout := time.After(3 * time.Second)
	go workerTimeout(timeout, &wg)

	// Даем горутинам поработать
	time.Sleep(2 * time.Second)

	// Останавливаем горутины разными способами
	fmt.Println("\nОстанавливаем горутины...")

	// 1. Устанавливаем флаг
	flagStop = true

	// 2. Закрываем канал
	close(chanStop)

	// 3. Отменяем контекст
	cancel()

	// 4. runtime.Goexit() уже сработал внутри горутины
	// 5. Таймаут сработает автоматически через 3 секунды

	// Ждем завершения всех горутин
	wg.Wait()
	fmt.Println("Все горутины остановлены")
}
