package repositories

import (
	"backend/internal/domain"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type profileRepository struct {
	db *sqlx.DB
}

// NewProfileRepository создаёт новый репозиторий для профилей.
func NewProfileRepository(db *sqlx.DB) ProfileRepository {
	return &profileRepository{db: db}
}

func (r *profileRepository) GetUserProfile(userID int) (*domain.UserProfile, error) {
	var profile domain.UserProfile
	query := "SELECT * FROM user_profiles WHERE user_id = $1"
	if err := r.db.Get(&profile, query, userID); err != nil {
		return nil, fmt.Errorf("failed to fetch profile: %w", err)
	}
	return &profile, nil
}

func (r *profileRepository) UpdateUserProfile(profile *domain.UserProfile) error {
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
