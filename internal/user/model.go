package user

import (
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Login     string    `json:"login"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` //hash
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}
