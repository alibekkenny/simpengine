package repository

import (
	"context"
	"time"

	model "github.com/alibekkenny/simpengine/internal/romantic_event/model"
)

type RomaticEventRepository interface {
	CreateRomanticEvent(ctx context.Context, eventDate time.Time, title, description string, status model.RomanticEventStatus, simpTargetID, userID int64) (int64, error)
	UpdateRomanticEvent(ctx context.Context, id int64, eventDate time.Time, title, description string, simpTargetID, userID int64) error
	DeleteRomanticEvent(ctx context.Context, id int64, userID int64) error
	FindByIDAndUserID(ctx context.Context, id, userID int64) (*model.RomanticEvent, error)
	FindAllByUserID(ctx context.Context, userID int64) ([]*model.RomanticEvent, error)
	UpdateStatusAndToken(ctx context.Context, id int64, userID int64, status model.RomanticEventStatus, token string) error
	UpdateStatus(ctx context.Context, id int64, userID int64, status model.RomanticEventStatus) error
	FindByPublicToken(ctx context.Context, token string) (*model.RomanticEvent, error)
}
