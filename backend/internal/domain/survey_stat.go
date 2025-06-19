package domain

import (
	"time"
)

// Статистика опросов
type SurveyStat struct {
	ID             int     `json:"id" db:"id"`
	SurveyID       int     `json:"survey_id" db:"survey_id"`
	ViewsCount     int     `json:"views_count" db:"views_count"`
	CompletionRate float64 `json:"completion_rate" db:"completion_rate"`
}

// Определяем структуры данных
type SurveyStats struct {
	StartedInterviews     int             `json:"started_interviews" db:"started_interviews"`
	CompletedInterviews   int             `json:"completed_interviews" db:"completed_interviews"`
	AverageCompletionTime float64         `json:"average_completion_time" db:"-"` // Новое поле для среднего времени
	Questions             []QuestionStats `json:"questions" db:"-"`
	InterviewTimes        []InterviewTime `json:"-" db:"-"`
}

type InterviewTime struct {
	StartTime time.Time `json:"start_time" db:"start_time"`
	EndTime   time.Time `json:"end_time" db:"end_time"`
}

type QuestionStats struct {
	ID          int           `json:"id" db:"id"`
	Label       string        `json:"label" db:"label"`
	Type        QuestionType  `json:"type" db:"type"`
	Options     []OptionStats `json:"options,omitempty" db:"-"`
	Answers     []string      `json:"answers" db:"-"`
	ExtraParams ExtraParams   `json:"extra_params" db:"extra_params"` // Типизированные параметры
}

type OptionStats struct {
	ID    int    `json:"id" db:"id"`
	Label string `json:"label" db:"label"`
}
