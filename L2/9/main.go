package main

import (
	"errors"
	"strconv"
	"unicode"
)

// Unpack распаковывает строку с повторяющимися символами
func Unpack(s string) (string, error) {
	var result []rune
	runes := []rune(s)
	n := len(runes)

	if n == 0 {
		return "", nil
	}

	// Проверка, если строка начинается с цифры
	if unicode.IsDigit(runes[0]) {
		return "", errors.New("некорректная строка: начинается с цифры")
	}

	i := 0
	for i < n {
		current := runes[i]

		// Обработка escape-последовательностей
		if current == '\\' {
			if i+1 < n {
				result = append(result, runes[i+1])
				i += 2
				continue
			} else {
				// Одиночный обратный слеш в конце строки
				result = append(result, '\\')
				i++
				continue
			}
		}

		// Если текущий символ не цифра и не escape
		if !unicode.IsDigit(current) {
			// Проверяем следующий символ
			if i+1 < n && unicode.IsDigit(runes[i+1]) {
				// Извлекаем число
				j := i + 1
				for j < n && unicode.IsDigit(runes[j]) {
					j++
				}

				countStr := string(runes[i+1 : j])
				count, err := strconv.Atoi(countStr)
				if err != nil {
					return "", err
				}

				// Добавляем символ count раз
				for k := 0; k < count; k++ {
					result = append(result, current)
				}

				i = j
			} else {
				// Просто добавляем символ один раз
				result = append(result, current)
				i++
			}
		} else {
			// Цифра без escape - ошибка
			return "", errors.New("некорректная строка: цифра без escape")
		}
	}

	return string(result), nil
}
