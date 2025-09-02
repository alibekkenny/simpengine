package repository

import (
	"context"

	"github.com/alibekkenny/simpengine/internal/media/model"
)

type MediaRepository interface {
	Create(ctx context.Context, objectName string, originalName string, mimeType string, size int64, userID int64) (int64, error)
	FindByIDAndUserID(ctx context.Context, id int64, userID int64) (*model.Media, error)
	Delete(ctx context.Context, id int64, userID int64) error
}
