package main

import (
	"fmt"
)

func groupTemperatures(temps []float64) map[int][]float64 {
	groups := make(map[int][]float64)

	for _, temp := range temps {
		// Определяем ключ группы:
		// Для temp = -25.4 → key = -25.4 / 10 = -2.54 → int(-2.54) = -2 → -2 * 10 = -20
		// Для temp = 13.0 → key = 13.0 / 10 = 1.3 → int(1.3) = 1 → 1 * 10 = 10
		key := int(temp/10) * 10
		groups[key] = append(groups[key], temp)
	}

	return groups
}

func main() {
	temperatures := []float64{-25.4, -27.0, 13.0, 19.0, 15.5, 24.5, -21.0, 32.5}
	grouped := groupTemperatures(temperatures)

	// Выводим результат
	for key, values := range grouped {
		fmt.Printf("%d: %v\n", key, values)
	}
}
