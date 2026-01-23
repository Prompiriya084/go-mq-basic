package config

import "os"

func Load() Config {
	return Config{
		HTTPPort: getEnv("HTTP_PORT", "8080"),
		Database: DatabaseConfig{
			DB_Host:     getEnv("DB_Host", ""),
			DB_Port:     getEnv("DB_Port", ""),
			DB_Name:     getEnv("DB_Name", ""),
			DB_Username: getEnv("DB_Username", ""),
			DB_Password: getEnv("DB_Password", ""),
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
