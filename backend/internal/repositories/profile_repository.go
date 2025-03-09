package repositories

import (
	"backend/internal/models"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// ProfileRepository определяет методы для работы с профилем пользователя.
type ProfileRepository interface {
	GetUserProfile(userID int) (*models.UserProfile, error)
	UpdateUserProfile(profile *models.UserProfile) error
}

type profileRepository struct {
	db *sqlx.DB
}

// NewProfileRepository создаёт новый репозиторий для профилей.
func NewProfileRepository(db *sqlx.DB) ProfileRepository {
	return &profileRepository{db: db}
}

func (r *profileRepository) GetUserProfile(userID int) (*models.UserProfile, error) {
	var profile models.UserProfile
	query := "SELECT * FROM user_profiles WHERE user_id = $1"
	if err := r.db.Get(&profile, query, userID); err != nil {
		return nil, fmt.Errorf("failed to fetch profile: %w", err)
	}
	return &profile, nil
}

func (r *profileRepository) UpdateUserProfile(profile *models.UserProfile) error {
	query := `UPDATE user_profiles SET first_name = $1, last_name = $2, birth_date = $3, phone_number = $4, lang = $5 WHERE user_id = $6`
	result, err := r.db.Exec(query, profile.FirstName, profile.LastName, profile.BirthDate, profile.PhoneNumber, profile.Lang, profile.UserID)
	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no profile updated for user_id: %d", profile.UserID)
	}
	return nil
}
