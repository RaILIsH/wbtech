package main

import (
	"fmt"
	"strconv"
)

// setBit устанавливает i-й бит (считая с 0) в value (0 или 1)
func setBit(num int64, i uint, value int) int64 {
	if value == 1 {
		return num | (1 << i) // Установка бита в 1
	} else {
		return num &^ (1 << i) // Установка бита в 0 (AND NOT)
	}
}

// printBinary печатает число в 8-битном двоичном формате
func printBinary(num int64) string {
	return fmt.Sprintf("%08s", strconv.FormatInt(num, 2))
}

func main() {
	var num int64 = 5 // 00000101 (биты: 7 6 5 4 3 2 1 0)

	fmt.Printf("Исходное число: %d (%s)\n", num, printBinary(num))

	num = setBit(num, 1, 1)
	fmt.Printf("После установки 1-го бита в 1: %d (%s)\n", num, printBinary(num))

	num = setBit(num, 0, 0)
	fmt.Printf("После установки 0-го бита в 0: %d (%s)\n", num, printBinary(num))
}
