package domain

import "time"

type User struct {
	ID        int       `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"` // "-" скрывает поле в JSON
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Профили пользователей
type UserProfile struct {
	ID          int        `json:"id" db:"id"`
	UserID      int        `json:"user_id" db:"user_id"`
	FirstName   string     `json:"first_name" db:"first_name"`
	LastName    string     `json:"last_name" db:"last_name"`
	BirthDate   *time.Time `json:"birth_date,omitempty" db:"birth_date"` // Указатель, чтобы можно было null
	PhoneNumber string     `json:"phone_number,omitempty" db:"phone_number"`
	Lang        string     `json:"lang" db:"lang"`
}

// Роли пользователей
type UserRole struct {
	ID     int      `json:"id" db:"id"`
	UserID int      `json:"user_id" db:"user_id"`
	Roles  []string `json:"roles" db:"roles"` // Массив строк для ролей
}
