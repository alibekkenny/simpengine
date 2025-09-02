package model

import "time"

type Media struct {
	ID           int64  `json:"id"`
	OriginalName int64  `json:"original_name"`
	Size         int64  `json:"size"`
	MimeType     string `json:"mime_type"`
	ObjectName   string `json:"-"`
	UserID       int64  `json:"-"`
	CreatedAt    time.Time
}
