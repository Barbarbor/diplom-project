package dbutils

import (
	"backend/internal/models"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func IsUserExists(db *sqlx.DB, email string) bool {
	var exists bool
	err := db.Get(&exists, "SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)", email)
	return err == nil && exists
}

func IsValidUserCredentials(db *sqlx.DB, email, password string) bool {
	var user models.User
	err := db.Get(&user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return false
	}

	// Сравниваем пароли
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}
