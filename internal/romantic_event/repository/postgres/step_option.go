package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	rmodel "github.com/alibekkenny/simpengine/internal/romantic_event/model"
	"github.com/alibekkenny/simpengine/internal/shared/db"
	"github.com/alibekkenny/simpengine/internal/shared/model"
	"github.com/lib/pq"
)

type EventStepOptionRepository struct {
	db *sql.DB
}

func NewEventStepOptionRepository(db *sql.DB) EventStepOptionRepository {
	return EventStepOptionRepository{db: db}
}

func (r EventStepOptionRepository) CreateEventStepOption(ctx context.Context, label string, imgID int64, eventStepID int64) (int64, error) {
	var id int64
	stmt := `INSERT INTO event_step_options(label, img_id, event_step_id)
	VALUES($1, $2, $3) RETURNING id`

	if err := r.db.QueryRowContext(ctx, stmt, label, imgID, eventStepID).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r EventStepOptionRepository) UpdateEventStepOption(ctx context.Context, id int64, label string, imgID int64) error {
	stmt := `UPDATE event_step_options
	SET label = $1, img_id = $2
	WHERE id = $3`

	rows, err := r.db.ExecContext(ctx, stmt, label, imgID, id)
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

func (r EventStepOptionRepository) DeleteEventStepOption(ctx context.Context, id int64) error {
	stmt := `DELETE FROM event_step_options
	WHERE id = $1`

	rows, err := r.db.ExecContext(ctx, stmt, id)
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

func (r EventStepOptionRepository) FindAllByEventStepID(ctx context.Context, stepID int64) ([]*rmodel.EventStepOption, error) {
	stmt := `SELECT id, label, img_id, event_step_id FROM event_step_options
	WHERE event_step_id = $1`

	rows, err := r.db.QueryContext(ctx, stmt, stepID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	options := []*rmodel.EventStepOption{}

	for rows.Next() {
		var option rmodel.EventStepOption
		if err := rows.Scan(&option.ID, &option.Label, &option.ImgID, &option.EventStepID); err != nil {
			return nil, err
		}

		options = append(options, &option)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return options, nil
}

func (r EventStepOptionRepository) FindAllByEventStepIDs(ctx context.Context, stepIDs []int64) (map[int64][]*rmodel.EventStepOption, error) {
	if len(stepIDs) == 0 {
		return map[int64][]*rmodel.EventStepOption{}, nil
	}

	stmt := `SELECT id, label, img_id, event_step_id FROM event_step_options
	WHERE event_step_id = ANY($1)`

	rows, err := r.db.QueryContext(ctx, stmt, pq.Array(stepIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	optionsByStep := make(map[int64][]*rmodel.EventStepOption)

	for rows.Next() {
		var option rmodel.EventStepOption
		if err := rows.Scan(&option.ID, &option.Label, &option.ImgID, &option.EventStepID); err != nil {
			return nil, err
		}

		optionsByStep[option.EventStepID] = append(optionsByStep[option.EventStepID], &option)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return optionsByStep, nil
}

func (r EventStepOptionRepository) FindAllByUserID(ctx context.Context, userID int64) ([]*rmodel.EventStepOption, error) {
	stmt := `SELECT o.id, o.label, o.img_id, o.event_step_id FROM event_step_options o
	JOIN event_steps s ON s.id = o.event_step_id
    JOIN romantic_events e ON e.id = s.event_id
	WHERE e.user_id = $1`

	rows, err := r.db.QueryContext(ctx, stmt, userID)
	if err != nil {
		return nil, db.MapPQError(err)
	}

	defer rows.Close()

	options := []*rmodel.EventStepOption{}
	for rows.Next() {
		var option rmodel.EventStepOption
		if err := rows.Scan(&option.ID, &option.Label, &option.ImgID, &option.EventStepID); err != nil {
			return nil, db.MapPQError(err)
		}

		options = append(options, &option)
	}

	if err := rows.Err(); err != nil {
		return nil, db.MapPQError(err)
	}

	return options, nil
}

func (r EventStepOptionRepository) CreateEventStepOptionMany(ctx context.Context, options []*rmodel.EventStepOption, stepID int64) ([]*rmodel.EventStepOption, error) {
	stmt := `INSERT INTO event_step_options(label, img_id, event_step_id)
			VALUES %s
			RETURNING id, label, img_id, event_step_id`

	var (
		values []string
		args   []interface{}
	)

	argIndex := 1
	for _, option := range options {
		values = append(values, fmt.Sprintf("($%d, $%d, $%d)", argIndex, argIndex+1, argIndex+2))
		args = append(args, option.Label, option.ImgID, stepID)

		argIndex += 3
	}

	query := fmt.Sprintf(stmt, strings.Join(values, ","))
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, db.MapPQError(err)
	}

	defer rows.Close()

	created := []*rmodel.EventStepOption{}
	for rows.Next() {
		o := rmodel.EventStepOption{}
		if err := rows.Scan(&o.ID, &o.Label, &o.ImgID, &o.EventStepID); err != nil {
			return nil, db.MapPQError(err)
		}

		created = append(created, &o)
	}

	if err := rows.Err(); err != nil {
		return nil, db.MapPQError(err)
	}

	return created, nil
}
