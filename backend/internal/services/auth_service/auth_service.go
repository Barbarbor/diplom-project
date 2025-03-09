package auth

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"fmt"

	"backend/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

// AuthService определяет бизнес-логику для аутентификации и регистрации.
type AuthService struct {
	repo repositories.AuthRepository
}

// NewAuthService создаёт новый сервис для аутентификации.
func NewAuthService(repo repositories.AuthRepository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) RegisterUser(email, password string) (int, error) {
	return s.repo.CreateUser(email, password)
}

func (s *AuthService) AuthenticateUser(email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return "", fmt.Errorf("incorrect username or password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", fmt.Errorf("incorrect username or password")
	}
	// Генерируем токен с помощью jwt пакета
	return jwt.GenerateToken(user.ID, user.Email)
}

func (s *AuthService) GetUser(userID int) (*models.User, error) {
	return s.repo.GetUserByID(userID)
}
