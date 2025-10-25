package romanticevent

import (
	"time"

	"github.com/alibekkenny/simpengine/internal/romantic_event/model"
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
	PublishedAt  *time.Time         `json:"published_at"`
	PublicToken  *string            `json:"public_token"`
	SimpTargetID int64              `json:"simp_target_id"`
	Steps        []*model.EventStep `json:"steps"`
}

type PublishRomanticEventResponseDTO struct {
	Status model.RomanticEventStatus `json:"status"`
	Token  string                    `json:"token"`
}

type EventStepRequestDTO struct {
	Title       string                 `json:"title" validate:"required,min=3"`
	Description string                 `json:"description" validate:"required"`
	StepOrder   int32                  `json:"step_order" validate:"required"`
	Options     []StepOptionRequestDTO `json:"options" validate:"required"`
}

type EventStepsRequestDTO struct {
	Steps []EventStepRequestDTO `json:"steps" validate:"required"`
}

type EventStepResponseDTO struct {
	ID          int64                   `json:"id"`
	Title       string                  `json:"title"`
	Description string                  `json:"description"`
	StepOrder   int32                   `json:"step_order"`
	Options     []StepOptionResponseDTO `json:"options"`
}

type EventStepsResponseDTO struct {
	Steps []EventStepResponseDTO `json:"steps"`
}

type StepOptionRequestDTO struct {
	Label string `json:"label" validate:"required,min=3"`
	ImgID int64  `json:"img_id" validate:"required"`
}

type StepOptionResponseDTO struct {
	ID    int64  `json:"id"`
	Label string `json:"label"`
	ImgID int64  `json:"img_id"`
}

type TemplateEventStepRequestDTO struct {
	Title       string `json:"title" validate:"required,min=3"`
	Description string `json:"description" validate:"required"`
}

type TemplateEventStepResponseDTO struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ViewTemplateEventStepResponseDTO struct {
	ID          int64                    `json:"id"`
	Title       string                   `json:"title"`
	Description string                   `json:"description"`
	Options     []*StepOptionResponseDTO `json:"options"`
}
