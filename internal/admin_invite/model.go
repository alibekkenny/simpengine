package admininvite

import (
	"time"
)

type AdminInvite struct {
	ID        int
	Token     string
	CreatedAt time.Time
	ExpiresAt time.Time
	UsedAt    time.Time
	CreatedBy int
	UsedBy    int
}
