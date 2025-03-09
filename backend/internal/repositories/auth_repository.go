package repositories

import (
	"backend/internal/domain"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type authRepository struct {
	db *sqlx.DB
}

// NewAuthRepository создаёт новый репозиторий для работы с пользователями.
func NewAuthRepository(db *sqlx.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) CreateUser(email, password string) (int, error) {
	// Проверяем, существует ли пользователь
	var exists bool
	err := r.db.Get(&exists, "SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)", email)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, fmt.Errorf("user already exists")
	}

	// Хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	var userID int
	query := "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id"
	if err := r.db.QueryRow(query, email, string(hashedPassword)).Scan(&userID); err != nil {
		return 0, err
	}

	return userID, nil
}

func (r *authRepository) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	query := "SELECT * FROM users WHERE email = $1"
	if err := r.db.Get(&user, query, email); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) GetUserByID(userID int) (*domain.User, error) {
	var user domain.User
	query := "SELECT * FROM users WHERE id = $1"
	if err := r.db.Get(&user, query, userID); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}
