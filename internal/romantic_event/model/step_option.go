package model

type EventStepOption struct {
	ID          int64  `json:"id"`
	Label       string `json:"label"`
	Description string `json:"description"`
	ImgID       string `json:"img_id"`
	EventStepID int64  `json:"-"`
}
