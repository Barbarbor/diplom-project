package survey

import (
	"backend/internal/domain"
	"backend/internal/repositories"
	"backend/pkg/i18n"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"sync"
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
	titleTemplate, err := i18n.TWithData("survey.service.defaultTitle", map[string]interface{}{"Date": now.Format("02.01.2006")})
	if err != nil {
		return nil, fmt.Errorf("failed to process title template: %w", err)
	}
	hash, err := GenerateRandomHash(15)
	if err != nil {
		return nil, fmt.Errorf("failed to generate survey hash: %w", err)
	}
	state := domain.SurveyStateDraft

	surveyID, err := s.surveyRepo.CreateSurvey(titleTemplate, authorID, hash, state, now)
	if err != nil {
		return nil, err
	}

	survey := &domain.Survey{
		ID:        surveyID,
		Title:     titleTemplate,
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

// GetQuestionsForSurvey возвращает вопросы из временной таблицы
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
			// Парсим экстра параметры
			extraParams, err := domain.ParseExtraParams(row.QExtraParams, domain.QuestionType(row.QType))
			if err != nil {
				return nil, err
			}

			q = &domain.SurveyQuestionTemp{
				ID:                 row.QID,
				QuestionOriginalID: row.QQuestionOriginalId,
				SurveyID:           row.QSurveyID,
				QuestionState:      *row.QQuestionState,
				Label:              row.QLabel,
				Type:               domain.QuestionType(row.QType),
				QuestionOrder:      row.QOrder,
				RawExtraParams:     row.QExtraParams, // Сохраняем сырые данные (если нужно для отладки)
				ExtraParams:        extraParams,      // Присваиваем типизированные параметры
				CreatedAt:          row.QCreatedAt.Time,
				UpdatedAt:          row.QUpdatedAt.Time,
				Options:            []domain.OptionTemp{},
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
			}
			q.Options = append(q.Options, option)
		}
	}

	return questions, nil
}

