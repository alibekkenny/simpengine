package model

type EventStep struct {
	ID          int64              `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	StepOrder   int32              `json:"step_order"`
	EventID     int64              `json:"-"`
	Options     []*EventStepOption `json:"options"`
}

type EventStepChoice struct {
	ID        int64   `json:"id"`
	EventID   int64   `json:"event_id"`
	StepID    int64   `json:"step_id"`
	OptionIDs []int64 `json:"option_ids"`
}
