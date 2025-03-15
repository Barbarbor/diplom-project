package repositories

import (
	"backend/internal/domain"
	"time"
)

// AuthRepository определяет методы для работы с пользователями
type AuthRepository interface {
	CreateUser(email, password string) (int, error)
	GetUserByEmail(email string) (*domain.User, error)
	GetUserByID(userID int) (*domain.User, error)
}

// ProfileRepository определяет методы для работы с профилем пользователя.
type ProfileRepository interface {
	GetUserProfile(userID int) (*domain.UserProfile, error)
	UpdateUserProfile(profile *domain.UserProfile) error
}

// SurveyRepository определяет методы для работы с опросами в БД.
type SurveyRepository interface {
	CreateSurvey(title string, authorID int, hash string, state domain.SurveyState, now time.Time) (int, error)
	GetSurveyByHash(hash string) (*domain.Survey, string, error)
	CheckUserAccess(userID int, surveyID int) (bool, error)           // Добавляем этот метод
	GetSurveysByAuthor(authorID int) ([]*domain.SurveySummary, error) // Новый метод
}

// QuestionRepository определяет методы для работы с вопросами в БД
type QuestionRepository interface {
	CreateQuestion(question *domain.SurveyQuestionTemp) error
	GetQuestionsBySurveyID(surveyID int) ([]*domain.SurveyQuestion, error)
	GetOptionsByQuestionID(questionID int) ([]domain.Option, error)
	GetQuestionOptionRows(surveyID int) ([]QuestionOptionRow, error)
}
