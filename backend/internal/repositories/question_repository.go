// TODO: сделать так, чтобы при изменении у не новых вопросов состояние менялось на CHANGED

package repositories

import (
	"backend/internal/domain"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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
	QExtraParams        json.RawMessage       `db:"q_extra_params"`
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

	extraParamsJSON, err := json.Marshal(question.ExtraParams)
	if err != nil {
		return fmt.Errorf("failed to marshal extra_params: %w", err)
	}

	query := `
		INSERT INTO survey_questions_temp (
			question_original_id, survey_id, label, type, question_order, question_state, extra_params, created_at, updated_at
		) VALUES (
			NULL, $1, $2, $3, $4, 'NEW', $5, NOW(), NOW()
		)
		RETURNING id, question_order`
	err = tx.QueryRow(query, question.SurveyID, question.Label, question.Type, maxOrder+1, extraParamsJSON).
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
			q.extra_params as q_extra_params,
			o.id AS o_id,
			o.question_id AS o_question_id,
			o.option_original_id AS o_original_id,
			o.label AS o_label,
			o.option_order as o_order,
			o.option_state AS o_state,
			o.created_at AS o_created_at,
			o.updated_at AS o_updated_at
		FROM survey_questions_temp q
		LEFT JOIN survey_options_temp o ON q.id = o.question_id AND o.option_state != 'DELETED'
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
func (r *questionRepository) GetSurveyQuestionsWithOptionsAndAnswers(
	surveyID int,
	interviewID string,
	isDemo bool,
) ([]QuestionOptionAnswerRow, error) {
	var query string
	var args []interface{}

	if isDemo {
		query = `
            SELECT 
                q.id AS q_id,
                q.survey_id AS q_survey_id,
                q.label AS q_label,
                q.type AS q_type,
                q.question_order AS q_order,
                q.extra_params AS q_extra_params,
                o.id AS o_id,
                o.question_id AS o_question_id,
                o.label AS o_label,
                o.option_order AS o_order,
                a.answer AS a_answer
            FROM survey_questions_temp q
            LEFT JOIN survey_options_temp o ON q.id = o.question_id AND o.option_state != 'DELETED'
            LEFT JOIN survey_answers a ON q.id = a.question_id AND a.interview_id = $2
            WHERE q.survey_id = $1 AND q.question_state != 'DELETED'
            ORDER BY q.question_order, o.option_order
        `
		args = []interface{}{surveyID, interviewID}
	} else {
		query = `
            SELECT 
                q.id AS q_id,
                q.survey_id AS q_survey_id,
                q.label AS q_label,
                q.type AS q_type,
                q.question_order AS q_order,
                q.extra_params AS q_extra_params,
                o.id AS o_id,
                o.question_id AS o_question_id,
                o.label AS o_label,
                o.option_order AS o_order,
                a.answer AS a_answer
            FROM survey_questions q
            LEFT JOIN survey_options o ON q.id = o.question_id
            LEFT JOIN survey_answers a ON q.id = a.question_id AND a.interview_id = $2
            WHERE q.survey_id = $1
            ORDER BY q.question_order, o.option_order
        `
		args = []interface{}{surveyID, interviewID}
	}

	var rows []QuestionOptionAnswerRow
	err := r.db.Select(&rows, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query survey questions with options and answers: %w", err)
	}

	// Отладка: выведем количество строк и наличие опций
	for _, row := range rows {
		if row.OID.Valid {
			log.Printf("Row with QID %d has option OID %d, Label: %s", row.QID, row.OID.Int64, row.OLabel.String)
		} else {
			log.Printf("Row with QID %d has no options", row.QID)
		}
	}

	return rows, nil
}

// QuestionOptionAnswerRow — структура для хранения данных из запроса
type QuestionOptionAnswerRow struct {
	QID          int             `db:"q_id"`
	QSurveyID    int             `db:"q_survey_id"`
	QLabel       string          `db:"q_label"`
	QType        string          `db:"q_type"`
	QOrder       int             `db:"q_order"`
	QExtraParams json.RawMessage `db:"q_extra_params"`
	OID          sql.NullInt64   `db:"o_id"`
	OQuestionID  sql.NullInt64   `db:"o_question_id"`
	OLabel       sql.NullString  `db:"o_label"`
	OOrder       *int            `db:"o_order"`
	AAnswer      *string         `db:"a_answer"`
}

// UpdateQuestionType обновляет тип вопроса в таблице survey_questions_temp и возвращает обновленный вопрос.
// Параметры:
// - questionID: идентификатор вопроса,
// - newType: новый тип вопроса,
// - currentState: текущее состояние вопроса, полученное из контекста.
func (r *questionRepository) UpdateQuestionType(questionID int, newType domain.QuestionType, currentState string) (*domain.SurveyQuestionTemp, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	now := time.Now()

	if currentState == "NEW" {
		// Удаляем все опции для этого вопроса из временной таблицы
		delQuery := `DELETE FROM survey_options_temp WHERE question_id = $1`
		if _, err := tx.Exec(delQuery, questionID); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to delete options: %w", err)
		}

		// Обновляем тип вопроса, оставляя состояние как NEW
		updateQuery := `
			UPDATE survey_questions_temp
			SET type = $1, updated_at = $2
			WHERE id = $3`
		if _, err := tx.Exec(updateQuery, newType, now, questionID); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update question type (NEW): %w", err)
		}
	} else if currentState == "ACTUAL" || currentState == "CHANGED" {
		// Обновляем опции: помечаем их как DELETED
		updOptionsQuery := `
			UPDATE survey_options_temp
			SET option_state = 'DELETED', updated_at = $1
			WHERE question_id = $2`
		if _, err := tx.Exec(updOptionsQuery, now, questionID); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update options state: %w", err)
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
			return nil, fmt.Errorf("failed to update question type (ACTUAL/CHANGED): %w", err)
		}
	} else {
		tx.Rollback()
		return nil, fmt.Errorf("update not allowed for question_state: %s", currentState)
	}

	// Извлекаем обновленный вопрос
	var updatedQuestion domain.SurveyQuestionTemp
	selectQuery := `
		SELECT id, survey_id, label, type, question_order, extra_params, question_state, created_at, updated_at
		FROM survey_questions_temp
		WHERE id = $1`
	err = tx.Get(&updatedQuestion, selectQuery, questionID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to fetch updated question: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &updatedQuestion, nil
}

