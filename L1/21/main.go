package main

import (
	"fmt"
	"math"
)

// FahrenheitThermometer - старый интерфейс, который мы не можем изменить
type FahrenheitThermometer interface {
	GetFahrenheitTemperature() float64
}

// OldFahrenheitSensor - конкретная реализация старого интерфейса
type OldFahrenheitSensor struct{}

func (s *OldFahrenheitSensor) GetFahrenheitTemperature() float64 {
	// В реальности здесь было бы получение данных с датчика
	return 98.6 // пример температуры тела человека
}

// CelsiusThermometer - новый интерфейс, который ожидает клиент
type CelsiusThermometer interface {
	GetCelsiusTemperature() float64
}

// FahrenheitToCelsiusAdapter - адаптер, преобразующий старый интерфейс в новый
type FahrenheitToCelsiusAdapter struct {
	fahrenheitThermometer FahrenheitThermometer
}

func (a *FahrenheitToCelsiusAdapter) GetCelsiusTemperature() float64 {
	fahr := a.fahrenheitThermometer.GetFahrenheitTemperature()
	// Конвертируем Фаренгейты в Цельсии
	celsius := (fahr - 32) * 5 / 9
	// Округляем до одного знака после запятой
	return math.Round(celsius*10) / 10
}

func main() {
	// Создаем старый датчик
	oldSensor := &OldFahrenheitSensor{}

	// Создаем адаптер, передавая ему старый датчик
	adapter := &FahrenheitToCelsiusAdapter{
		fahrenheitThermometer: oldSensor,
	}

	// Клиентский код работает с новым интерфейсом
	fmt.Printf("Температура: %.1f°C\n", adapter.GetCelsiusTemperature())
}
