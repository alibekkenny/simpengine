package user

import (
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Login     string    `json:"login"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` //hash
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)
