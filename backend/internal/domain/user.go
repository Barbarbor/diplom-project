package domain

import (
	"fmt"
	"strings"
	"time"
)

type User struct {
	ID        int       `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"` // "-" скрывает поле в JSON
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type CustomTime struct {
	Time *time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" || s == "null" {
		ct.Time = nil
		return nil
	}
	t, err := time.Parse("2006-01-02", s) // Adjust format to match your input
	if err != nil {
		return err
	}
	ct.Time = &t
	return nil
}

// Scan implements the sql.Scanner interface for database retrieval
func (ct *CustomTime) Scan(value interface{}) error {
	if value == nil {
		ct.Time = nil
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		ct.Time = &v
	case []byte:
		if len(v) == 0 {
			ct.Time = nil
			return nil
		}
		t, err := time.Parse("2006-01-02", string(v)) // Match the expected date format
		if err != nil {
			return err
		}
		ct.Time = &t
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	return nil
}

type UserProfile struct {
	ID          int         `json:"id" db:"id"`
	UserID      int         `json:"user_id" db:"user_id"`
	FirstName   *string     `json:"first_name,omitempty" db:"first_name"`
	LastName    *string     `json:"last_name,omitempty" db:"last_name"`
	BirthDate   *CustomTime `json:"birth_date,omitempty" db:"birth_date"` // Use CustomTime
	PhoneNumber *string     `json:"phone_number,omitempty" db:"phone_number"`
	Lang        string      `json:"lang" db:"lang"`
}

// Роли пользователей
type UserRole struct {
	ID     int      `json:"id" db:"id"`
	UserID int      `json:"user_id" db:"user_id"`
	Roles  []string `json:"roles" db:"roles"` // Массив строк для ролей
}
