package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port int
}

func Load() *Config {
	portStr := os.Getenv("CALENDAR_PORT")
	if portStr == "" {
		portStr = "8080"
	}

	port, _ := strconv.Atoi(portStr)
	if port == 0 {
		port = 8080
	}

	return &Config{Port: port}
}
