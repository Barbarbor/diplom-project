package domain

import "github.com/lib/pq"

type SurveyRole string

const (
	Read SurveyRole = "read"
	Edit SurveyRole = "edit"
)

// Роли в опросах
type SurveyRoles struct {
	ID       int            `json:"id" db:"id"`
	SurveyID int            `json:"survey_id" db:"survey_id"`
	UserID   int            `json:"user_id" db:"user_id"`
	Roles    pq.StringArray `json:"roles" db:"roles"`
}

type Roles struct {
	ID     int            `json:"id" db:"id"`
	Roles  pq.StringArray `json:"roles" db:"roles"`
	UserID int            `json:"user_id" db:"user_id"`
}
