package profile

import (
	models "backend/internal/domain"
	"backend/internal/repositories"
	"context"
	"fmt"
)

// ProfileService определяет бизнес-логику для профилей.
type ProfileService struct {
	repo repositories.ProfileRepository
}

// NewProfileService создаёт новый сервис профилей.
func NewProfileService(repo repositories.ProfileRepository) *ProfileService {
	return &ProfileService{repo: repo}
}

func (s *ProfileService) GetUserProfile(ctx context.Context, userID int) (*models.UserProfile, error) {
	profile, err := s.repo.GetUserProfile(userID)

	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}
	return profile, nil
}

func (s *ProfileService) UpdateUserProfile(ctx context.Context, profile *models.UserProfile) error {
	if err := s.repo.UpdateUserProfile(profile); err != nil {
		return fmt.Errorf("failed to update user profile: %w", err)
	}
	return nil
}
