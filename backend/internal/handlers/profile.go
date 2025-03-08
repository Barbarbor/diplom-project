package handlers

import (
	"backend/internal/models"
	"backend/pkg/profile"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// GetProfileHandler получает профиль пользователя
func GetProfileHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Пример получения идентификатора пользователя из токена
		userIDStr := c.GetString("userID")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userID"})
			return
		}

		ctx := c.Request.Context()
		userProfile, err := profile.GetUserProfile(ctx, db, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"profile": userProfile})
	}
}

// UpdateProfileHandler обновляет профиль пользователя
func UpdateProfileHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userProfile models.UserProfile
		if err := c.ShouldBindJSON(&userProfile); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Пример получения идентификатора пользователя из токена
		userIDStr := c.GetString("userID")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userID"})
			return
		}
		userProfile.UserID = userID

		ctx := c.Request.Context()
		err = profile.UpdateUserProfile(ctx, db, &userProfile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
	}
}
