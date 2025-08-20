package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"wget/internal/downloader"
	"wget/internal/parser"
	"wget/internal/storage"
)

func main() {
	url := flag.String("url", "", "URL для загрузки")
	depth := flag.Int("depth", 1, "Глубина рекурсии")
	workers := flag.Int("workers", 5, "Количество одновременных загрузок")
	output := flag.String("output", "downloads", "Выходная директория")
	flag.Parse()

	if *url == "" {
		fmt.Println("Использование: wget -url <URL> [-depth <глубина>] [-workers <потоки>] [-output <директория>]")
		os.Exit(1)
	}

	// Создаем директорию для загрузок
	if err := os.MkdirAll(*output, 0755); err != nil {
		log.Fatalf("Ошибка создания директории: %v", err)
	}

	// Инициализируем компоненты
	storage := storage.NewStorage(*output)
	parser := parser.NewParser()
	dl := downloader.NewDownloader(*workers, storage, parser)

	// Запускаем загрузку
	start := time.Now()
	fmt.Printf("Начинаем загрузку %s (глубина: %d)\n", *url, *depth)

	err := dl.DownloadRecursive(*url, *depth, *url)
	if err != nil {
		log.Fatalf("Ошибка загрузки: %v", err)
	}

	fmt.Printf("Загрузка завершена за %v. Файлы сохранены в: %s\n",
		time.Since(start), *output)
}
