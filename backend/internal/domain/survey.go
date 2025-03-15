package domain

import (
	"time"
)

type SurveyState string

const (
	SurveyStateDraft  SurveyState = "DRAFT"
	SurveyStateActive SurveyState = "ACTIVE"
)

type SurveyAction string

const (
	CreateAction SurveyAction = "create"
	UpdateAction SurveyAction = "update"
	DeleteAction SurveyAction = "delete"
	PassAction   SurveyAction = "pass"
)

// Опросы
type Survey struct {
	ID        int         `json:"id" db:"id"`
	Title     string      `json:"title" db:"title"`
	AuthorID  int         `json:"author_id" db:"author_id"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt time.Time   `json:"updated_at" db:"updated_at"`
	Hash      string      `json:"hash" db:"hash"`
	State     SurveyState `json:"state" db:"state"`
}

type SurveysTemp struct {
	SurveyOriginalID int       `json:"survey_original_id" db:"survey_original_id"`
	Title            string    `json:"title" db:"title"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// SurveySummary представляет краткую информацию об опросе.
type SurveySummary struct {
	Title     string      `json:"title" db:"title"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt time.Time   `json:"updated_at" db:"updated_at"`
	Hash      string      `json:"hash" db:"hash"`
	State     SurveyState `json:"state" db:"state"`
}

// Действия с опросами
type SurveyActionLog struct {
	ID         int          `json:"id" db:"id"`
	Action     SurveyAction `json:"action" db:"action"`
	SurveyID   *int         `json:"survey_id,omitempty" db:"survey_id"`
	UserID     *int         `json:"user_id,omitempty" db:"user_id"`
	Body       interface{}  `json:"body,omitempty" db:"body"`
	ActionTime time.Time    `json:"action_time" db:"action_time"`
}
