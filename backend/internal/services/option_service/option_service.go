package option

import (
	"backend/internal/domain"
	"backend/internal/repositories"
	"fmt"
)

// OptionService отвечает за бизнес-логику работы с вариантами ответа.
type OptionService struct {
	repo repositories.OptionRepository
	// При необходимости можно добавить репозиторий вопросов для получения типа вопроса
}

// NewOptionService создает новый сервис для опций.
func NewOptionService(repo repositories.OptionRepository) *OptionService {
	return &OptionService{repo: repo}
}

// CreateOption создает новый вариант ответа для заданного вопроса.
// Здесь label генерируется по умолчанию в зависимости от типа вопроса.
func (s *OptionService) CreateOption(questionID int, questionType domain.QuestionType) (*domain.OptionTemp, error) {
	// Генерируем дефолтный лейбл. Можно расширить логику для разных типов вопросов.
	defaultLabel := "Вариант ответа"
	if questionType == domain.SingleChoice || questionType == domain.MultiChoice {
		defaultLabel = "Вариант ответа"
	} else {
		// Если тип не поддерживает опции, можно вернуть ошибку.
		return nil, fmt.Errorf("question type %s does not support options", questionType)
	}

	option := &domain.OptionTemp{
		QuestionID: questionID,
		Label:      defaultLabel,
		// Опция будет создана с option_original_id = NULL, а состояние - 'NEW'
		OptionState:      "NEW",
		OptionOriginalID: nil,
	}

	if err := s.repo.CreateOption(option); err != nil {
		return nil, err
	}
	return option, nil
}
func (s *OptionService) UpdateOptionOrder(option *domain.OptionTemp, newOrder int) error {
	return s.repo.UpdateOptionOrder(option.ID, newOrder, option.OptionOrder, option.QuestionID)
}

func (s *OptionService) UpdateOptionLabel(option *domain.OptionTemp, newLabel string, questionID int) error {
	return s.repo.UpdateOptionLabel(option.ID, newLabel, questionID)
}

func (s *OptionService) DeleteOption(option *domain.OptionTemp, questionID int) error {
	return s.repo.DeleteOption(option.ID, questionID)
}
