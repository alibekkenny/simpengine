package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/alibekkenny/simpengine/internal/media/model"
	"github.com/alibekkenny/simpengine/internal/shared/db"
	shared_model "github.com/alibekkenny/simpengine/internal/shared/model"
)

type MediaRepository struct {
	db *sql.DB
}

func NewMediaRepository(db *sql.DB) MediaRepository {
	return MediaRepository{db: db}
}

func (r MediaRepository) Create(ctx context.Context, objectName string, originalName string, mimeType string, size int64, userID int64) (int64, error) {
	var id int64
	stmt := `INSERT INTO media(object_name, original_name, mime_type, size, user_id)
	VALUES($1, $2, $3, $4, $5) RETURNING id`

	if err := r.db.QueryRowContext(ctx, stmt, objectName, originalName, mimeType, size, userID).Scan(&id); err != nil {
		return 0, db.MapPQError(err)
	}

	return id, nil
}

func (r MediaRepository) FindByIDAndUserID(ctx context.Context, id int64, userID int64) (*model.Media, error) {
	var media model.Media
	stmt := `SELECT id, object_name, original_name, mime_type, size, user_id FROM media WHERE id = $1 AND user_id = $2`

	if err := r.db.QueryRowContext(ctx, stmt, id, userID).Scan(&media.ID, &media.ObjectName, &media.OriginalName, &media.MimeType, &media.Size, &media.UserID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared_model.ErrNoRecord
		}

		return nil, db.MapPQError(err)
	}

	return &media, nil
}

func (r MediaRepository) Delete(ctx context.Context, id int64, userID int64) error {
	stmt := `DELETE FROM media WHERE id = $1 AND user_id = $2`

	rows, err := r.db.ExecContext(ctx, stmt, id, userID)
	if err != nil {
		return db.MapPQError(err)
	}

	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return db.MapPQError(err)
	}
	if rowsAffected == 0 {
		return shared_model.ErrNoRecord
	}

	return nil
}
