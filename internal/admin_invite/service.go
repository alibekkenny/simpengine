package admininvite

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type AdminInviteService struct {
	repo AdminInviteRepository
}

func NewAdminInviteService(repo AdminInviteRepository) *AdminInviteService {
	return &AdminInviteService{repo: repo}
}

func (s *AdminInviteService) GenerateInvite(ctx context.Context) (string, error) {
	token := uuid.New().String()
	invite := AdminInvite{
		Token:     token,
		ExpiresAt: time.Now().Add(2 * 24 * time.Hour),
	}
	_, err := s.repo.Create(ctx, invite)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AdminInviteService) ValidateInvite(ctx context.Context, token string) error {
	invite, err := s.repo.GetInviteByToken(ctx, token)
	if err != nil {
		return err
	}
	if invite.UsedBy != 0 || time.Now().After(invite.ExpiresAt) {
		return errors.New("token invalid or expired")
	}

	return nil
}

func (s *AdminInviteService) UseInvite(ctx context.Context, token string, usedById int) error {
	return s.repo.MarkInviteAsUsedByToken(ctx, token, usedById)
}
