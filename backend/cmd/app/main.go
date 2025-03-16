package main

import (
	api "backend/internal/api"
	"backend/internal/api/handlers"
	"backend/internal/api/middleware"
	"backend/internal/config"
	"backend/internal/db"
	"backend/internal/repositories"
	auth "backend/internal/services/auth_service"
	option "backend/internal/services/option_service"
	profile "backend/internal/services/profile_service"
	question "backend/internal/services/question_service"
	survey "backend/internal/services/survey_service"

	// "backend/pkg/redisclient"
	// "fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	// Подключаем Redis
	// redisclient.Init(cfg.RedisAddr, "", 0)
	// pong, err := redisclient.Client.Ping(redisclient.Ctx).Result()
	// if err != nil {
	// 	log.Fatalf("Failed to connect to Redis: %v", err)
	// }
	// fmt.Println("Redis connected:", pong)

	// Подключение к PostgreSQL
	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Инициализация репозиториев
	authRepo := repositories.NewAuthRepository(database)
	profileRepo := repositories.NewProfileRepository(database)
	surveyRepo := repositories.NewSurveyRepository(database)
	questionRepo := repositories.NewQuestionRepository(database)
	optionRepo := repositories.NewOptionRepository(database)

	// Инициализация сервисов
	authService := auth.NewAuthService(authRepo)
	profileService := profile.NewProfileService(profileRepo)
	surveyService := survey.NewSurveyService(surveyRepo, questionRepo)
	questionService := question.NewQuestionService(questionRepo)
	optionService := option.NewOptionService(optionRepo)

	// Инициализация хэндлеров
	authHandler := handlers.NewAuthHandler(authService)
	profileHandler := handlers.NewProfileHandler(profileService)
	surveyHandler := handlers.NewSurveyHandler(surveyService)
	questionHandler := handlers.NewQuestionHandler(questionService)
	optionHandler := handlers.NewOptionHandler(optionService)

	surveyAccessMiddleware := middleware.SurveyAccessMiddleware(surveyRepo)
	questionMiddleware := middleware.QuestionMiddleware(questionRepo)

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
	api.RegisterRoutes(router, authHandler, profileHandler, surveyHandler, questionHandler, optionHandler, surveyAccessMiddleware, questionMiddleware)

	// Запуск сервера
	log.Printf("Starting server on %s", cfg.ServerPort)
	router.Run(cfg.ServerPort)
}
