package repositories

import (
	"backend/internal/domain"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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
	// Вставляем запись в таблицу survey_stats
	queryStats := `
		INSERT INTO survey_stats (survey_id, started_interviews, completed_interviews)
		VALUES ($1, $2, $3)`
	if _, err := tx.Exec(queryStats, surveyID, 0, 0); err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to create survey stats record: %w", err)
	}

	// Вставляем запись в таблицу survey_roles с ролями ['read', 'edit']
	queryRoles := `
        INSERT INTO survey_roles (survey_id, user_id, roles)
        VALUES ($1, $2, $3)`
	roles := []string{"read", "edit"}
	if _, err := tx.Exec(queryRoles, surveyID, authorID, pq.Array(roles)); err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to create survey roles record: %w", err)
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

func (r *surveyRepository) GetSurveyIdByHash(hash string, isDemo bool) (int, error) {
	var surveyID int
	var query string

	if isDemo {
		query = `
			SELECT s.id
			FROM surveys s
			WHERE s.hash = $1`
	} else {
		query = `
			SELECT s.id
			FROM surveys s
			WHERE s.hash = $1 AND s.state = 'ACTIVE'`
	}

	if err := r.db.QueryRow(query, hash).Scan(&surveyID); err != nil {
		fmt.Print("surveyid", surveyID)
		return -1, err
	}

	return surveyID, nil
}
func (r *surveyRepository) CheckUserAccess(userID int, surveyID int) (bool, error) {
	var count int
	query := `
        SELECT COUNT(*) 
        FROM surveys s
        LEFT JOIN survey_roles sr ON s.id = sr.survey_id AND sr.user_id = $2
        LEFT JOIN roles r ON r.user_id = $2
        WHERE s.id = $1 
        AND (sr.user_id = $2 AND 'edit' = ANY(sr.roles)
             OR r.user_id = $2 AND ('admin' = ANY(r.roles) OR 'moderator' = ANY(r.roles)))`
	err := r.db.QueryRow(query, surveyID, userID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check user access: %w", err)
	}
	return count > 0, nil
}
func (r *surveyRepository) GetSurveysByAuthor(authorID int) ([]*domain.SurveySummary, error) {
	var summaries []*domain.SurveySummary
	query := `
        SELECT st.title, s.created_at, s.updated_at, s.hash, s.state, ss.completed_interviews
        FROM surveys s
        JOIN survey_roles sr ON s.id = sr.survey_id
        JOIN survey_stats ss ON s.id = ss.survey_id
        JOIN surveys_temp st ON s.id = st.survey_original_id -- Assuming hash links surveys and surveys_temp
        WHERE sr.user_id = $1
        AND 'read' = ANY(sr.roles)
        ORDER BY s.created_at DESC`
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
			   extra_params   = tmp.extra_params,
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
		QuestionID int    `db:"question_original_id"` // Теперь это ID из survey_questions
		Label      string `db:"label"`
		Order      int    `db:"option_order"`
	}
	var toCreateOpts []newOpt
	if err := tx.Select(&toCreateOpts, `
        SELECT o.id, q.question_original_id, o.label, o.option_order
          FROM survey_options_temp o
          JOIN survey_questions_temp q ON o.question_id = q.id
         WHERE q.survey_id = $1
           AND o.option_state = 'NEW'
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

	// 7) Удаляем опции, помеченные на удаление (option_state = DELETED)
	if _, err := tx.Exec(`
        DELETE FROM survey_options
         WHERE id IN (
             SELECT option_original_id
               FROM survey_options_temp
              WHERE option_state = 'DELETED'
         )
    `); err != nil {
		return fmt.Errorf("delete real options with state DELETED: %w", err)
	}
	if _, err := tx.Exec(`
        DELETE FROM survey_options_temp
         WHERE option_state = 'DELETED'
    `); err != nil {
		return fmt.Errorf("delete temp options with state DELETED: %w", err)
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

// FinishInterview updates the interview status and increments completed_interviews if not a demo
func (r *surveyRepository) FinishInterview(interviewID string, endTime time.Time, isDemo bool) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Обновляем статус и end_time в survey_interviews
	query := "UPDATE survey_interviews SET status = 'completed', end_time = $1 WHERE id = $2"
	_, err = tx.Exec(query, endTime, interviewID)
	if err != nil {
		return err
	}

	// Инкремент completed_interviews только если isDemo = false
	if !isDemo {
		_, err = tx.Exec(`
			UPDATE survey_stats
			SET completed_interviews = completed_interviews + 1
			WHERE survey_id = (
				SELECT survey_id FROM survey_interviews WHERE id = $1
			)
		`, interviewID)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
func (r *surveyRepository) GetSurveyStats(surveyID int) (*domain.SurveyStats, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Шаг 1: Получаем общую статистику
	var stats domain.SurveyStats
	err = tx.Get(&stats, `
		SELECT started_interviews, completed_interviews
		FROM survey_stats
		WHERE survey_id = $1
	`, surveyID)
	if err != nil {
		if err == sql.ErrNoRows {
			// If no stats exist, initialize with zeros
			stats.StartedInterviews = 0
			stats.CompletedInterviews = 0
		} else {
			return nil, fmt.Errorf("failed to get survey stats: %w", err)
		}
	}

	// Шаг 2: Получаем список вопросов с raw extra_params
	var tempQuestions []struct {
		ID          int                 `db:"id"`
		Label       string              `db:"label"`
		Type        domain.QuestionType `db:"type"`
		ExtraParams json.RawMessage     `db:"extra_params"` // Use json.RawMessage for raw DB data
	}
	err = tx.Select(&tempQuestions, `
		SELECT id, label, type, extra_params
		FROM survey_questions
		WHERE survey_id = $1
		ORDER BY question_order
	`, surveyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get survey questions: %w", err)
	}

	// Шаг 3: Парсим extra_params и собираем финальный список вопросов
	var questions []domain.QuestionStats
	for _, tempQ := range tempQuestions {
		extraParams, err := domain.ParseExtraParams(tempQ.ExtraParams, tempQ.Type)
		if err != nil {
			// Handle parsing error (e.g., log it or use a default value)
			extraParams = domain.ConsentExtraParams{} // Default fallback
		}

		// Создаём новый QuestionStats
		question := domain.QuestionStats{
			ID:          tempQ.ID,
			Label:       tempQ.Label,
			Type:        tempQ.Type,
			ExtraParams: extraParams,
		}

		// Получаем опции для вопросов с опциями
		if question.Type == "single_choice" || question.Type == "multi_choice" {
			var options []domain.OptionStats
			err = tx.Select(&options, `
				SELECT id, label
				FROM survey_options
				WHERE question_id = $1
				ORDER BY option_order
			`, question.ID)
			if err != nil && err != sql.ErrNoRows {
				return nil, fmt.Errorf("failed to get options for question %d: %w", question.ID, err)
			}
			question.Options = options
		}

		// Получаем ответы на вопросы только из завершённых интервью
		var answers []string
		err = tx.Select(&answers, `
			SELECT sa.answer
			FROM survey_answers sa
			JOIN survey_interviews si ON sa.interview_id = si.id
			WHERE sa.question_id = $1
			AND si.status = 'completed'
		`, question.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				// No answers found, set to nil
				question.Answers = nil
			} else {
				return nil, fmt.Errorf("failed to get answers for question %d: %w", question.ID, err)
			}
		} else {
			question.Answers = answers
		}

		questions = append(questions, question)
	}

	// Шаг 4: Получаем времена интервью
	var interviewTimes []domain.InterviewTime
	err = tx.Select(&interviewTimes, `
		SELECT start_time, end_time
		FROM survey_interviews
		WHERE survey_id = $1
		AND status = 'completed'
	`, surveyID)
	if err != nil {
		if err == sql.ErrNoRows {
			// No interview times, leave as empty slice
			interviewTimes = []domain.InterviewTime{}
		} else {
			return nil, fmt.Errorf("failed to get interview times for survey %d: %w", surveyID, err)
		}
	}

	// Собираем результат
	stats.Questions = questions
	stats.InterviewTimes = interviewTimes

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &stats, nil
}

// GetAccessEmails retrieves emails of users with 'edit' access for a survey
func (r *surveyRepository) GetAccessEmails(surveyID int) ([]string, error) {
	var emails []string
	query := `
		SELECT DISTINCT u.email
		FROM survey_roles sr
		JOIN users u ON sr.user_id = u.id
		WHERE sr.survey_id = $1 AND $2 = ANY(sr.roles)
	`
	err := r.db.Select(&emails, query, surveyID, domain.Edit)
	if err != nil {
		return nil, fmt.Errorf("failed to get access emails: %w", err)
	}
	return emails, nil
}

// AddEditAccess adds 'edit' access to a user for the survey
func (r *surveyRepository) AddEditAccess(surveyID int, email string) error {
	user, err := r.GetUserByEmail(email)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user with email %s not found", email)
	}

	var existingRoles []domain.SurveyRole
	err = r.db.Get(&existingRoles, `
		SELECT roles FROM survey_roles WHERE survey_id = $1 AND user_id = $2
	`, surveyID, user.ID)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to get existing roles: %w", err)
	}

	var roles []domain.SurveyRole
	if err == sql.ErrNoRows {
		roles = []domain.SurveyRole{domain.Edit}
	} else {
		roles = existingRoles
		hasEdit := false
		for _, r := range roles {
			if r == domain.Edit {
				hasEdit = true
				break
			}
		}
		if !hasEdit {
			roles = append(roles, domain.Edit)
		}
	}

	query := `
		INSERT INTO survey_roles (survey_id, user_id, roles)
		VALUES ($1, $2, $3)
		ON CONFLICT (survey_id, user_id) DO UPDATE
		SET roles = $3
	`
	_, err = r.db.Exec(query, surveyID, user.ID, pq.Array(roles))
	if err != nil {
		// Check if the error is due to a missing constraint (though this should be resolved by the migration)
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" { // Unique violation
			return fmt.Errorf("duplicate entry for survey %d and user %d: %w", surveyID, user.ID, err)
		}
		return fmt.Errorf("failed to update survey roles: %w", err)
	}
	return nil
}

// RemoveEditAccess removes 'edit' access from a user for the survey
func (r *surveyRepository) RemoveEditAccess(surveyID int, email string, currentUserID int, surveyAuthorEmail string) error {
	user, err := r.GetUserByEmail(email)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user with email %s not found", email)
	}

	if user.ID == currentUserID {
		return fmt.Errorf("cannot remove self from survey access")
	}

	creatorUser, err := r.GetUserByEmail(surveyAuthorEmail)
	if err != nil {
		return err
	}
	if creatorUser != nil && creatorUser.ID == user.ID {
		return fmt.Errorf("cannot remove survey creator from access")
	}

	var roles domain.SurveyRoles
	err = r.db.Get(&roles, `
        SELECT id, survey_id, user_id, roles FROM survey_roles WHERE survey_id = $1 AND user_id = $2
    `, surveyID, user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil // No roles to remove
		}
		return fmt.Errorf("failed to get user roles: %w", err)
	}

	// Since roles.Roles is now pq.StringArray (i.e., []string), filter out "edit"
	var newRoles []string
	for _, r := range roles.Roles {
		if r != string(domain.Edit) {
			newRoles = append(newRoles, r)
		}
	}

	if len(newRoles) == 0 {
		_, err = r.db.Exec(`DELETE FROM survey_roles WHERE id = $1`, roles.ID)
	} else {
		_, err = r.db.Exec(`
            UPDATE survey_roles SET roles = $1 WHERE id = $2
        `, pq.Array(newRoles), roles.ID)
	}
	if err != nil {
		return fmt.Errorf("failed to remove edit access: %w", err)
	}
	return nil
}

// GetUserByEmail retrieves a user by email
func (r *surveyRepository) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, email FROM users WHERE email = $1`
	err := r.db.Get(&user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}
