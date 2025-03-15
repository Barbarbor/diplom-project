package domain

import (
	"errors"
	"time"
)

var ErrInvalidQuestionType = errors.New("invalid question type")

type QuestionType string

const (
	SingleChoice QuestionType = "single_choice"
	MultiChoice  QuestionType = "multi_choice"
)

// QuestionState – состояние временного вопроса относительно постоянной сущности.
type QuestionState string

const (
	QuestionStateActual  QuestionState = "ACTUAL"
	QuestionStateNew     QuestionState = "NEW"
	QuestionStateChanged QuestionState = "CHANGED"
	QuestionStateDeleted QuestionState = "DELETED"
)

// SurveyQuestion представляет вопрос опроса.
// Теперь добавляем поле QuestionOrder и удаляем Options.
type SurveyQuestion struct {
	ID            int          `json:"id" db:"id"`
	SurveyID      int          `json:"survey_id" db:"survey_id"`
	Label         string       `json:"label" db:"label"`
	Type          QuestionType `json:"type" db:"type"`
	QuestionOrder int          `json:"order" db:"question_order"`
	Options       []Option     `json:"options,omitempty" db:"-"`
	CreatedAt     time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at" db:"updated_at"`
}

// SurveyQuestionsTemp - временная таблица для вопросов.
type SurveyQuestionTemp struct {
	ID                 int           `json:"id" db:"id"`
	QuestionOriginalID *int          `json:"question_original_id" db:"question_original_id"`
	SurveyID           int           `json:"survey_id" db:"survey_id"`
	Label              string        `json:"label" db:"label"`
	Type               QuestionType  `json:"type" db:"type"`
	QuestionOrder      int           `json:"order" db:"question_order"`
	Options            []OptionTemp  `json:"options,omitempty" db:"-"`
	QuestionState      QuestionState `json:"question_state" db:"question_state"`
	CreatedAt          time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time     `json:"updated_at" db:"updated_at"`
}
