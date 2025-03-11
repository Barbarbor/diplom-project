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
	surveyRepo   repositories.SurveyRepository
	questionRepo repositories.QuestionRepository
}

// SurveyWithCreator объединяет опрос с email создателя.
type SurveyWithCreator struct {
	Survey       *domain.Survey `json:"survey"`
	CreatorEmail string         `json:"creator"`
}

// NewSurveyService создаёт новый экземпляр сервиса, внедряя репозиторий.
func NewSurveyService(surveyRepo repositories.SurveyRepository, questionRepo repositories.QuestionRepository) *SurveyService {
	return &SurveyService{
		surveyRepo:   surveyRepo,
		questionRepo: questionRepo,
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

func (s *SurveyService) GetSurveyByHash(hash string) (*SurveyWithCreator, error) {
	survey, email, err := s.surveyRepo.GetSurveyByHash(hash)
	if err != nil {
		return nil, err
	}

	return &SurveyWithCreator{
		Survey:       survey,
		CreatorEmail: email,
	}, nil
}

// GetQuestionsForSurvey получает «сырые» данные из репозитория и группирует их в срез вопросов.
func (s *SurveyService) GetQuestionsForSurvey(surveyID int) ([]*domain.SurveyQuestion, error) {
	rows, err := s.questionRepo.GetQuestionOptionRows(surveyID)
	if err != nil {
		return nil, err
	}

	// Группировка по вопросу.
	questionsMap := make(map[int]*domain.SurveyQuestion)
	var questions []*domain.SurveyQuestion

	// Если у вопроса тип поддерживает опции, они заполняются; иначе оставляем пустым.
	for _, row := range rows {
		q, exists := questionsMap[row.QID]
		if !exists {
			q = &domain.SurveyQuestion{
				ID:            row.QID,
				SurveyID:      row.QSurveyID,
				Label:         row.QLabel,
				Type:          domain.QuestionType(row.QType),
				QuestionOrder: row.QOrder,
				// Поле Options помечено как транзитное (db:"-")
				Options: []domain.Option{},
			}
			questionsMap[row.QID] = q
			questions = append(questions, q)
		}

		// Заполняем опции только для типов, где опции имеют смысл.
		if (q.Type == domain.SingleChoice || q.Type == domain.MultiChoice) && row.OptionID.Valid {
			option := domain.Option{
				ID:         int(row.OptionID.Int64),
				QuestionID: int(row.OptionQuestionID.Int64),
				Label:      row.OptionLabel.String,
				CreatedAt:  row.OptionCreatedAt.Time,
				UpdatedAt:  row.OptionUpdatedAt.Time,
			}
			q.Options = append(q.Options, option)
		}
	}

	// Если требуется параллелизация обработки (но в данном случае она уже достаточно быстрая),
	// можно обернуть логику группировки в горутины. Но здесь она проходит последовательно.
	return questions, nil
}
func (s *SurveyService) GetSurveysByAuthor(authorID int) ([]*domain.SurveySummary, error) {
	return s.surveyRepo.GetSurveysByAuthor(authorID)
}
