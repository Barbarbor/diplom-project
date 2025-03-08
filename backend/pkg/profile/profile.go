package profile

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"backend/internal/models"

	"github.com/jmoiron/sqlx"
)

func GetUserProfile(ctx context.Context, db *sqlx.DB, userID int) (*models.UserProfile, error) {
	var profile models.UserProfile
	err := db.GetContext(ctx, &profile, "SELECT * FROM profiles WHERE id = $1", userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("profile not found for userID: %d", userID)
		}
		return nil, fmt.Errorf("failed to fetch profile: %w", err)
	}
	return &profile, nil
}

func UpdateUserProfile(ctx context.Context, db *sqlx.DB, profile *models.UserProfile) error {
	result, err := db.ExecContext(
		ctx,
		"UPDATE profiles SET first_name=$1, last_name=$2, birth_date=$3, phone_number=$4, language=$5 WHERE id=$6",
		profile.FirstName, profile.LastName, profile.BirthDate, profile.PhoneNumber, profile.Lang, profile.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no profile updated for userID: %s", profile.ID)
	}

	return nil
}
