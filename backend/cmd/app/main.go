package main

import (
	api "backend/internal/api"
	"backend/internal/config"
	"backend/internal/db"

	"backend/internal/repositories"
	auth "backend/internal/services/auth_service"
	profile "backend/internal/services/profile_service"
	survey "backend/internal/services/survey_service"
	"log"
	"time"

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

	// Инициализация репозиториев
	authRepo := repositories.NewAuthRepository(database)
	profileRepo := repositories.NewProfileRepository(database)
	surveyRepo := repositories.NewSurveyRepository(database)

	// Инициализация сервисов
	authService := auth.NewAuthService(authRepo)
	profileService := profile.NewProfileService(profileRepo)
	surveyService := survey.NewSurveyService(surveyRepo)

	// Инициализация хэндлеров
	authHandler := api.NewAuthHandler(authService)
	profileHandler := api.NewProfileHandler(profileService)
	surveyHandler := api.NewSurveyHandler(surveyService)

	// Инициализация роутера
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Укажите домен вашего фронтенда
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}))
	// Регистрация маршрутов
	api.RegisterRoutes(router, authHandler, profileHandler, surveyHandler)

	// Запуск сервера
	log.Printf("Starting server on %s", cfg.ServerPort)
	router.Run(cfg.ServerPort)
}
