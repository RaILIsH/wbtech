package main

import (
	"fmt"
)

func reverseString(input string) string {
	// Преобразуем строку в срез рун
	runes := []rune(input)

	// Переворачиваем срез рун
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	// Возвращаем перевёрнутую строку
	return string(runes)
}

func main() {
	input := "главрыба"
	reversed := reverseString(input)
	fmt.Println(reversed) // Выведет: абырвалг

	// Тест с emoji
	fmt.Println(reverseString("Hello, World! 😊")) // Выведет: 😊 !dlroW ,olleH
}
