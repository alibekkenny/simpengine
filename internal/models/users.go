package models

import (
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id           int
	Login        string
	Email        string
	PasswordHash string
	Role         string
	CreatedAt    time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(login string, email string, password string, role string) (int, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	var id int
	stmt := `INSERT INTO users(login, email, password_hash, role, created_at) 
			VALUES($1, $2, $3, $4, NOW()) RETURNING id`

	err = m.DB.QueryRow(stmt, login, email, string(passwordHash), role).Scan(&id)
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *UserModel) Get(id int) (*User, error) {
	user := &User{}

	stmt := `SELECT id, login, email, password_hash, role, created_at FROM users WHERE id = $1`
	err := m.DB.QueryRow(stmt, id).Scan(&user.Id, &user.Login, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return user, nil
}

func (m *UserModel) Latest(user *User) ([]*User, error) {
	stmt := `SELECT id, login, email, password_hash, role, created_at FROM users 
			ORDER BY created_at DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// We defer rows.Close() to ensure the sql.Rows resultset is
	// always properly closed before the Latest() method returns.
	// This defer // statement should come *after* you check for an error from the Query() // method.
	// Otherwise, if Query() returns an error, you'll get a panic trying to close a nil resultset.
	defer rows.Close()

	users := []*User{}
	for rows.Next() {
		u := &User{}
		err = rows.Scan(&u.Id, &u.Login, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
