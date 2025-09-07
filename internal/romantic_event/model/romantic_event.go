package model

import "time"

type RomanticEventStatus string

const (
	StatusDraft     RomanticEventStatus = "draft"
	StatusPublished RomanticEventStatus = "published"
	StatusAccepted  RomanticEventStatus = "accepted"
	StatusRejected  RomanticEventStatus = "rejected"
	StatusConfirmed RomanticEventStatus = "confirmed"
	StatusArchived  RomanticEventStatus = "archived"
)

type RomanticEvent struct {
	ID           int64               `json:"id"`
	EventDate    time.Time           `json:"event_date"`
	Title        string              `json:"title"`
	Status       RomanticEventStatus `json:"status"`
	Description  string              `json:"description"`
	PublicToken  string              `json:"public_token"`
	PublishedAt  time.Time           `json:"published_at"`
	SimpTargetID int64               `json:"simp_target_id"`
	UserID       int64               `json:"-"`
	Steps        []*EventStep        `json:"-"`
}
