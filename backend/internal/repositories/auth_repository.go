package repositories

import (
	"backend/internal/domain"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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

	// Начинаем транзакцию
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// Вставляем пользователя
	var userID int
	query := "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id"
	if err := tx.QueryRow(query, email, string(hashedPassword)).Scan(&userID); err != nil {
		return 0, err
	}

	// Вставляем запись в roles с ролью "user"
	rolesQuery := "INSERT INTO roles (user_id, roles) VALUES ($1, $2)"
	_, err = tx.Exec(rolesQuery, userID, pq.Array([]string{"user"}))
	if err != nil {
		return 0, fmt.Errorf("failed to insert roles: %w", err)
	}

	// Вставляем запись в user_profiles (только user_id)
	profileQuery := "INSERT INTO user_profiles (user_id,lang) VALUES ($1, 'ru')"
	_, err = tx.Exec(profileQuery, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user profile: %w", err)
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
