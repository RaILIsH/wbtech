package main

import (
	"fmt"
	"time"
)

func Sleep(duration time.Duration) {
	<-time.After(duration) // Блокируемся, пока не придет сигнал из канала
}

func main() {
	fmt.Println("Начало сна...")
	Sleep(2 * time.Second)
	fmt.Println("Прошло 2 секунды!")
}
