package services

import (
	"backend/internal/models"

	"github.com/jmoiron/sqlx"
)

func GetSurveys(db *sqlx.DB) ([]models.Survey, error) {
	var surveys []models.Survey
	err := db.Select(&surveys, "SELECT * FROM surveys")
	return surveys, err
}
