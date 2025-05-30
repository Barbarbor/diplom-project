package repositories

import (
	"backend/internal/domain"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type interviewRepository struct {
	db *sqlx.DB
}

func NewInterviewRepository(db *sqlx.DB) InterviewRepository {
	return &interviewRepository{db: db}
}

// InterviewExists checks if an interview with the given ID exists
func (r *interviewRepository) InterviewExists(interviewID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM survey_interviews WHERE id = $1)", interviewID).Scan(&exists)
	return exists, err
}

// CreateInterview inserts a new interview into the database and increments started_interviews
func (r *interviewRepository) CreateInterview(interview *domain.SurveyInterview) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Вставляем запись в survey_interviews
	_, err = tx.Exec(`
		INSERT INTO survey_interviews (id, survey_id, user_id, status, start_time, end_time, is_demo)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, interview.ID, interview.SurveyID, interview.UserID, interview.Status, interview.StartTime, interview.EndTime, interview.IsDemo)
	if err != nil {
		fmt.Println("ERRR", err)
		return err
	}

	// Инкремент started_interviews только если isDemo = false
	if !interview.IsDemo {
		_, err = tx.Exec(`
			UPDATE survey_stats
			SET started_interviews = started_interviews + 1
			WHERE survey_id = $1
		`, interview.SurveyID)
		if err != nil {
			fmt.Println("ERRR updating started_interviews:", err)
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// GetInterviewByID fetches an interview by its ID
func (r *interviewRepository) GetInterviewByID(interviewID string) (*domain.SurveyInterview, error) {
	var interview domain.SurveyInterview
	var userID sql.NullInt64
	var endTime sql.NullTime

	err := r.db.QueryRow(`
		SELECT id, survey_id, user_id, status, start_time, end_time, is_demo
		FROM survey_interviews
		WHERE id = $1
	`, interviewID).Scan(
		&interview.ID,
		&interview.SurveyID,
		&userID,
		&interview.Status,
		&interview.StartTime,
		&endTime,
		&interview.IsDemo,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrInterviewNotFound
	}
	if err != nil {
		return nil, err
	}

	// Handle nullable fields
	if userID.Valid {
		userIDInt := int(userID.Int64)
		interview.UserID = &userIDInt
	}
	if endTime.Valid {
		interview.EndTime = &endTime.Time
	}

	return &interview, nil
}
