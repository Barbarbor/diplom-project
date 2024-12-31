package services

import (
	"backend/internal/models"
	"backend/internal/utils"
	"errors"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser регистрирует нового пользователя
func RegisterUser(db *sqlx.DB, email, password string) error {
	// Проверяем, существует ли пользователь
	var exists bool
	err := db.Get(&exists, "SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)", email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("user already exists")
	}

	// Хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Создаем нового пользователя
	_, err = db.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", email, string(hashedPassword))
	return err
}

// AuthenticateUser аутентифицирует пользователя и возвращает JWT токен
func AuthenticateUser(db *sqlx.DB, email, password string) (string, error) {
	var user models.User

	// Получаем пользователя из базы данных
	err := db.Get(&user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return "", errors.New("user not found")
	}

	// Сравниваем пароли
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid password")
	}

	// Генерируем токен
	return utils.GenerateToken(user.ID, user.Email)
}
