package handlers

import (
	models "backend/internal/domain"
	profile "backend/internal/services/profile_service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ProfileHandler структурирует обработку запросов для профилей.
type ProfileHandler struct {
	profileService *profile.ProfileService
}

// NewProfileHandler создаёт новый обработчик профилей.
func NewProfileHandler(profileService *profile.ProfileService) *ProfileHandler {
	return &ProfileHandler{profileService: profileService}
}

// GetProfile возвращает профиль пользователя.
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	// Извлекаем user_id из контекста (например, middleware установил его)
	userIDStr := c.GetString("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	profile, err := h.profileService.GetUserProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"profile": profile})
}

// UpdateProfile обновляет профиль пользователя.
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	var updatedProfile models.UserProfile
	if err := c.ShouldBindJSON(&updatedProfile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr := c.GetString("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}
	updatedProfile.UserID = userID

	if err := h.profileService.UpdateUserProfile(c.Request.Context(), &updatedProfile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}
