package domain

// Роли в опросах
type SurveyRole struct {
	ID       int      `json:"id" db:"id"`
	SurveyID int      `json:"survey_id" db:"survey_id"`
	UserID   int      `json:"user_id" db:"user_id"`
	Roles    []string `json:"roles" db:"roles"`
}
