package config

type Config struct {
	PostgreSQL struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
	}
	Kafka struct {
		Broker  string
		Topic   string
		GroupID string
	}
	Server struct {
		Port string
	}
}

func Load() *Config {
	cfg := &Config{}

	// PostgreSQL configuration
	cfg.PostgreSQL.Host = "localhost"
	cfg.PostgreSQL.Port = "5432"
	cfg.PostgreSQL.User = "order_user"
	cfg.PostgreSQL.Password = "order_password"
	cfg.PostgreSQL.DBName = "orders_db"

	// Kafka configuration
	cfg.Kafka.Broker = "localhost:9092"
	cfg.Kafka.Topic = "orders"
	cfg.Kafka.GroupID = "order-service-group"

	// Server configuration
	cfg.Server.Port = "8081"

	return cfg
}
