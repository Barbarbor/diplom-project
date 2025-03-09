package auth

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"fmt"

	"backend/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

// AuthService определяет бизнес-логику для аутентификации и регистрации.
type AuthService interface {
	RegisterUser(email, password string) (int, error)
	AuthenticateUser(email, password string) (string, error)
	GetUser(userID int) (*models.User, error)
}

type authService struct {
	repo repositories.AuthRepository
}

// NewAuthService создаёт новый сервис для аутентификации.
func NewAuthService(repo repositories.AuthRepository) AuthService {
	return &authService{
		repo: repo,
	}
}

func (s *authService) RegisterUser(email, password string) (int, error) {
	return s.repo.CreateUser(email, password)
}

func (s *authService) AuthenticateUser(email, password string) (string, error) {
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

func (s *authService) GetUser(userID int) (*models.User, error) {
	return s.repo.GetUserByID(userID)
}
