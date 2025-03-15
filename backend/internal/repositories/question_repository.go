package repositories

import (
	"backend/internal/domain"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// QuestionOptionRow представляет одну строку результата запроса с JOIN опросов и опций.
type QuestionOptionRow struct {
	QID                 int                   `db:"q_id"`
	QSurveyID           int                   `db:"q_survey_id"`
	QLabel              string                `db:"q_label"`
	QType               string                `db:"q_type"`
	QQuestionOriginalId *int                  `db:"q_original_id"`
	QOrder              int                   `db:"q_order"`
	QQuestionState      *domain.QuestionState `db:"q_state"`
	OptionID            sql.NullInt64         `db:"o_id"`
	OptionOriginalId    *int                  `db:"o_original_id"`
	OptionQuestionID    sql.NullInt64         `db:"o_question_id"`
	OptionLabel         sql.NullString        `db:"o_label"`
	OptionState         *domain.OptionState   `db:"o_state"`
	OptionCreatedAt     sql.NullTime          `db:"o_created_at"`
	OptionUpdatedAt     sql.NullTime          `db:"o_updated_at"`
}

// QuestionRepository отвечает за операции с вопросами в БД
type questionRepository struct {
	db *sqlx.DB
}

// NewQuestionRepository создаёт новый экземпляр репозитория
func NewQuestionRepository(db *sqlx.DB) QuestionRepository {
	return &questionRepository{db: db}
}
func (r *questionRepository) CreateQuestion(question *domain.SurveyQuestionTemp) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	// Определяем максимальное значение question_order для вопросов в таблице survey_questions_temp
	var maxOrder int
	queryOrder := `SELECT COALESCE(MAX(question_order), 0) FROM survey_questions_temp WHERE survey_id = $1`
	err = tx.Get(&maxOrder, queryOrder, question.SurveyID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Вставляем новый вопрос в таблицу survey_questions_temp
	// Поле question_original_id задается как NULL, а question_state устанавливается в 'NEW'
	query := `
		INSERT INTO survey_questions_temp (question_original_id, survey_id, label, type, question_order, question_state)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, question_order`
	err = tx.QueryRow(query, question.QuestionOriginalID, question.SurveyID, question.Label, question.Type, maxOrder+1, question.QuestionState).
		Scan(&question.ID, &question.QuestionOrder)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// GetQuestionsBySurveyID возвращает список вопросов для опроса из временной таблицы.
func (r *questionRepository) GetQuestionsBySurveyID(surveyID int) ([]*domain.SurveyQuestion, error) {
	var questions []*domain.SurveyQuestion
	query := `
		SELECT id, survey_id, label, type, question_order, question_state
		FROM survey_questions_temp
		WHERE survey_id = $1
		ORDER BY question_order`
	err := r.db.Select(&questions, query, surveyID)
	if err != nil {
		return nil, err
	}
	return questions, nil
}

// GetOptionsByQuestionID возвращает список опций для заданного вопроса из временной таблицы.
func (r *questionRepository) GetOptionsByQuestionID(questionID int) ([]domain.Option, error) {
	var options []domain.Option
	query := `
		SELECT id, question_id, label, created_at, updated_at, option_state
		FROM survey_options_temp
		WHERE question_id = $1`
	err := r.db.Select(&options, query, questionID)
	if err != nil {
		return nil, err
	}
	return options, nil
}

// GetQuestionOptionRows выбирает данные вопросов из временной таблицы
// и джойнит опции из временной таблицы survey_options_temp.
func (r *questionRepository) GetQuestionOptionRows(surveyID int) ([]QuestionOptionRow, error) {
	query := `
		SELECT 
			q.id AS q_id,
			q.survey_id AS q_survey_id,
			q.label AS q_label,
			q.type AS q_type,
			q.question_order AS q_order,
			q.question_state AS q_state,
			q.question_original_id AS q_original_id,
			o.id AS o_id,
			o.question_id AS o_question_id,
			o.option_original_id AS o_original_id,
			o.label AS o_label,
			o.created_at AS o_created_at,
			o.updated_at AS o_updated_at,
			o.option_state AS o_state
		FROM survey_questions_temp q
		LEFT JOIN survey_options_temp o ON q.id = o.question_id
		WHERE q.survey_id = $1
		ORDER BY q.question_order, o.id`
	var rows []QuestionOptionRow
	err := r.db.Select(&rows, query, surveyID)
	if err != nil {
		return nil, fmt.Errorf("failed to query question options: %w", err)
	}
	return rows, nil
}
