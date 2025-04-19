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
		TempID int    `db:"id"`
		Label  string `db:"label"`
		Type   string `db:"type"`
		Order  int    `db:"question_order"`
	}
	var toCreateQs []newQ
	if err := tx.Select(&toCreateQs, `
        SELECT id, label, type, question_order
          FROM survey_questions_temp
         WHERE survey_id = $1 AND question_state = 'NEW'
    `, surveyID); err != nil {
		return fmt.Errorf("select new questions: %w", err)
	}
	for _, q := range toCreateQs {
		var newID int
		if err := tx.QueryRow(`
            INSERT INTO survey_questions
                        (survey_id, label, type, question_order, created_at, updated_at)
                 VALUES ($1,      $2,    $3,   $4,            NOW(),      NOW())
             RETURNING id
        `, surveyID, q.Label, q.Type, q.Order).Scan(&newID); err != nil {
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

	// 3) Обновляем «изменённые» (state = CHANGED или ACTUAL) вопросы
	//    (например, label/type/order)
	if _, err := tx.Exec(`
        UPDATE survey_questions AS real
           SET label          = tmp.label,
               type           = tmp.type,
               question_order = tmp.question_order,
               updated_at     = NOW()
          FROM survey_questions_temp AS tmp
         WHERE tmp.question_original_id = real.id
           AND tmp.survey_id = $1
           AND tmp.question_state IN ('ACTUAL', 'CHANGED')
    `, surveyID); err != nil {
		return fmt.Errorf("sync changed questions: %w", err)
	}
	// приводим все их temp‑состояния к ACTUAL
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
	//   сначала опции, потом вопросы
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
	// и теперь сами temp‑записи
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

	// 5) (опционально) можно почистить surveys_temp, если больше не нужен

	return tx.Commit()
}

func (r *surveyRepository) RestoreSurvey(tx *sqlx.Tx, surveyID int) error {
	// 1. Сначала восстанавливаем title и updated_at в основной таблице из temp
	_, err := tx.Exec(`
		UPDATE surveys s
		   SET title = t.title,
		       updated_at = t.updated_at
		  FROM surveys_temp t
		 WHERE t.survey_original_id = s.id
		   AND s.id = $1
	`, surveyID)
	if err != nil {
		return fmt.Errorf("failed to restore survey metadata: %w", err)
	}
	return nil
}
