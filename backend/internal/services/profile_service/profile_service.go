package profile

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"context"
	"fmt"
)

// ProfileService определяет бизнес-логику для профилей.
type ProfileService interface {
	GetUserProfile(ctx context.Context, userID int) (*models.UserProfile, error)
	UpdateUserProfile(ctx context.Context, profile *models.UserProfile) error
}

type profileService struct {
	repo repositories.ProfileRepository
}

// NewProfileService создаёт новый сервис профилей.
func NewProfileService(repo repositories.ProfileRepository) ProfileService {
	return &profileService{repo: repo}
}

func (s *profileService) GetUserProfile(ctx context.Context, userID int) (*models.UserProfile, error) {
	profile, err := s.repo.GetUserProfile(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}
	return profile, nil
}

func (s *profileService) UpdateUserProfile(ctx context.Context, profile *models.UserProfile) error {
	if err := s.repo.UpdateUserProfile(profile); err != nil {
		return fmt.Errorf("failed to update user profile: %w", err)
	}
	return nil
}
