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

func CreateSurvey(db *sqlx.DB, survey models.Survey) error {
	_, err := db.Exec("INSERT INTO surveys (title, description, content) VALUES ($1, $2, $3)", survey.Title, survey.Description, survey.Content)
	return err
}
