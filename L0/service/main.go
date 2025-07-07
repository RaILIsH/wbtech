package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"service/service/config"
	"service/service/internal/cache"
	"service/service/internal/db"
	"service/service/internal/kafka"
	"service/service/internal/web"
	"syscall"
	"time"
)

func main() {
	cfg := config.Load()

	// Инициализация базы данных
	dbConn, err := db.New(
		cfg.PostgreSQL.Host,
		cfg.PostgreSQL.Port,
		cfg.PostgreSQL.User,
		cfg.PostgreSQL.Password,
		cfg.PostgreSQL.DBName,
	)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	// Инициализация кеша
	cache := cache.New()

	// Загрузка существующих заказов в кеш
	if orders, err := dbConn.GetAllOrders(); err == nil {
		for _, order := range orders {
			cache.Set(&order)
		}
		log.Printf("Loaded %d orders into cache", len(orders))
	} else {
		log.Printf("Failed to load orders into cache: %v", err)
	}

	// Создаем контекст с обработкой прерываний
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запускаем Kafka consumer в отдельной горутине
	go kafka.Consume(ctx, cfg.Kafka.Broker, cfg.Kafka.Topic, cfg.Kafka.GroupID, dbConn, cache)

	// Запускаем HTTP сервер
	go web.StartServer(cfg.Server.Port, cache)

	// Ожидаем сигналы завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down gracefully...")
	cancel()
	time.Sleep(2 * time.Second) // Даем время на завершение операций
	log.Println("Service stopped")
}
