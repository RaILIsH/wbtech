package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	// Парсим флаги
	column := flag.Int("k", 1, "sort by column number")
	numeric := flag.Bool("n", false, "sort numerically")
	reverse := flag.Bool("r", false, "reverse the result")
	unique := flag.Bool("u", false, "output only unique lines")
	flag.Parse()

	// Получаем имя файла из аргументов
	var filename string
	if flag.NArg() > 0 {
		filename = flag.Arg(0)
	}

	// Читаем строки из файла или STDIN
	lines, err := readLines(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Сортируем строки
	sortLines(lines, *column, *numeric, *reverse, *unique)

	// Выводим результат
	for _, line := range lines {
		fmt.Println(line)
	}
}

// readLines читает строки из файла или STDIN
func readLines(filename string) ([]string, error) {
	var file *os.File
	var err error

	if filename == "" {
		file = os.Stdin
	} else {
		file, err = os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer file.Close()
	}

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// sortLines сортирует строки согласно параметрам
func sortLines(lines []string, column int, numeric, reverse, unique bool) {
	// Уникальные строки
	if unique {
		lines = removeDuplicates(lines)
	}

	// Сортируем
	sort.Slice(lines, func(i, j int) bool {
		key1 := getSortKey(lines[i], column)
		key2 := getSortKey(lines[j], column)

		var result bool
		if numeric {
			num1, err1 := strconv.ParseFloat(key1, 64)
			num2, err2 := strconv.ParseFloat(key2, 64)

			// Если оба числа - сравниваем как числа
			if err1 == nil && err2 == nil {
				result = num1 < num2
			} else if err1 != nil && err2 != nil {
				// Если оба не числа - сравниваем как строки
				result = key1 < key2
			} else {
				// Числа всегда идут перед не-числами
				result = err1 == nil
			}
		} else {
			result = key1 < key2
		}

		if reverse {
			return !result
		}
		return result
	})
}

// getSortKey извлекает ключ для сортировки из строки
func getSortKey(line string, column int) string {
	if column <= 1 {
		return line
	}

	// Разделяем по табуляции
	parts := strings.Split(line, "\t")
	if column-1 < len(parts) {
		return parts[column-1]
	}

	return ""
}

// removeDuplicates удаляет повторяющиеся строки
func removeDuplicates(lines []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, line := range lines {
		if !seen[line] {
			seen[line] = true
			result = append(result, line)
		}
	}

	return result
}
