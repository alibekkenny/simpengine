package model

type EventStepOption struct {
	ID          int64  `json:"id"`
	Label       string `json:"label"`
	ImgID       int64  `json:"img_id"`
	EventStepID int64  `json:"-"`
}
