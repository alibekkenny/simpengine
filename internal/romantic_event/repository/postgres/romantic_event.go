package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	rmodel "github.com/alibekkenny/simpengine/internal/romantic_event/model"
	"github.com/alibekkenny/simpengine/internal/shared/model"
)

type RomanticEventRepository struct {
	db *sql.DB
}

func NewRomanticEventRepository(db *sql.DB) RomanticEventRepository {
	return RomanticEventRepository{db: db}
}

// CreateRomanticEvent(ctx context.Context, eventDate time.Time, title, description string, simpTargetID, userID int64) (int64, error)
func (r RomanticEventRepository) CreateRomanticEvent(ctx context.Context, eventDate time.Time, title, description string, simpTargetID, userID int64) (int64, error) {
	var id int64
	stmt := `INSERT INTO romantic_events(event_date, title, description, simp_target_id, user_id)
	VALUES($1, $2, $3, $4, $5) RETURNING id`

	if err := r.db.QueryRowContext(ctx, stmt, eventDate, title, description, simpTargetID, userID).Scan(&id); err != nil {
		return 0, err
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
		return err
	}
	if rowsAfected == 0 {
		return model.ErrNoRecord
	}

	return nil
}

// FindByIDAndUserID(ctx context.Context, id, userID int64) (*RomanticEvent, error)
func (r RomanticEventRepository) FindByIDAndUserID(ctx context.Context, id, userID int64) (*rmodel.RomanticEvent, error) {
	var event rmodel.RomanticEvent
	stmt := `SELECT id, event_date, title, description, simp_target_id, user_id FROM romantic_events WHERE id = $1 AND user_id = $2`

	if err := r.db.QueryRowContext(ctx, stmt, id, userID).Scan(&event.ID, &event.EventDate, &event.Title, &event.Description, &event.SimpTargetID, &event.UserID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNoRecord
		}
		return nil, err
	}

	return &event, nil
}

// FindAllByUserID(ctx context.Context, userID int64) ([]*RomanticEvent, error)
func (r RomanticEventRepository) FindAllByUserID(ctx context.Context, userID int64) ([]*rmodel.RomanticEvent, error) {
	stmt := `SELECT id, event_date, title, description, simp_target_id, user_id FROM romantic_events WHERE user_id = $1`

	rows, err := r.db.QueryContext(ctx, stmt, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []*rmodel.RomanticEvent{}

	for rows.Next() {
		var event rmodel.RomanticEvent
		if err := rows.Scan(&event.ID, &event.EventDate, &event.Title, &event.Description, &event.SimpTargetID, &event.UserID); err != nil {
			return nil, err
		}

		events = append(events, &event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
