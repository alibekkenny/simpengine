package model

type EventStepOption struct {
	ID          int64  `json:"id"`
	Label       string `json:"label"`
	Description string `json:"description"`
	ImgID       int64  `json:"img_id,string,omitempty"`
	EventStepID int64  `json:"-"`
}
