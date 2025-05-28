package question

import (
	"backend/internal/domain"
	"backend/internal/repositories"
	"backend/pkg/i18n"
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
	// Устанавливаем ExtraParams в зависимости от типа вопроса
	switch questionType {
	case domain.Consent:
		// Для Consent экстра параметры пустые
		question.ExtraParams = domain.ConsentExtraParams{}
	case domain.SingleChoice:
		question.ExtraParams = domain.SingleChoiceExtraParams{Required: true}
	case domain.MultiChoice:
		question.ExtraParams = domain.MultiChoiceExtraParams{Required: true}
	case domain.Email:
		question.ExtraParams = domain.EmailExtraParams{Required: true}
	case domain.Rating:
		question.ExtraParams = domain.RatingExtraParams{Required: true}
	case domain.Date:
		question.ExtraParams = domain.DateExtraParams{Required: true}
	case domain.ShortText:
		question.ExtraParams = domain.TextExtraParams{Required: true}
	case domain.LongText:
		question.ExtraParams = domain.TextExtraParams{Required: true}
	case domain.Number:
		question.ExtraParams = domain.NumberExtraParams{Required: true}
	default:
		return nil, domain.ErrInvalidQuestionType
	}
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

// RestoreQuestion восстанавливает вопрос и возвращает его данные
func (s *QuestionService) RestoreQuestion(id int, surveyID int) (*domain.SurveyQuestionTemp, error) {
	// Выполняем восстановление
	err := s.repo.RestoreQuestion(id)
	if err != nil {
		return nil, err
	}
	// Получаем восстановленный вопрос из базы данных
	question, err := s.repo.GetQuestionByID(id, surveyID)
	if err != nil {
		return nil, err
	}
	return question, nil
}
func (s *QuestionService) UpdateQuestionExtraParams(questionID int, params map[string]interface{}) error {
	return s.repo.UpdateQuestionExtraParams(questionID, params)
}
