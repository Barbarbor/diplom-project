package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	DatabaseURL string
	ServerPort  string
}

var JwtSecretKey = os.Getenv("JWT_SECRET_KEY")

func LoadConfig() Config {

	// Возвращаем конфигурацию
	config := Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		ServerPort:  os.Getenv("SERVER_PORT"),
	}

	// Устанавливаем значение по умолчанию для порта
	if config.ServerPort == "" {
		config.ServerPort = ":8080" // Значение по умолчанию
	}

	return config
}
