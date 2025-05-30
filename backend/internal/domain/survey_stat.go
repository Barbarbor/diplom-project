package domain

// Статистика опросов
type SurveyStat struct {
	ID             int     `json:"id" db:"id"`
	SurveyID       int     `json:"survey_id" db:"survey_id"`
	ViewsCount     int     `json:"views_count" db:"views_count"`
	CompletionRate float64 `json:"completion_rate" db:"completion_rate"`
}

// Определяем структуры данных
type SurveyStats struct {
	StartedInterviews   int             `json:"started_interviews" db:"started_interviews"`
	CompletedInterviews int             `json:"completed_interviews" db:"completed_interviews"`
	Questions           []QuestionStats `json:"questions" db:"-"`
}

type QuestionStats struct {
	ID      int           `json:"id" db:"id"`
	Label   string        `json:"label" db:"label"`
	Type    QuestionType  `json:"type" db:"type"`
	Options []OptionStats `json:"options,omitempty" db:"-"`
	Answers []string      `json:"answers" db:"-"`
}

type OptionStats struct {
	ID    int    `json:"id" db:"id"`
	Label string `json:"label" db:"label"`
}