// GetSurveyQuestionsWithAnswers возвращает вопросы с опциями и ответами
func (s *SurveyService) GetSurveyQuestionsWithAnswers(
	surveyID int,
	interviewID string,
	isDemo bool,
) ([]interface{}, error) {
	// Получаем данные из репозитория
	rows, err := s.questionRepo.GetSurveyQuestionsWithOptionsAndAnswers(surveyID, interviewID, isDemo)
	if err != nil {
		return nil, err
	}

	// Карта для группировки вопросов
	questionsMap := make(map[int]interface{})
	var questions []interface{}

	for _, row := range rows {
		var q interface{}
		if isDemo {
			// Для демо-режима используем SurveyQuestionTemp
			tempQ := &domain.SurveyQuestionTemp{
				ID:             row.QID,
				SurveyID:       row.QSurveyID,
				Label:          row.QLabel,
				Type:           domain.QuestionType(row.QType),
				QuestionOrder:  row.QOrder,
				RawExtraParams: row.QExtraParams,
				Options:        []domain.OptionTemp{},
			}
			// Парсим экстра параметры
			tempQ.ExtraParams, err = domain.ParseExtraParams(tempQ.RawExtraParams, tempQ.Type)
			if err != nil {
				return nil, err
			}
			q = tempQ
		} else {
			// Для обычного режима используем SurveyQuestion
			regularQ := &domain.SurveyQuestion{
				ID:             row.QID,
				SurveyID:       row.QSurveyID,
				Label:          row.QLabel,
				Type:           domain.QuestionType(row.QType),
				QuestionOrder:  row.QOrder,
				RawExtraParams: row.QExtraParams,
				Options:        []domain.Option{},
			}
			// Парсим экстра параметры
			regularQ.ExtraParams, err = domain.ParseExtraParams(regularQ.RawExtraParams, regularQ.Type)
			if err != nil {
				return nil, err
			}
			q = regularQ
		}

		// Добавляем вопрос в карту и массив, если его еще нет
		if _, exists := questionsMap[row.QID]; !exists {
			questionsMap[row.QID] = q
			questions = append(questions, q)
		} else {
			q = questionsMap[row.QID]
		}

		// Добавляем опции, если они есть
		if row.OID.Valid && row.OLabel.Valid && row.OOrder != nil {
			if isDemo {
				tempQ := q.(*domain.SurveyQuestionTemp)
				option := domain.OptionTemp{
					ID:          int(row.OID.Int64),
					QuestionID:  int(row.OQuestionID.Int64),
					Label:       row.OLabel.String,
					OptionOrder: *row.OOrder,
				}
				tempQ.Options = append(tempQ.Options, option)
			} else {
				regularQ := q.(*domain.SurveyQuestion)
				option := domain.Option{
					ID:          int(row.OID.Int64),
					QuestionID:  int(row.OQuestionID.Int64),
					Label:       row.OLabel.String,
					OptionOrder: *row.OOrder,
				}
				regularQ.Options = append(regularQ.Options, option)
			}
		}

		// Добавляем ответ, если он есть
		if row.AAnswer != nil {
			if isDemo {
				tempQ := q.(*domain.SurveyQuestionTemp)
				tempQ.Answer = row.AAnswer
			} else {
				regularQ := q.(*domain.SurveyQuestion)
				regularQ.Answer = row.AAnswer
			}
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
} // Восстанавливаем опрос целиком по его ID
func (s *SurveyService) RestoreSurveyByID(surveyID int) (err error) {
	tx, err := s.surveyRepo.BeginTx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	if err = s.surveyRepo.UpdateSurveyTitleTx(tx, surveyID); err != nil {
		return err
	}

	tempIDs, err := s.questionRepo.GetTempQuestionIDsBySurveyIDTx(tx, surveyID)
	if err != nil {
		return fmt.Errorf("fetch temp question ids: %w", err)
	}

	// Параллельное выполнение с отдельными транзакциями
	errCh := make(chan error, len(tempIDs))
	var wg sync.WaitGroup
	for _, qID := range tempIDs {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// Вызываем публичный метод, который сам открывает транзакцию
			if e := s.questionRepo.RestoreQuestion(id); e != nil {
				errCh <- fmt.Errorf("restore question %d: %w", id, e)
			}
		}(qID)
	}
	wg.Wait()
	close(errCh)
	if e, ok := <-errCh; ok {
		return e
	}

	return nil
}

// UpdateQuestionAnswer обновляет или создает ответ на вопрос
func (s *SurveyService) UpdateQuestionAnswer(interviewID string, questionID int, answer string, isDemo bool) error {
	// Проверяем, существует ли вопрос
	exists, err := s.questionRepo.QuestionExists(questionID, isDemo)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrQuestionNotFound
	}

	// Проверяем, существует ли ответ
	answerExists, err := s.questionRepo.AnswerExists(interviewID, questionID)
	if err != nil {
		return err
	}

	if answerExists {
		// Обновляем существующий ответ
		return s.questionRepo.UpdateAnswer(interviewID, questionID, answer)
	} else {
		// Создаем новый ответ
		return s.questionRepo.CreateAnswer(interviewID, questionID, answer)
	}
}

