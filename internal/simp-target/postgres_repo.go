package simptarget

import (
	"context"
	"database/sql"
	"errors"

	"github.com/alibekkenny/simpengine/internal/shared/model"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPosgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateSimpTarget(ctx context.Context, name, description string, userID int64) (int64, error) {
	var id int64
	stmt := `INSERT INTO simp_targets(name, description, user_id)
	VALUES($1, $2, $3) RETURNING id`

	err := r.db.QueryRowContext(ctx, stmt, name, description, userID).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *PostgresRepository) UpdateSimpTarget(ctx context.Context, id int64, name, description string, userID int64) error {
	stmt := `UPDATE simp_target SET name = $1, DESCRIPTION = $2 WHERE id = $3 AND user_id = $4`

	row, err := r.db.ExecContext(ctx, stmt, name, description, id, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := row.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return model.ErrNoRecord
	}

	return nil
}

func (r *PostgresRepository) DeleteSimpTarget(ctx context.Context, id int64, userID int64) error {
	stmt := `DELETE FROM simp_targets WHERE id = $1 AND user_id = $2`

	row, err := r.db.ExecContext(ctx, stmt, id, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := row.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return model.ErrNoRecord
	}

	return nil
}

func (r *PostgresRepository) FindAllByUserID(ctx context.Context, userID int64) ([]*SimpTarget, error) {
	stmt := `SELECT id, name, description FROM simp_targets WHERE user_id = $1`
	result, err := r.db.QueryContext(ctx, stmt, userID)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	simpTargets := []*SimpTarget{}

	for result.Next() {
		var simpTarget SimpTarget
		err := result.Scan(&simpTarget.ID, &simpTarget.Name, &simpTarget.Description)
		if err != nil {
			return nil, err
		}
		simpTargets = append(simpTargets, &simpTarget)
	}

	return simpTargets, nil
}

func (r *PostgresRepository) FindByIDAndUserID(ctx context.Context, id, userID int64) (*SimpTarget, error) {
	var target SimpTarget
	stmt := `SELECT id, name, description FROM simp_targets WHERE id = $1 AND user_id = $2`

	err := r.db.QueryRowContext(ctx, stmt, id, userID).Scan(&target.ID, &target.Name, &target.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNoRecord
		}
		return nil, err
	}

	return &target, nil
}

func (r *PostgresRepository) FindByID(ctx context.Context, id int64) (*SimpTarget, error) {
	var target SimpTarget
	stmt := `SELECT id, name, description FROM simp_targets WHERE id = $1`

	err := r.db.QueryRowContext(ctx, stmt, id).Scan(&target.ID, &target.Name, &target.Description)
	if err != nil {
		return nil, err
	}

	return &target, nil
}
