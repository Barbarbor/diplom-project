package repositories

import (
	"backend/internal/models"
	"time"
)

// AuthRepository определяет методы для работы с пользователями
type AuthRepository interface {
	CreateUser(email, password string) (int, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(userID int) (*models.User, error)
}

// ProfileRepository определяет методы для работы с профилем пользователя.
type ProfileRepository interface {
	GetUserProfile(userID int) (*models.UserProfile, error)
	UpdateUserProfile(profile *models.UserProfile) error
}

// SurveyRepository определяет методы для работы с опросами в БД.
type SurveyRepository interface {
	CreateSurvey(title string, authorID int, hash string, state models.SurveyState, now time.Time) (int, error)
	GetSurveyByHash(hash string) (*models.Survey, string, error)
}
