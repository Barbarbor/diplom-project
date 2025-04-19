package repositories

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// DBRunner defines the interface for executing SQL queries, satisfied by both *sqlx.DB and *sqlx.Tx
type DBRunner interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

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

// --- универсальный метод для логического удаления ---
func deleteEntity(tx *sqlx.Tx, table, stateField string, entityID int) error {
	// получаем текущее состояние
	var state string
	if err := tx.Get(&state, fmt.Sprintf("SELECT %s FROM %s WHERE id=$1", stateField, table), entityID); err != nil {
		return fmt.Errorf("failed to get state: %w", err)
	}
	switch state {
	case "NEW":
		// для NEW — удаляем совсем
		if _, err := tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE id=$1", table), entityID); err != nil {
			return fmt.Errorf("failed to delete new entity: %w", err)
		}
	default:
		// для ACTUAL/CHANGED — помечаем DELETED
		if _, err := tx.Exec(
			fmt.Sprintf("UPDATE %s SET %s='DELETED', updated_at=NOW() WHERE id=$1", table, stateField),
			entityID,
		); err != nil {
			return fmt.Errorf("failed to mark entity deleted: %w", err)
		}
	}
	return nil
}

// --- универсальный метод для смены label ---
func updateEntityLabel(tx *sqlx.Tx, table, labelField, stateField string, entityID int, newLabel string) error {
	if _, err := tx.Exec(
		fmt.Sprintf("UPDATE %s SET %s=$1, updated_at=NOW() WHERE id=$2", table, labelField),
		newLabel, entityID,
	); err != nil {
		return fmt.Errorf("failed to update label: %w", err)
	}
	// если было ACTUAL — ставим CHANGED
	if _, err := tx.Exec(
		fmt.Sprintf("UPDATE %s SET %s='CHANGED' WHERE id=$1 AND %s='ACTUAL'", table, stateField, stateField),
		entityID,
	); err != nil {
		return fmt.Errorf("failed to update state: %w", err)
	}
	return nil
}

// updateActualState updates the state to 'CHANGED' if it is currently 'ACTUAL'
func updateActualState(runner DBRunner, table string, stateField string, id int) error {
	// Construct the SQL query with dynamic table and state field names
	query := fmt.Sprintf(
		"UPDATE %s SET %s = 'CHANGED' WHERE id = $1 AND %s = 'ACTUAL'",
		table, stateField, stateField,
	)
	// Execute the query with the ID parameter
	_, err := runner.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to update state to CHANGED in table %s: %w", table, err)
	}
	return nil
}
