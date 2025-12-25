package config

type Config struct {
	HTTPPort string
	Database DatabaseConfig
	RabbitMQ RabbitMQConfig
}
