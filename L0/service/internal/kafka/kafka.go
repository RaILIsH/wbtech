package kafka

import (
	"context"
	"encoding/json"
	"log"
	"service/service/internal/cache"
	"service/service/internal/db"
	"service/service/internal/models"

	"github.com/segmentio/kafka-go"
)

func Consume(ctx context.Context, broker, topic, groupID string, db *db.Database, cache *cache.Cache) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	defer r.Close()

	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			log.Printf("Kafka error: %v", err)
			continue
		}

		var order models.Order
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Printf("Failed to parse order: %v", err)
			continue
		}

		if err := db.SaveOrder(&order); err != nil {
			log.Printf("DB save error: %v", err)
			continue
		}

		cache.Set(&order)
		log.Printf("Order saved: %s", order.OrderUID)
	}
}
