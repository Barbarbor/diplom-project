package survey

import (
	"backend/internal/domain"
	"backend/internal/repositories"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

// SurveyService определяет бизнес-логику для опросов.
type SurveyService struct {
	surveyRepo repositories.SurveyRepository
}

// NewSurveyService создаёт новый экземпляр сервиса, внедряя репозиторий.
func NewSurveyService(repo repositories.SurveyRepository) *SurveyService {
	return &SurveyService{
		surveyRepo: repo,
	}
}

// GenerateRandomHash генерирует случайную строку длиной n символов из набора допустимых символов.
func GenerateRandomHash(n int) (string, error) {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	hash := make([]byte, n)
	for i := range hash {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letterBytes))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		hash[i] = letterBytes[num.Int64()]
	}
	return string(hash), nil
}

func (s *SurveyService) CreateSurvey(authorID int) (*domain.Survey, error) {
	now := time.Now()
	title := fmt.Sprintf("Опрос от %s", now.Format("02.01.2006"))
	hash, err := GenerateRandomHash(15)
	if err != nil {
		return nil, fmt.Errorf("failed to generate survey hash: %w", err)
	}
	state := domain.SurveyStateDraft

	surveyID, err := s.surveyRepo.CreateSurvey(title, authorID, hash, state, now)
	if err != nil {
		return nil, err
	}

	survey := &domain.Survey{
		ID:        surveyID,
		Title:     title,
		AuthorID:  authorID,
		CreatedAt: now,
		UpdatedAt: now,
		Hash:      hash,
		State:     state,
	}
	return survey, nil
}

func (s *SurveyService) GetSurveyByHash(hash string) (*domain.SurveyWithCreator, error) {
	survey, email, err := s.surveyRepo.GetSurveyByHash(hash)
	if err != nil {
		return nil, err
	}

	return &domain.SurveyWithCreator{
		Survey:       survey,
		CreatorEmail: email,
	}, nil
}

func (s *SurveyService) GetSurveysByAuthor(authorID int) ([]*domain.SurveySummary, error) {
	return s.surveyRepo.GetSurveysByAuthor(authorID)
}
