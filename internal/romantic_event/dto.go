package romanticevent

import (
	"time"

	"github.com/alibekkenny/simpengine/internal/romantic_event/model"
	"github.com/google/uuid"
)

type RomanticEventRequestDTO struct {
	EventDate    time.Time `json:"event_date" validate:"required"`
	Title        string    `json:"title" validate:"required,min=3"`
	Description  string    `json:"description" validate:"required"`
	SimpTargetID int64     `json:"simp_target_id" validate:"required"`
}

type RomanticEventResponseDTO struct {
	ID           int64              `json:"id"`
	EventDate    time.Time          `json:"event_date"`
	Title        string             `json:"title"`
	Description  string             `json:"description"`
	SimpTargetID int64              `json:"simp_target_id"`
	Steps        []*model.EventStep `json:"steps"`
}

type EventStepRequestDTO struct {
	Title       string `json:"title" validate:"required,min=3"`
	Description string `json:"description" validate:"required"`
	StepOrder   int32  `json:"step_order" validate:"required"`
}

type StepOptionRequestDTO struct {
	Label string    `json:"label" validate:"required,min=3"`
	ImgID uuid.UUID `json:"img_id" validate:"required"`
}