// ValidateAnswer валидирует ответ на основе типа вопроса и экстра параметров
func (s *SurveyService) ValidateAnswer(question interface{}, answer *string, interviewID string, isDemo bool) error {
	var qType domain.QuestionType
	var extraParams domain.ExtraParams

	que := question.(*domain.SurveyQuestion)
	// Определяем тип вопроса и извлекаем экстра параметры
	if q, ok := question.(*domain.SurveyQuestion); ok {
		qType = q.Type
		extraParams = q.ExtraParams
	} else if tempQ, ok := question.(*domain.SurveyQuestionTemp); ok {
		qType = tempQ.Type
		extraParams = tempQ.ExtraParams
	} else {
		return errors.New("invalid question type")
	}

	// Проверка required
	switch params := extraParams.(type) {
	case domain.SingleChoiceExtraParams:
		if params.Required && (answer == nil || *answer == "") {
			return errors.New("single choice answer is required")
		}
	case domain.MultiChoiceExtraParams:
		if params.Required && (answer == nil || *answer == "") {
			return errors.New("multi choice answer is required")
		}
	case domain.EmailExtraParams:
		if params.Required && (answer == nil || *answer == "") {
			return errors.New("email answer is required")
		}
	case domain.RatingExtraParams:
		if params.Required && (answer == nil || *answer == "") {
			return errors.New("rating answer is required")
		}
	case domain.DateExtraParams:
		if params.Required && (answer == nil || *answer == "") {
			return errors.New("date answer is required")
		}
	case domain.TextExtraParams:
		if params.Required && (answer == nil || *answer == "") {
			return errors.New("text answer is required")
		}
	case domain.NumberExtraParams:
		if params.Required && (answer == nil || *answer == "") {
			return errors.New("number answer is required")
		}
	}

	// Валидация в зависимости от типа вопроса
	switch qType {
	case domain.SingleChoice:
		if answer != nil {
			_, err := strconv.Atoi(*answer)
			if err != nil {
				return errors.New("invalid single_choice: expected number")
			}
		}
	case domain.MultiChoice:
		if answer != nil {
			var arr []int
			if err := json.Unmarshal([]byte(*answer), &arr); err != nil {
				return errors.New("invalid multi_choice: expected array of numbers")
			}
			if params, ok := extraParams.(domain.MultiChoiceExtraParams); ok {
				if params.MinAnswersCount > 0 && len(arr) < params.MinAnswersCount {
					return errors.New("too few answers for multi_choice")
				}
				if params.MaxAnswersCount > 0 && len(arr) > params.MaxAnswersCount {
					return errors.New("too many answers for multi_choice")
				}
			}
		}
	case domain.Consent:
		if answer == nil {

			// Создаём новый указатель с значением "false", если ответ отсутствует
			defaultAnswer := "false"
			answer = &defaultAnswer

			s.questionRepo.CreateAnswer(interviewID, que.ID, *answer)
		} else if *answer != "true" && *answer != "false" {
			return errors.New("invalid consent: expected 'true' or 'false'")
		}
	case domain.Email:
		if answer != nil {
			if !strings.Contains(*answer, "@") {
				return errors.New("invalid email format")
			}
		}
	case domain.Rating:
		if answer != nil {
			rating, err := strconv.Atoi(*answer)
			if err != nil {
				return errors.New("invalid rating: expected number")
			}
			if params, ok := extraParams.(domain.RatingExtraParams); ok && params.StarsCount > 0 {
				if rating < 1 || rating > params.StarsCount {
					return errors.New("rating out of range")
				}
			}
		}
	case domain.Date:
		if answer != nil {
			date, err := time.Parse("2006-01-02", *answer)
			if err != nil {
				return errors.New("invalid date format: expected yyyy-mm-dd")
			}
			if params, ok := extraParams.(domain.DateExtraParams); ok {
				if params.MinDate != "" {
					// Парсим minDate из формата ISO 8601
					minDate, err := time.Parse(time.RFC3339, params.MinDate)
					if err != nil {
						return fmt.Errorf("invalid minDate format: %w", err)
					}
					// Приводим к началу дня для корректного сравнения
					minDate = time.Date(minDate.Year(), minDate.Month(), minDate.Day(), 0, 0, 0, 0, minDate.Location())
					fmt.Println("mindate", minDate)
					if date.Before(minDate) {
						return errors.New("date is before minDate")
					}
				}
				if params.MaxDate != "" {
					// Парсим maxDate из формата ISO 8601
					maxDate, err := time.Parse(time.RFC3339, params.MaxDate)
					if err != nil {
						return fmt.Errorf("invalid maxDate format: %w", err)
					}
					// Приводим к началу дня для корректного сравнения
					maxDate = time.Date(maxDate.Year(), maxDate.Month(), maxDate.Day(), 0, 0, 0, 0, maxDate.Location())
					fmt.Println("maxdate", maxDate)
					fmt.Println("date", date)
					if date.After(maxDate) {
						return errors.New("date is after maxDate")
					}
				}
			}
		}
	case domain.ShortText, domain.LongText:
		if answer != nil {
			if params, ok := extraParams.(domain.TextExtraParams); ok && params.MaxLength > 0 {
				if len(*answer) > params.MaxLength {
					return errors.New("text exceeds maxLength")
				}
			}
		}
	case domain.Number:
		if answer != nil {
			num, err := strconv.ParseFloat(*answer, 64)
			if err != nil {
				return errors.New("invalid number: expected numeric value")
			}
			if params, ok := extraParams.(domain.NumberExtraParams); ok {
				if params.MinNumber > 0 && num < params.MinNumber {
					return errors.New("number is less than minNumber")
				}
				if params.MaxNumber > 0 && num > params.MaxNumber {
					return errors.New("number is greater than maxNumber")
				}
			}
		}
	default:
		return errors.New("unknown question type")
	}
	return nil
}

