package domain

import "time"

type SurveyStatus string

const (
	InProgress SurveyStatus = "in_progress"
	Completed  SurveyStatus = "completed"
)

// Прохождения опросов
type SurveyInterview struct {
	ID        int          `json:"id" db:"id"`
	UserID    int          `json:"user_id" db:"user_id"`
	SurveyID  int          `json:"survey_id" db:"survey_id"`
	Status    SurveyStatus `json:"status" db:"status"`
	StartTime time.Time    `json:"start_time" db:"start_time"`
	EndTime   *time.Time   `json:"end_time,omitempty" db:"end_time"`
}

// Ответы на вопросы
type SurveyAnswer struct {
	ID          int   `json:"id" db:"id"`
	InterviewID int   `json:"interview_id" db:"interview_id"`
	QuestionID  int   `json:"question_id" db:"question_id"`
	Options     []int `json:"options" db:"options"` // транзитное поле, его может не быть, если вопрос не подразумевает наличие опций
}
