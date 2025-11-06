package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	rmodel "github.com/alibekkenny/simpengine/internal/romantic_event/model"
	"github.com/alibekkenny/simpengine/internal/shared/db"
	"github.com/alibekkenny/simpengine/internal/shared/model"
)

type RomanticEventRepository struct {
	db *sql.DB
}

func NewRomanticEventRepository(db *sql.DB) RomanticEventRepository {
	return RomanticEventRepository{db: db}
}

// CreateRomanticEvent(ctx context.Context, eventDate time.Time, title, description string, simpTargetID, userID int64) (int64, error)
func (r RomanticEventRepository) CreateRomanticEvent(ctx context.Context, eventDate time.Time, title, description string, status rmodel.RomanticEventStatus, simpTargetID, userID int64) (int64, error) {
	var id int64
	stmt := `INSERT INTO romantic_events(event_date, title, description, status, simp_target_id, user_id)
	VALUES($1, $2, $3, $4, $5, $6) RETURNING id`

	if err := r.db.QueryRowContext(ctx, stmt, eventDate, title, description, status, simpTargetID, userID).Scan(&id); err != nil {
		return 0, db.MapPQError(err)
	}

	return id, nil
}

// UpdateRomanticEvent(ctx context.Context, eventDate time.Time, title, description string, simpTargetID, userID int64) error
func (r RomanticEventRepository) UpdateRomanticEvent(ctx context.Context, id int64, eventDate time.Time, title, description string, simpTargetID, userID int64) error {
	stmt := `UPDATE romantic_events
	SET event_date = $1, title = $2, description = $3, simp_target_id = $4, user_id = $5
	WHERE id = $6`

	rows, err := r.db.ExecContext(ctx, stmt, eventDate, title, description, simpTargetID, userID, id)
	if err != nil {
		return err
	}

	rowsAfected, err := rows.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAfected == 0 {
		return model.ErrNoRecord
	}

	return nil
}

func (r RomanticEventRepository) UpdateStatusAndToken(ctx context.Context, id int64, userID int64, status rmodel.RomanticEventStatus, token string) error {
	stmt := `UPDATE romantic_events
	SET status = $1, public_token = $2
	WHERE id = $3 AND user_id = $4`

	rows, err := r.db.ExecContext(ctx, stmt, status, token, id, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return model.ErrNoRecord
	}

	return nil
}

func (r RomanticEventRepository) UpdateStatus(ctx context.Context, id int64, userID int64, status rmodel.RomanticEventStatus) error {
	stmt := `UPDATE romantic_events
	SET status = $1
	WHERE id = $2 AND user_id = $3`

	rows, err := r.db.ExecContext(ctx, stmt, status, id, userID)
	if err != nil {
		return db.MapPQError(err)
	}

	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return db.MapPQError(err)
	}
	if rowsAffected == 0 {
		return model.ErrNoRecord
	}

	return nil
}

// DeleteRomanticEvent(ctx context.Context, id int64) error
func (r RomanticEventRepository) DeleteRomanticEvent(ctx context.Context, id int64, userID int64) error {
	stmt := `DELETE FROM romantic_events
	WHERE id = $1 AND user_id = $2`

	rows, err := r.db.ExecContext(ctx, stmt, id, userID)
	if err != nil {
		return err
	}

	rowsAfected, err := rows.RowsAffected()
	if err != nil {
		return db.MapPQError(err)
	}
	if rowsAfected == 0 {
		return model.ErrNoRecord
	}

	return nil
}

// FindByIDAndUserID(ctx context.Context, id, userID int64) (*RomanticEvent, error)
func (r RomanticEventRepository) FindByIDAndUserID(ctx context.Context, id, userID int64) (*rmodel.RomanticEvent, error) {
	var event rmodel.RomanticEvent
	stmt := `SELECT id, event_date, title, description, status, public_token, published_at, simp_target_id, user_id FROM romantic_events WHERE id = $1 AND user_id = $2`

	if err := r.db.QueryRowContext(ctx, stmt, id, userID).Scan(&event.ID, &event.EventDate, &event.Title, &event.Description, &event.Status, &event.PublicToken, &event.PublishedAt, &event.SimpTargetID, &event.UserID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNoRecord
		}
		return nil, db.MapPQError(err)
	}

	return &event, nil
}

// FindAllByUserID(ctx context.Context, userID int64) ([]*RomanticEvent, error)
func (r RomanticEventRepository) FindAllByUserID(ctx context.Context, userID int64) ([]*rmodel.RomanticEvent, error) {
	stmt := `SELECT id, event_date, title, description, status, public_token, published_at, simp_target_id, user_id FROM romantic_events WHERE user_id = $1`

	rows, err := r.db.QueryContext(ctx, stmt, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []*rmodel.RomanticEvent{}

	for rows.Next() {
		var event rmodel.RomanticEvent
		if err := rows.Scan(&event.ID, &event.EventDate, &event.Title, &event.Description, &event.Status, &event.PublicToken, &event.PublishedAt, &event.SimpTargetID, &event.UserID); err != nil {
			return nil, err
		}

		events = append(events, &event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (r RomanticEventRepository) FindByPublicToken(ctx context.Context, token string) (*rmodel.RomanticEvent, error) {
	var event rmodel.RomanticEvent
	stmt := `SELECT id, event_date, title, description, status FROM romantic_events WHERE public_token = $1`

	if err := r.db.QueryRowContext(ctx, stmt, token).Scan(&event.ID, &event.EventDate, &event.Title, &event.Description, &event.Status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNoRecord
		}
		return nil, db.MapPQError(err)
	}

	return &event, nil
}
