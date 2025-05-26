package question

import (
	"backend/internal/domain"
	"backend/internal/repositories"
	"backend/pkg/i18n"
	"encoding/json"
)

// QuestionService отвечает за бизнес-логику работы с вопросами
type QuestionService struct {
	repo repositories.QuestionRepository
}

// NewQuestionService создаёт новый сервис
func NewQuestionService(repo repositories.QuestionRepository) *QuestionService {
	return &QuestionService{repo: repo}
}

// CreateQuestion создаёт новый вопрос с дефолтными параметрами
func (s *QuestionService) CreateQuestion(surveyID int, questionType domain.QuestionType) (*domain.SurveyQuestionTemp, error) {
	// Предопределённые параметры для типов вопросов
	defaultQuestions := map[domain.QuestionType]*domain.SurveyQuestionTemp{
		domain.SingleChoice: {Label: i18n.T("question.service.defaultSingle"), Type: domain.SingleChoice},
		domain.MultiChoice:  {Label: i18n.T("question.service.defaultMulti"), Type: domain.MultiChoice},

		// Новые типы
		domain.Consent:   {Label: i18n.T("question.service.defaultConsent"), Type: domain.Consent},
		domain.Email:     {Label: i18n.T("question.service.defaultEmail"), Type: domain.Email},
		domain.Rating:    {Label: i18n.T("question.service.defaultRating"), Type: domain.Rating},
		domain.Date:      {Label: i18n.T("question.service.defaultDate"), Type: domain.Date},
		domain.ShortText: {Label: i18n.T("question.service.defaultShortText"), Type: domain.ShortText},
		domain.LongText:  {Label: i18n.T("question.service.defaultLongText"), Type: domain.LongText},
		domain.Number:    {Label: i18n.T("question.service.defaultNumber"), Type: domain.Number},
	}
	question, exists := defaultQuestions[questionType]
	if !exists {
		return nil, domain.ErrInvalidQuestionType
	}

	question.SurveyID = surveyID
	question.QuestionState = "NEW"
	question.QuestionOriginalID = nil
	question.ExtraParams = json.RawMessage(`{"required": true}`) // <-- вот ключевое
	// Сохраняем в БД
	err := s.repo.CreateQuestion(question)
	if err != nil {
		return nil, err
	}

	return question, nil
}

// UpdateQuestionType обновляет тип вопроса. currentState передается как параметр.
func (s *QuestionService) UpdateQuestionType(questionID int, newType domain.QuestionType, currentState string) (*domain.SurveyQuestionTemp, error) {
	return s.repo.UpdateQuestionType(questionID, newType, currentState)
}

// UpdateQuestionLabel обновляет только label вопроса.
func (s *QuestionService) UpdateQuestion(questionID int, newLabel string) error {
	return s.repo.UpdateQuestion(questionID, newLabel)
}

// UpdateQuestionOrder обновляет порядок вопроса.
// Параметры currentOrder и surveyID должны быть извлечены до вызова этой функции.
func (s *QuestionService) UpdateQuestionOrder(questionID, newOrder, currentOrder, surveyID int) error {
	return s.repo.UpdateQuestionOrder(questionID, newOrder, currentOrder, surveyID)
}

func (s *QuestionService) DeleteQuestion(id int) error {
	return s.repo.DeleteQuestion(id)
}

func (s *QuestionService) RestoreQuestion(id int) error {
	return s.repo.RestoreQuestion(id)
}

func (s *QuestionService) UpdateQuestionExtraParams(questionID int, params map[string]interface{}) error {
	return s.repo.UpdateQuestionExtraParams(questionID, params)
}
