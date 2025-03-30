package repositories

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// updateEntityOrder – универсальная функция для изменения порядкового номера в таблице.
// table: имя таблицы (например, QuestionTable или OptionTable)
// fkField: имя поля внешнего ключа (например, QuestionFKField или OptionFKField)
// orderField: имя поля порядка (например, QuestionOrderField или OptionOrderField)
// stateField: имя поля состояния (например, QuestionStateField или OptionStateField)
// entityID: идентификатор изменяемой записи
// newOrder: новое значение порядка, которое хотят установить
// currentOrder: текущее значение порядка (полученное из контекста)
// fkValue: значение внешнего ключа (например, surveyID для вопросов или questionID для опций)
// maxOrder: максимальный порядковый номер у сущности
func updateEntityOrder(tx *sqlx.Tx, table, fkField, orderField, stateField string, entityID, newOrder, currentOrder, fkValue int) error {
	// В зависимости от направления перемещения изменяем порядковые номера остальных записей.
	if newOrder < currentOrder {
		// Перемещение вверх: увеличиваем order для записей с order от newOrder до currentOrder - 1
		queryUp := fmt.Sprintf(`
			UPDATE %s
			SET %s = %s + 1, updated_at = NOW()
			WHERE %s = $1 AND %s >= $2 AND %s < $3 AND %s != 'DELETED'`,
			table, orderField, orderField, fkField, orderField, orderField, stateField)
		if _, err := tx.Exec(queryUp, fkValue, newOrder, currentOrder); err != nil {
			return fmt.Errorf("failed to shift entities up: %w", err)
		}
	} else if newOrder > currentOrder {
		// Перемещение вниз: уменьшаем order для записей с order от currentOrder+1 до newOrder
		queryDown := fmt.Sprintf(`
			UPDATE %s
			SET %s = %s - 1, updated_at = NOW()
			WHERE %s = $1 AND %s > $2 AND %s <= $3 AND %s != 'DELETED'`,
			table, orderField, orderField, fkField, orderField, orderField, stateField)
		if _, err := tx.Exec(queryDown, fkValue, currentOrder, newOrder); err != nil {
			return fmt.Errorf("failed to shift entities down: %w", err)
		}
	}
	// Обновляем порядковый номер целевой записи
	finalQuery := fmt.Sprintf(`
		UPDATE %s
		SET %s = $1, updated_at = NOW()
		WHERE id = $2`, table, orderField)
	if _, err := tx.Exec(finalQuery, newOrder, entityID); err != nil {
		return fmt.Errorf("failed to update target entity order: %w", err)
	}
	return nil
}
