package repositories

import (
	"backend/internal/domain"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type surveyRepository struct {
	db *sqlx.DB
}

// NewSurveyRepository создаёт новый экземпляр репозитория с внедрённой БД.
func NewSurveyRepository(db *sqlx.DB) SurveyRepository {
	return &surveyRepository{db: db}
}

func (r *surveyRepository) CreateSurvey(title string, authorID int, hash string, state domain.SurveyState, now time.Time) (int, error) {
	var surveyID int
	query := `
		INSERT INTO surveys (title, author_id, hash, state, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $5) RETURNING id`
	if err := r.db.QueryRow(query, title, authorID, hash, state, now).Scan(&surveyID); err != nil {
		return 0, fmt.Errorf("failed to create survey: %w", err)
	}
	return surveyID, nil
}

func (r *surveyRepository) GetSurveyByHash(hash string) (*domain.Survey, string, error) {
	var survey domain.Survey
	var email string
	query := `
		SELECT s.id, s.title, s.created_at, s.updated_at, s.hash, s.state, u.email
		FROM surveys s
		JOIN users u ON s.author_id = u.id
		WHERE s.hash = $1`
	if err := r.db.QueryRow(query, hash).Scan(
		&survey.ID,
		&survey.Title,
		&survey.CreatedAt,
		&survey.UpdatedAt,
		&survey.Hash,
		&survey.State,
		&email,
	); err != nil {
		return nil, "", err
	}
	return &survey, email, nil
}

func (r *surveyRepository) CheckUserAccess(userID int, surveyID int) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM surveys WHERE id = $1 AND author_id = $2`
	err := r.db.QueryRow(query, surveyID, userID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check user access: %w", err)
	}
	return count > 0, nil
}

func (r *surveyRepository) GetSurveysByAuthor(authorID int) ([]*domain.SurveySummary, error) {
	var summaries []*domain.SurveySummary
	query := `
		SELECT title, created_at, updated_at, hash, state
		FROM surveys
		WHERE author_id = $1
		ORDER BY created_at DESC`
	if err := r.db.Select(&summaries, query, authorID); err != nil {
		return nil, fmt.Errorf("failed to fetch surveys: %w", err)
	}
	return summaries, nil
}
