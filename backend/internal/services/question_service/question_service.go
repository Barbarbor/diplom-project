package question

import (
	"backend/internal/domain"
	"backend/internal/repositories"
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
func (s *QuestionService) CreateQuestion(surveyID int, questionType domain.QuestionType) (*domain.SurveyQuestion, error) {
	// Предопределённые параметры для типов вопросов
	defaultQuestions := map[domain.QuestionType]*domain.SurveyQuestion{
		domain.SingleChoice: {Label: "Выберите один вариант", Type: domain.SingleChoice, Options: []domain.Option{{ID: 1, Label: "Вариант 1"}, {ID: 2, Label: "Вариант 2"}}},
		domain.MultiChoice:  {Label: "Выберите несколько вариантов", Type: domain.MultiChoice, Options: []domain.Option{{ID: 1, Label: "Вариант A"}, {ID: 2, Label: "Вариант B"}}},
	}

	question, exists := defaultQuestions[questionType]
	if !exists {
		return nil, domain.ErrInvalidQuestionType
	}

	question.SurveyID = surveyID

	// Сохраняем в БД
	err := s.repo.CreateQuestion(question)
	if err != nil {
		return nil, err
	}

	return question, nil
}
