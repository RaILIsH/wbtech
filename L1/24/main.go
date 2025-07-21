package main

import (
	"fmt"
	"math"
)

// Point представляет точку на плоскости с инкапсулированными координатами x и y.
type Point struct {
	x float64
	y float64
}

// NewPoint создает новую точку с заданными координатами.
func NewPoint(x, y float64) Point {
	return Point{x: x, y: y}
}

// Distance вычисляет расстояние между текущей точкой и другой точкой.
func (p Point) Distance(other Point) float64 {
	dx := p.x - other.x
	dy := p.y - other.y
	return math.Sqrt(dx*dx + dy*dy)
}

func main() {
	// Создаем две точки
	point1 := NewPoint(1.0, 2.0)
	point2 := NewPoint(4.0, 6.0)

	// Вычисляем расстояние между ними
	distance := point1.Distance(point2)

	// Выводим результат
	fmt.Printf("Расстояние между точками: %.2f\n", distance)
}
