package repository

import (
	"context"
	"time"

	model "github.com/alibekkenny/simpengine/internal/romantic_event/model"
)

type EventViewRepository interface {
	InsertView(ctx context.Context, v model.EventView) error
	LastViewAt(ctx context.Context, eventID int64, visitorID string) (*time.Time, error)
	StatsByEventID(ctx context.Context, eventID int64, recentLimit int) (model.EventViewStats, error)
}
