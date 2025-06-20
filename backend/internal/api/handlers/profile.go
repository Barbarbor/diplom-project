package handlers

import (
	"backend/internal/domain"
	profile "backend/internal/services/profile_service"
	"fmt"
	"net/http"

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
	userIDInterface, _ := c.Get("user_id")
	userID, ok := userIDInterface.(int)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}
	fmt.Print("HERE!")
	profile, err := h.profileService.GetUserProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Print("OR HERE")
	c.JSON(http.StatusOK, gin.H{"profile": profile})
}
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	var updatedProfile domain.UserProfile
	if err := c.ShouldBindJSON(&updatedProfile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDInterface, _ := c.Get("user_id")
	userID, ok := userIDInterface.(int)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}
	updatedProfile.UserID = userID

	// Preprocess and validate fields
	if updatedProfile.FirstName != nil && *updatedProfile.FirstName == "" {
		updatedProfile.FirstName = nil
	}
	if updatedProfile.LastName != nil && *updatedProfile.LastName == "" {
		updatedProfile.LastName = nil
	}
	if updatedProfile.PhoneNumber != nil && *updatedProfile.PhoneNumber == "" {
		updatedProfile.PhoneNumber = nil
	}
	if updatedProfile.Lang == "" {
		updatedProfile.Lang = ""
	}

	if err := h.profileService.UpdateUserProfile(c.Request.Context(), &updatedProfile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}
