package config

import "database/sql"

type Config struct {
	JWTSecret []byte
	DB        *sql.DB
}
