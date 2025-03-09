package repositories

import (
	"backend/internal/domain"
	"encoding/json"

	"github.com/jmoiron/sqlx"
)

// QuestionRepository отвечает за операции с вопросами в БД
type questionRepository struct {
	db *sqlx.DB
}

// NewQuestionRepository создаёт новый экземпляр репозитория
func NewQuestionRepository(db *sqlx.DB) QuestionRepository {
	return &questionRepository{db: db}
}

// CreateQuestion добавляет новый вопрос в опрос
func (r *questionRepository) CreateQuestion(question *domain.SurveyQuestion) error {
	optionsJSON, err := json.Marshal(question.Options)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO survey_questions (survey_id, label, type, options)
		VALUES ($1, $2, $3, $4) RETURNING id`
	return r.db.QueryRow(query, question.SurveyID, question.Label, question.Type, optionsJSON).Scan(&question.ID)
}
