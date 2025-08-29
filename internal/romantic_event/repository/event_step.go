package repository

import (
	"context"

	model "github.com/alibekkenny/simpengine/internal/romantic_event/model"
)

type EventStepRepository interface {
	CreateEventStep(ctx context.Context, title, description string, stepOrder int32, eventID int64) (int64, error)
	UpdateEventStep(ctx context.Context, id int64, title, description string, stepOrder int32) error
	DeleteEventStep(ctx context.Context, id int64) error
	FindAllByEventID(ctx context.Context, eventID int64) ([]*model.EventStep, error)
	FindByIDandEventID(ctx context.Context, id int64, eventID int64) (*model.EventStep, error)
}
