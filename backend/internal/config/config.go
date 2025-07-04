package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	DatabaseURL string
	ServerPort  string
	RedisAddr   string
}

var JwtSecretKey = os.Getenv("JWT_SECRET_KEY")

func LoadConfig() Config {
	config := Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		ServerPort:  os.Getenv("SERVER_PORT"),
		RedisAddr:   os.Getenv("REDIS_ADDR"),
	}

	if config.ServerPort == "" {
		config.ServerPort = ":8000"
	}

	return config
}
