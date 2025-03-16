package domain

import "time"

// OptionState – состояние временной опции относительно постоянной сущности.
type OptionState string

const (
	OptionStateActual  OptionState = "ACTUAL"
	OptionStateNew     OptionState = "NEW"
	OptionStateChanged OptionState = "CHANGED"
	OptionStateDeleted OptionState = "DELETED"
)

// Option представляет опцию вопроса, которая хранится в таблице survey_options
type Option struct {
	ID          int       `json:"id" db:"id"`
	QuestionID  int       `json:"question_id" db:"question_id"`
	Label       string    `json:"label" db:"label"`
	OptionOrder int       `json:"option_order" db:"option_order"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// SurveyOptionsTemp - временная таблица для опций.
type OptionTemp struct {
	ID               int         `json:"id" db:"id"`
	OptionOriginalID *int        `json:"option_original_id" db:"option_original_id"`
	QuestionID       int         `json:"question_id" db:"question_id"`
	Label            string      `json:"label" db:"label"`
	OptionOrder      int         `json:"option_order" db:"option_order"`
	OptionState      OptionState `json:"option_state" db:"option_state"`
	CreatedAt        time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at" db:"updated_at"`
}