// UpdateQuestionLabel обновляет только label для вопроса в survey_questions_temp.
func (r *questionRepository) UpdateQuestion(questionID int, newLabel string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := updateEntityLabel(tx, QuestionTable, QuestionLabelField, QuestionStateField, questionID, newLabel, nil); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

// UpdateQuestionOrder обновляет порядок вопроса в таблице survey_questions_temp.
// Параметры currentOrder и surveyID получаются ранее (например, из контекста).
func (r *questionRepository) UpdateQuestionOrder(questionID int, newOrder, currentOrder, surveyID int) error {

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
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

func (r *questionRepository) DeleteQuestion(questionID int) error {
	// если NEW — удаляем, иначе ставим state = DELETED
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	return deleteEntity(tx, QuestionTable, QuestionFKField, QuestionOrderField, QuestionStateField, questionID, nil)
}

// Public API: сам открывает транзакцию и вызывает приватную логику
func (r *questionRepository) RestoreQuestion(questionTempID int) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	return r.restoreQuestionTx(tx, questionTempID)
}

func (r *questionRepository) RestoreQuestionTx(tx *sqlx.Tx, questionTempID int) error {
	return r.restoreQuestionTx(tx, questionTempID)
}

// Здесь — ваша существующая логика в точности, только вынесена в приватный метод
func (r *questionRepository) restoreQuestionTx(tx *sqlx.Tx, questionTempID int) error {
	// 1) Вытаскиваем temp-вопрос
	var temp struct {
		ID                 int           `db:"id"`
		QuestionOriginalID sql.NullInt64 `db:"question_original_id"`
	}
	if err := tx.Get(&temp, `
		SELECT id, question_original_id
		FROM survey_questions_temp
		WHERE id = $1`, questionTempID,
	); err != nil {
		return fmt.Errorf("temp question lookup: %w", err)
	}

	// 2) Если это новая запись — просто удаляем
	if !temp.QuestionOriginalID.Valid {
		if _, err := tx.Exec(`DELETE FROM survey_questions_temp WHERE id = $1`, questionTempID); err != nil {
			return fmt.Errorf("delete new temp question: %w", err)
		}
		return nil
	}

	origID := int(temp.QuestionOriginalID.Int64)

	var orig struct {
		ID            int             `db:"id"`
		Label         string          `db:"label"`
		Type          string          `db:"type"`
		QuestionOrder int             `db:"question_order"`
		ExtraParams   json.RawMessage `db:"extra_params"`
	}
	if err := tx.Get(&orig, `
		SELECT id, label, type, question_order, extra_params
		  FROM survey_questions
		 WHERE id = $1`, origID,
	); err != nil {
		return fmt.Errorf("original question lookup: %w", err)
	}

	// 4) Обновляем temp-запись, включая extra_params
	if _, err := tx.Exec(`
		UPDATE survey_questions_temp
		   SET label         = $1,
			   type          = $2,
			   question_order= $3,
			   extra_params  = $4,
			   question_state= 'ACTUAL',
			   updated_at    = NOW()
		 WHERE id = $5`,
		orig.Label, orig.Type, orig.QuestionOrder,
		orig.ExtraParams,
		questionTempID,
	); err != nil {
		return fmt.Errorf("reset temp question: %w", err)
	}

	// 5) Удаляем все temp-опции
	if _, err := tx.Exec(`
		DELETE FROM survey_options_temp
		WHERE question_id = $1`, questionTempID,
	); err != nil {
		return fmt.Errorf("clear temp options: %w", err)
	}

	// 6) Копируем оригинальные опции в temp
	rows, err := tx.Queryx(`
		SELECT id, label, option_order
		FROM survey_options
		WHERE question_id = $1`, origID)
	if err != nil {
		return fmt.Errorf("select original options: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var o struct {
			ID          int    `db:"id"`
			Label       string `db:"label"`
			OptionOrder int    `db:"option_order"`
		}
		if err := rows.StructScan(&o); err != nil {
			return err
		}
		if _, err := tx.Exec(`
			INSERT INTO survey_options_temp
			  (option_original_id, question_id, label, option_order, option_state, created_at, updated_at)
			VALUES ($1, $2, $3, $4, 'ACTUAL', NOW(), NOW())
		`, o.ID, questionTempID, o.Label, o.OptionOrder); err != nil {
			return fmt.Errorf("insert temp option: %w", err)
		}
	}

	return nil
}

// internal/repositories/question_repository.go
// Добавляем метод для получения списка temp-question‑ID по опросу
func (r *questionRepository) GetTempQuestionIDsBySurveyIDTx(
	tx *sqlx.Tx, surveyID int,
) ([]int, error) {
	var ids []int
	err := tx.Select(&ids, `
		SELECT id FROM survey_questions_temp
		WHERE survey_id = $1`, surveyID)
	return ids, err
}

// UpdateQuestionExtraParams сливает переданные params в поле extra_params JSONB
func (r *questionRepository) UpdateQuestionExtraParams(questionID int, params map[string]interface{}) error {
	// 1) маршалим params в JSON
	raw, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("marshal extra_params: %w", err)
	}
	// 2) обновляем через JSONB-конкатенацию `||`
	query := `
      UPDATE survey_questions_temp
      SET extra_params = extra_params || $1::jsonb,
          updated_at = NOW()
      WHERE id = $2`
	if _, err := r.db.Exec(query, raw, questionID); err != nil {
		return fmt.Errorf("update extra_params: %w", err)
	}
	return nil
}

