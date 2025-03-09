package handlers

import (
	"backend/internal/middleware"
	"backend/internal/repositories"
	auth "backend/internal/services/auth_service"
	profile "backend/internal/services/profile_service"
	survey "backend/internal/services/survey_service"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RegisterRoutes(router *gin.Engine, db *sqlx.DB) {
	authRepo := repositories.NewAuthRepository(db)
	profileRepo := repositories.NewProfileRepository(db)
	surveyRepo := repositories.NewSurveyRepository(db)
	// Инициализируем сервисы, внедряя репозитории
	authService := auth.NewAuthService(authRepo)
	profileService := profile.NewProfileService(profileRepo)
	surveyService := survey.NewSurveyService(surveyRepo)
	// Инициализируем контроллеры, внедряя сервисы
	authHandler := NewAuthHandler(authService)
	profileHandler := NewProfileHandler(profileService)
	surveyHandler := NewSurveyHandler(surveyService)

	api := router.Group("/api")
	{
		// Authorization routes
		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/register", authHandler.Register)
			authRoutes.POST("/login", authHandler.Login)
			authRoutes.GET("/user", middleware.AuthMiddleware(), authHandler.GetUser)
		}

		// Profile routes
		profileRoutes := api.Group("/profile", middleware.AuthMiddleware())
		{
			profileRoutes.GET("", profileHandler.GetProfile)
			profileRoutes.PUT("", profileHandler.UpdateProfile)
		}

		// Survey routes
		surveyRoutes := api.Group("/surveys", middleware.AuthMiddleware())
		{
			surveyRoutes.POST("", surveyHandler.CreateSurvey)
			surveyRoutes.GET("/:hash", surveyHandler.GetSurvey)
		}
	}
}
