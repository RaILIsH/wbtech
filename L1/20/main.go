package main

import (
	"fmt"
	"strings"
)

func reverseWords(s string) string {
	// Разбиваем строку на слова
	words := strings.Fields(s)

	// Переворачиваем слова на месте
	for i, j := 0, len(words)-1; i < j; i, j = i+1, j-1 {
		words[i], words[j] = words[j], words[i]
	}

	// Собираем строку обратно
	return strings.Join(words, " ")
}

func main() {
	input := "snow dog sun"
	output := reverseWords(input)
	fmt.Printf("Вход:  «%s»\n", input)
	fmt.Printf("Выход: «%s»\n", output)
}
