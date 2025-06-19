package interview

import (
	"backend/internal/domain"
	"backend/internal/repositories"
	"time"
)

type InterviewService struct {
	repo       repositories.InterviewRepository
	surveyRepo repositories.SurveyRepository
}

func NewInterviewService(repo repositories.InterviewRepository, surveyRepo repositories.SurveyRepository) *InterviewService {
	return &InterviewService{repo: repo, surveyRepo: surveyRepo}
}

func (s *InterviewService) StartInterview(hash, interviewID string, isDemo bool) error {
	// Check if the survey exists by hash
	surveyID, err := s.surveyRepo.GetSurveyIdByHash(hash, isDemo)

	if err != nil {
		return domain.ErrSurveyNotFound
	}
	// Check if an interview with this ID already exists
	exists, err := s.repo.InterviewExists(interviewID)
	if err != nil {
		return err
	}
	if exists {
		return domain.ErrInterviewAlreadyExists
	}

	// Create a new interview
	interview := &domain.SurveyInterview{
		ID:        interviewID,
		SurveyID:  surveyID,
		UserID:    nil,
		Status:    "in_progress",
		StartTime: time.Now(),
		EndTime:   nil,
		IsDemo:    isDemo,
	}

	return s.repo.CreateInterview(interview)
}
