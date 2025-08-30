package repository

import (
	"context"

	model "github.com/alibekkenny/simpengine/internal/romantic_event/model"
	"github.com/google/uuid"
)

type EventStepOptionRepository interface {
	CreateEventStepOption(ctx context.Context, label string, imgID uuid.UUID, eventStepID int64) (int64, error)
	UpdateEventStepOption(ctx context.Context, id int64, label string, imgID uuid.UUID) error
	DeleteEventStepOption(ctx context.Context, id int64) error
	FindAllByEventStepID(ctx context.Context, stepID int64) ([]*model.EventStepOption, error)
	FindAllByEventStepIDs(ctx context.Context, stepIDs []int64) (map[int64][]*model.EventStepOption, error)
}
