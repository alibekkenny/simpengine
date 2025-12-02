package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	rmodel "github.com/alibekkenny/simpengine/internal/romantic_event/model"
	"github.com/alibekkenny/simpengine/internal/shared/db"
	"github.com/alibekkenny/simpengine/internal/shared/model"
	"github.com/lib/pq"
)

type EventStepRepository struct {
	db *sql.DB
}

func NewEventStepRepository(db *sql.DB) EventStepRepository {
	return EventStepRepository{db: db}
}

// CreateEventStep(ctx context.Context, title, description string, eventID int64) (int64, error)
func (r EventStepRepository) CreateEventStep(ctx context.Context, title, descripion string, stepOrder int32, eventID int64) (int64, error) {
	var id int64
	stmt := `INSERT INTO event_steps(title, description, step_order, event_id) 
	VALUES($1, $2, $3, $4) RETURNING id`

	if err := r.db.QueryRowContext(ctx, stmt, title, descripion, stepOrder, eventID).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

// UpdateEventStep(ctx context.Context, id int64, title, description string) error
func (r EventStepRepository) UpdateEventStep(ctx context.Context, id int64, title, description string, stepOrder int32) error {
	stmt := `UPDATE event_steps SET title = $1, description = $2, step_order = $3 WHERE id = $4`

	rows, err := r.db.ExecContext(ctx, stmt, title, description, stepOrder, id)
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

// DeleteEventStep(ctx context.Context, id int64) error
func (r EventStepRepository) DeleteEventStep(ctx context.Context, id int64) error {
	stmt := `DELETE FROM event_steps WHERE id = $1`

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

// FindAllByEventID(ctx context.Context, eventID int64) ([]*models.EventStep, error)
func (r EventStepRepository) FindAllByEventID(ctx context.Context, eventID int64) ([]*rmodel.EventStep, error) {
	stmt := `SELECT id, title, description, event_id, step_order FROM event_steps WHERE event_id = $1 ORDER BY step_order ASC`

	rows, err := r.db.QueryContext(ctx, stmt, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	steps := []*rmodel.EventStep{}

	for rows.Next() {
		var step rmodel.EventStep
		if err := rows.Scan(&step.ID, &step.Title, &step.Description, &step.EventID, &step.StepOrder); err != nil {
			return nil, err
		}

		steps = append(steps, &step)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return steps, nil
}

func (r EventStepRepository) FindByIDandEventID(ctx context.Context, id, eventID int64) (*rmodel.EventStep, error) {
	var step rmodel.EventStep
	stmt := `SELECT id, title, description, event_id, step_order FROM event_steps WHERE id = $1 AND event_id = $2`

	if err := r.db.QueryRowContext(ctx, stmt, id, eventID).Scan(&step.ID, &step.Title, &step.Description, &step.EventID, &step.StepOrder); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNoRecord
		}
		return nil, err
	}

	return &step, nil
}

func (r EventStepRepository) CreateTemplateEventStep(ctx context.Context, title, description string) (int64, error) {
	var id int64
	stmt := `INSERT INTO event_steps(title, description)
	VALUES($1, $2) RETURNING id`

	if err := r.db.QueryRowContext(ctx, stmt, title, description).Scan(&id); err != nil {
		return 0, db.MapPQError(err)
	}

	return id, nil
}

func (r EventStepRepository) UpdateTemplateEventStep(ctx context.Context, id int64, title, description string) error {
	stmt := `UPDATE event_steps SET title = $1, description = $2 WHERE id = $3`

	rows, err := r.db.ExecContext(ctx, stmt, title, description, id)
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

func (r EventStepRepository) FindByID(ctx context.Context, id int64) (*rmodel.EventStep, error) {
	var step rmodel.EventStep
	stmt := `SELECT id, title, description FROM event_steps 
	WHERE id = $1 AND event_id IS NULL`

	if err := r.db.QueryRowContext(ctx, stmt, id).Scan(&step.ID, &step.Title, &step.Description); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNoRecord
		}
		return nil, err
	}

	return &step, nil
}

func (r EventStepRepository) FindAllTemplates(ctx context.Context) ([]*rmodel.EventStep, error) {
	stmt := `SELECT id, title, description FROM event_steps WHERE event_id IS NULL`

	rows, err := r.db.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	steps := []*rmodel.EventStep{}

	for rows.Next() {
		var step rmodel.EventStep
		if err := rows.Scan(&step.ID, &step.Title, &step.Description); err != nil {
			return nil, err
		}

		steps = append(steps, &step)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return steps, nil
}

func (r EventStepRepository) CreateEventStepMany(ctx context.Context, steps []*rmodel.EventStep, eventID int64) ([]*rmodel.EventStep, error) {
	stmt := `INSERT INTO event_steps(title, description, step_order, event_id)
			VALUES %s
			RETURNING id, title, description, step_order, event_id`

	var (
		values []string
		args   []interface{}
	)
	argIndex := 1

	for _, step := range steps {
		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d)", argIndex, argIndex+1, argIndex+2, argIndex+3))
		args = append(args, step.Title, step.Description, step.StepOrder, eventID)

		argIndex += 4
	}

	query := fmt.Sprintf(stmt, strings.Join(values, ","))
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, db.MapPQError(err)
	}
	defer rows.Close()

	created := []*rmodel.EventStep{}
	for rows.Next() {
		var s rmodel.EventStep
		if err := rows.Scan(&s.ID, &s.Title, &s.Description, &s.StepOrder, &s.EventID); err != nil {
			return nil, db.MapPQError(err)
		}

		created = append(created, &s)
	}

	if rows.Err() != nil {
		return nil, db.MapPQError(rows.Err())
	}

	return created, nil
}

func (r EventStepRepository) CreateAnswers(ctx context.Context, choices []*rmodel.EventStepChoice, eventID int64) error {
	stmt := `INSERT INTO romantic_event_choices(event_id, step_id, options_ids)
			VALUES %s`

	var (
		values []string
		args   []interface{}
	)
	argIndex := 1

	for _, step := range choices {
		values = append(values, fmt.Sprintf("($%d, $%d, $%d)", argIndex, argIndex+1, argIndex+2))
		args = append(args, eventID, step.StepID, pq.Array(step.OptionIDs))

		argIndex += 3
	}

	query := fmt.Sprintf(stmt, strings.Join(values, ","))
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return db.MapPQError(err)
	}

	if rows.Err() != nil {
		return db.MapPQError(rows.Err())
	}

	return nil
}
