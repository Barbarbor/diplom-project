package repositories

import (
	"backend/internal/domain"
	"encoding/json"
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
	tx, err := r.db.Beginx()
	if err != nil {
		return 0, err
	}

	// Вставляем запись в таблицу surveys
	var surveyID int
	query := `
		INSERT INTO surveys (title, author_id, hash, state, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $5) RETURNING id`
	if err := tx.QueryRow(query, title, authorID, hash, state, now).Scan(&surveyID); err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to create survey: %w", err)
	}

	// Одновременно вставляем запись в таблицу surveys_temp
	// В surveys_temp поле survey_original_id является как первичным, так и внешним ключом к surveys.
	queryTemp := `
		INSERT INTO surveys_temp (survey_original_id, title, created_at, updated_at)
		VALUES ($1, $2, $3, $4)`
	if _, err := tx.Exec(queryTemp, surveyID, title, now, now); err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to create survey temp record: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return surveyID, nil
}
func (r *surveyRepository) GetSurveyByHash(hash string) (*domain.Survey, string, error) {
	var survey domain.Survey
	var tempTitle string
	var tempUpdatedAt time.Time
	var email string

	query := `
		SELECT s.id, st.title, s.created_at, st.updated_at, s.hash, s.state, u.email
		FROM surveys s
		JOIN users u ON s.author_id = u.id
		JOIN surveys_temp st ON st.survey_original_id = s.id
		WHERE s.hash = $1`
	if err := r.db.QueryRow(query, hash).Scan(
		&survey.ID,
		&tempTitle,
		&survey.CreatedAt,
		&tempUpdatedAt,
		&survey.Hash,
		&survey.State,
		&email,
	); err != nil {
		return nil, "", err
	}

	// Переопределяем title и updated_at из surveys_temp
	survey.Title = tempTitle
	survey.UpdatedAt = tempUpdatedAt

	return &survey, email, nil
}

