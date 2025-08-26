package simptarget

import "context"

type SimpTargetRepository interface {
	CreateSimpTarget(ctx context.Context, name, description string, userID int64) (int64, error)
	UpdateSimpTarget(ctx context.Context, id int64, name, description string, userID int64) error
	DeleteSimpTarget(ctx context.Context, id int64, userID int64) error
	FindAllByUserID(ctx context.Context, userID int64) ([]*SimpTarget, error)
	FindByIDAndUserID(ctx context.Context, id, userID int64) (*SimpTarget, error)
	FindByID(ctx context.Context, id int64) (*SimpTarget, error)
}
