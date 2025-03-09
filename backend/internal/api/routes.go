package api

import (
	"backend/internal/handlers"
	"backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, authHandler *handlers.AuthHandler, profileHandler *handlers.ProfileHandler, surveyHandler *handlers.SurveyHandler) {
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
