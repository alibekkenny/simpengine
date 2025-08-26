package simptarget

import (
	"context"
	"errors"
	"fmt"

	"github.com/alibekkenny/simpengine/internal/auth"
	"github.com/alibekkenny/simpengine/internal/shared/model"
)

type SimpTargetService struct {
	repo SimpTargetRepository
}

func NewSimpTargetService(r SimpTargetRepository) *SimpTargetService {
	return &SimpTargetService{repo: r}
}

func (s *SimpTargetService) CreateSimpTarget(ctx context.Context, name, description string) (int64, error) {
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return 0, model.ErrInvalidCredentials
	}

	id, err := s.repo.CreateSimpTarget(ctx, name, description, userID)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	return id, nil
}

func (s *SimpTargetService) UpdateSimpTarget(ctx context.Context, id int64, name, description string) error {
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return model.ErrInvalidCredentials
	}

	err := s.repo.UpdateSimpTarget(ctx, id, name, description, userID)
	if err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			return fmt.Errorf("%w: simp target not found", model.ErrNoRecord)
		}
		return fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	return nil
}

func (s *SimpTargetService) DeleteSimpTarget(ctx context.Context, id int64) error {
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return model.ErrInvalidCredentials
	}

	err := s.repo.DeleteSimpTarget(ctx, id, userID)
	if err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			return fmt.Errorf("%w: simp target not found", model.ErrNoRecord)
		}
		return fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	return nil
}

func (s *SimpTargetService) GetSimpTargetsByUserID(ctx context.Context) ([]*SimpTarget, error) {
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return nil, model.ErrInvalidCredentials
	}
	simpTargets, err := s.repo.FindAllByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	return simpTargets, nil
}

func (s *SimpTargetService) GetSimpTargetByIDAndUser(ctx context.Context, id int64) (*SimpTarget, error) {
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return nil, model.ErrInvalidCredentials
	}

	simpTarget, err := s.repo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		fmt.Println(err, errors.Is(err, model.ErrNoRecord))
		if errors.Is(err, model.ErrNoRecord) {
			return nil, fmt.Errorf("%w: simp target not found", model.ErrNoRecord)
		}
		return nil, fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	return simpTarget, nil
}
