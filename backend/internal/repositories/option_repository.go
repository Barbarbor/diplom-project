package repositories

import (
	"backend/internal/domain"

	"fmt"

	"github.com/jmoiron/sqlx"
)

type optionRepository struct {
	db *sqlx.DB
}

func NewOptionRepository(db *sqlx.DB) OptionRepository {
	return &optionRepository{db: db}
}

func (r *optionRepository) GetMaxOptionOrder(questionID int) (int, error) {
	var maxOrder int
	query := `SELECT COALESCE(MAX(option_order), 0) FROM survey_options_temp WHERE question_id = $1`
	err := r.db.Get(&maxOrder, query, questionID)
	if err != nil {
		return 0, fmt.Errorf("failed to get max option order: %w", err)
	}
	return maxOrder, nil
}

func (r *optionRepository) CreateOption(option *domain.OptionTemp) error {
	// Начинаем транзакцию
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Получаем следующий порядковый номер для опции
	maxOrder, err := r.GetMaxOptionOrder(option.QuestionID)
	if err != nil {
		return err
	}
	nextOrder := maxOrder + 1

	// Вставляем новую опцию в таблицу survey_options_temp.
	// Предполагается, что поле option_original_id передается как NULL.
	query := `
		INSERT INTO survey_options_temp (option_original_id, question_id, label, option_order, option_state, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, option_order`
	if err := tx.QueryRow(query, option.OptionOriginalID, option.QuestionID, option.Label, nextOrder, option.OptionState).Scan(&option.ID, &option.OptionOrder); err != nil {
		return fmt.Errorf("failed to insert option: %w", err)
	}

	// Обновляем question_state на CHANGED, если текущий статус ACTUAL
	updateQuery := `
		UPDATE survey_questions_temp
		SET question_state = 'CHANGED',
			updated_at = NOW()
		WHERE id = $1 AND question_state = 'ACTUAL'`
	result, err := tx.Exec(updateQuery, option.QuestionID)
	if err != nil {
		return fmt.Errorf("failed to update question state: %w", err)
	}

	// Проверяем, было ли обновление (если rowsAffected == 0, значит статус не был ACTUAL)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		// Это не ошибка, просто статус уже не ACTUAL, можно продолжить
	}

	return nil
}

func (r *optionRepository) GetOptionById(questionID, optionID int) (*domain.OptionTemp, error) {
	var o domain.OptionTemp
	query := fmt.Sprintf(
		`SELECT id, option_original_id, question_id, label, option_order, option_state, created_at, updated_at
		 FROM %s WHERE id=$1 AND question_id=$2`, OptionTable)
	if err := r.db.Get(&o, query, optionID, questionID); err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *optionRepository) UpdateOptionOrder(optionID, newOrder, currentOrder, questionID int) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := updateEntityOrder(tx, OptionTable, OptionFKField, OptionOrderField, OptionStateField,
		optionID, newOrder, currentOrder, questionID); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (r *optionRepository) UpdateOptionLabel(optionID int, newLabel string, questionID int) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := updateEntityLabel(tx, OptionTable, OptionLabelField, OptionStateField, optionID, newLabel, &questionID); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (r *optionRepository) DeleteOption(optionID int, questionID int) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := deleteEntity(tx, OptionTable, OptionFKField, OptionOrderField, OptionStateField, optionID, &questionID); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
