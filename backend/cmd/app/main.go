package main

import (
	"backend/internal/config"
	"backend/internal/db"
	"backend/internal/handlers"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	// Подключение к базе данных
	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Инициализация роутера
	router := gin.Default()
	router.Use(cors.Default())
	// Регистрация маршрутов
	handlers.RegisterRoutes(router, database)

	// Запуск сервера
	log.Printf("Starting server on %s", cfg.ServerPort)
	router.Run(cfg.ServerPort)
}
