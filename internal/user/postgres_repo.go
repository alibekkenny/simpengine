package user

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

func (r *PostgresRepository) Create(ctx context.Context, user *User) (int, error) {
	var id int
	stmt := `INSERT INTO users(login, email, password_hash, role, created_at) 
			VALUES($1, $2, $3, $4, NOW()) RETURNING id`

	err := r.db.QueryRowContext(ctx, stmt, user.Login, user.Email, user.Password, user.Role).Scan(&id)
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *PostgresRepository) FindById(ctx context.Context, id int) (*User, error) {
	user := &User{}
	stmt := `SELECT id, login, email, password_hash, role, created_at FROM users WHERE id = $1`

	err := r.db.QueryRowContext(ctx, stmt, id).Scan(&user.ID, &user.Login, &user.Email, &user.Password, &user.Role, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return user, nil
}

func (r *PostgresRepository) FindByLogin(ctx context.Context, login string) (*User, error) {
	user := &User{}
	stmt := `SELECT id, login, email, password_hash, role, created_at FROM users WHERE login = $1`

	err := r.db.QueryRowContext(ctx, stmt, login).Scan(&user.ID, &user.Login, &user.Email, &user.Password, &user.Role, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return user, nil
}

func (r *PostgresRepository) ExistsByEmailOrLogin(ctx context.Context, email, login string) (bool, error) {
	var exists bool
	stmt := `SELECT EXISTS (
		SELECT 1 FROM users WHERE email = $1 OR login = $2
	)`

	err := r.db.QueryRowContext(ctx, stmt, email, login).Scan(&exists)
	return exists, err
}
