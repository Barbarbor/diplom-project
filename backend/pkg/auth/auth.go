package auth

import (
	"backend/internal/models"
	"backend/pkg/jwt"
	"errors"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser регистрирует нового пользователя
func RegisterUser(db *sqlx.DB, email, password string) (int, error) {
	// Проверяем, существует ли пользователь
	var exists bool
	err := db.Get(&exists, "SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)", email)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, errors.New("user already exists")
	}

	// Хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	// Создаем нового пользователя и возвращаем его ID
	var userID int
	err = db.QueryRow(
		"INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id",
		email, string(hashedPassword),
	).Scan(&userID)
	if err != nil {
		return 0, err
	}

	// Создаем профиль пользователя
	_, err = db.Exec("INSERT INTO user_profiles (user_id) VALUES ($1)", userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

// AuthenticateUser аутентифицирует пользователя и возвращает JWT токен
func AuthenticateUser(db *sqlx.DB, email, password string) (string, error) {
	var user models.User

	// Получаем пользователя из базы данных
	err := db.Get(&user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return "", errors.New("incorrect username or password")
	}

	// Сравниваем пароли
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("incorrect username or password")
	}

	// Генерируем токен
	return jwt.GenerateToken(user.ID, user.Email)
}
