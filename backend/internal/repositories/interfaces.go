package repositories

import (
	"backend/internal/domain"
	"time"

	"github.com/jmoiron/sqlx"
)

// AuthRepository определяет методы для работы с пользователями
type AuthRepository interface {
	CreateUser(email, password string) (int, error)
	GetUserByEmail(email string) (*domain.User, error)
	GetUserByID(userID int) (*domain.User, error)
}

// ProfileRepository определяет методы для работы с профилем пользователя.
type ProfileRepository interface {
	GetUserProfile(userID int) (*domain.UserProfile, error)
	UpdateUserProfile(profile *domain.UserProfile) error
}

// SurveyRepository определяет методы для работы с опросами в БД.
type SurveyRepository interface {
	CreateSurvey(title string, authorID int, hash string, state domain.SurveyState, now time.Time) (int, error)
	GetSurveyByHash(hash string) (*domain.Survey, string, error)
	GetSurveyIdByHash(hash string) (int, error)
	CheckUserAccess(userID int, surveyID int) (bool, error)           // Добавляем этот метод
	GetSurveysByAuthor(authorID int) ([]*domain.SurveySummary, error) // Новый метод
	PublishSurvey(surveyID int) error
	UpdateSurveyTitle(surveyID int, newTitle string) error
	// Для RestoreSurvey
	BeginTx() (*sqlx.Tx, error)
	UpdateSurveyTitleTx(tx *sqlx.Tx, surveyID int) error

	FinishInterview(interviewID string, endTime time.Time) error
}

// QuestionRepository определяет методы для работы с вопросами в БД
type QuestionRepository interface {
	CreateQuestion(question *domain.SurveyQuestionTemp) error

	GetQuestionMaxOrder(surveyID int) (int, error)
	GetQuestionByID(questionID int, surveyID int) (*domain.SurveyQuestionTemp, error)
	GetQuestionsBySurveyID(surveyID int) ([]*domain.SurveyQuestionTemp, error)
	GetOptionsByQuestionID(questionID int) ([]domain.OptionTemp, error)
	GetQuestionOptionRows(surveyID int) ([]QuestionOptionRow, error)
	GetSurveyQuestionsWithOptionsAndAnswers(
		surveyID int,
		interviewID string,
		isDemo bool,
	) ([]QuestionOptionAnswerRow, error)

	UpdateQuestion(questionID int, newLabel string) error
	UpdateQuestionType(questionID int, newType domain.QuestionType, currentState string) (*domain.SurveyQuestionTemp, error)
	UpdateQuestionOrder(questionID int, newOrder, currentOrder, surveyID int) error
	UpdateQuestionExtraParams(questionID int, params map[string]interface{}) error

	DeleteQuestion(questionID int) error
	RestoreQuestion(questionTempID int) error

	GetTempQuestionIDsBySurveyIDTx(
		tx *sqlx.Tx, surveyID int,
	) ([]int, error)
	RestoreQuestionTx(tx *sqlx.Tx, questionTempID int) error
	restoreQuestionTx(tx *sqlx.Tx, questionTempID int) error

	QuestionExists(questionID int, isDemo bool) (bool, error)
	AnswerExists(interviewID string, questionID int) (bool, error)
	CreateAnswer(interviewID string, questionID int, answer string) error
	UpdateAnswer(interviewID string, questionID int, answer string) error
}

// OptionRepository описывает операции с опциями.
type OptionRepository interface {
	// CreateOption создает новую опцию в таблице survey_options_temp.
	CreateOption(option *domain.OptionTemp) error
	// GetMaxOptionOrder возвращает максимальное значение option_order для заданного вопроса.
	GetMaxOptionOrder(questionID int) (int, error)
	GetOptionById(questionID, optionID int) (*domain.OptionTemp, error)
	UpdateOptionOrder(optionID, newOrder, currentOrder, questionID int) error
	UpdateOptionLabel(optionID int, newLabel string) error
	DeleteOption(optionID int) error
}

type InterviewRepository interface {
	InterviewExists(interviewID string) (bool, error)
	CreateInterview(interview *domain.SurveyInterview) error
	GetInterviewByID(interviewID string) (*domain.SurveyInterview, error)
}
