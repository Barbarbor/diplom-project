package services

import (
	"backend/internal/models"
	"crypto/rand"

	"fmt"

	"time"

	"math/big"

	"github.com/jmoiron/sqlx"
)

// Допустимые символы: латинские буквы (верхний и нижний регистр)
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// GenerateRandomHash генерирует случайную строку длиной n символов из letterBytes
func GenerateRandomHash(n int) (string, error) {
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

// CreateSurvey создает новый опрос с дефолтными значениями
func CreateSurvey(db *sqlx.DB, authorID int) (*models.Survey, error) {
	now := time.Now()
	title := fmt.Sprintf("Опрос от %s", now.Format("02.01.2006"))
	hash, err := GenerateRandomHash(15)
	if err != nil {
		return nil, fmt.Errorf("failed to generate survey hash: %w", err)
	}
	state := models.SurveyStateDraft

	var surveyID int
	query := `
		INSERT INTO surveys (title, author_id, hash, state, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id`
	if err := db.QueryRow(query, title, authorID, hash, state).Scan(&surveyID); err != nil {
		return nil, fmt.Errorf("failed to create survey: %w", err)
	}

	survey := &models.Survey{
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

// GetSurveyByHash возвращает опрос по hash вместе с email автора
func GetSurveyByHash(db *sqlx.DB, hash string) (*models.Survey, string, error) {
	// Выполняем JOIN с таблицей пользователей, чтобы получить email автора
	var survey models.Survey
	var email string
	query := `
		SELECT s.id, s.title, s.created_at, s.updated_at, s.hash, s.state, u.email
		FROM surveys s
		JOIN users u ON s.author_id = u.id
		WHERE s.hash = $1`
	if err := db.QueryRowx(query, hash).Scan(&survey.ID, &survey.Title, &survey.CreatedAt, &survey.UpdatedAt, &survey.Hash, &survey.State, &email); err != nil {
		return nil, "", err
	}
	return &survey, email, nil
}
