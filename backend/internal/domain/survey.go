package domain

import (
	"errors"
	"time"
)

var ErrInvalidQuestionType = errors.New("invalid question type")

type SurveyState string

const (
	SurveyStateDraft  SurveyState = "DRAFT"
	SurveyStateActive SurveyState = "ACTIVE"
)

type QuestionType string

const (
	SingleChoice QuestionType = "single_choice"
	MultiChoice  QuestionType = "multi_choice"
)

type SurveyStatus string

const (
	InProgress SurveyStatus = "in_progress"
	Completed  SurveyStatus = "completed"
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
	State     SurveyState `json:"state" db:"state"` // Новое поле
}

// SurveyWithCreator объединяет опрос с email создателя.
type SurveyWithCreator struct {
	Survey       *Survey `json:"survey"`
	CreatorEmail string  `json:"creator"`
}

// SurveySummary представляет краткую информацию об опросе.
type SurveySummary struct {
	Title     string      `json:"title" db:"title"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt time.Time   `json:"updated_at" db:"updated_at"`
	Hash      string      `json:"hash" db:"hash"`
	State     SurveyState `json:"state" db:"state"`
}

// Вопросы опросов
type SurveyQuestion struct {
	ID       int          `json:"id" db:"id"`
	SurveyID int          `json:"survey_id" db:"survey_id"`
	Label    string       `json:"label" db:"label"`
	Type     QuestionType `json:"type" db:"type"`       // Используем строгий тип
	Options  []Option     `json:"options" db:"options"` // JSON массив с опциями
}

// Опции вопроса
type Option struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}

// Прохождения опросов
type SurveyInterview struct {
	ID        int          `json:"id" db:"id"`
	UserID    int          `json:"user_id" db:"user_id"`
	SurveyID  int          `json:"survey_id" db:"survey_id"`
	Status    SurveyStatus `json:"status" db:"status"` // Используем строгий тип
	StartTime time.Time    `json:"start_time" db:"start_time"`
	EndTime   *time.Time   `json:"end_time,omitempty" db:"end_time"`
}

// Ответы на вопросы
type SurveyAnswer struct {
	ID          int   `json:"id" db:"id"`
	InterviewID int   `json:"interview_id" db:"interview_id"`
	QuestionID  int   `json:"question_id" db:"question_id"`
	Options     []int `json:"options" db:"options"` // Массив выбранных ID опций
}

// Статистика опросов
type SurveyStat struct {
	ID             int     `json:"id" db:"id"`
	SurveyID       int     `json:"survey_id" db:"survey_id"`
	ViewsCount     int     `json:"views_count" db:"views_count"`
	CompletionRate float64 `json:"completion_rate" db:"completion_rate"`
}

// Роли в опросах
type SurveyRole struct {
	ID       int      `json:"id" db:"id"`
	SurveyID int      `json:"survey_id" db:"survey_id"`
	UserID   int      `json:"user_id" db:"user_id"`
	Roles    []string `json:"roles" db:"roles"` // Массив строк для ролей
}

// Действия с опросами
type SurveyActionLog struct {
	ID         int          `json:"id" db:"id"`
	Action     SurveyAction `json:"action" db:"action"`                 // Используем строгий тип
	SurveyID   *int         `json:"survey_id,omitempty" db:"survey_id"` // Может быть null
	UserID     *int         `json:"user_id,omitempty" db:"user_id"`     // Может быть null
	Body       interface{}  `json:"body,omitempty" db:"body"`           // JSON с изменениями
	ActionTime time.Time    `json:"action_time" db:"action_time"`
}
