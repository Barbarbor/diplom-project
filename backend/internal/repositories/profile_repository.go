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
	var birthDateArg interface{} // Use interface{} to handle nil or *time.Time
	if profile.BirthDate != nil && profile.BirthDate.Time != nil {
		birthDateArg = profile.BirthDate.Time
	} else {
		birthDateArg = nil
	}

	query := `
		UPDATE user_profiles 
		SET 
			first_name = COALESCE(NULLIF($1, ''), first_name),
			last_name = COALESCE(NULLIF($2, ''), last_name),
			birth_date = COALESCE($3, birth_date),
			phone_number = COALESCE(NULLIF($4, ''), phone_number),
			lang = COALESCE(NULLIF($5, ''), lang)
		WHERE user_id = $6`
	result, err := r.db.Exec(query,
		profile.FirstName,
		profile.LastName,
		birthDateArg,
		profile.PhoneNumber,
		profile.Lang,
		profile.UserID)
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
