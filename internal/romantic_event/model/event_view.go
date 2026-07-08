package model

import "time"

type EventView struct {
	ID        int64
	EventID   int64
	VisitorID string
	Device    string
	OS        string
	Browser   string
	IP        string
	CreatedAt time.Time
}

type EventViewSummary struct {
	Device   string    `json:"device"`
	OS       string    `json:"os"`
	Browser  string    `json:"browser"`
	OpenedAt time.Time `json:"opened_at"`
}

type EventViewStats struct {
	Views        int                `json:"views"`
	Opens        int                `json:"opens"`
	LastOpenedAt *time.Time         `json:"last_opened_at"`
	RecentOpens  []EventViewSummary `json:"recent_opens"`
}
