package main

import (
	"sort"
	"strings"
)

func findAnagrams(words []string) map[string][]string {
	// Приводим все слова к нижнему регистру
	lowerWords := make([]string, len(words))
	for i, word := range words {
		lowerWords[i] = strings.ToLower(word)
	}

	// Создаем map для группировки анаграмм
	anagramGroups := make(map[string][]string)

	// Проходим по всем словам
	for _, word := range lowerWords {
		// Сортируем буквы слова для создания ключа
		runes := []rune(word)
		sort.Slice(runes, func(i, j int) bool {
			return runes[i] < runes[j]
		})
		key := string(runes)

		// Добавляем слово в соответствующую группу
		anagramGroups[key] = append(anagramGroups[key], word)
	}

	// Создаем результирующий map
	result := make(map[string][]string)

	// Обрабатываем группы анаграмм
	for _, group := range anagramGroups {
		// Пропускаем группы из одного слова
		if len(group) <= 1 {
			continue
		}

		// Удаляем дубликаты и сортируем
		uniqueSorted := removeDuplicatesAndSort(group)

		// Используем первое слово в отсортированном списке как ключ
		result[uniqueSorted[0]] = uniqueSorted
	}

	return result
}

// Функция для удаления дубликатов и сортировки
func removeDuplicatesAndSort(words []string) []string {
	// Удаляем дубликаты
	unique := make(map[string]bool)
	for _, word := range words {
		unique[word] = true
	}

	// Создаем срез уникальных слов
	result := make([]string, 0, len(unique))
	for word := range unique {
		result = append(result, word)
	}

	// Сортируем по возрастанию
	sort.Strings(result)
	return result
}
