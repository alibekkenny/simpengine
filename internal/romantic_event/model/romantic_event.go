package model

import "time"

type RomanticEvent struct {
	ID           int64        `json:"id"`
	EventDate    time.Time    `json:"event_date"`
	Title        string       `json:"title"`
	Description  string       `json:"description"`
	SimpTargetID int64        `json:"simp_target_id"`
	UserID       int64        `json:"-"`
	Steps        []*EventStep `json:"-"`
}