// FinishInterview завершает интервью после валидации
func (s *SurveyService) FinishInterview(surveyID int, interviewID string, isDemo bool) error {
	// Получаем вопросы и ответы
	questions, err := s.GetSurveyQuestionsWithAnswers(surveyID, interviewID, isDemo)
	if err != nil {
		return err
	}

	// Валидируем каждый вопрос
	for _, question := range questions {
		var answer *string

		// Приводим question к нужному типу
		if q, ok := question.(*domain.SurveyQuestion); ok {
			answer = q.Answer
		} else if tempQ, ok := question.(*domain.SurveyQuestionTemp); ok {
			answer = tempQ.Answer
		} else {
			return errors.New("invalid question type")
		}

		// Валидируем вопрос
		if err := s.ValidateAnswer(question, answer, interviewID, isDemo); err != nil {
			return err
		}
	}

	// Завершаем интервью
	if err := s.surveyRepo.FinishInterview(interviewID, time.Now(), isDemo); err != nil {
		return err
	}

	return nil
}

func (s *SurveyService) GetSurveyStats(surveyID int) (*domain.SurveyStats, error) {
	stats, err := s.surveyRepo.GetSurveyStats(surveyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get survey stats: %w", err)
	}

	// Вычисляем среднее время прохождения анкеты
	var totalDuration time.Duration
	var completedCount int
	for _, interviewTime := range stats.InterviewTimes {
		if !interviewTime.EndTime.IsZero() && !interviewTime.StartTime.IsZero() {
			totalDuration += interviewTime.EndTime.Sub(interviewTime.StartTime)
			completedCount++
		}
	}

	// Вычисляем среднее время в секундах
	var averageCompletionTime float64
	if completedCount > 0 {
		averageCompletionTime = totalDuration.Seconds() / float64(completedCount)
	}

	// Устанавливаем среднее время в статистику
	stats.AverageCompletionTime = averageCompletionTime

	return stats, nil
}

func (s *SurveyService) AddEditAccess(surveyID int, email string) error {
	return s.surveyRepo.AddEditAccess(surveyID, email)
}

// RemoveEditAccess removes 'edit' access from a user for the survey
func (s *SurveyService) RemoveEditAccess(surveyID int, email string, currentUserID int, surveyAuthorEmail string) error {
	return s.surveyRepo.RemoveEditAccess(surveyID, email, currentUserID, surveyAuthorEmail)
}

// GetAccessList retrieves the list of users with 'edit' access to the survey
func (s *SurveyService) GetAccessList(surveyID int) ([]string, error) {
	return s.surveyRepo.GetAccessEmails(surveyID)
}