func (r *surveyRepository) GetSurveyIdByHash(hash string) (int, error) {
	var surveyID int
	query := `
		SELECT s.id
		FROM surveys s
		WHERE s.hash = $1`
	if err := r.db.QueryRow(query, hash).Scan(
		&surveyID,
	); err != nil {
		fmt.Print("surveyid", surveyID)
		return -1, err
	}

	return surveyID, nil
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

func (r *surveyRepository) UpdateSurveyTitle(surveyID int, newTitle string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Обновляем только title и updated_at в temp-таблице
	if _, err := tx.Exec(`
        UPDATE surveys_temp
           SET title      = $1,
               updated_at = NOW()
         WHERE survey_original_id = $2
    `, newTitle, surveyID); err != nil {
		return fmt.Errorf("update surveys_temp title: %w", err)
	}

	return tx.Commit()
}
func (r *surveyRepository) PublishSurvey(surveyID int) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1) Обновляем основной заголовок и переводим в ACTIVE
	if _, err := tx.Exec(`
        UPDATE surveys
           SET title = (SELECT title FROM surveys_temp WHERE survey_original_id = $1),
               state = 'ACTIVE'
         WHERE id = $1 AND state = 'DRAFT'
    `, surveyID); err != nil {
		return fmt.Errorf("publish surveys: %w", err)
	}

	// 2) Вставляем «новые» вопросы (state = NEW)
	type newQ struct {
		TempID      int             `db:"id"`
		Label       string          `db:"label"`
		Type        string          `db:"type"`
		Order       int             `db:"question_order"`
		ExtraParams json.RawMessage `db:"extra_params"`
	}
	var toCreateQs []newQ
	if err := tx.Select(&toCreateQs, `
        SELECT id, label, type, question_order, extra_params
          FROM survey_questions_temp
         WHERE survey_id = $1 AND question_state = 'NEW'
    `, surveyID); err != nil {
		return fmt.Errorf("select new questions: %w", err)
	}
	for _, q := range toCreateQs {
		var newID int
		if err := tx.QueryRow(`
            INSERT INTO survey_questions
                        (survey_id, label, type, question_order, extra_params, created_at, updated_at)
                 VALUES ($1,      $2,    $3::question_type_enum,   $4,  $5,          NOW(),      NOW())
             RETURNING id
        `, surveyID, q.Label, q.Type, q.Order, q.ExtraParams).Scan(&newID); err != nil {
			return fmt.Errorf("insert question #%d: %w", q.TempID, err)
		}
		// Наводим порядок в temp: привязываем original и переводим в ACTUAL
		if _, err := tx.Exec(`
            UPDATE survey_questions_temp
               SET question_original_id = $1,
                   question_state       = 'ACTUAL',
                   updated_at           = NOW()
             WHERE id = $2
        `, newID, q.TempID); err != nil {
			return fmt.Errorf("tag new temp question #%d: %w", q.TempID, err)
		}
	}

	// 3) Обновляем «изменённые» вопросы (state = CHANGED или ACTUAL)
	if _, err := tx.Exec(`
        UPDATE survey_questions AS real
           SET label          = tmp.label,
               type           = tmp.type::question_type_enum,
               question_order = tmp.question_order,
               updated_at     = NOW()
          FROM survey_questions_temp AS tmp
         WHERE tmp.question_original_id = real.id
           AND tmp.survey_id = $1
           AND tmp.question_state IN ('ACTUAL', 'CHANGED')
    `, surveyID); err != nil {
		return fmt.Errorf("sync changed questions: %w", err)
	}
	// Приводим все их temp-состояния к ACTUAL
	if _, err := tx.Exec(`
        UPDATE survey_questions_temp
           SET question_state = 'ACTUAL',
               updated_at     = NOW()
         WHERE survey_id = $1
           AND question_state IN ('ACTUAL', 'CHANGED')
    `, surveyID); err != nil {
		return fmt.Errorf("reset temp states questions: %w", err)
	}

	// 4) Удаляем «помеченные на удаление» (state = DELETED)
	// Сначала опции, потом вопросы
	if _, err := tx.Exec(`
        DELETE FROM survey_options
         WHERE question_id IN (
             SELECT question_original_id
               FROM survey_questions_temp
              WHERE survey_id = $1
                AND question_state = 'DELETED'
         )
    `, surveyID); err != nil {
		return fmt.Errorf("delete real options: %w", err)
	}
	if _, err := tx.Exec(`
        DELETE FROM survey_questions
         WHERE id IN (
             SELECT question_original_id
               FROM survey_questions_temp
              WHERE survey_id = $1
                AND question_state = 'DELETED'
         )
    `, surveyID); err != nil {
		return fmt.Errorf("delete real questions: %w", err)
	}
	// Удаляем temp-записи
	if _, err := tx.Exec(`
        DELETE FROM survey_options_temp
         WHERE question_id IN (
             SELECT id
               FROM survey_questions_temp
              WHERE survey_id = $1
                AND question_state = 'DELETED'
         )
    `, surveyID); err != nil {
		return fmt.Errorf("delete temp options: %w", err)
	}
	if _, err := tx.Exec(`
        DELETE FROM survey_questions_temp
         WHERE survey_id = $1
           AND question_state = 'DELETED'
    `, surveyID); err != nil {
		return fmt.Errorf("delete temp questions: %w", err)
	}

	// 5) Вставляем «новые» опции (state = NEW)
	type newOpt struct {
		TempID     int    `db:"id"`
		QuestionID int    `db:"question_id"`
		Label      string `db:"label"`
		Order      int    `db:"option_order"`
	}
	var toCreateOpts []newOpt
	if err := tx.Select(&toCreateOpts, `
        SELECT id, question_id, label, option_order
          FROM survey_options_temp
         WHERE question_id IN (
             SELECT id FROM survey_questions_temp WHERE survey_id = $1
         ) AND option_state = 'NEW'
    `, surveyID); err != nil {
		return fmt.Errorf("select new options: %w", err)
	}
	for _, opt := range toCreateOpts {
		var newID int
		if err := tx.QueryRow(`
            INSERT INTO survey_options
                        (question_id, label, option_order, created_at, updated_at)
                 VALUES ($1,          $2,    $3,            NOW(),      NOW())
             RETURNING id
        `, opt.QuestionID, opt.Label, opt.Order).Scan(&newID); err != nil {
			return fmt.Errorf("insert option #%d: %w", opt.TempID, err)
		}
		// Обновляем temp: привязываем original и переводим в ACTUAL
		if _, err := tx.Exec(`
            UPDATE survey_options_temp
               SET option_original_id = $1,
                   option_state       = 'ACTUAL',
                   updated_at         = NOW()
             WHERE id = $2
        `, newID, opt.TempID); err != nil {
			return fmt.Errorf("tag new temp option #%d: %w", opt.TempID, err)
		}
	}

	// 6) Обновляем «изменённые» опции (state = ACTUAL или CHANGED)
	if _, err := tx.Exec(`
        UPDATE survey_options AS real
           SET label        = tmp.label,
               option_order = tmp.option_order,
               updated_at   = NOW()
          FROM survey_options_temp AS tmp
         WHERE tmp.option_original_id = real.id
           AND tmp.option_state IN ('ACTUAL', 'CHANGED')
    `); err != nil {
		return fmt.Errorf("sync changed options: %w", err)
	}
	// Приводим все их temp-состояния к ACTUAL
	if _, err := tx.Exec(`
        UPDATE survey_options_temp
           SET option_state = 'ACTUAL',
               updated_at   = NOW()
         WHERE option_state IN ('ACTUAL', 'CHANGED')
    `); err != nil {
		return fmt.Errorf("reset temp states options: %w", err)
	}

	return tx.Commit()
}

// Меняем title = surveys_temp.title и updated_at
func (r *surveyRepository) UpdateSurveyTitleTx(tx *sqlx.Tx, surveyID int) error {
	_, err := tx.Exec(`
        UPDATE surveys s
        SET title      = st.title,
            updated_at = NOW()
        FROM surveys_temp st
        WHERE s.id = st.survey_original_id
          AND st.survey_original_id = $1
    `, surveyID)
	if err != nil {
		return fmt.Errorf("update survey title: %w", err)
	}
	return nil
}
func (r *surveyRepository) BeginTx() (*sqlx.Tx, error) {
	return r.db.Beginx()
}

func (r *surveyRepository) FinishInterview(interviewID string, endTime time.Time) error {
	query := "UPDATE survey_interviews SET status = 'completed', end_time = $1 WHERE id = $2"
	_, err := r.db.Exec(query, endTime, interviewID)
	return err
}
