package romanticevent

import (
	"time"

	"github.com/google/uuid"
)

type RomanticEventRequestDTO struct {
	EventDate    time.Time `json:"event_date" validate:"required"`
	Title        string    `json:"title" validate:"required,min=3"`
	Description  string    `json:"description" validate:"required"`
	SimpTargetID int64     `json:"simp_target_id" validate:"required"`
}

type EventStepRequestDTO struct {
	Title       string `json:"title" validate:"required,min=3"`
	Description string `json:"description" validate:"required"`
	StepOrder   int32  `json:"step_order" validate:"required"`
}

type StepOptionRequestDTO struct {
	Label       string    `json:"label" validate:"required,min=3"`
	Description string    `json:"description" validate:"required"`
	ImgID       uuid.UUID `json:"img_id" validate:"required"`
}
