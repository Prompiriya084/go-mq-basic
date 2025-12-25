package config

import "os"

func Load() Config {
	return Config{
		HTTPPort: getEnv("HTTP_PORT", "8080"),
		Database: DatabaseConfig{
			DSN: getEnv("DB_DSN", ""),
		},
		RabbitMQ: RabbitMQConfig{
			URL: getEnv("RABBITMQ_URL", ""),
		},
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
