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
	}
	question, exists := defaultQuestions[questionType]
	if !exists {
		return nil, domain.ErrInvalidQuestionType
	}

	question.SurveyID = surveyID
	question.QuestionState = "NEW"
	question.QuestionOriginalID = nil
	// Сохраняем в БД
	err := s.repo.CreateQuestion(question)
	if err != nil {
		return nil, err
	}

	return question, nil
}

// UpdateQuestionType обновляет тип вопроса. currentState передается как параметр.
func (s *QuestionService) UpdateQuestionType(questionID int, newType domain.QuestionType, currentState string) error {
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
