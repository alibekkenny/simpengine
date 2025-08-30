package model

type EventStepOption struct {
	ID          int64  `json:"id"`
	Label       string `json:"label"`
	ImgID       string `json:"img_id"`
	EventStepID int64  `json:"-"`
}