func (r *questionRepository) QuestionExists(questionID int, isDemo bool) (bool, error) {
	var query string
	if isDemo {
		query = "SELECT EXISTS(SELECT 1 FROM survey_questions_temp WHERE id = $1)"
	} else {
		query = "SELECT EXISTS(SELECT 1 FROM survey_questions WHERE id = $1)"
	}

	var exists bool
	err := r.db.QueryRow(query, questionID).Scan(&exists)
	return exists, err
}

func (r *questionRepository) AnswerExists(interviewID string, questionID int) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM survey_answers WHERE interview_id = $1 AND question_id = $2)"
	var exists bool
	err := r.db.QueryRow(query, interviewID, questionID).Scan(&exists)
	return exists, err
}

func (r *questionRepository) CreateAnswer(interviewID string, questionID int, answer string) error {
	query := "INSERT INTO survey_answers (interview_id, question_id, answer) VALUES ($1, $2, $3)"
	_, err := r.db.Exec(query, interviewID, questionID, answer)

	return err
}

func (r *questionRepository) UpdateAnswer(interviewID string, questionID int, answer string) error {
	query := "UPDATE survey_answers SET answer = $1 WHERE interview_id = $2 AND question_id = $3"
	_, err := r.db.Exec(query, answer, interviewID, questionID)
	return err
}
