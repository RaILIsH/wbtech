package main

import (
	"fmt"
	"math/big"
)

func main() {
	// Инициализация больших чисел
	a := new(big.Int)
	b := new(big.Int)

	// Установка значений (больше 2^20)
	a.SetString("1000000000000000000", 10)
	b.SetString("500000000000000000", 10)

	// Выполнение операций
	sum := new(big.Int).Add(a, b)
	diff := new(big.Int).Sub(a, b)
	product := new(big.Int).Mul(a, b)
	quotient := new(big.Int).Div(a, b)

	// Вывод результатов
	fmt.Printf("a = %s\n", a.String())
	fmt.Printf("b = %s\n", b.String())
	fmt.Printf("a + b = %s\n", sum.String())
	fmt.Printf("a - b = %s\n", diff.String())
	fmt.Printf("a * b = %s\n", product.String())
	fmt.Printf("a / b = %s\n", quotient.String())
}
