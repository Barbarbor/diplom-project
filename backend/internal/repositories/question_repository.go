// TODO: сделать так, чтобы при изменении у не новых вопросов состояние менялось на CHANGED

package repositories

import (
	"backend/internal/domain"
	"database/sql"
	"fmt"
	"time"

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
	QCreatedAt          sql.NullTime          `db:"q_created_at"`
	QUpdatedAt          sql.NullTime          `db:"q_updated_at"`
	OptionID            sql.NullInt64         `db:"o_id"`
	OptionOriginalId    *int                  `db:"o_original_id"`
	OptionQuestionID    sql.NullInt64         `db:"o_question_id"`
	OptionLabel         sql.NullString        `db:"o_label"`
	OptionOrder         *int                  `db:"o_order"`
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

// GetQuestionMaxOrder возвращает максимальное значение question_order для вопросов в survey_questions_temp,
// учитывая только те вопросы, у которых question_state != 'DELETED'.
func (r *questionRepository) GetQuestionMaxOrder(surveyID int) (int, error) {
	var maxOrder int
	query := `SELECT COALESCE(MAX(question_order), 0) FROM survey_questions_temp WHERE survey_id = $1 AND question_state != 'DELETED'`
	err := r.db.Get(&maxOrder, query, surveyID)
	if err != nil {
		return 0, fmt.Errorf("failed to get max question order: %w", err)
	}
	return maxOrder, nil
}

// CreateQuestion создает новый вопрос в таблице survey_questions_temp.
// Поле question_original_id устанавливается в NULL, а question_state в 'NEW'.
func (r *questionRepository) CreateQuestion(question *domain.SurveyQuestionTemp) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	// Получаем максимальное значение question_order (учитывая только вопросы, не удаленные)
	maxOrder, err := r.GetQuestionMaxOrder(question.SurveyID)
	if err != nil {
		tx.Rollback()
		return err
	}

	query := `
		INSERT INTO survey_questions_temp (question_original_id, survey_id, label, type, question_order, question_state, created_at, updated_at)
		VALUES (NULL, $1, $2, $3, $4, 'NEW', NOW(), NOW())
		RETURNING id, question_order`
	err = tx.QueryRow(query, question.SurveyID, question.Label, question.Type, maxOrder+1).
		Scan(&question.ID, &question.QuestionOrder)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create question: %w", err)
	}
	return tx.Commit()
}

// GetQuestionsBySurveyID возвращает список вопросов для опроса из временной таблицы.
func (r *questionRepository) GetQuestionsBySurveyID(surveyID int) ([]*domain.SurveyQuestionTemp, error) {
	var questions []*domain.SurveyQuestionTemp
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
func (r *questionRepository) GetOptionsByQuestionID(questionID int) ([]domain.OptionTemp, error) {
	var options []domain.OptionTemp
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

// GetQuestionByID возвращает вопрос из временной таблицы по его ID.
func (r *questionRepository) GetQuestionByID(questionID int, surveyID int) (*domain.SurveyQuestionTemp, error) {
	var question domain.SurveyQuestionTemp
	query := `
		SELECT *
		FROM survey_questions_temp
		WHERE id = $1
		AND survey_id = $2`
	if err := r.db.Get(&question, query, questionID, surveyID); err != nil {
		return nil, fmt.Errorf("failed to get question by id: %w", err)
	}
	return &question, nil
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
			q.created_at as q_created_at,
			q.updated_at as q_updated_at,
			o.id AS o_id,
			o.question_id AS o_question_id,
			o.option_original_id AS o_original_id,
			o.label AS o_label,
			o.option_order as o_order,
			o.option_state AS o_state,
			o.created_at AS o_created_at,
			o.updated_at AS o_updated_at
		FROM survey_questions_temp q
		LEFT JOIN survey_options_temp o ON q.id = o.question_id
		WHERE q.survey_id = $1
		ORDER BY q.question_order, o.option_order`
	var rows []QuestionOptionRow

	err := r.db.Select(&rows, query, surveyID)
	fmt.Print(err)
	if err != nil {
		return nil, fmt.Errorf("failed to query question options: %w", err)
	}

	return rows, nil
}

// UpdateQuestionType обновляет тип вопроса в таблице survey_questions_temp.
// Параметры:
// - questionID: идентификатор вопроса,
// - newType: новый тип вопроса,
// - currentState: текущее состояние вопроса, полученное из контекста.
func (r *questionRepository) UpdateQuestionType(questionID int, newType domain.QuestionType, currentState string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	now := time.Now()

	if currentState == "NEW" {
		// Удаляем все опции для этого вопроса из временной таблицы
		delQuery := `DELETE FROM survey_options_temp WHERE question_id = $1`
		if _, err := tx.Exec(delQuery, questionID); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to delete options: %w", err)
		}

		// Обновляем тип вопроса, оставляя состояние как NEW
		updateQuery := `
			UPDATE survey_questions_temp
			SET type = $1, updated_at = $2
			WHERE id = $3`
		if _, err := tx.Exec(updateQuery, newType, now, questionID); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update question type (NEW): %w", err)
		}
	} else if currentState == "ACTUAL" || currentState == "CHANGED" {
		// Обновляем опции: помечаем их как DELETED
		updOptionsQuery := `
			UPDATE survey_options_temp
			SET option_state = 'DELETED', updated_at = $1
			WHERE question_id = $2`
		if _, err := tx.Exec(updOptionsQuery, now, questionID); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update options state: %w", err)
		}

		// Если состояние было ACTUAL, меняем его на CHANGED, иначе оставляем CHANGED
		newState := currentState
		if currentState == "ACTUAL" {
			newState = "CHANGED"
		}

		// Обновляем вопрос: меняем тип, состояние и updated_at
		updQuestionQuery := `
			UPDATE survey_questions_temp
			SET type = $1, question_state = $2, updated_at = $3
			WHERE id = $4`
		if _, err := tx.Exec(updQuestionQuery, newType, newState, now, questionID); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update question type (ACTUAL/CHANGED): %w", err)
		}
	} else {
		tx.Rollback()
		return fmt.Errorf("update not allowed for question_state: %s", currentState)
	}

	return tx.Commit()
}

// UpdateQuestionLabel обновляет только label для вопроса в survey_questions_temp.
func (r *questionRepository) UpdateQuestion(questionID int, newLabel string) error {
	query := `
		UPDATE survey_questions_temp
		SET label = $1, updated_at = NOW()
		WHERE id = $2`
	_, err := r.db.Exec(query, newLabel, questionID)
	if err != nil {
		return fmt.Errorf("failed to update question label: %w", err)
	}
	return nil
}

// UpdateQuestionOrder обновляет порядок вопроса в таблице survey_questions_temp.
// Параметры currentOrder и surveyID получаются ранее (например, из контекста).
func (r *questionRepository) UpdateQuestionOrder(questionID int, newOrder, currentOrder, surveyID int) error {

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	// Получаем максимальное значение question_order (учитывая только вопросы, не удаленные)
	maxOrder, err := r.GetQuestionMaxOrder(surveyID)
	if err != nil {
		tx.Rollback()
		return err
	}
	if newOrder > maxOrder {
		return fmt.Errorf("new order value can`t be more than max order value")
	}
	if err := updateEntityOrder(tx, QuestionTable, QuestionFKField, QuestionOrderField, QuestionStateField, questionID, newOrder, currentOrder, surveyID); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
