package validators

import (
	"backend/pkg/dbutils"
	"fmt"
	"regexp"

	"github.com/jmoiron/sqlx"
)

func GetRegisterValidationErrors(email, password string, db *sqlx.DB) []map[string]string {
	var validationErrors []map[string]string

	// Валидация email
	if !isValidEmail(email) {
		validationErrors = append(validationErrors, map[string]string{
			"formField": "email",
			"message":   "Invalid email format",
		})
	} else if dbutils.IsUserExists(db, email) {
		validationErrors = append(validationErrors, map[string]string{
			"formField": "email",
			"message":   "User already exists",
		})
	}

	// Валидация пароля
	passwordErrors := validatePassword(password)
	if passwordErrors != "" {
		validationErrors = append(validationErrors, map[string]string{
			"formField": "password",
			"message":   passwordErrors,
		})
	}

	return validationErrors
}

func GetLoginValidationErrors(email, password string, db *sqlx.DB) []map[string]string {
	// Проверяем существование пользователя и корректность пароля
	if !dbutils.IsValidUserCredentials(db, email, password) {
		return []map[string]string{
			{
				"formField": "general",
				"message":   "Invalid email or password",
			},
		}
	}

	return nil
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)
	return re.MatchString(email)
}

func validatePassword(password string) string {
	var errors []string

	if len(password) < 9 {
		errors = append(errors, "be at least 9 characters long")
	}
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		errors = append(errors, "contain at least one number")
	}
	if !regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>_+&=-]`).MatchString(password) {
		errors = append(errors, "contain at least one special character")
	}

	if len(errors) > 0 {
		return "Password should: " + fmt.Sprint(errors)
	}
	return ""
}
