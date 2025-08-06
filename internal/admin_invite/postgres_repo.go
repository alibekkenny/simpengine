package admininvite

import (
	"context"
	"database/sql"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPosgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, invite AdminInvite) (int, error) {
	var id int
	stmt := `INSERT INTO admin_invite_tokens(token, created_at, expires_at, created_by) 
			VALUES($1, now(), $2, $3) RETURNING id`

	err := r.db.QueryRowContext(ctx, stmt, invite.Token, invite.ExpiresAt, invite.CreatedBy).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *PostgresRepository) GetInviteByToken(ctx context.Context, token string) (*AdminInvite, error) {
	invite := &AdminInvite{}
	stmt := `SELECT id, token, created_at, expires_at, used_at, created_by, used_by FROM admin_invite_tokens WHERE token = $1`

	err := r.db.QueryRowContext(ctx, stmt, token).Scan(&invite.ID, &invite.Token, &invite.CreatedAt, &invite.ExpiresAt, &invite.UsedAt, &invite.CreatedBy, &invite.UsedBy)
	if err != nil {
		return nil, err
	}

	return invite, nil
}

func (r *PostgresRepository) MarkInviteAsUsedByToken(ctx context.Context, token string, usedById int) error {
	var id int
	stmt := `SELECT id FROM admin_invite_token WHERE token = $1 AND used_at IS NULL`
	err := r.db.QueryRowContext(ctx, stmt, token).Scan(&id)
	if err != nil {
		return err
	}

	stmt = `UPDATE admin_invite_tokens SET usedBy = $1, used_at = now() WHERE id = $2`
	_, err = r.db.ExecContext(ctx, stmt, usedById, id)
	return err
}
