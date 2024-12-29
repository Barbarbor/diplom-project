package config

import "os"

type Config struct {
	DatabaseURL string
	ServerPort  string
}

func LoadConfig() Config {
	return Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		ServerPort:  os.Getenv("SERVER_PORT"),
	}
}
