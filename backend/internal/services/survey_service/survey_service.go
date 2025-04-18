package survey

import (
	"backend/internal/domain"
	"backend/internal/repositories"
	"backend/pkg/i18n"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/jmoiron/sqlx"
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
	titleTemplate := i18n.T("survey.service.defaultTitle")
	title := fmt.Sprintf(titleTemplate, now.Format("02.01.2006"))
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
func (s *SurveyService) GetQuestionsForSurvey(surveyID int) ([]*domain.SurveyQuestionTemp, error) {
	// Получаем "сырые" строки из временной таблицы
	rows, err := s.questionRepo.GetQuestionOptionRows(surveyID)
	if err != nil {
		return nil, err
	}

	// Группируем данные по вопросу
	questionsMap := make(map[int]*domain.SurveyQuestionTemp)
	var questions []*domain.SurveyQuestionTemp

	for _, row := range rows {
		q, exists := questionsMap[row.QID]
		if !exists {
			q = &domain.SurveyQuestionTemp{
				ID:                 row.QID,
				QuestionOriginalID: row.QQuestionOriginalId,
				SurveyID:           row.QSurveyID,
				QuestionState:      *row.QQuestionState,
				Label:              row.QLabel,
				Type:               domain.QuestionType(row.QType),
				QuestionOrder:      row.QOrder,
				CreatedAt:          row.QCreatedAt.Time,
				UpdatedAt:          row.QUpdatedAt.Time,
				// Если нужно, можно сохранить состояние, например:
				// QuestionState: domain.QuestionState(row.QState),
				Options: []domain.OptionTemp{},
			}
			questionsMap[row.QID] = q
			questions = append(questions, q)
		}

		// Если вопрос поддерживает опции и опция существует, добавляем ее
		if (q.Type == domain.SingleChoice || q.Type == domain.MultiChoice) && row.OptionID.Valid {
			option := domain.OptionTemp{
				ID:               int(row.OptionID.Int64),
				QuestionID:       int(row.OptionQuestionID.Int64),
				OptionOriginalID: row.OptionOriginalId,
				OptionOrder:      *row.OptionOrder,
				OptionState:      *row.OptionState,
				Label:            row.OptionLabel.String,
				CreatedAt:        row.OptionCreatedAt.Time,
				UpdatedAt:        row.OptionUpdatedAt.Time,
				// Если требуется, можно установить состояние опции:
				// OptionState: domain.OptionState(row.OptionState.String),
			}
			q.Options = append(q.Options, option)
		}
	}

	return questions, nil
}
func (s *SurveyService) GetSurveysByAuthor(authorID int) ([]*domain.SurveySummary, error) {
	return s.surveyRepo.GetSurveysByAuthor(authorID)
}

func (s *SurveyService) UpdateSurvey(surveyID int, newTitle string) error {
	return s.surveyRepo.UpdateSurveyTitle(surveyID, newTitle)
}

func (s *SurveyService) PublishSurvey(surveyID int) error {
	return s.surveyRepo.PublishSurvey(surveyID)
}

func (s *SurveyService) RestoreSurvey(db *sqlx.DB, surveyID int) (err error) {
	tx, err := db.Beginx()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// 1) Найти все temp‑вопросы
	var tempQIDs []int
	if err = tx.Select(&tempQIDs, `
        SELECT id FROM survey_questions_temp WHERE survey_id = $1
    `, surveyID); err != nil {
		return fmt.Errorf("list temp questions: %w", err)
	}

	// 2) Для каждого вопроса вызвать репозиторий
	for _, qID := range tempQIDs {
		if err = s.questionRepo.RestoreQuestion(tx, qID); err != nil {
			return fmt.Errorf("restore question %d: %w", qID, err)
		}
	}

	// 3) Восстановить metadata опроса
	if err = s.surveyRepo.RestoreSurvey(tx, surveyID); err != nil {
		return fmt.Errorf("restore survey metadata: %w", err)
	}

	return
}
