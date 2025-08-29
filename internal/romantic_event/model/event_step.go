package model

type EventStep struct {
	ID          int64              `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	StepOrder   int                `json:"step_order"`
	EventID     int64              `json:"-"`
	Options     []*EventStepOption `json:"options"`
}
