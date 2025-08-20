package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// FieldSelector представляет набор выбранных полей
type FieldSelector struct {
	fields map[int]bool
}

// NewFieldSelector создает новый FieldSelector из строки спецификации
func NewFieldSelector(fieldSpec string) (*FieldSelector, error) {
	fs := &FieldSelector{fields: make(map[int]bool)}

	if fieldSpec == "" {
		return fs, nil
	}

	parts := strings.Split(fieldSpec, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.Contains(part, "-") {
			// Обработка диапазона
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("неверный формат диапазона: %s", part)
			}

			start, err := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
			if err != nil {
				return nil, fmt.Errorf("неверное начальное значение диапазона: %s", rangeParts[0])
			}

			end, err := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
			if err != nil {
				return nil, fmt.Errorf("неверное конечное значение диапазона: %s", rangeParts[1])
			}

			if start > end {
				return nil, fmt.Errorf("неверный диапазон: начало (%d) больше конца (%d)", start, end)
			}

			for i := start; i <= end; i++ {
				if i > 0 {
					fs.fields[i] = true
				}
			}
		} else {
			// Обработка отдельного поля
			fieldNum, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("неверный номер поля: %s", part)
			}
			if fieldNum > 0 {
				fs.fields[fieldNum] = true
			}
		}
	}

	return fs, nil
}

// ShouldInclude проверяет, должно ли быть включено поле с указанным номером
func (fs *FieldSelector) ShouldInclude(fieldNum int) bool {
	if len(fs.fields) == 0 {
		return true // Если поля не указаны, включаем все
	}
	return fs.fields[fieldNum]
}

// HasFields проверяет, есть ли выбранные поля
func (fs *FieldSelector) HasFields() bool {
	return len(fs.fields) > 0
}

func main() {
	// Парсинг флагов
	fieldSpec := flag.String("f", "", "список полей для вывода (например: 1,3-5)")
	delimiter := flag.String("d", "\t", "разделитель полей")
	separated := flag.Bool("s", false, "выводить только строки с разделителем")
	flag.Parse()

	// Валидация флагов
	if *fieldSpec == "" {
		fmt.Fprintln(os.Stderr, "Ошибка: необходимо указать поля с помощью флага -f")
		os.Exit(1)
	}

	// Создание селектора полей
	fieldSelector, err := NewFieldSelector(*fieldSpec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка парсинга полей: %v\n", err)
		os.Exit(1)
	}

	// Обработка входных данных
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		// Разделение строки
		fields := strings.Split(line, *delimiter)

		// Проверка на наличие разделителя (если указан флаг -s)
		if *separated && len(fields) == 1 {
			continue
		}

		// Формирование выходной строки
		var outputFields []string
		for i := 1; i <= len(fields); i++ {
			if fieldSelector.ShouldInclude(i) {
				outputFields = append(outputFields, fields[i-1])
			}
		}

		// Вывод результата
		if len(outputFields) > 0 {
			fmt.Println(strings.Join(outputFields, *delimiter))
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка чтения: %v\n", err)
		os.Exit(1)
	}
}
