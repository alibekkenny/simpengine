package repository

import (
	"context"

	model "github.com/alibekkenny/simpengine/internal/romantic_event/model"
)

type EventStepOptionRepository interface {
	CreateEventStepOption(ctx context.Context, label string, imgID int64, eventStepID int64) (int64, error)
	UpdateEventStepOption(ctx context.Context, id int64, label string, imgID int64) error
	DeleteEventStepOption(ctx context.Context, id int64) error
	FindAllByEventStepID(ctx context.Context, stepID int64) ([]*model.EventStepOption, error)
	FindAllByEventStepIDs(ctx context.Context, stepIDs []int64) (map[int64][]*model.EventStepOption, error)
	FindAllByUserID(ctx context.Context, userID int64) ([]*model.EventStepOption, error)
	CreateEventStepOptionMany(ctx context.Context, options []*model.EventStepOption, stepID int64) ([]*model.EventStepOption, error)
}
